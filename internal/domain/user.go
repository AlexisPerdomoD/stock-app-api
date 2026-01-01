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

func (u User) GetID() uint64 {
	return u.ID
}

type UserRepository interface {
	/*
		Returns a user by its id, if not found, returns nil
	*/
	GetByID(ctx context.Context, id uint64) (*User, error)

	/*
		Returns a user by its id including password field, if not found, returns nil
	*/
	GetByIDWithPassword(ctx context.Context, id uint64) (*User, error)

	/*
		Returns a user by its username, if not found, returns nil
	*/
	GetByUsername(ctx context.Context, username string) (*User, error)

	/*
		Returns a user by its username(incliding password field), if not found, returns nil
	*/
	GetByUsernameWithPassword(ctx context.Context, username string) (*User, error)

	/*
		Saves a user and maps id and CreatedAt fields, password field is remaped to nil always.

		- returns error if nil arguments are passed.

		- any persistence constraints violated by any argument returns an error (e.g unique indexes like username).
	*/
	Save(ctx context.Context, args *User) error

	/*
		returns true if the user has the stock saved, false otherwise.
	*/
	HasUserStock(ctx context.Context, userID uint64, stockID uint64) (bool, error)

	/*
		Registers the user stock association to save the stock as a favorite.

		- any persistence constraints violated by any argument returns an error (e.g unique constraint or not existing user or stock).
	*/
	RegisterUserStock(ctx context.Context, userID uint64, stockID uint64) error

	/*
		Removes the user stock association to the stock as a favorite.

		- If register does not exists this is no-op.
	*/
	RemoveUserStock(ctx context.Context, userID uint64, stockID uint64) error
}
