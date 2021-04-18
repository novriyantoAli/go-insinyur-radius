package udp

import (
	"context"
	"insinyur-radius/domain"
	"math"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

// SchedulerHandler ...
type SchedulerHandler struct {
	ucase domain.SchedulerUsecase
}

// NewSchedulerHandler ...
func NewSchedulerHandler( uc domain.SchedulerUsecase) {
	handler := &SchedulerHandler{ucase: uc}

	// scheduler 2 minutes
	go func() {
		for true {
			handler.TaskTwoMinutes()

			time.Sleep(2 * time.Minute)
		}
	}()

	// scheduler 1 day
	go func(){
		for true {
			handler.TaskOneDay()

			time.Sleep(5 * time.Minute)
			// time.Sleep(24 * time.Hour)
		}
	}()
}

// TaskTwoMinutes ....
func (uc *SchedulerHandler) TaskTwoMinutes() {
	ctx := context.Background()

	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	// date, _ := time.ParseInLocation(layoutFormat, "04 Oct 2020 12:23:05", wita)

	// logrus.Error(date)
	// logrus.Error(time.Now().AddDate(0, 0, -1).Format("02 Jan 2006"))
	// // fmt.Println(date)
	// // fmt.Println(time.Now().AddDate(0, 0, -1).Format("02 Jan 2006"))

	// duration := time.Now().Sub(date)
	// logrus.Error(duration)
	// // fmt.Println(duration)

	// get all users with expiration today
	radcheck, username, err := uc.ucase.GetUsers(ctx, false)
	if err != nil {
		logrus.Error(err)
		return
	}

	if username == "" {
		logrus.Error("username not found")
		return
	}

	// get all user online and ready for kicking
	radacct, err := uc.ucase.GetOnlineUsers(ctx, username)

	for _, value := range radcheck {
		// check if date and time now
		date, _ := time.ParseInLocation(layoutFormat, *value.Value, wita)
		duration := time.Now().Sub(date)
		if math.Signbit(duration.Seconds()) == false {
			for _, ra := range radacct {
				if *value.Username == *ra.Username {
					cmd := exec.Command("sh", "-c", `echo "Acct-Session-Id=`+*ra.Acctsessionid+`,User-Name=`+*ra.Username+`,NAS-IP-Address=`+*ra.Nasipaddress+`,Framed-IP-Address=`+*ra.Framedipaddress+`" | radclient -x `+*ra.Nasipaddress+`:3799 disconnect '`+*ra.Secret+`'`)
					err := cmd.Run()
					if err != nil {
						logrus.Error(err)
					}
				}
			}
		}
	}
}

// TaskOneDay ...
func (uc *SchedulerHandler) TaskOneDay() {
	ctx := context.Background()

	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	// date, _ := time.ParseInLocation(layoutFormat, "04 Oct 2020 12:23:05", wita)

	// logrus.Error(date)
	// logrus.Error(time.Now().AddDate(0, 0, -1).Format("02 Jan 2006"))
	// // fmt.Println(date)
	// // fmt.Println(time.Now().AddDate(0, 0, -1).Format("02 Jan 2006"))

	// duration := time.Now().Sub(date)
	// logrus.Error(duration)
	// // fmt.Println(duration)

	// get all users with expiration today
	radcheck, username, err := uc.ucase.GetUsers(ctx, true)
	if err != nil {
		logrus.Error(err)
		return
	}

	if username == "" {
		logrus.Error("username not found")
		return
	}

	// get all user online and ready for kicking
	radacct, err := uc.ucase.GetOnlineUsers(ctx, username)

	for _, value := range radcheck {
		// check if date and time now
		date, _ := time.ParseInLocation(layoutFormat, *value.Value, wita)
		duration := time.Now().Sub(date)
		if math.Signbit(duration.Seconds()) == false {
			for _, ra := range radacct {
				if *value.Username == *ra.Username {
					cmd := exec.Command("sh", "-c", `echo "Acct-Session-Id=`+*ra.Acctsessionid+`,User-Name=`+*ra.Username+`,NAS-IP-Address=`+*ra.Nasipaddress+`,Framed-IP-Address=`+*ra.Framedipaddress+`" | radclient -x `+*ra.Nasipaddress+`:3799 disconnect '`+*ra.Secret+`'`)
					err := cmd.Run()
					if err != nil {
						logrus.Error(err)
					}
				}
			}
			// delete expire user
			err := uc.ucase.DeleteExpireUsers(ctx, *value.Username)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
