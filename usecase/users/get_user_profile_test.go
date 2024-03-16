package users

import (
	"context"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GetUserProfileTestSuite struct {
	suite.Suite

	gomock *gomock.Controller
	repo   *repository.MockUserRepository

	usecase usecase.UserUsecases

	getUserInput  repository.GetUserInput
	getUserOutput repository.GetUserOutput

	input  usecase.GetUserProfileInput
	output usecase.GetUserProfileOutput

	ctx     context.Context
	mockErr error
}

func TestGetUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserProfileTestSuite))
}

func (s *GetUserProfileTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.repo = repository.NewMockUserRepository(s.gomock)

	s.usecase = NewUserUsecases(NewUserUsecasesOptions{UserRepo: s.repo})

	s.getUserInput = repository.GetUserInput{ID: 123}
	s.getUserOutput = repository.GetUserOutput{
		ID:                   123,
		PhoneNo:              "+6281215183300",
		FullName:             "John Smith",
		PasswordHash:         []byte("password-hash-here"),
		SuccessfulLoginCount: 2,
	}

	s.input = usecase.GetUserProfileInput{UserID: 123}
	s.output = usecase.GetUserProfileOutput{
		UserID:               123,
		PhoneNo:              "+6281215183300",
		FullName:             "John Smith",
		SuccessfulLoginCount: 2,
	}

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *GetUserProfileTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *GetUserProfileTestSuite) TestRepositoryError() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(repository.GetUserOutput{}, s.mockErr)

	out, err := s.usecase.GetUserProfile(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *GetUserProfileTestSuite) TestUserNotFound() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(repository.GetUserOutput{}, repository.ErrorRecordNotFound)

	out, err := s.usecase.GetUserProfile(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserNotFoundError)
}

func (s *GetUserProfileTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(s.getUserOutput, nil)

	out, err := s.usecase.GetUserProfile(s.ctx, s.input)

	a.Empty(err)
	a.Equal(s.output, out)
}
