package handler

import (
	"context"
	"fmt"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type UserHandlerTestSuite struct {
	suite.Suite

	gomock  *gomock.Controller
	usecase *usecase.MockUserUsecases

	handler *Server

	ctx     context.Context
	mockErr error
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) SetupTest() {
	s.gomock = gomock.NewController(s.T())
	s.usecase = usecase.NewMockUserUsecases(s.gomock)

	s.handler = NewServer(NewServerOptions{UserUsecase: s.usecase})

	s.ctx = context.Background()
	s.mockErr = fmt.Errorf("simulated error")
}

func (s *UserHandlerTestSuite) TearDownTest() {
	s.gomock.Finish()
}

func (s *UserHandlerTestSuite) TestLoginUserInvalidJson() {
	a := assert.New(s.T())

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/session", strings.NewReader("NOT-A-JSON"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.LoginUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusBadRequest, rec.Code)
	a.Equal(`{"error":"invalid JSON Body"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestLoginUserInvalidCredentials() {
	a := assert.New(s.T())

	s.usecase.EXPECT().LoginUser(s.ctx, usecase.LoginUserInput{
		PhoneNo:  "+62812141733",
		Password: "SomeP@ssw0rdHere",
	}).Return(usecase.LoginUserOutput{}, usecase.UserInvalidLogin)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/session", strings.NewReader(`{"phone_no":"+62812141733","password":"SomeP@ssw0rdHere"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.LoginUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusBadRequest, rec.Code)
	a.Equal(`{"error":"invalid phone number or password"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestLoginUserSuccess() {
	a := assert.New(s.T())

	s.usecase.EXPECT().LoginUser(s.ctx, usecase.LoginUserInput{
		PhoneNo:  "+62812141733",
		Password: "SomeP@ssw0rdHere",
	}).Return(usecase.LoginUserOutput{JwtToken: "jwt-token"}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user/session", strings.NewReader(`{"phone_no":"+62812141733","password":"SomeP@ssw0rdHere"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.LoginUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusOK, rec.Code)
	a.Equal(`{"jwt_token":"jwt-token"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestRegisterUserInternalError() {
	a := assert.New(s.T())

	s.usecase.EXPECT().RegisterUser(s.ctx, usecase.RegisterUserInput{
		PhoneNo:  "+62812141733",
		FullName: "John Smith",
		Password: "SomeP@ssw0rdHere",
	}).Return(usecase.RegisterUserOutput{}, s.mockErr)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(`{"phone_no":"+62812141733","full_name":"John Smith","password":"SomeP@ssw0rdHere"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.RegisterUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusInternalServerError, rec.Code)
	a.Equal(`{"error":"simulated error"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestRegisterUserBadRequest() {
	a := assert.New(s.T())

	s.usecase.EXPECT().RegisterUser(s.ctx, usecase.RegisterUserInput{
		PhoneNo:  "not-a-phone-number",
		FullName: "A",
		Password: "simplepassword",
	}).Return(usecase.RegisterUserOutput{}, usecase.NewValidationError(map[string][]error{
		"phone_no":  {fmt.Errorf("phone_no too long"), fmt.Errorf("phone_no must only consist of digits")},
		"full_name": {fmt.Errorf("full_name short")},
	}))

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(`{"phone_no":"not-a-phone-number","full_name":"A","password":"simplepassword"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.RegisterUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusBadRequest, rec.Code)
	a.Contains(rec.Body.String(), `{"error":"phone_no too long","field":"phone_no"},{"error":"phone_no must only consist of digits","field":"phone_no"}`)
	a.Contains(rec.Body.String(), `{"error":"full_name short","field":"full_name"}`)
}

func (s *UserHandlerTestSuite) TestRegisterUserSuccess() {
	a := assert.New(s.T())

	s.usecase.EXPECT().RegisterUser(s.ctx, usecase.RegisterUserInput{
		PhoneNo:  "+62812141733",
		FullName: "John Smith",
		Password: "SomeP@ssw0rdHere",
	}).Return(usecase.RegisterUserOutput{UserID: 123}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(`{"phone_no":"+62812141733","full_name":"John Smith","password":"SomeP@ssw0rdHere"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.RegisterUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusOK, rec.Code)
	a.Equal(`{"user_id":123}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestGetUserWithoutAuth() {
	a := assert.New(s.T())

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.GetUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusForbidden, rec.Code)
	a.Equal(`{"error":"invalid / expired token, please login again"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestGetUserInvalidToken() {
	a := assert.New(s.T())

	s.usecase.EXPECT().ValidateUserToken(s.ctx, usecase.ValidateUserTokenInput{
		JwtToken: "jwt-token",
	}).Return(usecase.ValidateUserTokenOutput{}, usecase.UserInvalidToken)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(""))
	req.Header.Set("Authorization", "Bearer jwt-token")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.GetUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusForbidden, rec.Code)
	a.Equal(`{"error":"invalid / expired token, please login again"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestGetUserSuccess() {
	a := assert.New(s.T())

	s.usecase.EXPECT().ValidateUserToken(s.ctx, usecase.ValidateUserTokenInput{
		JwtToken: "jwt-token",
	}).Return(usecase.ValidateUserTokenOutput{UserID: 123}, nil)
	s.usecase.EXPECT().GetUserProfile(s.ctx, usecase.GetUserProfileInput{
		UserID: 123,
	}).Return(usecase.GetUserProfileOutput{
		UserID:               123,
		PhoneNo:              "+62812141733",
		FullName:             "John Smith",
		SuccessfulLoginCount: 42,
	}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(""))
	req.Header.Set("Authorization", "Bearer jwt-token")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.GetUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusOK, rec.Code)
	a.Equal(`{"full_name":"John Smith","phone_no":"+62812141733","successful_login_count":42,"user_id":123}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestUpdateUserWithBrokenAuth() {
	a := assert.New(s.T())

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/user", strings.NewReader(`{"phone_no":"+628121417","full_name":"Alex Smith"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Basic basic-auth")
	rec := httptest.NewRecorder()
	err := s.handler.UpdateUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusForbidden, rec.Code)
	a.Equal(`{"error":"invalid / expired token, please login again"}`, strings.TrimSpace(rec.Body.String()))
}

func (s *UserHandlerTestSuite) TestUpdateUserSuccess() {
	a := assert.New(s.T())

	phoneNo := "+62812141722"
	fullName := "Alex Smith"
	s.usecase.EXPECT().ValidateUserToken(s.ctx, usecase.ValidateUserTokenInput{
		JwtToken: "jwt-token",
	}).Return(usecase.ValidateUserTokenOutput{UserID: 123}, nil)
	s.usecase.EXPECT().UpdateUserProfile(s.ctx, usecase.UpdateUserProfileInput{
		UserID:   123,
		PhoneNo:  &phoneNo,
		FullName: &fullName,
	}).Return(usecase.UpdateUserProfileOutput{}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/user", strings.NewReader(`{"phone_no":"+62812141722","full_name":"Alex Smith"}`))
	req.Header.Set("Authorization", "Bearer jwt-token")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := s.handler.UpdateUser(e.NewContext(req, rec))

	a.Empty(err)
	a.Equal(http.StatusOK, rec.Code)
	a.Equal(`{"message":"profile updated"}`, strings.TrimSpace(rec.Body.String()))
}
