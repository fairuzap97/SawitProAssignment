// This file contains the repository implementation layer.
package users

import (
	"database/sql"
	"github.com/SawitProRecruitment/UserService/repository"

	_ "github.com/lib/pq"
)

// userRepository is a postgresSQL implementation of repository.UserRepository
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(opts repository.NewRepositoryOptions) (repository.UserRepository, error) {
	db, err := sql.Open("postgres", opts.Dsn)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		db: db,
	}, nil
}
