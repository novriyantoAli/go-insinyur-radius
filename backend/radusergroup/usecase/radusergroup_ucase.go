package usecase

import (
	"context"
	"fmt"
	"insinyur-radius/domain"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type radusergroupUsecase struct {
	Timeout        time.Duration
	Repository     domain.RadusergroupRepository
	RepoGroupCheck domain.RadgroupcheckRepository
	RepoGroupReply domain.RadgroupreplyRepository
	Radcheck       domain.RadcheckRepository
}

// NewRadusergroupUsecase ...
func NewRadusergroupUsecase(t time.Duration, r domain.RadusergroupRepository, rgc domain.RadgroupcheckRepository, rgr domain.RadgroupreplyRepository, rcr domain.RadcheckRepository) domain.RadusergroupUsecase {
	return &radusergroupUsecase{Timeout: t, Repository: r, RepoGroupCheck: rgc, RepoGroupReply: rgr, Radcheck: rcr}
}

// Fetch ...
func (uc *radusergroupUsecase) Fetch(c context.Context) (res []domain.Radusergroup, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (uc *radusergroupUsecase) Get(c context.Context, username string) (res domain.Radusergroup, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	res, err = uc.Repository.Get(ctx, username)
	if err != nil {
		return domain.Radusergroup{}, err
	}

	return
}

func (uc *radusergroupUsecase) LoadProfile(c context.Context, profile string) (res domain.Profile, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	rs := strings.Split(profile, "_")

	if len(rs) <= 1 {
		err = fmt.Errorf("profile name is invalid")
		return
	}

	radusergroup, err := uc.Repository.Get(ctx, profile)
	if err != nil {
		return
	}

	res.ProfileName = *radusergroup.Username
	res.Priority = *radusergroup.Priority

	// now search radgroupcheck and insert all you need to profile
	// search limit session
	limitSession := "Simultaneous-Use"
	op := ":="
	value := ""
	var id int64 = 0
	radgroupcheck := domain.Radgroupcheck{}
	radgroupcheck.ID = &id
	radgroupcheck.Groupname = radusergroup.Groupname
	radgroupcheck.Attribute = &limitSession
	radgroupcheck.OP = &op
	radgroupcheck.Value = &value

	limit, _ := uc.RepoGroupCheck.Find(ctx, radgroupcheck)
	if limit != (domain.Radgroupcheck{}) {
		intres, err := strconv.Atoi(*limit.Value)
		if err != nil {
			return domain.Profile{}, err
		}
		res.UseLimitSession = true
		res.LimitSession = intres
	}

	if strings.Trim(rs[0], " ") == "pppoe" {
		// now search radgroupreply and insert all you need to profile
		// search MIR Upload
		// override op value
		op = "="

		mirUpload := "Mikrotik-Rate-Limit"
		radgroupreply := domain.Radgroupreply{}
		radgroupreply.ID = &id
		radgroupreply.Groupname = radusergroup.Groupname
		radgroupreply.Attribute = &mirUpload
		radgroupreply.OP = &op
		radgroupreply.Value = &value

		instance, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)

		if instance != (domain.Radgroupreply{}) {
			rsRateLimit := strings.Split(*instance.Value, "/")
			if len(rsRateLimit) > 0 {
				upload := rsRateLimit[0]
				download := rsRateLimit[1]

				speedUpload := upload[:len(upload)-1]
				speedDownload := download[:len(download)-1]

				int64Upload, err := strconv.ParseInt(speedUpload, 10, 64)
				if err != nil {
					return domain.Profile{}, err
				}

				int64Download, err := strconv.ParseInt(speedDownload, 10, 64)
				if err != nil {
					return domain.Profile{}, err
				}
				res.UseLimitSpeed = true
				res.LimitMIRUpload = int64Upload * 1000
				res.LimitMIRDownload = int64Download * 1000
			}
		}

		fremedPool := "Framed-Pool"
		radgroupreply.Attribute = &fremedPool

		instance2, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)
		if instance2 != (domain.Radgroupreply{}){
			res.PoolName = *instance2.Value
		}

	} else {
		// now search radgroupreply and insert all you need to profile
		// search CIR Upload
		cirUpload := "WISPr-Bandwidth-Min-Up"
		radgroupreply := domain.Radgroupreply{}
		radgroupreply.ID = &id
		radgroupreply.Groupname = radusergroup.Groupname
		radgroupreply.Attribute = &cirUpload
		radgroupreply.OP = &op
		radgroupreply.Value = &value

		CIRupload, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)
		if CIRupload != (domain.Radgroupreply{}) {
			int64res, err := strconv.ParseInt(*CIRupload.Value, 10, 64)
			if err != nil {
				return domain.Profile{}, err
			}
			res.UseLimitSpeed = true
			res.LimitCIRUpload = int64res
		}

		// search MIR Upload
		mirUpload := "WISPr-Bandwidth-Max-Up"
		radgroupreply = domain.Radgroupreply{}
		radgroupreply.ID = &id
		radgroupreply.Groupname = radusergroup.Groupname
		radgroupreply.Attribute = &mirUpload
		radgroupreply.OP = &op
		radgroupreply.Value = &value

		MIRupload, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)
		if MIRupload != (domain.Radgroupreply{}) {
			int64res, err := strconv.ParseInt(*MIRupload.Value, 10, 64)
			if err != nil {
				return domain.Profile{}, err
			}
			res.UseLimitSpeed = true
			res.LimitMIRUpload = int64res
		}

		// search CIR Download
		cirDownload := "WISPr-Bandwidth-Min-Down"
		radgroupreply = domain.Radgroupreply{}
		radgroupreply.ID = &id
		radgroupreply.Groupname = radusergroup.Groupname
		radgroupreply.Attribute = &cirDownload
		radgroupreply.OP = &op
		radgroupreply.Value = &value

		CIRdownload, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)
		if CIRdownload != (domain.Radgroupreply{}) {
			int64res, err := strconv.ParseInt(*CIRdownload.Value, 10, 64)
			if err != nil {
				return domain.Profile{}, err
			}
			res.UseLimitSpeed = true
			res.LimitCIRDownload = int64res
		}

		// search MIR Download
		mirDownload := "WISPr-Bandwidth-Max-Down"
		radgroupreply = domain.Radgroupreply{}
		radgroupreply.ID = &id
		radgroupreply.Groupname = radusergroup.Groupname
		radgroupreply.Attribute = &mirDownload
		radgroupreply.OP = &op
		radgroupreply.Value = &value

		MIRdownload, _ := uc.RepoGroupReply.Find(ctx, radgroupreply)
		if MIRdownload != (domain.Radgroupreply{}) {
			int64res, err := strconv.ParseInt(*MIRdownload.Value, 10, 64)
			if err != nil {
				return domain.Profile{}, err
			}
			res.UseLimitSpeed = true
			res.LimitMIRDownload = int64res
		}
	}

	res.ProfileName = rs[1]
	res.PrefixName = rs[0]
	return
}

func (uc *radusergroupUsecase) Save(c context.Context, m *domain.Radusergroup) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Insert(ctx, *m)

	return

}

func (uc *radusergroupUsecase) SaveProfile(c context.Context, m *domain.Profile) (result domain.Radusergroup, err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	withPrefixName := m.PrefixName + "_" + m.ProfileName

	poolNameValidate := strings.Trim(m.PoolName, "")
	poolNameArr := strings.Split(poolNameValidate, " ")
	m.PoolName = strings.Join(poolNameArr, "")

	errToLog := uc.Repository.Delete(ctx, withPrefixName)
	if errToLog != nil {
		logrus.Errorln(errToLog)
	}

	errToLog = uc.RepoGroupCheck.DeleteWithGroupname(ctx, withPrefixName)
	if errToLog != nil {
		logrus.Warn(errToLog)
	}

	errToLog = uc.RepoGroupReply.DeleteWithGroupname(ctx, withPrefixName)
	if errToLog != nil {
		logrus.Warn(errToLog)
	}

	if m.UseLimitSpeed == true {
		// CIR upload
		switch m.LimitCIRUploadUnit {
		case "b":
			logrus.Debug("original selected, not item are change")
		case "kb":
			m.LimitCIRUpload = m.LimitCIRUpload * 1000
		case "mb":
			m.LimitCIRUpload = (m.LimitCIRUpload * 1000) * 1000
		default:
			return domain.Radusergroup{}, fmt.Errorf("cir upload unit undefined")
		}

		// MIR upload
		switch m.LimitMIRUploadUnit {
		case "b":
			logrus.Debug("original selected, not item are change")
		case "kb":
			m.LimitMIRUpload = m.LimitMIRUpload * 1000
		case "mb":
			m.LimitMIRUpload = (m.LimitMIRUpload * 1000) * 1000
		default:
			return domain.Radusergroup{}, fmt.Errorf("mir upload unit undefined")
		}

		// CIR Download
		switch m.LimitCIRDownloadUnit {
		case "b":
			logrus.Debug("original selected, not item are change")
		case "kb":
			m.LimitCIRDownload = m.LimitCIRDownload * 1000
		case "mb":
			m.LimitCIRDownload = (m.LimitCIRDownload * 1000) * 1000
		default:
			return domain.Radusergroup{}, fmt.Errorf("cir download unit undefined")
		}

		// MIR Download
		switch m.LimitMIRDownloadUnit {
		case "b":
			logrus.Debug("original selected, not item are change")
		case "kb":
			m.LimitMIRDownload = m.LimitMIRDownload * 1000
		case "mb":
			m.LimitMIRDownload = (m.LimitMIRDownload * 1000) * 1000
		default:
			return domain.Radusergroup{}, fmt.Errorf("mir download unit undefined")
		}
	}

	err = uc.Repository.SaveProfile(ctx, *m)

	if err != nil {
		return domain.Radusergroup{}, err
	}

	result, err = uc.Repository.Get(ctx, withPrefixName)

	return
}

func (uc *radusergroupUsecase) Delete(c context.Context, username string) (err error) {
	ctx, cancel := context.WithTimeout(c, uc.Timeout)
	defer cancel()

	err = uc.Repository.Delete(ctx, username)

	specAttribute := "User-Profile"
	radcheckSpec := domain.Radcheck{}
	radcheckSpec.Attribute = &specAttribute
	radcheckSpec.Value = &username

	radcheck, _ := uc.Radcheck.Search(ctx, radcheckSpec)

	if len(radcheck) > 0 {
		for _, item := range radcheck {
			uc.Radcheck.DeleteWithUsername(ctx, *item.Username)
		}
	}

	return
}
