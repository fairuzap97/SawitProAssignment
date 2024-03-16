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

type CreateUserTestSuite struct {
	suite.Suite

	dbMock sqlmock.Sqlmock
	repo   repository.UserRepository

	input repository.CreateUserInput
	ctx   context.Context
}

func TestCreateUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}

func (s *CreateUserTestSuite) SetupTest() {
	repo := &userRepository{}
	repo.db, s.dbMock, _ = sqlmock.New()
	s.repo = repo

	s.input = repository.CreateUserInput{
		PhoneNo:      "+6281315184400",
		FullName:     "John Smith",
		PasswordHash: []byte("random-hash-string-plus-salt"),
	}
	s.ctx = context.Background()
}

func (s *CreateUserTestSuite) TearDownTest() {
	if err := s.dbMock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (s *CreateUserTestSuite) TestDatabaseError() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(createUserQuery)).WithArgs(s.input.PhoneNo, s.input.FullName, s.input.PasswordHash).
		WillReturnError(&pq.Error{Message: "some error message here"})

	res, err := s.repo.CreateUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorContains(err, "pq: some error message here")
}

func (s *CreateUserTestSuite) TestRecordConflictError() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(createUserQuery)).WithArgs(s.input.PhoneNo, s.input.FullName, s.input.PasswordHash).
		WillReturnError(&pq.Error{Message: "some error message here", Code: "23505"})

	res, err := s.repo.CreateUser(s.ctx, s.input)
	a.Empty(res)
	a.ErrorIs(err, repository.ErrorRecordConflict)
}

func (s *CreateUserTestSuite) TestSuccess() {
	a := assert.New(s.T())

	s.dbMock.ExpectQuery(regexp.QuoteMeta(createUserQuery)).WithArgs(s.input.PhoneNo, s.input.FullName, s.input.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	res, err := s.repo.CreateUser(s.ctx, s.input)
	a.Empty(err)
	a.Equal(uint64(123), res.ID)
}
