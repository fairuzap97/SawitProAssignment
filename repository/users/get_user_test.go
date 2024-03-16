package users

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type GetUserTestSuite struct {
	suite.Suite

	dbMock sqlmock.Sqlmock
	repo   repository.UserRepository

	input  repository.GetUserInput
	output repository.GetUserOutput
	ctx    context.Context
}

func TestGetUserTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserTestSuite))
}

func (s *GetUserTestSuite) SetupTest() {
	repo := &userRepository{}
	repo.db, s.dbMock, _ = sqlmock.New()
	s.repo = repo

	s.input = repository.GetUserInput{
		PhoneNo: "+6281315184400",
	}
	s.output = repository.GetUserOutput{
		ID:                   12,
		PhoneNo:              "+6281215182299",
		FullName:             "John Smith",
		PasswordHash:         []byte("random-salted-password-hash"),
		SuccessfulLoginCount: 2,
	}
	s.ctx = context.Background()
}

func (s *GetUserTestSuite) TearDownTest() {
	if err := s.dbMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *GetUserTestSuite) TestGetByPhoneNoDatabaseError() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByPhoneNoQuery)).WithArgs(s.input.PhoneNo).
		WillReturnError(&pq.Error{Message: "some error message here"})

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorContains(err, "pq: some error message here")
}

func (s *GetUserTestSuite) TestGetByPhoneNoNotFound() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByPhoneNoQuery)).WithArgs(s.input.PhoneNo).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_no", "full_name", "password_hash", "successful_login_count"}))

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorIs(err, repository.ErrorRecordNotFound)
}

func (s *GetUserTestSuite) TestSuccessGetByPhoneNo() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByPhoneNoQuery)).WithArgs(s.input.PhoneNo).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "phone_no", "full_name", "password_hash", "successful_login_count"}).
				AddRow(s.output.ID, s.output.PhoneNo, s.output.FullName, s.output.PasswordHash, s.output.SuccessfulLoginCount),
		)

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(err)
	a.Equal(s.output, res)
}

func (s *GetUserTestSuite) TestGetByIDDatabaseError() {
	a := assert.New(s.T())

	s.input.PhoneNo = ""
	s.input.ID = 123
	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByIDQuery)).WithArgs(s.input.ID).
		WillReturnError(&pq.Error{Message: "some error message here"})

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorContains(err, "pq: some error message here")
}

func (s *GetUserTestSuite) TestGetByIDNotFound() {
	a := assert.New(s.T())

	s.input.PhoneNo = ""
	s.input.ID = 123
	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByIDQuery)).WithArgs(s.input.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_no", "full_name", "password_hash", "successful_login_count"}))

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorIs(err, repository.ErrorRecordNotFound)
}

func (s *GetUserTestSuite) TestSuccessGetByID() {
	a := assert.New(s.T())

	s.input.PhoneNo = ""
	s.input.ID = 123
	s.dbMock.ExpectQuery(regexp.QuoteMeta(getUserByIDQuery)).WithArgs(s.input.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "phone_no", "full_name", "password_hash", "successful_login_count"}).
				AddRow(s.output.ID, s.output.PhoneNo, s.output.FullName, s.output.PasswordHash, s.output.SuccessfulLoginCount),
		)

	res, err := s.repo.GetUser(s.ctx, s.input)
	a.Empty(err)
	a.Equal(s.output, res)
}
