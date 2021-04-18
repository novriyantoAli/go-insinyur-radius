package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlRadcheckRepository struct {
	Conn *sql.DB
}

// NewMysqlRadcheckRepository ...
func NewMysqlRadcheckRepository(conn *sql.DB) domain.RadcheckRepository {
	return &mysqlRadcheckRepository{Conn: conn}
}

func (m *mysqlRadcheckRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radcheck, err error) {
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

	result = make([]domain.Radcheck, 0)
	for rows.Next() {
		t := domain.Radcheck{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Attribute,
			&t.OP,
			&t.Value,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRadcheckRepository) Fetch(ctx context.Context) (res []domain.Radcheck, err error) {
	query := `SELECT * FROM radcheck`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRadcheckRepository) FetchWithUsername(ctx context.Context, username string) (res []domain.Radcheck, err error) {
	query := `SELECT * FROM radcheck WHERE username = ?`

	res, err = m.fetch(ctx, query, username)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRadcheckRepository) FetchWithValueExpiration(ctx context.Context, delete bool) (res []domain.Radcheck, err error) {
	query := `SELECT id, username, attribute, op, value FROM radcheck WHERE attribute='Expiration' AND STR_TO_DATE(value, "%d %b %Y") <= CURDATE()`
	
	if delete == true {
		query += " - INTERVAL 7 DAY"
	}

	res, err = m.fetch(ctx, query)

	return
}

func (m *mysqlRadcheckRepository) Search(ctx context.Context, radcheck domain.Radcheck)(res []domain.Radcheck, err error){
	query := `SELECT * FROM radcheck`
	whereUsing := false

	args := make([]interface{}, 0)

	if radcheck.Username != nil {
		if whereUsing == false {
			query += "WHERE "
			whereUsing = true
		} else {
			query += "AND "
		}
		query += "username = ? "
		args = append(args, *radcheck.Username)
	}

	if radcheck.Attribute != nil {
		if whereUsing == false {
			query += "WHERE "
			whereUsing = true
		} else {
			query += "AND "
		}
		query += "attribute = ? "
		args = append(args, *radcheck.Attribute)
	}

	if radcheck.OP != nil {
		if whereUsing == false {
			query += "WHERE "
			whereUsing = true
		} else {
			query += "AND "
		}
		query += "op = ? "
		args = append(args, *radcheck.OP)
	}

	if radcheck.Value != nil {
		if whereUsing == false {
			query += "WHERE "
			whereUsing = true
		} else {
			query += "AND "
		}
		query += "value = ? "
		args = append(args, *radcheck.Value)
	}

	res, err = m.fetch(ctx, query, args...)

	return
}
func(m *mysqlRadcheckRepository) Find(ctx context.Context, radcheck domain.Radcheck)(res domain.Radcheck, err error){
	query := `SELECT * FROM radcheck WHERE username = ? AND attribute = ? LIMIT 1`

	resArr, err := m.fetch(ctx, query, *radcheck.Username, *radcheck.Attribute)
	if err != nil {
		return domain.Radcheck{}, err
	}

	if len(resArr) <= 0 {
		return domain.Radcheck{}, nil
	}

	res = resArr[0]

	return
}

func (m *mysqlRadcheckRepository) DeleteWithUsername(ctx context.Context, username string) (err error) {
	query := "DELETE FROM radcheck WHERE username = ?"

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

	if rowsAffected < 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return
	}

	return
}
