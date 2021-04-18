package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

// mysqlRepository ...
type mysqlRepository struct {
	Conn *sql.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB) domain.PaymentRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(c context.Context, query string, args ...interface{}) (res []domain.Payment, er error) {
	rows, err := m.Conn.QueryContext(c, query, args...)
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

	res = make([]domain.Payment, 0)
	for rows.Next() {
		t := domain.Payment{}
		err = rows.Scan(
			&t.ID,
			&t.NoInvoice,
			&t.Nominal,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (m *mysqlRepository) Search(ctx context.Context, spec *domain.Payment) (res []domain.Payment, err error) {
	query := ` SELECT * FROM payment WHERE no_invoice = ?`
	args := make([]interface{}, 0)
	args = append(args, *spec.NoInvoice)

	if spec.Nominal != nil {
		query += ` AND nominal = ?`
		args = append(args, *spec.Nominal)
	}

	res, err = m.fetch(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRepository) CountPage(ctx context.Context) (res int64, err error) {
	query := ` SELECT COUNT(id) as total_page FROM payment `
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
	}

	return
}

func (m *mysqlRepository) Fetch(ctx context.Context, id int64, limit int64) (res []domain.Payment, err error) {
	query := `SELECT * FROM payment WHERE id <= ? ORDER BY id DESC LIMIT ? `

	res, err = m.fetch(ctx, query, id, limit)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRepository) Find(ctx context.Context, spec *domain.Payment) (res []domain.Payment, err error) {
	query := "SELECT * FROM payment"
	args := make([]interface{}, 0)

	first := false

	if spec.ID != nil {
		if first == false {
			query += " WHERE "
			first = true
		} else {
			query += " AND "
		}
		query += " id = ? "
		args = append(args, *spec.ID)
	}

	if spec.NoInvoice != nil {
		if first == false {
			query += " WHERE "
			first = true
		} else {
			query += " AND "
		}
		query += " no_invoice = ? "
		args = append(args, *spec.NoInvoice)
	}

	if spec.Nominal != nil {
		if first == false {
			query += " WHERE "
			first = true
		} else {
			query += " AND "
		}
		query += " nominal = ? "
		args = append(args, *spec.Nominal)
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) Get(ctx context.Context, id int64) (res domain.Payment, err error) {
	query := `SELECT * FROM payment WHERE id = ? `
	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Payment{}, err
	}

	if len(resArr) > 0 {
		res = resArr[0]

		return
	}

	return domain.Payment{}, nil
}

func (m *mysqlRepository) Update(ctx context.Context, payment *domain.Payment) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	query := "UPDATE payment SET no_invoice = ?, nominal = ?"
	_, err = tx.ExecContext(ctx, query, *payment.NoInvoice, *payment.Nominal)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (m *mysqlRepository) Save(ctx context.Context, payment *domain.Payment) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	query := "INSERT INTO payment(no_invoice, nominal) VALUES(?,?)"
	res, err := tx.ExecContext(ctx, query, *payment.NoInvoice, *payment.Nominal)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	payment.ID = &lastID

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (m *mysqlRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM payment WHERE id = ?"

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
