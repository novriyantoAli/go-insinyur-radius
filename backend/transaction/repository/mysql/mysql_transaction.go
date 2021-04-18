package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"insinyur-radius/domain"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type mysqlTransactionRepository struct {
	Conn *sql.DB
}

// NewMysqlTransactionRepository ...
func NewMysqlTransactionRepository(conn *sql.DB) domain.TransactionRepository {
	return &mysqlTransactionRepository{Conn: conn}
}

func (m *mysqlTransactionRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Transaction, err error) {
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

	result = make([]domain.Transaction, 0)
	for rows.Next() {
		t := domain.Transaction{}
		err = rows.Scan(
			&t.ID,
			&t.IDReseller,
			&t.NameReseller,
			&t.TransactionCode,
			&t.Status,
			&t.Value,
			&t.Information,
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

func (m *mysqlTransactionRepository) Fetch(ctx context.Context) (res []domain.Transaction, err error) {
	query := `SELECT transaction.id, transaction.id_reseller, reseller.name, transaction.transaction_code, transaction.status, transaction.value, transaction.information, transaction.created_at FROM transaction INNER JOIN reseller ON transaction.id_reseller = reseller.id`

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlTransactionRepository) Report(ctx context.Context, dateStart string, dateEnd string) (res []domain.Transaction, err error) {
	query := `SELECT transaction.id, transaction.id_reseller, reseller.name, transaction.transaction_code, transaction.status, transaction.value, transaction.information, transaction.created_at FROM transaction INNER JOIN reseller ON transaction.id_reseller = reseller.id WHERE (transaction.created_at BETWEEN ? AND ?) ORDER BY transaction.id_reseller, transaction.status`

	res, err = m.fetch(ctx, query, dateStart, dateEnd)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlTransactionRepository) FetchWithIDReseller(ctx context.Context, idReseller int64) (res []domain.Transaction, err error) {
	query := `SELECT transaction.id, transaction.id_reseller, reseller.name, transaction.transaction_code, transaction.status, transaction.value, transaction.information, transaction.created_at FROM transaction INNER JOIN reseller ON transaction.id_reseller = reseller.id WHERE id_reseller = ?`

	res, err = m.fetch(ctx, query, idReseller)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlTransactionRepository) GetWithTransactionCode(ctx context.Context, code string) (res domain.Transaction, err error) {
	query := `SELECT transaction.id, transaction.id_reseller, reseller.name, transaction.transaction_code, transaction.status, transaction.value, transaction.information, transaction.created_at FROM transaction INNER JOIN reseller ON transaction.id_reseller = reseller.id WHERE transaction_code = ?`

	resArr, err := m.fetch(ctx, query, code)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(resArr) > 0 {
		res = resArr[0]
		return
	}

	return
}

func (m *mysqlTransactionRepository) RefillBalance(ctx context.Context, transaction domain.Transaction, message domain.Message) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// transaction
	query := "INSERT transaction SET id_reseller = ?, status = ?, transaction_code = ?, value = ?"
	_, err = tx.ExecContext(ctx, query, *transaction.IDReseller, *transaction.Status, *transaction.TransactionCode, *transaction.Value)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
	// invoice
	query = "INSERT INTO invoice(no_invoice, type, name, nominal, collector) VALUES(?,?,(SELECT name FROM reseller WHERE id = ? LIMIT 1),?,?)"
	_, err = tx.ExecContext(
		ctx,
		query,
		uniqTime, "postpaid", *transaction.IDReseller, *transaction.Value, "no",
	)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// message
	query = "INSERT INTO message(chat_id, message_id, message) VALUES((SELECT chat_id FROM reseller WHERE id = ? LIMIT 1),(SELECT message_id FROM reseller WHERE id = ? LIMIT 1),?)"
	_, err = tx.ExecContext(ctx, query, *transaction.IDReseller, *transaction.IDReseller, *message.Message)

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

	return
}

func (m *mysqlTransactionRepository) ResellerRefillTransaction(ctx context.Context, transaction domain.Transaction, customer string, expiration string) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// update radcheck
	query := "UPDATE radcheck SET value = ? WHERE username = ? AND attribute = 'Expiration'"
	_, err = tx.ExecContext(ctx, query, expiration, *transaction.Information)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// transaction
	query = "INSERT INTO transaction(id_reseller, transaction_code, status, value, information) VALUES(?,?,?,?,?)"
	_, err = tx.ExecContext(
		ctx,
		query,
		*transaction.IDReseller, *transaction.TransactionCode, *transaction.Status, *transaction.Value, *transaction.Information,
	)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// invoice
	query = "DELETE FROM invoice WHERE no_invoice = ?"
	_, err = tx.ExecContext(ctx, query, customer)

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

	return
}

func (m *mysqlTransactionRepository) ResellerTransaction(ctx context.Context, transaction *domain.Transaction, idPackage int64, profile string) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// and save radpackage, radcheck, transaction

	// radpackage
	query := "INSERT INTO radpackage(id_package, username) VALUES(?,?)"
	_, err = tx.ExecContext(ctx, query, idPackage, *transaction.Information)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// radcheck
	query = "INSERT INTO radcheck(username, attribute, op, value) VALUES(?,?,?,?),(?,?,?,?)"
	_, err = tx.ExecContext(
		ctx,
		query,
		*transaction.Information, "Cleartext-Password", ":=", *transaction.Information,
		*transaction.Information, "User-Profile", ":=", profile,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// transaction
	query = "INSERT INTO transaction(id_reseller, transaction_code, status, value, information) VALUES(?,?,?,?,?)"
	_, err = tx.ExecContext(
		ctx,
		query,
		*transaction.IDReseller, *transaction.TransactionCode, *transaction.Status, *transaction.Value, *transaction.Information,
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

	trc, _ := m.GetWithTransactionCode(ctx, *transaction.TransactionCode)
	if trc != (domain.Transaction{}) {
		transaction.CreatedAt = trc.CreatedAt
	}

	return nil
}

func (m *mysqlTransactionRepository) Insert(ctx context.Context, p1 *domain.Transaction) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// transaction
	query := "INSERT transaction SET id_reseller = ?, transaction_code = ?, status = ?, value = ?, information = ?"
	res, err := tx.ExecContext(ctx, query, p1.IDReseller, p1.TransactionCode, p1.Status, p1.Value, p1.Information)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// invoice
	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
	query = "INSERT invoice SET no_invoice = ?, type = ?, name = (SELECT name FROM reseller WHERE id = ?), nominal = ?, collector = 'no'"
	_, err = tx.ExecContext(
		ctx, query,
		uniqTime, "postpaid", *p1.IDReseller, *p1.Value,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	lastID, err := res.LastInsertId()
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

	// query := `INSERT transaction SET id_reseller = ?, transaction_code = ?, status = ?, value = ?, information = ?`
	// stmt, err := m.Conn.PrepareContext(ctx, query)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// res, err := stmt.ExecContext(ctx, p1.IDReseller, p1.TransactionCode, p1.Status, p1.Value, p1.Information)
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }

	// lastID, err := res.LastInsertId()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }
	p1.ID = &lastID

	return
}

func (m *mysqlTransactionRepository) Update(ctx context.Context, p1 *domain.Transaction) (err error) {
	query := `UPDATE transaction SET id_reseller = ?, status = ?, value = ?, information = ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return
	}

	res, err := stmt.ExecContext(ctx, p1.IDReseller, p1.Status, p1.Value, p1.Information, p1.ID)
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

func (m *mysqlTransactionRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM reseller WHERE id = ?"

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
