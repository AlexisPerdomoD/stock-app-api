package entity

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
