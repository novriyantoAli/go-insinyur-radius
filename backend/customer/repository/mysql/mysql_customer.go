package mysql

import (
	"context"
	"database/sql"
	"insinyur-radius/domain"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB) domain.CustomerRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Customer, err error) {
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

	// ID         *int64  `json:"id"`
	// IDPackage  *int64  `json:"id_package"`
	// Name       *string `json:"name"`
	// Username   *string `json:"username"`
	// Password   *string `json:"password"`
	// Type       *string `json:"type"`
	// Profile    *string `json:"profile"`
	// Expiration *string `json:"expiration"`
	// CreatedAt  *string `json:"created_at"`

	result = make([]domain.Customer, 0)
	for rows.Next() {
		t := domain.Customer{}
		err = rows.Scan(
			&t.ID,
			&t.IDPackage,
			&t.Name,
			&t.Username,
			&t.Password,
			&t.Type,
			&t.Profile,
			&t.CreatedAt,
			&t.Expiration,
		)

		// ID         *int64  `json:"id"`
		// IDPackage  *int64  `json:"id_package"`
		// Name       *string `json:"name"`
		// Username   *string `json:"username"`
		// Password   *string `json:"password"`
		// Type       *string `json:"type"`
		// Profile    *string `json:"profile"`
		// Expiration *string `json:"expiration"`
		// CreatedAt  *string `json:"created_at"`

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) CountPage(ctx context.Context) (res int, err error) {
	query := `SELECT * FROM customer`
	resArr, err := m.fetch(ctx, query)
	if err != nil {
		return 0, err
	}

	return len(resArr), nil
}

func (m *mysqlRepository) Fetch(ctx context.Context, id int64, limit int64) (res []domain.Customer, err error) {
	query := `
			SELECT 
				customer.*, radcheck.value as expiration 
			FROM 
				customer 
			LEFT JOIN 
				radcheck 
			ON 
				radcheck.username = customer.username 
			AND
				radcheck.attribute = 'Expiration' 
			 
			WHERE 

				customer.id < ? ORDER BY customer.id DESC LIMIT ?`

	res, err = m.fetch(ctx, query, id, limit)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRepository) GetUsername(ctx context.Context, username string) (res domain.Customer, err error) {
	query := `
	SELECT 
		customer.*, radcheck.value as expiration 
	FROM 
		customer 
	LEFT JOIN 
		radcheck 
	ON 
		radcheck.username = customer.username 
	AND
		radcheck.attribute = 'Expiration'	
			WHERE 
				customer.username = ? `

	resArr, err := m.fetch(ctx, query, username)
	if err != nil {
		return domain.Customer{}, err
	}

	if len(resArr) <= 0 {
		return domain.Customer{}, nil
	}

	res = resArr[0]

	return
}

func (m *mysqlRepository) Get(ctx context.Context, id int64) (res domain.Customer, err error) {
	query := `
			SELECT 
				customer.*, radcheck.value as expiration 
			FROM 
				customer 
			LEFT JOIN 
				radcheck 
			ON 
				radcheck.username = customer.username 
				AND radcheck.attribute = 'Expiration' 
			WHERE 
				customer.id = ? `

	resArr, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Customer{}, err
	}

	if len(resArr) <= 0 {
		return domain.Customer{}, nil
	}

	res = resArr[0]

	return
}

func (m *mysqlRepository) Update(ctx context.Context, username string, customer domain.Customer) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	args := make([]interface{}, 0)

	// and save radpackage, radcheck, invoice

	first := true
	// radpackage
	query := "UPDATE customer SET "
	if customer.IDPackage != nil {
		if first == true {
			first = false
			query += "id_package = ?"
		} else {
			query += ", id_package = ?"
		}
		args = append(args, *customer.IDPackage)
	}

	if customer.Name != nil {
		if first == true {
			first = false
			query += "name = ?"
		} else {
			query += ", name = ?"
		}
		args = append(args, *customer.Name)
	}

	if customer.Username != nil {
		if first == true {
			first = false
			query += "username = ?"
		} else {
			query += ", username = ?"
		}
		args = append(args, *customer.Username)
	}

	if customer.Password != nil {
		if first == true {
			first = false
			query += "password = ?"
		} else {
			query += ", password = ?"
		}
		args = append(args, *customer.Password)
	}

	if customer.Profile != nil {
		if first == true {
			first = false
			query += "profile = ?"
		} else {
			query += ", profile = ?"
		}
		args = append(args, *customer.Profile)
	}

	if customer.Type != nil {
		if first == true {
			first = false
			query += "type = ?"
		} else {
			query += ", type = ?"
		}
		args = append(args, *customer.Type)
	}

	query += " WHERE username = ? "
	args = append(args, username)

	_, err = tx.ExecContext(ctx, query)
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

func (m *mysqlRepository) Refill(ctx context.Context, customer domain.Customer, expiration string) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// and save radpackage, radcheck, invoice

	// radpackage
	query := "UPDATE radcheck SET expiration = ? WHERE username = ? "
	_, err = tx.ExecContext(ctx, query, expiration, *customer.Username)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
	query = "INSERT INTO invoice(no_invoice, type, name, nominal) VALUES(?,?,?,((SELECT price FROM package WHERE id = ?) + (SELECT margin FROM package WHERE id = ?)))"
	_, err = tx.ExecContext(ctx, query,
		uniqTime,
		*customer.Type,
		*customer.Username,
		*customer.IDPackage,
		*customer.IDPackage,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	if *customer.Type == "prepaid" {
		query = "INSERT INTO payment(no_invoice, nominal) VALUES(?,((SELECT price FROM package WHERE id = ?) + (SELECT margin FROM package WHERE id = ?)))"
		_, err = tx.ExecContext(
			ctx,
			query,
			uniqTime,
			*customer.IDPackage, *customer.IDPackage,
		)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return
}

func (m *mysqlRepository) Insert(ctx context.Context, customer domain.Customer) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// and save radpackage, radcheck, invoice

	// radpackage
	query := "INSERT INTO radpackage(id_package, username) VALUES(?,?)"
	_, err = tx.ExecContext(ctx, query, *customer.IDPackage, *customer.Username)
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
		*customer.Username, "Cleartext-Password", ":=", *customer.Password,
		*customer.Username, "User-Profile", ":=", *customer.Profile,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
	// invoice
	query = "INSERT INTO invoice(no_invoice, type, name, nominal) VALUES(?,?,?,((SELECT price FROM package WHERE id = ?) + (SELECT margin FROM package WHERE id = ?)))"
	_, err = tx.ExecContext(
		ctx,
		query,
		uniqTime,
		*customer.Type,
		*customer.Name,
		*customer.IDPackage, *customer.IDPackage,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	if *customer.Type == "prepaid" {
		query = "INSERT INTO payment(no_invoice, nominal) VALUES(?,((SELECT price FROM package WHERE id = ?) + (SELECT margin FROM package WHERE id = ?)))"
		_, err = tx.ExecContext(
			ctx,
			query,
			uniqTime,
			*customer.IDPackage, *customer.IDPackage,
		)
		if err != nil {
			logrus.Error(err)

			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return
}

func (m *mysqlRepository) Delete(ctx context.Context, id int64) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// and delete radcheck, customer

	// radpackage
	query := "DELETE FROM radcheck WHERE username = (SELECT username FROM customer WHERE id = ?)"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// customer
	query = "DELETE FROM customer WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, id)
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
