package users

import (
	"context"
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
)

func (u *userUsecases) UpdateUserProfile(ctx context.Context, input usecase.UpdateUserProfileInput) (output usecase.UpdateUserProfileOutput, err error) {
	if validationErrors := validateUpdateUserPayload(input); len(validationErrors) > 0 {
		err = usecase.NewValidationError(validationErrors)
		return
	}

	var updatePayload = repository.UpdateUserInput{ID: input.UserID}
	if input.FullName != nil {
		updatePayload.FullName = *input.FullName
	}
	if input.PhoneNo != nil {
		updatePayload.PhoneNo = *input.PhoneNo
	}
	if _, err = u.userRepo.UpdateUser(ctx, updatePayload); err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			err = usecase.UserNotFoundError
		} else if errors.Is(err, repository.ErrorRecordConflict) {
			err = usecase.UserConflictError
		}
		return
	}
	return usecase.UpdateUserProfileOutput{}, nil
}

func validateUpdateUserPayload(input usecase.UpdateUserProfileInput) map[string][]error {
	var validationErrors = map[string][]error{}
	if input.FullName != nil {
		if errs := validateUserFullName(*input.FullName); len(errs) > 0 {
			validationErrors["full_name"] = errs
		}
	}
	if input.PhoneNo != nil {
		if errs := validateUserPhoneNo(*input.PhoneNo); len(errs) > 0 {
			validationErrors["phone_no"] = errs
		}
	}
	return validationErrors
}
