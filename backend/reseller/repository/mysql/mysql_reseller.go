package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

type mysqlResellerRepository struct {
	Conn *sql.DB
}

// NewMysqlResellerRepository ...
func NewMysqlResellerRepository(conn *sql.DB) domain.ResellerRepository {
	return &mysqlResellerRepository{Conn: conn}
}

func (m *mysqlResellerRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Reseller, err error) {
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

	result = make([]domain.Reseller, 0)
	for rows.Next() {
		t := domain.Reseller{}
		err = rows.Scan(
			&t.ID,
			&t.TelegramID,
			&t.ChatID,
			&t.Name,
			&t.RegisterCode,
			&t.Active,
			&t.CreatedAt,
			&t.StatusTransaction,
			&t.DateTransaction,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlResellerRepository) Fetch(ctx context.Context) (res []domain.Reseller, err error) {
	query := `
	SELECT 
		reseller.*, 
		transaction.status as status_transaction, 
		transaction.created_at as date_transaction 
	FROM 
		reseller 
	LEFT JOIN 
		transaction 
		ON 
			reseller.id = transaction.id_reseller 
		AND 
			transaction.created_at = (SELECT MAX(created_at) as mx from transaction WHERE id_reseller = reseller.id) 
	ORDER BY 
		transaction.created_at 
	DESC
	`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlResellerRepository) Get(ctx context.Context, id int64) (res domain.Reseller, err error) {
	query := `	SELECT 
	reseller.*, 
	transaction.status as status_transaction, 
	transaction.created_at as date_transaction 
FROM 
	reseller 
LEFT JOIN 
	transaction 
	ON 
		reseller.id = transaction.id_reseller 
	AND 
		transaction.created_at = (SELECT MAX(created_at) as mx from transaction WHERE id_reseller = reseller.id) 
		WHERE reseller.id = ?`

	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Reseller{}, err
	}

	if len(resArr) <= 0 {
		return domain.Reseller{}, err
	}

	return resArr[0], err
}

func (m *mysqlResellerRepository) Insert(ctx context.Context, p1 *domain.Reseller) (err error) {
	query := `INSERT reseller SET name = ?, telegram_id = ?, chat_id = ?, active = ?, register_code = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, p1.Name, p1.TelegramID, p1.ChatID, p1.Active, p1.RegisterCode)
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

func (m *mysqlResellerRepository) Update(ctx context.Context, p1 *domain.Reseller) (err error) {
	query := `UPDATE reseller SET name = ?, telegram_id = ?, chat_id = ?, active = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, p1.Name, p1.TelegramID, p1.ChatID, p1.Active, p1.ID)
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

func (m *mysqlResellerRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM reseller WHERE id = ?"

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
