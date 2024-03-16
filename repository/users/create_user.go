package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/lib/pq"
)

const (
	createUserQuery = `INSERT INTO "users" (phone_no, full_name, password_hash, successful_login_count) VALUES ($1, $2, $3, 0) RETURNING id;`
)

func (u *userRepository) CreateUser(ctx context.Context, input repository.CreateUserInput) (output repository.CreateUserOutput, err error) {
	var result *sql.Rows
	if result, err = u.db.QueryContext(ctx, createUserQuery, input.PhoneNo, input.FullName, input.PasswordHash); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique constraint violation
				err = repository.ErrorRecordConflict
			}
		}
		return
	}

	result.Next()
	if err = result.Scan(&output.ID); err != nil {
		return
	}
	return output, nil
}
