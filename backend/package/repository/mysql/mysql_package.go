package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlPackageRepository struct {
	Conn *sql.DB
}

// NewMysqlPackageRepository ...
func NewMysqlPackageRepository(conn *sql.DB) domain.PackageRepository {
	return &mysqlPackageRepository{Conn: conn}
}

func (m *mysqlPackageRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Package, err error) {
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

	result = make([]domain.Package, 0)
	for rows.Next() {
		t := domain.Package{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.ValidityValue,
			&t.ValidityUnit,
			&t.Price,
			&t.Margin,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlPackageRepository) Fetch(ctx context.Context) (res []domain.Package, err error) {
	query := `SELECT * FROM package`

	res, err = m.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlPackageRepository) Get(ctx context.Context, id int64) (res domain.Package, err error) {
	query := `SELECT * FROM package WHERE id = ?`

	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		logrus.Error(err)
		return domain.Package{}, err
	}

	res = resArr[0]

	return
}

func (m *mysqlPackageRepository) Insert(ctx context.Context, p1 *domain.Package) (err error) {
	query := `INSERT package SET name = ?, validity_value = ?, validity_unit = ?, price = ?, margin = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, *p1.Name, *p1.ValidityValue, *p1.ValidityUnit, *p1.Price, *p1.Margin)
	if err != nil {
		logrus.Error(err)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logrus.Error(err)
		return
	}
	p1.ID = &lastID

	return
}

func (m *mysqlPackageRepository) Update(ctx context.Context, p1 *domain.Package) (err error) {
	query := `UPDATE package SET name = ?, validity_value = ?, validity_unit = ?, price = ?, margin = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, p1.Name, p1.ValidityValue, p1.ValidityUnit, p1.Price, p1.Margin, p1.ID)
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

func (m *mysqlPackageRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM package WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		logrus.Error(err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return
	}

	return
}
