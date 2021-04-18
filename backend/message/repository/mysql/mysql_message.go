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
func NewMysqlRepository(conn *sql.DB) domain.MessageRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Message, err error) {
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

	result = make([]domain.Message, 0)
	for rows.Next() {
		t := domain.Message{}
		err = rows.Scan(
			&t.ID,
			&t.ChatID,
			&t.MessageID,
			&t.Received,
			&t.Message,
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

func (m *mysqlRepository) Update(ctx context.Context, message domain.Message) (err error) {
	args := make([]interface{}, 0)

	query := `UPDATE message SET chat_id = ?, message_id = ?`
	// To parameter not null
	args = append(args, *message.ChatID)
	args = append(args, *message.MessageID)

	if message.Received != nil && *message.Received != "" {
		query += ", received = ? "
		args = append(args, *message.Received)
	}

	if message.Message != nil && *message.Message != "" {
		query += ", message = ? "
		args = append(args, *message.Message)
	}

	query += " WHERE id = ? "
	args = append(args, *message.ID)

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, args...)
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

func (m *mysqlRepository) Insert(ctx context.Context, message domain.Message) (err error) {
	query := `INSERT message SET chat_id = ?, message_id = ?, message = ? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, message.ChatID, message.MessageID, message.Message)
	if err != nil {
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}

func (m *mysqlRepository) Find(ctx context.Context, spec domain.Message) (res []domain.Message, err error) {
	query := `SELECT * FROM message `

	args := make([]interface{}, 0)

	whereSelected := false
	if spec.ID != nil && *spec.ID != 0 {
		query += ` WHERE id = ? `
		whereSelected = true
		args = append(args, *spec.ID)
	}

	if spec.ChatID != nil && *spec.ChatID != 0 {
		if whereSelected == false {
			query += ` WHERE chat_id = ? `
			whereSelected = true
		} else {
			query += ` AND chat_id = ? `
		}
		args = append(args, *spec.ChatID)
	}

	if spec.MessageID != nil && *spec.MessageID != 0 {
		if whereSelected == false {
			query += ` WHERE message_id = ? `
			whereSelected = true
		} else {
			query += ` AND message_id = ? `
		}
		args = append(args, *spec.MessageID)
	}

	if spec.Received != nil && *spec.Received != "" {
		if whereSelected == false {
			query += ` WHERE received = ? `
			whereSelected = true
		} else {
			query += `AND received = ? `
		}
		args = append(args, *spec.Received)
	}

	if spec.Message != nil && *spec.Message != "" {
		if whereSelected == false {
			query += ` WHERE message = ? `
			whereSelected = true
		} else {
			query += ` AND message = ? `
		}
		args = append(args, *spec.Message)
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return
}
