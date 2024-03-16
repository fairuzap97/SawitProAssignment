package users

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

type LoginUserTestSuite struct {
	suite.Suite

	gomock *gomock.Controller
	repo   *repository.MockUserRepository

	jwtSecret *rsa.PrivateKey
	usecase   usecase.UserUsecases

	getUserInput  repository.GetUserInput
	getUserOutput repository.GetUserOutput

	updateUserInput  repository.UpdateUserInput
	updateUserOutput repository.UpdateUserOutput

	input usecase.LoginUserInput

	ctx     context.Context
	mockErr error
}

func TestLoginUserTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUserTestSuite))
}

func (s *LoginUserTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.repo = repository.NewMockUserRepository(s.gomock)

	s.jwtSecret, _ = rsa.GenerateKey(rand.Reader, 1024)
	s.usecase = NewUserUsecases(NewUserUsecasesOptions{
		UserRepo:  s.repo,
		JwtSecret: s.jwtSecret,
		JwtTtl:    time.Minute * 5,
	})

	passwdHash, _ := bcrypt.GenerateFromPassword([]byte("SomeVal1dPassw@rd"), bcrypt.DefaultCost)
	s.getUserInput = repository.GetUserInput{PhoneNo: "+62812151833"}
	s.getUserOutput = repository.GetUserOutput{
		ID:                   123,
		PhoneNo:              "+62812151833",
		FullName:             "John Smith",
		PasswordHash:         passwdHash,
		SuccessfulLoginCount: 2,
	}

	s.updateUserInput = repository.UpdateUserInput{ID: 123, SuccessfulLoginCount: 3}
	s.updateUserOutput = repository.UpdateUserOutput{}

	s.input = usecase.LoginUserInput{
		PhoneNo:  "+62812151833",
		Password: "SomeVal1dPassw@rd",
	}

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *LoginUserTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *LoginUserTestSuite) TestRepositoryError() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(repository.GetUserOutput{}, s.mockErr)

	out, err := s.usecase.LoginUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *LoginUserTestSuite) TestUserNotFound() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(repository.GetUserOutput{}, repository.ErrorRecordNotFound)

	out, err := s.usecase.LoginUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidLogin)
}

func (s *LoginUserTestSuite) TestInvalidPassword() {
	a := assert.New(s.T())

	s.input.Password = "invalid-password"
	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(s.getUserOutput, nil)

	out, err := s.usecase.LoginUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidLogin)
}

func (s *LoginUserTestSuite) TestFailedUpdate() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(s.getUserOutput, nil)
	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(repository.UpdateUserOutput{}, s.mockErr)

	out, err := s.usecase.LoginUser(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, s.mockErr)
}

func (s *LoginUserTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.repo.EXPECT().GetUser(s.ctx, s.getUserInput).Return(s.getUserOutput, nil)
	s.repo.EXPECT().UpdateUser(s.ctx, s.updateUserInput).Return(s.updateUserOutput, nil)

	out, err := s.usecase.LoginUser(s.ctx, s.input)

	a.Empty(err)
	parsedToken, err := jwt.Parse(out.JwtToken, func(token *jwt.Token) (interface{}, error) {
		return &s.jwtSecret.PublicKey, nil
	})
	a.Empty(err)
	a.True(parsedToken.Valid)
	subj, err := parsedToken.Claims.GetSubject()
	a.Empty(err)
	a.Equal("123", subj)
	exp, err := parsedToken.Claims.GetExpirationTime()
	a.Empty(err)
	a.True(time.Now().Before(exp.Time))
}
