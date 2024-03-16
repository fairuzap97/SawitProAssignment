package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/lib/pq"
	"strings"
)

func (u *userRepository) UpdateUser(ctx context.Context, input repository.UpdateUserInput) (output repository.UpdateUserOutput, err error) {
	var updates []string
	var params []interface{}
	id := 1
	if input.FullName != "" {
		updates = append(updates, fmt.Sprintf("full_name=$%d", id))
		params = append(params, input.FullName)
		id += 1
	}
	if input.PhoneNo != "" {
		updates = append(updates, fmt.Sprintf("phone_no=$%d", id))
		params = append(params, input.PhoneNo)
		id += 1
	}
	if input.SuccessfulLoginCount != 0 {
		updates = append(updates, fmt.Sprintf("successful_login_count=$%d", id))
		params = append(params, input.SuccessfulLoginCount)
		id += 1
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(updates, ", "), id)
	params = append(params, input.ID)
	var result sql.Result
	if result, err = u.db.ExecContext(ctx, query, params...); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique constraint violation
				err = repository.ErrorRecordConflict
			}
		}
		return
	}

	if affected, _ := result.RowsAffected(); affected <= 0 {
		err = repository.ErrorRecordNotFound
		return
	}

	return output, nil
}
