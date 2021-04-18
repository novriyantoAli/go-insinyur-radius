package domain

import (
	"context"
	"time"
)

// Message ...
type Message struct {
	ID        *int64
	ChatID    *int64
	MessageID *int64
	Received  *string
	Message   *string
	CreatedAt *time.Time
}

// MessageRepository ...
type MessageRepository interface {
	Find(ctx context.Context, spec Message) (res []Message, err error)
	Insert(ctx context.Context, message Message)(err error)
	Update(ctx context.Context, message Message)(err error)
}

// MessageUsecase ...
type MessageUsecase interface {
	Find(ctx context.Context, spec Message) (res []Message, err error)
	Insert(ctx context.Context, message Message)(err error)
	Update(ctx context.Context, message Message)(err error)
}
