package telegram

import (
	"context"
	"fmt"
	"insinyur-radius/domain"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yudapc/go-rupiah"
)

// for the new feature add customer repository and collector repository

func contains(v string, a []domain.Menu) bool {
	for _, m := range a {
		if *m.Name == v {
			return true
		}
	}
	return false
}

func sendMessage(message string, chatID int64, messageID int64, ucme domain.MessageRepository) {
	messaging := domain.Message{}
	messaging.ChatID = &chatID
	messaging.MessageID = &messageID
	messaging.Message = &message

	err := ucme.Insert(context.Background(), messaging)
	if err != nil {
		logrus.Error(err)
	}
}

func validReseller(telegramID int64, ucr domain.ResellerUsecase) (res domain.Reseller, mes string) {
	res, err := ucr.GetWithTelegramID(context.Background(), telegramID)
	if err != nil {
		logrus.Error(err)
		mes = "sepertinya telah terjadi kesalahan, silahkan hubungi beberapa saat lagi..."
		return
	}

	if res == (domain.Reseller{}) {
		mes = "sepertinya anda belum terdaftar sebagai reseller, silahkan lakukan pendataran terlebih dahulu..."
		return
	}

	if *res.Active == "no" {
		mes = "data pendaftaran anda belum ter verifikasi, berikut adalah nomor pendaftaran anda:\n"
		mes += *res.RegisterCode
		mes += "\ntunjukkan kode kepada admin untuk memverifikasi pendaftaran anda"

		res = domain.Reseller{}
		return
	}

	return
}

func register(message string, tSplit string, telegramID int64, chatID int64, ucr domain.ResellerUsecase) (mes string) {
	reseller, _ := ucr.GetWithTelegramID(context.Background(), telegramID)
	if reseller != (domain.Reseller{}) {
		mes = "maaf, perangkat anda telah terdaftar sebelumnya dengan nomor registrasi sebagai berikut:\n\n\n" + *reseller.RegisterCode
		return
	}

	spString := strings.Split(message, tSplit)
	if len(spString) != 3 {
		mes = "maaf, format yang anda masukkan salah...\nformat: \n\n" + tSplit + "daftar" + tSplit + "nama_reseller\n\n"
	} else {
		t := time.Now()

		defaultNo := "no"
		uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
		reseller := domain.Reseller{}
		reseller.Name = &spString[2]
		reseller.TelegramID = &telegramID
		reseller.ChatID = &chatID
		reseller.Active = &defaultNo
		reseller.RegisterCode = &uniqTime

		mes = "perangkat anda telah ditambahkan sebagai reseller, berikut adalah kode registrasi anda\n\n" + *reseller.RegisterCode

		err := ucr.Insert(context.Background(), &reseller)
		if err != nil {
			mes = "sepertinya telah terjadi kesalahan, silahkan hubungi beberapa saat lagi..."
		}
	}
	return
}

func menu(tSplit string, balanceMenu string, telegramID int64, menu []domain.Menu, ucr domain.ResellerUsecase) (mes string) {
	reseller, mess := validReseller(telegramID, ucr)
	if reseller == (domain.Reseller{}) {
		mes = mess
		return
	}

	mes = "daftar menu: \n\n"
	mes += tSplit + balanceMenu + "\n"
	for _, value := range menu {
		mes += *value.Name + "\n"
	}

	return
}

func dynamicMenu(message string, telegramID int64, menus []domain.Menu, uct domain.TransactionUsecase, ucr domain.ResellerUsecase) (mes string) {
	reseller, mess := validReseller(telegramID, ucr)
	if reseller == (domain.Reseller{}) {
		mes = mess
		return
	}

	mes = "sepertinya masa promo untuk pilihan ini telah berakhir, silahkan pilih menu yang lain"
	for _, menu := range menus {
		if message == *menu.Name {
			trc, err := uct.ResellerTransaction(context.Background(), *reseller.ID, *menu.IDPackage, *menu.Profile)
			if err != nil {
				logrus.Error(err)

				mes = "transaksi gagal dilakukan, hubungi admin jika terjadi masalah..."

				switch err {
				case domain.ErrBalanceRequired:
					mes = "saldo anda tidak mencukupi..."
				}
				return
			}

			mes = "Detail Transaksi:\n\n"
			mes += "status:\n\t" + *trc.Status + "\n\n"
			mes += "kode transaksi:\n\t" + *trc.TransactionCode + "\n\n"
			mes += "harga:\n\t" + strconv.FormatInt(*trc.Value, 10) + "\n\n"
			mes += "kode voucher:\n\t" + *trc.Information + "\n\n"
			time := trc.CreatedAt.Format("2006-01-02 15:04:05")
			mes += "tanggal transaksi:\n\t" + time + "\n\n"

			break
		}
	}

	return

}

func balance(message string, telegramID int64, uct domain.TransactionUsecase, ucr domain.ResellerUsecase) (mes string) {
	reseller, mess := validReseller(telegramID, ucr)
	if reseller == (domain.Reseller{}) {
		mes = mess
		return
	}

	balance, err := uct.Balance(context.Background(), *reseller.ID)
	mes = "sisa saldo anda : \n\n" + rupiah.FormatRupiah(float64(balance)) //strconv.FormatInt(balance, 10)
	if err != nil {
		mes = "sepertinya telah terjadi kesalahan, silahkan hubungi beberapa saat lagi..."
	}
	return
}

func refill(message string, tSplit string, telegramID int64, uct domain.TransactionUsecase, ucr domain.ResellerUsecase) (mes string) {
	rs := strings.Split(message, tSplit)

	// timeNow := time.Now()

	// fmt.Println("TIME NOW: ",timeNow.String())



	mes = "pastikan keyword yang anda masukkan benar [/refill/nomor_invoice] contoh: /refill/1234567890"
	if len(rs) == 3 {
		reseller, mess := validReseller(telegramID, ucr)
		if reseller == (domain.Reseller{}) {
			mes = mess
			return
		}

		mesS, err := uct.ResellerRefillTransaction(context.Background(), *reseller.ID, rs[2])
		if err != nil {
			logrus.Error(err)
			mes = "sepertinya telah terjadi kesalahan, silahkan hubungi beberapa saat lagi..."
			if err == domain.ErrNotAccordingSpecifications {
				mes = "aktivasi user belum dilakukan, hubungi admin untuk melakukan verifikasi data..."
			}
			return mes
		}
		mes = mesS
	}

	return
}

// NewHandler ...
func NewHandler(ucr domain.ResellerUsecase, uct domain.TransactionUsecase, ucm domain.MenuUsecase, ucme domain.MessageUsecase) {
	var mtx sync.Mutex

LoopReconnect:
	for true {
		bot, err := tgbotapi.NewBotAPI(viper.GetString("telegram.token"))
		if err != nil {
			logrus.Error(err)

			time.Sleep(5 * time.Second)
			continue LoopReconnect
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			logrus.Error(err)
		}

	LoopMessage:
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue LoopMessage
			}

			// trim whitespace in start and end line
			message := fmt.Sprintf("%s", update.Message.Text)
			message = strings.Trim(message, " ")

			// trim all white space
			arrMessage := strings.Split(message, " ")
			message = strings.Join(arrMessage, "")

			// split message
			tSplit := viper.GetString(`telegram.split`)
			resultSplit := strings.Split(message, tSplit)

			if len(resultSplit) > 1 && string(message[0]) == tSplit {
				listMenu, err := ucm.Fetch(context.Background())
				if err != nil {
					mtx.Lock()

					logrus.Error(err)

					mes := "sepertinya telah terjadi kesalahan, silahkan hubungi beberapa saat lagi..."
					sendMessage(mes, update.Message.Chat.ID, int64(update.Message.MessageID), ucme)

					mtx.Unlock()

					continue LoopMessage
				}

				switch {
				case resultSplit[1] == viper.GetString(`telegram.register`):
					mtx.Lock()
					sendMessage(
						register(message, tSplit, int64(update.Message.From.ID), update.Message.Chat.ID, ucr),
						update.Message.Chat.ID, int64(update.Message.MessageID), ucme,
					)
					mtx.Unlock()
					break
				case resultSplit[1] == viper.GetString(`telegram.customer`):
					mtx.Lock()

					mtx.Unlock()
					break
				case resultSplit[1] == viper.GetString(`telegram.refill`):
					mtx.Lock()
					sendMessage(
						refill(message, tSplit, int64(update.Message.From.ID), uct, ucr),
						update.Message.Chat.ID, int64(update.Message.MessageID), ucme,
					)
					mtx.Unlock()
					break
				case resultSplit[1] == viper.GetString(`telegram.balance`):
					mtx.Lock()
					sendMessage(
						balance(message, int64(update.Message.From.ID), uct, ucr),
						update.Message.Chat.ID, int64(update.Message.MessageID), ucme,
					)
					mtx.Unlock()
					break
				case resultSplit[1] == viper.GetString(`telegram.menu`):
					mtx.Lock()
					sendMessage(
						menu(tSplit, viper.GetString(`telegram.balance`), int64(update.Message.From.ID), listMenu, ucr),
						update.Message.Chat.ID, int64(update.Message.MessageID), ucme,
					)
					mtx.Unlock()
					break
				case contains((tSplit + resultSplit[1]), listMenu):
					mtx.Lock()
					sendMessage(
						dynamicMenu(message, int64(update.Message.From.ID), listMenu, uct, ucr),
						update.Message.Chat.ID, int64(update.Message.MessageID), ucme,
					)
					mtx.Unlock()
					break
				}
			}
		}
	}
}
