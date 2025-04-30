package domain

import "context"

type User struct {
	ID       uint
	UserName string
	Password []byte
	Active   bool
}

type UserRepository interface {
	Get(ctx context.Context, id uint) (*User, error)

	GetByUsername(ctx context.Context, username string) (*User, error)

	Create(ctx context.Context, args *User) error

	RegisterUserStock(ctx context.Context, userID uint, stockID uint) error

	RemoveUserStock(ctx context.Context, userID uint, stockID uint) error
}
