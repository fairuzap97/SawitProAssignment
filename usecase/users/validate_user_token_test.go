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
	"testing"
	"time"
)

type ValidateUserTokenTestSuite struct {
	suite.Suite

	gomock *gomock.Controller
	repo   *repository.MockUserRepository

	jwtSecret *rsa.PrivateKey
	usecase   usecase.UserUsecases

	input  usecase.ValidateUserTokenInput
	output usecase.ValidateUserTokenOutput

	ctx     context.Context
	mockErr error
}

func TestValidateUserTokenTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateUserTokenTestSuite))
}

func (s *ValidateUserTokenTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.repo = repository.NewMockUserRepository(s.gomock)

	s.jwtSecret, _ = rsa.GenerateKey(rand.Reader, 1024)
	s.usecase = NewUserUsecases(NewUserUsecasesOptions{
		UserRepo:  s.repo,
		JwtSecret: s.jwtSecret,
		JwtTtl:    time.Minute * 5,
	})

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "123",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	jwtToken, _ := token.SignedString(s.jwtSecret)
	s.input = usecase.ValidateUserTokenInput{JwtToken: jwtToken}
	s.output = usecase.ValidateUserTokenOutput{UserID: 123}

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *ValidateUserTokenTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *ValidateUserTokenTestSuite) TestInvalidTokenSignature() {
	a := assert.New(s.T())

	otherSeret, _ := rsa.GenerateKey(rand.Reader, 1024)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "123",
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})
	s.input.JwtToken, _ = token.SignedString(otherSeret)
	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidToken)
}

func (s *ValidateUserTokenTestSuite) TestTokenMissingExp() {
	a := assert.New(s.T())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "123",
	})
	s.input.JwtToken, _ = token.SignedString(s.jwtSecret)
	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidToken)
}

func (s *ValidateUserTokenTestSuite) TestTokenExpired() {
	a := assert.New(s.T())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "123",
		"exp": time.Now().Add(time.Minute * -5).Unix(),
	})
	s.input.JwtToken, _ = token.SignedString(s.jwtSecret)
	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidToken)
}

func (s *ValidateUserTokenTestSuite) TestTokenMissingSub() {
	a := assert.New(s.T())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * -5).Unix(),
	})
	s.input.JwtToken, _ = token.SignedString(s.jwtSecret)
	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidToken)
}

func (s *ValidateUserTokenTestSuite) TestTokenSubjectInvalid() {
	a := assert.New(s.T())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "asdf",
		"exp": time.Now().Add(time.Minute * -5).Unix(),
	})
	s.input.JwtToken, _ = token.SignedString(s.jwtSecret)
	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(out)
	a.ErrorIs(err, usecase.UserInvalidToken)
}

func (s *ValidateUserTokenTestSuite) TestSuccess() {
	a := assert.New(s.T())

	out, err := s.usecase.ValidateUserToken(s.ctx, s.input)

	a.Empty(err)
	a.Equal(s.output, out)
}
