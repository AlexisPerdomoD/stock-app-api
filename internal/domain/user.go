package domain

import "context"

type User struct {
	ID       int
	UserName string
	Password string
	Active   bool
}

type CreateUserArgs struct {
	UserName string
	Password string
}

type UserRepository interface {
	Get(ctx context.Context, id int) (*User, error)

	GetByUsername(ctx context.Context, username string) (*User, error)

	Create(ctx context.Context, args *User) error
}
