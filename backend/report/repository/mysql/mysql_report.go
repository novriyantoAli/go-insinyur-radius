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
func NewMysqlRepository(conn *sql.DB) domain.ReportRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) getReportFinance(ctx context.Context, query string, args ...interface{}) (res domain.ReportFinance, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return domain.ReportFinance{}, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	for rows.Next() {
		err = rows.Scan(&res.Value)

		if err != nil {
			logrus.Error(err)
			return domain.ReportFinance{}, err
		}
	}

	if res.Value == nil  {
		var defaultValue int64 = 0
		res.Value = &defaultValue
	}

	return
}

func (m *mysqlRepository) ReportFinanceCurrentDay(ctx context.Context) (res []domain.ReportFinance, err error) {
	// Invoice
	invoice := "Invoice"
	query := `SELECT SUM(nominal) as value FROM invoice WHERE created_at = CURDATE()`
	result1, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	
	fmt.Println(result1)

	result1.Name = &invoice

	res = append(res, result1)

	payment := "Payment"
	query = `SELECT SUM(nominal) as value FROM payment WHERE created_at = CURDATE()`
	result2, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	result2.Name = &payment
	res = append(res, result2)

	credit := "Credit"
	v := *res[0].Value - *res[1].Value
	resultDebit := domain.ReportFinance{Name: &credit, Value: &v}
	res = append(res, resultDebit)

	return
}

func (m *mysqlRepository) ReportFinanceCurrentMonth(ctx context.Context) (res []domain.ReportFinance, err error) {
	// Invoice
	invoice := "Invoice"
	query := `SELECT SUM(nominal) as value FROM invoice WHERE MONTH(created_at) = MONTH(CURRENT_DATE()) AND YEAR(created_at) = YEAR(CURRENT_DATE())`
	result1, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	result1.Name = &invoice
	res = append(res, result1)

	payment := "Payment"
	query = `SELECT SUM(nominal) as value FROM payment WHERE MONTH(created_at) = MONTH(CURRENT_DATE()) AND YEAR(created_at) = YEAR(CURRENT_DATE())`
	result2, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	result2.Name = &payment
	res = append(res, result2)

	// debit
	credit := "Credit"
	v := *res[0].Value - *res[1].Value
	resultDebit := domain.ReportFinance{Name: &credit, Value: &v}
	res = append(res, resultDebit)

	return
}

func (m *mysqlRepository) ReportFinanceCurrentYear(ctx context.Context) (res []domain.ReportFinance, err error) {
	// Invoice
	invoice := "Invoice"
	query := `SELECT SUM(nominal) as value FROM invoice WHERE YEAR(created_at) = YEAR(CURRENT_DATE())`
	result, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	result.Name = &invoice
	res = append(res, result)

	payment := "Payment"
	query = `SELECT SUM(nominal) as value FROM payment WHERE YEAR(created_at) = YEAR(CURRENT_DATE())`
	result1, err := m.getReportFinance(ctx, query)
	if err != nil {
		return nil, err
	}
	result1.Name = &payment
	res = append(res, result1)

	// debit
	credit := "Credit"
	v := *res[0].Value - *res[1].Value
	resultDebit := domain.ReportFinance{Name: &credit, Value: &v}
	res = append(res, resultDebit)

	return
}
