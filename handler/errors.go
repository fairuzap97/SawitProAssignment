package handler

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/usecase"
)

var (
	JsonBodyInvalid = fmt.Errorf("invalid JSON Body")

	errorsToCodeMap = map[error]int{
		JsonBodyInvalid: 400,

		usecase.UserInvalidLogin:  400,
		usecase.UserInvalidToken:  403,
		usecase.UserConflictError: 409,
		usecase.UserNotFoundError: 404,
	}
)
