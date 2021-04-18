package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlRadgroupcheckRepository struct {
	Conn *sql.DB
}

// NewMysqlRadgroupcheckRepository ...
func NewMysqlRadgroupcheckRepository(conn *sql.DB) domain.RadgroupcheckRepository {
	return &mysqlRadgroupcheckRepository{Conn: conn}
}

func (m *mysqlRadgroupcheckRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radgroupcheck, err error) {
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

	result = make([]domain.Radgroupcheck, 0)
	for rows.Next() {
		t := domain.Radgroupcheck{}
		err = rows.Scan(
			&t.ID,
			&t.Groupname,
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

func (m *mysqlRadgroupcheckRepository) Find(ctx context.Context, spec domain.Radgroupcheck) (res domain.Radgroupcheck, err error) {
	args := make([]interface{}, 0)
	query := `SELECT * FROM radgroupcheck WHERE `
	first := false

	if *spec.ID != 0 {
		args = append(args, *spec.ID)
		query += `id = ? `
		first = true
	}

	if *spec.Groupname != "" {
		args = append(args, *spec.Groupname)
		if first == true {
			query += `AND groupname = ? `
		} else {
			first = true
			query += `groupname = ? `
		}
	}

	if *spec.Attribute != "" {
		args = append(args, *spec.Attribute)
		if first == true {
			query += `AND attribute = ? `
		} else {
			first = true
			query += `attribute = ? `
		}
	}

	if *spec.OP != "" {
		args = append(args, *spec.OP)
		if first == true {
			query += `AND op = ? `
		} else {
			first = true
			query += `op = ? `
		}
	}

	if *spec.Value != "" {
		args = append(args, *spec.Value)
		if first == true {
			query += `AND value = ? `
		} else {
			first = true
			query += `value = ? `
		}
	}

	resArr, err := m.fetch(ctx, query, args...)
	if err != nil {
		return domain.Radgroupcheck{}, err
	}

	res = resArr[0]

	return
}

func (m *mysqlRadgroupcheckRepository) DeleteWithGroupname(ctx context.Context, name string) (err error) {
	query := "DELETE FROM radgroupcheck WHERE groupname = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, name)
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
