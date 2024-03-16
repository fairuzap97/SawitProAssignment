package users

import (
	"context"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func (u *userUsecases) ValidateUserToken(ctx context.Context, input usecase.ValidateUserTokenInput) (output usecase.ValidateUserTokenOutput, err error) {
	parsedToken, err := jwt.Parse(input.JwtToken, func(token *jwt.Token) (interface{}, error) {
		return &u.jwtSecret.PublicKey, nil
	})
	if err != nil {
		err = usecase.UserInvalidToken
		return
	}

	// ensure token not expired
	exp, err := parsedToken.Claims.GetExpirationTime()
	if exp == nil || err != nil {
		err = usecase.UserInvalidToken
		return
	}
	if exp.Before(time.Now()) {
		err = usecase.UserInvalidToken
		return
	}

	// get userID from token subject
	subject, err := parsedToken.Claims.GetSubject()
	if err != nil {
		err = usecase.UserInvalidToken
		return
	}
	output.UserID, err = strconv.ParseUint(subject, 10, 64)
	if err != nil {
		err = usecase.UserInvalidToken
		return
	}

	return output, nil
}
