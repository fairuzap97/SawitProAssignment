package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UpdateUserProfileTestSuite struct {
	suite.Suite

	gomock *gomock.Controller
	repo   *repository.MockUserRepository

	usecase usecase.UserUsecases

	updateUserInput  repository.UpdateUserInput
	updateUserOutput repository.UpdateUserOutput

	input  usecase.UpdateUserProfileInput
	output usecase.UpdateUserProfileOutput

	ctx     context.Context
	mockErr error
}

func TestUpdateUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserProfileTestSuite))
}

func (s *UpdateUserProfileTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.repo = repository.NewMockUserRepository(s.gomock)

	s.usecase = NewUserUsecases(NewUserUsecasesOptions{UserRepo: s.repo})

	s.updateUserInput = repository.UpdateUserInput{
		ID:       123,
		PhoneNo:  "+62812151833",
		FullName: "John Smith",
	}
	s.updateUserOutput = repository.UpdateUserOutput{}

	s.input = usecase.UpdateUserProfileInput{
		UserID:   123,
		PhoneNo:  &s.updateUserInput.PhoneNo,
		FullName: &s.updateUserInput.FullName,
	}
	s.output = usecase.UpdateUserProfileOutput{}

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *UpdateUserProfileTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *UpdateUserProfileTestSuite) TestValidationFailed() {
	a := assert.New(s.T())

	badString := "no"
	s.input.PhoneNo = &badString
	s.input.FullName = &badString
	out, err := s.usecase.UpdateUserProfile(s.ctx, s.input)

	a.Empty(out)
	var validationError usecase.ValidationErrors
	a.True(errors.As(err, &validationError))
	validationErrors := validationError.GetErrors()
	a.Equal("phone_no must be between 10 and 13 characters long", validationErrors["phone_no"][0].Error())
	a.Equal("phone_no must start with indonesia country code (\"+62\")", validationErrors["phone_no"][1].Error())
	a.Equal("besides the country code, phone_no must only contain numbers", validationErrors["phone_no"][2].Error())
	a.Equal("full_name must be between 3 and 60 characters long", validationErrors["full_name"][0].Error())
}

func (s *UpdateUserProfileTestSuite) TestRepositoryError() {
	a := assert.New(s.T())

	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(repository.UpdateUserOutput{}, s.mockErr)

	out, err := s.usecase.UpdateUserProfile(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *UpdateUserProfileTestSuite) TestUserNotFound() {
	a := assert.New(s.T())

	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(repository.UpdateUserOutput{}, repository.ErrorRecordNotFound)

	out, err := s.usecase.UpdateUserProfile(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserNotFoundError)
}

func (s *UpdateUserProfileTestSuite) TestUserConflict() {
	a := assert.New(s.T())

	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(repository.UpdateUserOutput{}, repository.ErrorRecordConflict)

	out, err := s.usecase.UpdateUserProfile(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserConflictError)
}

func (s *UpdateUserProfileTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(s.updateUserOutput, nil)

	out, err := s.usecase.UpdateUserProfile(s.ctx, s.input)

	a.Empty(err)
	a.Equal(s.output, out)
}
