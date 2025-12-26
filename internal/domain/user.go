package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint64
	Username  string
	Firstname string
	Lastname  string
	Password  []byte
	Active    bool
	CreatedAt time.Time
}

type UserRepository interface {
	Get(ctx context.Context, id uint) (*User, error)

	GetByUsername(ctx context.Context, username string) (*User, error)

	Create(ctx context.Context, args *User) error

	RegisterUserStock(ctx context.Context, userID uint, stockID uint) error

	RemoveUserStock(ctx context.Context, userID uint, stockID uint) error
}
