package domain

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
	Get(id int) (*User, error)

	GetByUsername(username string) (*User, error)

	Create(args *User) error
}
