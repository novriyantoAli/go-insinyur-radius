package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_packageHandler "insinyur-radius/backend/package/delivery/http"
	_packageRepository "insinyur-radius/backend/package/repository/mysql"
	_packageUsecase "insinyur-radius/backend/package/usecase"

	_resellerHandler "insinyur-radius/backend/reseller/delivery/http"
	_telegramHandler "insinyur-radius/backend/telegram/delivery/bot"
	_resellerRepository "insinyur-radius/backend/reseller/repository/mysql"
	_resellerUsecase "insinyur-radius/backend/reseller/usecase"

	_transactionHandler "insinyur-radius/backend/transaction/delivery/http"
	_transactionRepository "insinyur-radius/backend/transaction/repository/mysql"
	_transactionUsecase "insinyur-radius/backend/transaction/usecase"

	_menuHandler "insinyur-radius/backend/menu/delivery/http"
	_menuRepository "insinyur-radius/backend/menu/repository/mysql"
	_menuUsecase "insinyur-radius/backend/menu/usecase"

	_radcheckRepository "insinyur-radius/backend/radcheck/repository/mysql"
	// _radcheckUsecase "insinyur-radius/backend/radcheck/usecase"

	_radusergroupHandler "insinyur-radius/backend/radusergroup/delivery/http"
	_radusergroupRepository "insinyur-radius/backend/radusergroup/repository/mysql"
	_radusergroupUsecase "insinyur-radius/backend/radusergroup/usecase"

	_radacctRepository "insinyur-radius/backend/radacct/repository/mysql"
	// _radacctUsecase "insinyur-radius/backend/radacct/usecase"

	_schedulerHandler "insinyur-radius/backend/scheduler/delivery/udp"
	_schedulerUsecase "insinyur-radius/backend/scheduler/usecase"

	_radgroupcheckHandler "insinyur-radius/backend/radgroupcheck/delivery/http"
	_radgroupcheckRepository "insinyur-radius/backend/radgroupcheck/repository/mysql"
	_radgroupcheckUsecase "insinyur-radius/backend/radgroupcheck/usecase"

	_radgroupreplyHandler "insinyur-radius/backend/radgroupreply/delivery/http"
	_radgroupreplyRepository "insinyur-radius/backend/radgroupreply/repository/mysql"
	_radgroupreplyUsecase "insinyur-radius/backend/radgroupreply/usecase"

	_usersHandler "insinyur-radius/backend/users/delivery/http"
	_usersRepository "insinyur-radius/backend/users/repository/mysql"
	_usersUsecase "insinyur-radius/backend/users/usecase"

	_messageHandler "insinyur-radius/backend/message/delivery/telegram"
	_messageRepository "insinyur-radius/backend/message/repository/mysql"
	_messageUsecase "insinyur-radius/backend/message/usecase"

	_nasHandler "insinyur-radius/backend/nas/delivery/http"
	_nasRepository "insinyur-radius/backend/nas/repository/mysql"
	_nasUsecase "insinyur-radius/backend/nas/usecase"

	_radpostauthHandler "insinyur-radius/backend/radpostauth/delivery/http"
	_radpostauthRepository "insinyur-radius/backend/radpostauth/repository/mysql"
	_radpostauthUsecase "insinyur-radius/backend/radpostauth/usecase"

	_invoiceHandler "insinyur-radius/backend/invoice/delivery/http"
	_invoiceRepository "insinyur-radius/backend/invoice/repository/mysql"
	_invoiceUsecase "insinyur-radius/backend/invoice/usecase"

	_paymentHandler "insinyur-radius/backend/payment/delivery/http"
	_paymentRepository "insinyur-radius/backend/payment/repository/mysql"
	_paymentUsecase "insinyur-radius/backend/payment/usecase"

	_reportHandler "insinyur-radius/backend/report/delivery/http"
	_reportRepository "insinyur-radius/backend/report/repository/mysql"
	_reportUsecase "insinyur-radius/backend/report/usecase"

	_customerRepository "insinyur-radius/backend/customer/repository/mysql"
)

type responseError struct {
	Message string `json:"error"`
}

func init() {

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	logrus.SetReportCaller(true)

	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN in DEBUG mode")
	}

	if viper.GetString(`administrator.mode`) == "bussiness" {
		dbConn := createDB()
		triggerName := viper.GetString(`administrator.triggerName`)
		// create all event after insert
		_, err = dbConn.Exec(`
		CREATE TRIGGER ` + triggerName + ` AFTER INSERT ON radacct FOR EACH ROW 
		
		BEGIN
		
		SET @expiration = (SELECT COUNT(*) FROM radcheck WHERE username = New.username AND attribute = 'Expiration'); 
		
		IF (@expiration = 0) THEN
			SET @validity_value = (SELECT package.validity_value FROM radpackage INNER JOIN package ON package_id = radpackage.package.id WHERE radpackage.username = New.username);
			SET @validity_unit = (SELECT package.validity_unit FROM radpackage INNER JOIN package ON package.id = radpackage.package_id WHERE radpackage.username = New.username);

			IF (@validity_unit = 'HOUR') THEN
				INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value HOUR), "%d %b %Y %H:%I:%S"));

			ELSEIF (@validity_unit = 'DAY') THEN
				INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value DAY), "%d %b %Y %H:%I:%S"));

			ELSEIF (@validity_unit = 'MONTH') THEN
				INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value MONTH), "%d %b %Y %H:%I:%S"));

			ELSEIF (@validity_unit = 'YEAR') THEN
				INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_unit YEAR), "%d %b %Y %H:%I:%S"));

			END IF;

		END IF;
		END;`)

		if err != nil {
			logrus.Error(err)
		}

		dbConn.Close()
	} else {
		dbConn := createDB()

		triggerName := viper.GetString("administrator.triggerName")

		// delete all event after insert
		_, err = dbConn.Exec(`DROP TRIGGER ` + triggerName + ` ;`)
		if err != nil {
			logrus.Error(err)
		}

		dbConn.Close()
	}
}

func createDB() *sql.DB {
	// set radacct to check if user logged in
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add(`parseTime`, "1")
	val.Add(`loc`, "Asia/Makassar")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func main() {

	// logging initialize
	f, err := os.OpenFile("goir.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)

	logrus.SetOutput(wrt)
	// log.SetOutput(wrt)

	// database initialize
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add(`parseTime`, "1")
	val.Add(`loc`, "Asia/Makassar")

	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		logrus.Fatalln(err)
		// log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		logrus.Fatalln(err)
		// log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			logrus.Fatalln(err)
			// log.Fatal(err)
		}
	}()

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// authenticate menu
	e.POST("/api/login", func(c echo.Context) error {

		logrus.Infoln("end point /api/login call")

		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" {
			return c.JSON(http.StatusBadRequest, responseError{Message: "username required..."})
		}

		if password == "" {
			return c.JSON(http.StatusBadRequest, responseError{Message: "password required..."})
		}

		if username != viper.GetString(`administrator.username`) {
			return c.JSON(http.StatusForbidden, responseError{Message: "username not valid..."})
		}

		if password != viper.GetString(`administrator.password`) {
			return c.JSON(http.StatusForbidden, responseError{Message: "password not valid..."})
		}

		logrus.Debugln("all information from user valid, sign in token for user")

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte(viper.GetString(`administrator.key`)))

		if err != nil {
			logrus.Error(err)
			return c.JSON(http.StatusInternalServerError, responseError{Message: err.Error()})
		}

		logrus.Debugln("sign in token complete, resulting token to user")

		return c.JSON(http.StatusOK, map[string]string{"token": t})
	})

	timeout := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// define all repository here
	packageRepository := _packageRepository.NewMysqlPackageRepository(dbConn)
	resellerRepository := _resellerRepository.NewMysqlResellerRepository(dbConn)
	transactionRepository := _transactionRepository.NewMysqlTransactionRepository(dbConn)
	menuRepository := _menuRepository.NewMysqlMenuRepository(dbConn)
	radcheckRepository := _radcheckRepository.NewMysqlRadcheckRepository(dbConn)
	radusergroupRepository := _radusergroupRepository.NewMysqlRadusergroupRepository(dbConn)
	radacctRepository := _radacctRepository.NewMysqlRadacctRepository(dbConn)
	radgroupcheckRepository := _radgroupcheckRepository.NewMysqlRadgroupcheckRepository(dbConn)
	radgroupreplyRepository := _radgroupreplyRepository.NewMysqlRadgroupreplyRepository(dbConn)
	usersRepository := _usersRepository.NewMysqlUsersRepository(dbConn)
	messageRepository := _messageRepository.NewMysqlRepository(dbConn)
	nasRepository := _nasRepository.NewMysqlRepository(dbConn)
	radpostauthRepository := _radpostauthRepository.NewMysqlRepository(dbConn)
	invoiceRepository := _invoiceRepository.NewMysqlRepository(dbConn)
	paymentRepository := _paymentRepository.NewMysqlRepository(dbConn)
	reportRepository := _reportRepository.NewMysqlRepository(dbConn)
	customerRepository := _customerRepository.NewMysqlRepository(dbConn)

	// define all usecase here
	messageUsecase := _messageUsecase.NewUsecase(timeout, messageRepository)
	packageUsecase := _packageUsecase.NewPackageUsecase(timeout, packageRepository)
	resellerUsecase := _resellerUsecase.NewResellerUsecase(timeout, resellerRepository)
	transactionUsecase := _transactionUsecase.NewTransactionUsecase(
		timeout,
		transactionRepository, radcheckRepository, packageRepository, 
		resellerRepository, messageRepository, customerRepository,
		invoiceRepository,
	)
	menuUsecase := _menuUsecase.NewMenuUsecase(timeout, menuRepository)
	radusergroupUsecase := _radusergroupUsecase.NewRadusergroupUsecase(
		timeout, radusergroupRepository, radgroupcheckRepository, radgroupreplyRepository,
		radcheckRepository,
	)
	schedulerUsecase := _schedulerUsecase.NewSchedulerUsecase(timeout, radcheckRepository, radacctRepository)
	radgroupcheckUsecase := _radgroupcheckUsecase.NewRadgroupcheckUsecase(timeout, radgroupcheckRepository)
	radgroupreplyUsecase := _radgroupreplyUsecase.NewRadgroupreplyUsecase(timeout, radgroupreplyRepository)
	usersUsecase := _usersUsecase.NewUsersUsecase(
		timeout, usersRepository, customerRepository,
	)
	nasUsecase := _nasUsecase.NewUsecase(timeout, nasRepository)
	radpostauthUsecase := _radpostauthUsecase.NewUsecase(timeout, radpostauthRepository)
	invoiceUsecase := _invoiceUsecase.NewUsecase(timeout, invoiceRepository)
	paymentUsecase := _paymentUsecase.NewUsecase(timeout, paymentRepository)
	reportUsecase := _reportUsecase.NewUsecase(timeout, reportRepository, usersRepository)

	// defined all handler here
	_packageHandler.NewPackageHandler(e, packageUsecase)
	_menuHandler.NewMenuHandler(e, menuUsecase)
	_radusergroupHandler.NewRadusergroupHandler(e, radusergroupUsecase)
	_resellerHandler.NewResellerUsecase(e, resellerUsecase, messageUsecase)
	_transactionHandler.NewTransactionHandler(e, transactionUsecase)
	_schedulerHandler.NewSchedulerHandler(schedulerUsecase)
	_radgroupcheckHandler.NewRadgroupcheckHandler(e, radgroupcheckUsecase)
	_radgroupreplyHandler.NewRadgroupreplyHandler(e, radgroupreplyUsecase)
	_usersHandler.NewHandler(e, usersUsecase)
	_nasHandler.NewHandler(e, nasUsecase)
	_radpostauthHandler.NewHandler(e, radpostauthUsecase)
	_invoiceHandler.NewHandler(e, invoiceUsecase)
	_paymentHandler.NewHandler(e, paymentUsecase)
	_reportHandler.NewHandler(e, reportUsecase)

	go _telegramHandler.NewHandler(
		resellerUsecase, transactionUsecase, menuUsecase, messageUsecase,
	)

	go _messageHandler.NewHandler(messageUsecase)

	logrus.Fatal(e.Start(viper.GetString("server.address")))
}
