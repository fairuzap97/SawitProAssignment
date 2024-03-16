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

type UpdateUserTestSuite struct {
	suite.Suite

	dbMock sqlmock.Sqlmock
	repo   repository.UserRepository

	input repository.UpdateUserInput
	ctx   context.Context
}

func TestUpdateUserTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserTestSuite))
}

func (s *UpdateUserTestSuite) SetupTest() {
	repo := &userRepository{}
	repo.db, s.dbMock, _ = sqlmock.New()
	s.repo = repo

	s.input = repository.UpdateUserInput{
		ID:                   123,
		PhoneNo:              "+6281315184400",
		FullName:             "John Smith",
		SuccessfulLoginCount: 2,
	}
	s.ctx = context.Background()
}

func (s *UpdateUserTestSuite) TearDownTest() {
	if err := s.dbMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *UpdateUserTestSuite) TestDatabaseError() {
	a := assert.New(s.T())

	s.input.FullName = ""
	s.input.PhoneNo = ""
	s.dbMock.ExpectExec(regexp.QuoteMeta("UPDATE users SET successful_login_count=$1 WHERE id=$2")).
		WithArgs(s.input.SuccessfulLoginCount, s.input.ID).
		WillReturnError(&pq.Error{Message: "some error message here"})

	_, err := s.repo.UpdateUser(s.ctx, s.input)
	a.ErrorContains(err, "pq: some error message here")
}

func (s *UpdateUserTestSuite) TestRecordConflictError() {
	a := assert.New(s.T())

	s.input.SuccessfulLoginCount = 0
	s.dbMock.ExpectExec(regexp.QuoteMeta("UPDATE users SET full_name=$1, phone_no=$2 WHERE id=$3")).
		WithArgs(s.input.FullName, s.input.PhoneNo, s.input.ID).
		WillReturnError(&pq.Error{Message: "some error message here", Code: "23505"})

	_, err := s.repo.UpdateUser(s.ctx, s.input)
	a.ErrorIs(err, repository.ErrorRecordConflict)
}

func (s *UpdateUserTestSuite) TestRecordNotFound() {
	a := assert.New(s.T())

	s.input.SuccessfulLoginCount = 0
	s.input.FullName = ""
	s.dbMock.ExpectExec(regexp.QuoteMeta("UPDATE users SET phone_no=$1 WHERE id=$2")).
		WithArgs(s.input.PhoneNo, s.input.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.repo.UpdateUser(s.ctx, s.input)
	a.ErrorIs(err, repository.ErrorRecordNotFound)
}

func (s *UpdateUserTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.input.SuccessfulLoginCount = 0
	s.input.FullName = ""
	s.dbMock.ExpectExec(regexp.QuoteMeta("UPDATE users SET phone_no=$1 WHERE id=$2")).
		WithArgs(s.input.PhoneNo, s.input.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	_, err := s.repo.UpdateUser(s.ctx, s.input)
	a.Empty(err)
}
