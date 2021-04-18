package mysql

import (
	"context"
	"database/sql"
	"insinyur-radius/domain"

	"github.com/sirupsen/logrus"
)

// mysqlUsersRepository ...
type mysqlUsersRepository struct {
	Conn *sql.DB
}

// NewMysqlUsersRepository ...
func NewMysqlUsersRepository(conn *sql.DB) domain.UsersRepository{
	return &mysqlUsersRepository{ Conn: conn }
}

func (m *mysqlUsersRepository) fetch(c context.Context, query string, args... interface{})(res []domain.Users, er error){
	rows, err := m.Conn.QueryContext(c, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func(){
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	res = make([]domain.Users, 0)
	for rows.Next(){
		t := domain.Users{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Password,
			&t.Expiration,
			&t.Profile,
			&t.PackageName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (m *mysqlUsersRepository) CountPage(ctx context.Context)(res int64, err error){
	query := `SELECT 
				COUNT(r1.username) as username
			FROM 
				radcheck r1 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = 'Expiration') r2 
				ON 
					r2.username = r1.username 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = 'User-Profile') r3 
				ON 
					r3.username = r1.username 
			WHERE 
				r1.attribute = 'Cleartext-Password' ORDER BY r1.id DESC`
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
			
	defer func(){
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()
			
	for rows.Next(){
		err = rows.Scan(&res)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
	}
		
	return
}

func (m *mysqlUsersRepository) Fetch(ctx context.Context, id int64, limit int64)(res []domain.Users, err error){
	query := `SELECT 
				r1.id, r1.username as username, r1.value as password, 
				r2.value as expiration, 
				r3.value as profile, 
				package.name as package
			FROM 
				radcheck r1 
			LEFT JOIN
				radpackage
				ON
					radpackage.username = r1.username
			LEFT JOIN
				package
				ON
					radpackage.id_package = package.id
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r2 
				ON 
					r2.username = r1.username 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r3 
				ON 
					r3.username = r1.username 
			WHERE 
				r1.attribute = ? `
	if id != 0 {
		query += "AND r1.id < ? "
	}

	query += "ORDER BY r1.id DESC LIMIT ? "

	if id != 0 {
		res, err = m.fetch(
			ctx, query, 
			"Expiration", "User-Profile", "Cleartext-Password", id, limit,
		)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = m.fetch(
			ctx, query, 
			"Expiration", "User-Profile", "Cleartext-Password", limit,
		)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (m *mysqlUsersRepository) Get(ctx context.Context, username string)(res domain.Users, err error){
	query := `SELECT 
				r1.id, r1.username as username, r1.value as password, 
				r2.value as expiration, 
				r3.value as profile,
				package.name as package
			FROM 
				radcheck r1 
			LEFT JOIN
				radpackage
				ON
					radpackage.username = r1.username
			LEFT JOIN
				package
				ON
					radpackage.id_package = package.id 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r2 
				ON 
					r2.username = r1.username 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r3 
				ON 
					r3.username = r1.username 
			WHERE 
				r1.attribute = ? AND r1.username = ? `
	resArr, err := m.fetch(ctx, query, "Expiration", "User-Profile", "Cleartext-Password", username)
	if err != nil {
		return domain.Users{}, err
	}

	if len(resArr) > 0 {
		res = resArr[0]
	} else {
		return domain.Users{}, err
	}
	
	
	return 
}

func (m *mysqlUsersRepository) ReportExpirationToday(ctx context.Context)(res []domain.Users, err error){
	query := `SELECT 
				r1.id, r1.username as username, r1.value as password, 
				r2.value as expiration, 
				r3.value as profile,
				package.name as package
			FROM 
				radcheck r1 
			LEFT JOIN
				radpackage
			ON
				radpackage.username = r1.username
			LEFT JOIN
				package
			ON
				radpackage.id_package = package.id 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r2 
			ON 
				r2.username = r1.username 
			LEFT JOIN 
				(SELECT * FROM radcheck WHERE attribute = ?) r3 
			ON 
				r3.username = r1.username 
			WHERE 
				r1.attribute = ? AND MONTH(STR_TO_DATE(r2.value, '%d %b %Y')) = MONTH(CURRENT_DATE())`
	
	res, err = m.fetch(
		ctx, query, 
		"Expiration", "User-Profile", "Cleartext-Password",
	)

	return
}

func (m *mysqlUsersRepository) UpdateUsers(ctx context.Context, user *domain.Users)(err error){
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		logrus.Error(err)
		return err
	}

	query := "UPDATE radpackage SET id_package = ? WHERE username = ?"
	_, err = tx.ExecContext(ctx, query, *user.Package, *user.Username)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	// radcheck
	query = "UPDATE radcheck SET value = ? WHERE attribute = ? AND username = ? "
	_, err = tx.ExecContext(
		ctx, query, *user.Password, "Cleartext-Password", *user.Username,
	)
	if err != nil {
		logrus.Error(err)
	
		tx.Rollback()
		return err
	}

	query = "UPDATE radcheck SET value = ? WHERE attribute = ? AND value = ? "
	_, err = tx.ExecContext(
		ctx, query, *user.Profile, "User-Profile", *user.Username,
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

func (m *mysqlUsersRepository) SaveUsers(ctx context.Context, users *domain.Users)(err error){
	tx, err := m.Conn.BeginTx(ctx, nil)

	// and save radpackage, radcheck, transaction

	// radpackage
	query := "INSERT INTO radpackage(id_package, username) VALUES(?,?)"
	_, err = tx.ExecContext(ctx, query, *users.Package, *users.Username)
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
		*users.Username, "Cleartext-Password", ":=", *users.Password,
		*users.Username, "User-Profile", ":=", *users.Profile,
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
