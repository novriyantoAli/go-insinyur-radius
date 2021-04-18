package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlMenuRepository struct {
	Conn *sql.DB
}

// NewMysqlMenuRepository ...
func NewMysqlMenuRepository(conn *sql.DB) domain.MenuRepository {
	return &mysqlMenuRepository{Conn: conn}
}

func (m *mysqlMenuRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Menu, err error) {
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

	result = make([]domain.Menu, 0)
	for rows.Next() {
		t := domain.Menu{}
		err = rows.Scan(
			&t.ID,
			&t.IDPackage,
			&t.Profile,
			&t.Name,
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

func (m *mysqlMenuRepository) Fetch(ctx context.Context) (res []domain.Menu, err error) {
	query := `SELECT * FROM menu`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlMenuRepository) Get(ctx context.Context, id int64) (res domain.Menu, err error) {
	query := `SELECT * FROm menu WHERE id = ?`
	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Menu{}, err
	}

	return resArr[0], err
}

func (m *mysqlMenuRepository) Insert(ctx context.Context, p1 *domain.Menu) (err error) {
	query := `INSERT menu SET name = ?, profile = ?, id_package = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, p1.Name, p1.Profile, p1.IDPackage)
	if err != nil {
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	p1.ID = &lastID

	return
}

func (m *mysqlMenuRepository) Update(ctx context.Context, p1 *domain.Menu) (err error) {
	query := `UPDATE menu SET name = ?, profile = ?, id_package = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, p1.Name, p1.Profile, p1.IDPackage, p1.ID)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *mysqlMenuRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM menu WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
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
