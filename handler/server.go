package handler

import (
	"errors"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Server struct {
	userUsecase usecase.UserUsecases
}

type NewServerOptions struct {
	UserUsecase usecase.UserUsecases
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		userUsecase: opts.UserUsecase,
	}
}

func (s *Server) getCurrentUser(ctx echo.Context) (userID uint64, err error) {
	authHeaders := ctx.Request().Header["Authorization"]
	if len(authHeaders) == 0 {
		return 0, usecase.UserInvalidToken
	}
	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, usecase.UserInvalidToken
	}
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := s.userUsecase.ValidateUserToken(ctx.Request().Context(), usecase.ValidateUserTokenInput{JwtToken: jwtToken})
	if err != nil {
		return 0, err
	}
	return token.UserID, nil
}

func renderError(ctx echo.Context, err error) error {
	var validationErrors usecase.ValidationErrors
	if errors.As(err, &validationErrors) {
		return renderValidationErrors(ctx, validationErrors)
	}

	var code int
	var ok bool
	resp := generated.ErrorResponse{Error: err.Error()}
	if code, ok = errorsToCodeMap[err]; !ok {
		code = 500
	}
	return ctx.JSON(code, resp)
}

func renderValidationErrors(ctx echo.Context, fieldErrors usecase.ValidationErrors) error {
	resp := generated.FieldErrorsResponse{Errors: nil}
	for field, errs := range fieldErrors.GetErrors() {
		for _, err := range errs {
			resp.Errors = append(resp.Errors, struct {
				Error string `json:"error"`
				Field string `json:"field"`
			}{Error: err.Error(), Field: field})
		}
	}
	return ctx.JSON(http.StatusBadRequest, resp)
}
