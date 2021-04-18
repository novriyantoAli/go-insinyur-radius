package usecase

import (
	"context"
	"insinyur-radius/domain"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
	"github.com/yudapc/go-rupiah"
)

type transactionUsecase struct {
	Timeout            time.Duration
	Repository         domain.TransactionRepository
	RepositoryRadcheck domain.RadcheckRepository
	RepositoryPackage  domain.PackageRepository
	RepositoryReseller domain.ResellerRepository
	RepositoryMessage  domain.MessageRepository
	RepositoryCustomer domain.CustomerRepository
	RepositoryInvoice  domain.InvoiceRepository
}

// NewTransactionUsecase ...
func NewTransactionUsecase(t time.Duration, r domain.TransactionRepository, rr domain.RadcheckRepository, rp domain.PackageRepository, rre domain.ResellerRepository, rm domain.MessageRepository, rc domain.CustomerRepository, ri domain.InvoiceRepository) domain.TransactionUsecase {
	return &transactionUsecase{Timeout: t, Repository: r, RepositoryRadcheck: rr, RepositoryPackage: rp, RepositoryReseller: rre, RepositoryMessage: rm, RepositoryCustomer: rc, RepositoryInvoice: ri}
}

func (uc *transactionUsecase) Fetch(c context.Context) (res []domain.Transaction, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *transactionUsecase) ResellerRefillTransaction(c context.Context, idReseller int64, noInvoice string) (res string, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	// check invoice if exists
	invoice, err := uc.RepositoryInvoice.Find(ctx, noInvoice)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if invoice == (domain.Invoice{}) {
		logrus.Error("item invoice not found")
		return "", domain.ErrNotFound
	}

	/**
	* check if customer has valid
	 */
	customer, err := uc.RepositoryCustomer.GetUsername(ctx, *invoice.Name)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if customer == (domain.Customer{}) {
		logrus.Error("item customer not found")
		return "", domain.ErrNotFound
	}
	/**
	* end check of customer valid
	 */

	/**
	* get package
	 */
	packages, err := uc.RepositoryPackage.Get(ctx, *customer.IDPackage)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if packages == (domain.Package{}) {
		logrus.Error("item package not found")
		return "", domain.ErrNotFound
	}
	// ==

	/**
	* check if reseller have a balance
	 */
	balance, err := uc.Balance(ctx, idReseller)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if balance < *packages.Price {
		logrus.Error("balance not equal or low")
		return "", domain.ErrBalanceRequired
	}
	/**
	* end of check balance
	 */

	/**
	* check if expiration have in radcheck
	 */
	radcheckArr, err := uc.RepositoryRadcheck.FetchWithUsername(ctx, *invoice.Name)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	if len(radcheckArr) < 3 {
		logrus.Error("user not activated")
		return "", domain.ErrNotAccordingSpecifications
	}
	var expiration *string
	for _, r := range radcheckArr {
		if strings.EqualFold(*r.Attribute, "expiration") {
			expiration = r.Value
			break
		}
	}

	if expiration == nil {
		logrus.Error("user not activated")
		return "", domain.ErrNotAccordingSpecifications
	}
	/**
	* end check expiration
	 */

	wita, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		logrus.Error(err)
		return "", domain.ErrInternalServerError
	}

	safeTransactionCode := ""
	for true {
		code, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		transaction, _ := uc.Repository.GetWithTransactionCode(ctx, code)
		if transaction == (domain.Transaction{}) {
			safeTransactionCode += code
			break
		}
	}

	layoutFormat := "02 Jan 2006 15:04:05"
	status := "out"

	transaction := domain.Transaction{}
	transaction.IDReseller = &idReseller
	transaction.TransactionCode = &safeTransactionCode
	transaction.Status = &status
	transaction.Value = packages.Price
	transaction.Information = invoice.Name

	timeNow := ""

	switch *packages.ValidityUnit {
	case "HOUR":
		// time now
		timeNow = time.Now().Local().Add(
			time.Hour*time.Duration(*packages.ValidityValue) +
				time.Minute*time.Duration(0) +
				time.Second*time.Duration(0)).In(wita).Format(layoutFormat)

		err = uc.Repository.ResellerRefillTransaction(ctx, transaction, *invoice.NoInvoice, timeNow)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		break
	case "DAY":
		// time now
		timeNow := time.Now().Local().AddDate(0, 0, int(*packages.ValidityValue)).In(wita).Format(layoutFormat)
		err = uc.Repository.ResellerRefillTransaction(ctx, transaction, *invoice.NoInvoice, timeNow)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		break
	case "MONTH":
		date, err := time.ParseInLocation(layoutFormat, *expiration, wita)
		if err != nil {
			return "", domain.ErrInternalServerError
		}
		timeNow := date.AddDate(0, 1, 0).In(wita).Format(layoutFormat)
		err = uc.Repository.ResellerRefillTransaction(ctx, transaction, *invoice.NoInvoice, timeNow)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		break
	case "YEAR":
		date, err := time.ParseInLocation(layoutFormat, *expiration, wita)
		if err != nil {
			return "", domain.ErrInternalServerError
		}
		timeNow := date.AddDate(1, 0, 0).In(wita).Format(layoutFormat)
		err = uc.Repository.ResellerRefillTransaction(ctx, transaction, *invoice.NoInvoice, timeNow)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		break
	}

	res = "====================\n"
	res += "====================\n"
	res += "Transaksi Berhasil..!\n"
	res += "====================\n"
	res += "Kode transaksi ~>\n\t"
	res += safeTransactionCode + "\n"
	res += "====================\n"
	res += "Nominal transaksi ~>\n\t"
	res += strconv.FormatInt(*packages.Price, 10) + "\n"
	res += "====================\n"
	res += "Nama Paket ~>\n\t"
	res += *packages.Name + "\n"
	res += "====================\n"
	res += "Tambahan Waktu ~>\n\t"
	res += strconv.FormatInt(*packages.ValidityValue, 10) + " " + *packages.ValidityUnit + "\n"
	res += "====================\n"
	res += "Nama pengguna ~>\n\t"
	res += *invoice.Name + "\n"
	res += "====================\n"
	res += "kadarluarsa ~>\n\t"
	res += timeNow + "\n"
	res += "====================\n"
	res += "====================\n"

	return

}
func (uc *transactionUsecase) ResellerTransaction(c context.Context, idReseller int64, idPackage int64, profile string) (res domain.Transaction, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res = domain.Transaction{}

	balance, err := uc.Balance(ctx, idReseller)
	if err != nil {
		return
	}

	pkg, _ := uc.RepositoryPackage.Get(ctx, idPackage)
	if pkg == (domain.Package{}) {
		err = domain.ErrNotFound
		return
	}

	if balance < *pkg.Price {
		err = domain.ErrBalanceRequired
		return
	}

	safeUsername := ""
	for true {
		// search random password
		username, err := password.Generate(8, 3, 0, true, true)
		if err != nil {
			logrus.Error(err)
			return domain.Transaction{}, err
		}
		radchecks, _ := uc.RepositoryRadcheck.FetchWithUsername(ctx, username)
		if len(radchecks) == 0 {
			safeUsername += username
			break
		}
	}

	safeTransactionCode := ""
	for true {
		code, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			logrus.Error(err)
			return domain.Transaction{}, err
		}
		transaction, _ := uc.Repository.GetWithTransactionCode(ctx, code)
		if transaction == (domain.Transaction{}) {
			safeTransactionCode += code
			break
		}
	}

	// before call, use package to detect price
	// call transaction repository
	// and save radpackage, radcheck, transaction

	defaultStatus := "OUT"

	res.IDReseller = &idReseller
	res.TransactionCode = &safeTransactionCode
	res.Information = &safeUsername
	res.Status = &defaultStatus
	res.Value = pkg.Price

	err = uc.Repository.ResellerTransaction(ctx, &res, idPackage, profile)

	if err != nil {
		return domain.Transaction{}, err
	}

	res.IDReseller = nil
	*res.Value = *pkg.Price + *pkg.Margin

	return
}

func (uc *transactionUsecase) Report(c context.Context, dateStart string, dateEnd string) (res []domain.Transaction, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Report(ctx, dateStart, dateEnd)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *transactionUsecase) Refill(c context.Context, idReseller int64, balance int64) (res domain.Transaction, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	safeTransactionCode := ""
	for true {
		code, err := password.Generate(64, 10, 0, false, true)
		if err != nil {
			logrus.Error(err)
			return domain.Transaction{}, err
		}
		transaction, _ := uc.Repository.GetWithTransactionCode(ctx, code)
		if transaction == (domain.Transaction{}) {
			safeTransactionCode = code
			break
		}
	}

	defStatus := "IN"
	res = domain.Transaction{}
	res.IDReseller = &idReseller
	res.Status = &defStatus
	res.TransactionCode = &safeTransactionCode
	res.Value = &balance

	mes := "saldo anda berhasil ditambahkan senilai " + rupiah.FormatRupiah(float64(balance))
	message := domain.Message{}
	message.Message = &mes

	err = uc.Repository.RefillBalance(ctx, res, message)
	if err != nil {
		return domain.Transaction{}, err
	}

	return
}

func (uc *transactionUsecase) Balance(c context.Context, idReseller int64) (res int64, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	resArr, err := uc.Repository.FetchWithIDReseller(ctx, idReseller)
	if err != nil {
		return 0, err
	}

	var in int64 = 0
	var out int64 = 0
	for _, value := range resArr {
		if strings.EqualFold("IN", *value.Status) == true {
			in += *value.Value
			continue
		}
		if strings.EqualFold("OUT", *value.Status) == true {
			out += *value.Value
			continue
		}
	}

	res = in - out

	return
}
