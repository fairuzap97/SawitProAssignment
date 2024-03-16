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

type RegisterUserTestSuite struct {
	suite.Suite

	gomock *gomock.Controller
	repo   *repository.MockUserRepository

	usecase usecase.UserUsecases

	createUserInput  repository.CreateUserInput
	createUserOutput repository.CreateUserOutput

	input  usecase.RegisterUserInput
	output usecase.RegisterUserOutput

	ctx     context.Context
	mockErr error
}

func TestRegisterUserTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterUserTestSuite))
}

func (s *RegisterUserTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.repo = repository.NewMockUserRepository(s.gomock)

	s.usecase = NewUserUsecases(NewUserUsecasesOptions{UserRepo: s.repo})
	generateBcryptHash = func(password string) ([]byte, error) {
		return []byte("password-hash"), nil
	}

	s.createUserInput = repository.CreateUserInput{
		PhoneNo:      "+62812151833",
		FullName:     "John Smith",
		PasswordHash: []byte("password-hash"),
	}
	s.createUserOutput = repository.CreateUserOutput{ID: 123}

	s.input = usecase.RegisterUserInput{
		PhoneNo:  "+62812151833",
		FullName: "John Smith",
		Password: "Sample-Val1d-Passw0rd",
	}
	s.output = usecase.RegisterUserOutput{UserID: 123}

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *RegisterUserTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *RegisterUserTestSuite) TestValidationFailed() {
	a := assert.New(s.T())

	s.input.PhoneNo = "no"
	s.input.FullName = "L"
	s.input.Password = "bad"
	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(out)
	var validationError usecase.ValidationErrors
	a.True(errors.As(err, &validationError))
	validationErrors := validationError.GetErrors()
	a.Equal("phone_no must be between 10 and 13 characters long", validationErrors["phone_no"][0].Error())
	a.Equal("phone_no must start with indonesia country code (\"+62\")", validationErrors["phone_no"][1].Error())
	a.Equal("besides the country code, phone_no must only contain numbers", validationErrors["phone_no"][2].Error())
	a.Equal("full_name must be between 3 and 60 characters long", validationErrors["full_name"][0].Error())
	a.Equal("password must be between 6 and 64 characters long", validationErrors["password"][0].Error())
	a.Equal("password must contains at least one number [0-9]", validationErrors["password"][1].Error())
	a.Equal("password must contains at least one capital letter [A-Z]", validationErrors["password"][2].Error())
	a.Equal("password must contains at least one non alphanumeric character", validationErrors["password"][3].Error())
}

func (s *RegisterUserTestSuite) TestValidationPartialFailed() {
	a := assert.New(s.T())

	s.input.PhoneNo = "+621419162299771122"
	s.input.FullName = "His Excellency John Smith Alexander Hamilton The III, Third of His Name"
	s.input.Password = "NotAnAwfulPasswordButMissingSomeStuffAsWellAsUnnecessaryLongLikeReallyWhyDidYouWriteThisLong"
	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(out)
	var validationError usecase.ValidationErrors
	a.True(errors.As(err, &validationError))
	validationErrors := validationError.GetErrors()
	a.Equal("phone_no must be between 10 and 13 characters long", validationErrors["phone_no"][0].Error())
	a.Equal("full_name must be between 3 and 60 characters long", validationErrors["full_name"][0].Error())
	a.Equal("password must be between 6 and 64 characters long", validationErrors["password"][0].Error())
	a.Equal("password must contains at least one number [0-9]", validationErrors["password"][1].Error())
	a.Equal("password must contains at least one non alphanumeric character", validationErrors["password"][2].Error())
}

func (s *RegisterUserTestSuite) TestBcryptError() {
	a := assert.New(s.T())

	generateBcryptHash = func(password string) ([]byte, error) {
		return nil, s.mockErr
	}
	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *RegisterUserTestSuite) TestRepositoryError() {
	a := assert.New(s.T())

	s.repo.EXPECT().CreateUser(s.ctx, s.createUserInput).Return(repository.CreateUserOutput{}, s.mockErr)

	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *RegisterUserTestSuite) TestUserConflict() {
	a := assert.New(s.T())

	s.repo.EXPECT().CreateUser(s.ctx, s.createUserInput).Return(repository.CreateUserOutput{}, repository.ErrorRecordConflict)

	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserConflictError)
}

func (s *RegisterUserTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.repo.EXPECT().CreateUser(s.ctx, s.createUserInput).Return(s.createUserOutput, nil)

	out, err := s.usecase.RegisterUser(s.ctx, s.input)

	a.Empty(err)
	a.Equal(s.output, out)
}
