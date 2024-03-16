package users

import (
	"context"
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
)

func (u *userUsecases) GetUserProfile(ctx context.Context, input usecase.GetUserProfileInput) (output usecase.GetUserProfileOutput, err error) {
	var resp repository.GetUserOutput
	if resp, err = u.userRepo.GetUser(ctx, repository.GetUserInput{ID: input.UserID}); err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			err = usecase.UserNotFoundError
		}
		return
	}
	return usecase.GetUserProfileOutput{
		UserID:               resp.ID,
		PhoneNo:              resp.PhoneNo,
		FullName:             resp.FullName,
		SuccessfulLoginCount: resp.SuccessfulLoginCount,
	}, nil
}
