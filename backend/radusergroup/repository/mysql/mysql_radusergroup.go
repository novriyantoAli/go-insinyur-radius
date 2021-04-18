package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"
	"strconv"

	"github.com/sirupsen/logrus"
)

type mysqlRadusergroupRepository struct {
	Conn *sql.DB
}

// NewMysqlRadusergroupRepository ...
func NewMysqlRadusergroupRepository(conn *sql.DB) domain.RadusergroupRepository {
	return &mysqlRadusergroupRepository{Conn: conn}
}

func (m *mysqlRadusergroupRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radusergroup, err error) {
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

	result = make([]domain.Radusergroup, 0)
	for rows.Next() {
		t := domain.Radusergroup{}
		err = rows.Scan(
			&t.Username,
			&t.Groupname,
			&t.Priority,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRadusergroupRepository) Fetch(ctx context.Context) (res []domain.Radusergroup, err error) {
	query := `SELECT * FROM radusergroup`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRadusergroupRepository) Get(ctx context.Context, username string)(res domain.Radusergroup, err error){
	query := `SELECT * FROM radusergroup WHERE username = ?`
	resArr, err := m.fetch(ctx, query, username)
	if err != nil {
		return domain.Radusergroup{}, err
	}

	if len(resArr) <= 0{
		err = fmt.Errorf("item not found")
		return domain.Radusergroup{}, err
	}

	res = resArr[0]

	return
}

func (m *mysqlRadusergroupRepository) SaveProfile(ctx context.Context, profile domain.Profile) (err error){
	tx, err := m.Conn.BeginTx(ctx, nil)

	prefixWithProfile := profile.PrefixName + "_" + profile.ProfileName
	// radusergroup
	query := "INSERT radusergroup SET username = ?, groupname = ?, priority = ?"
	_, err = tx.ExecContext(ctx, query, prefixWithProfile, prefixWithProfile, profile.Priority)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	if profile.UseLimitSession == true {
		// radgroupcheck
		query = "INSERT INTO radgroupcheck(groupname, attribute, op, value) VALUES(?,?,?,?)"
		_, err = tx.ExecContext(ctx, query, prefixWithProfile, "Simultaneous-Use", ":=", profile.LimitSession)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	}

	if profile.PrefixName == "pppoe" {
		query = "INSERT INTO radgroupcheck(groupname, attribute, op, value) VALUES(?,?,?,?)"
		_, err = tx.ExecContext(ctx, query,
			prefixWithProfile, "Framed-Protocol", "==", "PPP",
		)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	}

	if profile.PrefixName == "pppoe" {
		if profile.UseLimitSpeed == true {
			query = "INSERT INTO radgroupreply(groupname, attribute, op, value) VALUES(?,?,?,?)"
			_, err = tx.ExecContext(ctx, query,
				prefixWithProfile, "Mikrotik-Rate-Limit", "=", strconv.FormatInt((profile.LimitMIRUpload / 1000), 10)+"k/"+strconv.FormatInt((profile.LimitMIRDownload / 1000), 10)+"k",
			)
			if err != nil {
				logrus.Error(err)
	
				tx.Rollback()
				return err
			}
		}
		query = "INSERT INTO radgroupreply(groupname, attribute, op, value) VALUES(?,?,?,?)"
		_, err = tx.ExecContext(ctx, query,
			prefixWithProfile, "Framed-Pool", "=", profile.PoolName,
		)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	} else if profile.UseLimitSpeed == true && profile.PrefixName == "hotspot" {
		// radgroupreply
		query = "INSERT INTO radgroupreply(groupname, attribute, op, value) VALUES(?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?)"
		_, err = tx.ExecContext(
			ctx, query, 
			prefixWithProfile, "WISPr-Bandwidth-Min-Up", ":=", strconv.FormatInt(profile.LimitCIRUpload, 10),
			prefixWithProfile, "WISPr-Bandwidth-Max-Up", ":=", strconv.FormatInt(profile.LimitMIRUpload, 10),
			prefixWithProfile, "WISPr-Bandwidth-Min-Down", ":=", strconv.FormatInt(profile.LimitCIRDownload, 10),
			prefixWithProfile, "WISPr-Bandwidth-Max-Down", ":=", strconv.FormatInt(profile.LimitMIRDownload, 10),
		)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (m *mysqlRadusergroupRepository) Insert(ctx context.Context, rug domain.Radusergroup) (err error){
	query := `INSERT radusergroup SET username = ?, groupname = ?, priority = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, rug.Username, rug.Groupname, rug.Priority)
	if err != nil {
		return
	}

	afrows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if afrows <= 0 {
		err = fmt.Errorf("error, no rows has affected")
		return
	}

	return
}

func (m *mysqlRadusergroupRepository) Delete(ctx context.Context, username string)(err error){
	query := "DELETE FROM radusergroup WHERE username = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, username)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return
	}

	return
}