// This file contains types that are used in the repository layer.
package repository

type NewRepositoryOptions struct {
	Dsn string
}

type CreateUserInput struct {
	PhoneNo      string
	FullName     string
	PasswordHash []byte
}

type CreateUserOutput struct {
	ID uint64
}

type GetUserInput struct {
	ID      uint64
	PhoneNo string
}

type GetUserOutput struct {
	ID                   uint64
	PhoneNo              string
	FullName             string
	PasswordHash         []byte
	SuccessfulLoginCount uint64
}

type UpdateUserInput struct {
	ID                   uint64
	PhoneNo              string
	FullName             string
	SuccessfulLoginCount uint64
}

type UpdateUserOutput struct {
}
