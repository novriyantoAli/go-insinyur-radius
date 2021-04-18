package mysql

import (
	"context"
	"database/sql"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlRadacctRepository struct {
	Conn *sql.DB
}

// NewMysqlRadacctRepository ...
func NewMysqlRadacctRepository(conn *sql.DB) domain.RadacctRepository {
	return &mysqlRadacctRepository{Conn: conn}
}

func (m *mysqlRadacctRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radacct, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Radacct, 0)
	for rows.Next() {
		t := domain.Radacct{}
		err = rows.Scan(
			&t.Radacctid,
			&t.Acctsessionid,
			&t.Acctuniqueid,
			&t.Username,
			&t.Realm,
			&t.Nasipaddress,
			&t.Nasportid,
			&t.Nasporttype,
			&t.Acctstarttime,
			&t.Acctupdatetime,
			&t.Acctstoptime,
			&t.Acctinterval,
			&t.Acctsessiontime,
			&t.Acctauthentic,
			&t.ConnectinfoStart,
			&t.ConnectinfoStop,
			&t.Acctinputoctets,
			&t.Acctoutputoctets,
			&t.Calledstationid,
			&t.Callingstationid,
			&t.Acctterminatecause,
			&t.Servicetype,
			&t.Framedprotocol,
			&t.Framedipaddress,
			&t.Secret,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRadacctRepository) FetchWithUsernameBatch(ctx context.Context, usernameList string) (res []domain.Radacct, err error) {
	query := "SELECT radacct.*, nas.secret FROM radacct INNER JOIN nas ON nas.nasname = radacct.nasipaddress WHERE acctstoptime is NULL AND username IN(" + usernameList + ")"

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}
