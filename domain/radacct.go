package domain

import (
	"context"
	"time"
)

// Radacct ...
type Radacct struct {
	Radacctid          *int64     `json:"radacctid"`
	Acctsessionid      *string    `json:"acctsessionid"`
	Acctuniqueid       *string    `json:"acctuniqueid"`
	Username           *string    `json:"username"`
	Realm              *string    `json:"realm"`
	Nasipaddress       *string    `json:"nasipaddress"`
	Nasportid          *string    `json:"nasportid"`
	Nasporttype        *string    `json:"nasporttype"`
	Acctstarttime      *time.Time `json:"acctstarttime"`
	Acctupdatetime     *time.Time `json:"acctupdatetime"`
	Acctstoptime       *time.Time `json:"acctstoptime"`
	Acctinterval       *int64     `json:"acctinterval"`
	Acctsessiontime    *int64     `json:"acctsessiontime"`
	Acctauthentic      *string    `json:"acctauthentic"`
	ConnectinfoStart   *string    `json:"connectinfo_start"`
	ConnectinfoStop    *string    `json:"connectinfo_stop"`
	Acctinputoctets    *int64     `json:"acctinputoctets"`
	Acctoutputoctets   *int64     `json:"acctoutputoctets"`
	Calledstationid    *string    `json:"calledstationid"`
	Callingstationid   *string    `json:"callingstationid"`
	Acctterminatecause *string    `json:"acctterminatecause"`
	Servicetype        *string    `json:"servicetype"`
	Framedprotocol     *string    `json:"framedprotocol"`
	Framedipaddress    *string    `json:"framedipaddress"`
	Secret             *string    `json:"secret"`
}

// RadacctUsecase ...
type RadacctUsecase interface {
	FetchWithUsernameBatch(ctx context.Context, usernameList string)(res []Radacct, err error)
}

// RadacctRepository ...
type RadacctRepository interface {
	FetchWithUsernameBatch(ctx context.Context, usernameList string) (res []Radacct, err error)
}
