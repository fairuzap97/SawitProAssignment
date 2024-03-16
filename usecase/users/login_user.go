package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u *userUsecases) LoginUser(ctx context.Context, input usecase.LoginUserInput) (output usecase.LoginUserOutput, err error) {
	var usr repository.GetUserOutput
	if usr, err = u.userRepo.GetUser(ctx, repository.GetUserInput{PhoneNo: input.PhoneNo}); err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			err = usecase.UserInvalidLogin
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(input.Password)); err != nil {
		err = usecase.UserInvalidLogin
		return
	}

	updatePayload := repository.UpdateUserInput{ID: usr.ID, SuccessfulLoginCount: usr.SuccessfulLoginCount + 1}
	if _, err = u.userRepo.UpdateUser(ctx, updatePayload); err != nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", usr.ID),
		"exp": time.Now().Add(u.jwtTtl).Unix(),
	})
	if output.JwtToken, err = token.SignedString(u.jwtSecret); err != nil {
		return
	}

	return output, nil
}
