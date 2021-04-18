package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB) domain.NasRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.NAS, err error) {
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

	result = make([]domain.NAS, 0)
	for rows.Next() {
		t := domain.NAS{}
		err = rows.Scan(
			&t.ID,
			&t.Nasname,
			&t.Shortname,
			&t.Type,
			&t.Ports,
			&t.Secret,
			&t.Server,
			&t.Community,
			&t.Description,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) Get(ctx context.Context) (res domain.NAS, err error) {
	query := `SELECT * FROM nas `

	resArr, err := m.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return domain.NAS{}, err
	}

	if len(resArr) <= 0 {
		return domain.NAS{}, nil
	}

	res = resArr[0]

	return
}

func (m *mysqlRepository) Update(ctx context.Context, id int64, nas domain.NAS) (err error) {
	args := make([]interface{}, 0)

	query := `UPDATE nas SET nasname = ?`
	// To parameter not null
	args = append(args, *nas.Nasname)

	if nas.Shortname != nil && *nas.Shortname != "" {
		query += ", shortname = ? "
		args = append(args, *nas.Shortname)
	}

	if nas.Type != nil && *nas.Type != "" {
		query += ", type = ? "
		args = append(args, *nas.Type)
	}

	if nas.Ports != nil && *nas.Ports != "" {
		query += ", ports = ? "
		args = append(args, *nas.Ports)
	}

	if nas.Secret != nil && *nas.Secret != "" {
		query += ", secret = ? "
		args = append(args, *nas.Secret)
	}

	if nas.Server != nil && *nas.Server != "" {
		query += ", server = ? "
		args = append(args, *nas.Server)
	}

	if nas.Community != nil && *nas.Community != "" {
		query += ", community = ? "
		args = append(args, *nas.Community)
	}

	if nas.Description != nil && *nas.Description != "" {
		query += ", description = ? "
		args = append(args, *nas.Description)
	}

	query += " WHERE id = ? "
	args = append(args, id)

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		logrus.Error(err)
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *mysqlRepository) Insert(ctx context.Context, nas domain.NAS) (err error) {
	query := `INSERT nas SET nasname = ?, shortname = ?, type = ?, ports = ?, secret = ?, server = ?, community = ?, description = ? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, nas.Nasname, nas.Shortname, nas.Type, nas.Ports, nas.Secret, nas.Server, nas.Community, nas.Description)
	if err != nil {
		logrus.Error(err)
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		logrus.Error(err)
		return
	}

	return
}
