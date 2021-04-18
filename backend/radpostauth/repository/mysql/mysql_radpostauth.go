package mysql

import (
	"context"
	"database/sql"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB) domain.RadpostauthRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radpostauth, err error) {
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

	result = make([]domain.Radpostauth, 0)
	for rows.Next() {
		t := domain.Radpostauth{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Pass,
			&t.Reply,
			&t.Authdate,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) CountUsernamePage(ctx context.Context, username string) (res int64, err error) {
	query := `SELECT * FROM radpostauth WHERE username = ? ORDER BY id DESC`

	resArr, err := m.fetch(ctx, query, username)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return int64(len(resArr)), nil
}

func (m *mysqlRepository) Get(ctx context.Context, username string, id int64, limit int64) (res []domain.Radpostauth, err error) {
	args := make([]interface{}, 0)
	query := `SELECT * FROM radpostauth WHERE username = ? `

	args = append(args, username)

	if id != 0 {
		args = append(args, id)
		query += `AND id = ? `
	}

	query += " ORDER BY id DESC LIMIT ? "
	args = append(args, limit)

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}
