package users

import (
	"context"
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"golang.org/x/crypto/bcrypt"
)

var (
	// wrapper to bcrypt function call to make testing easier
	generateBcryptHash = func(password string) ([]byte, error) {
		return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	}
)

func (u *userUsecases) RegisterUser(ctx context.Context, input usecase.RegisterUserInput) (output usecase.RegisterUserOutput, err error) {
	if validationErrors := validateRegisterUserPayload(input); len(validationErrors) > 0 {
		err = usecase.NewValidationError(validationErrors)
		return
	}

	// bcrypt package implement password salting & hashing and output it into a single byte array to be stored
	var passwordHash []byte
	if passwordHash, err = generateBcryptHash(input.Password); err != nil {
		return
	}

	var resp repository.CreateUserOutput
	resp, err = u.userRepo.CreateUser(ctx, repository.CreateUserInput{
		PhoneNo:      input.PhoneNo,
		FullName:     input.FullName,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrorRecordConflict) {
			err = usecase.UserConflictError
		}
		return
	}

	return usecase.RegisterUserOutput{UserID: resp.ID}, nil
}

func validateRegisterUserPayload(input usecase.RegisterUserInput) map[string][]error {
	var validationErrors = map[string][]error{}
	if errs := validateUserFullName(input.FullName); len(errs) > 0 {
		validationErrors["full_name"] = errs
	}
	if errs := validateUserPhoneNo(input.PhoneNo); len(errs) > 0 {
		validationErrors["phone_no"] = errs
	}
	if errs := validateUserPassword(input.Password); len(errs) > 0 {
		validationErrors["password"] = errs
	}
	return validationErrors
}
