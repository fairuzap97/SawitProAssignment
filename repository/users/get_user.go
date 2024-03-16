package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
)

const (
	getUserByIDQuery      = `SELECT id, phone_no, full_name, password_hash, successful_login_count FROM users WHERE id=$1;`
	getUserByPhoneNoQuery = `SELECT id, phone_no, full_name, password_hash, successful_login_count FROM users WHERE phone_no=$1;`
)

func (u *userRepository) GetUser(ctx context.Context, input repository.GetUserInput) (output repository.GetUserOutput, err error) {
	var row *sql.Row
	if input.PhoneNo != "" {
		row = u.db.QueryRowContext(ctx, getUserByPhoneNoQuery, input.PhoneNo)
	} else {
		row = u.db.QueryRowContext(ctx, getUserByIDQuery, input.ID)
	}

	if err = row.Scan(&output.ID, &output.PhoneNo, &output.FullName, &output.PasswordHash, &output.SuccessfulLoginCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = repository.ErrorRecordNotFound
			return
		}
		return
	}

	return output, nil
}
