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
func NewMysqlRepository(conn *sql.DB) domain.InvoiceRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(c context.Context, query string, args ...interface{}) (res []domain.Invoice, er error) {
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

	res = make([]domain.Invoice, 0)
	for rows.Next() {
		t := domain.Invoice{}
		err = rows.Scan(
			&t.ID,
			&t.NoInvoice,
			&t.Type,
			&t.Name,
			&t.Nominal,
			&t.CreatedAt,
			&t.Collector,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (m *mysqlRepository) CountPage(ctx context.Context) (res int64, err error) {
	query := ` SELECT COUNT(id) as total_page FROM invoice `
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

func (m *mysqlRepository) Fetch(ctx context.Context, id int64, limit int64) (res []domain.Invoice, err error) {
	args := make([]interface{}, 0)

	query := `SELECT * FROM invoice  `
	if id != 0 {
		query += "WHERE id <= ? "
		args = append(args, id)
	}

	query += "ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) Find(ctx context.Context, noInvoice string)(res domain.Invoice, err error){
	query := "SELECT * FROM invoice WHERE no_invoice = ? "
	resArr, err := m.fetch(ctx, query, noInvoice)
	if err != nil {
		return domain.Invoice{}, err
	}

	if len(resArr) > 0 {
		res = resArr[0]
		return
	}

	return domain.Invoice{}, nil
}

func (m *mysqlRepository) Get(ctx context.Context, id int64) (res domain.Invoice, err error) {
	query := `SELECT * FROM invoice WHERE id = ? `
	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Invoice{}, err
	}

	if len(resArr) > 0 {
		res = resArr[0]

		return
	}

	return domain.Invoice{}, nil
}

func (m *mysqlRepository) Update(ctx context.Context, invoice *domain.Invoice) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	query := "UPDATE invoice SET no_invoice = ?, type = ?, name = ?, nominal = ?"
	_, err = tx.ExecContext(
		ctx,
		query,
		*invoice.NoInvoice, *invoice.Type, *invoice.Name, *invoice.Nominal,
	)

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

func (m *mysqlRepository) Save(ctx context.Context, invoice *domain.Invoice) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	query := "INSERT INTO invoice(no_invoice, type, name, nominal) VALUES(?,?,?,?)"
	res, err := tx.ExecContext(
		ctx,
		query,
		*invoice.NoInvoice, *invoice.Type, *invoice.Name, *invoice.Nominal,
	)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	invoice.ID = &lastID

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (m mysqlRepository) Collector(ctx context.Context)(res []domain.Invoice, err error){
	query := `
	SELECT 
		invoice.*, , payment.nominal as paymeny_value 
	FROM 
		invoice 
	
	LEFT JOIN 
		payment 
		ON 
		payment.no_invoice = invoice.no_invoice 
	WHERE 
		invoice.collector = ?
	AND 
		payment.no_invoice IS NULL
	`
	res, err = m.fetch(ctx, query, "yes")
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM invoice WHERE id = ?"

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
