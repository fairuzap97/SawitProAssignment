package handler

import (
	"encoding/json"
	"github.com/SawitProRecruitment/UserService/usecase"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

// Login API. Create new session for the target user.
// (POST /user/session)
func (s *Server) LoginUser(ctx echo.Context) error {
	var payload generated.LoginUserRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		err = JsonBodyInvalid
		return renderError(ctx, err)
	}

	result, err := s.userUsecase.LoginUser(ctx.Request().Context(), usecase.LoginUserInput{
		PhoneNo:  payload.PhoneNo,
		Password: payload.Password,
	})
	if err != nil {
		return renderError(ctx, err)
	}

	resp := generated.LoginUserResponse{JwtToken: result.JwtToken}
	return ctx.JSON(http.StatusOK, resp)
}

// Register new user with the provided phone number, full name, and password.
// (POST /user)
func (s *Server) RegisterUser(ctx echo.Context) error {
	var payload generated.RegisterUserRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		err = JsonBodyInvalid
		return renderError(ctx, err)
	}

	result, err := s.userUsecase.RegisterUser(ctx.Request().Context(), usecase.RegisterUserInput{
		PhoneNo:  payload.PhoneNo,
		FullName: payload.FullName,
		Password: payload.Password,
	})
	if err != nil {
		return renderError(ctx, err)
	}

	resp := generated.RegisterUserResponse{UserId: int(result.UserID)}
	return ctx.JSON(http.StatusOK, resp)
}

// Get logged-in user profile
// (GET /user)
func (s *Server) GetUser(ctx echo.Context) error {
	userID, err := s.getCurrentUser(ctx)
	if err != nil {
		return renderError(ctx, err)
	}

	result, err := s.userUsecase.GetUserProfile(ctx.Request().Context(), usecase.GetUserProfileInput{UserID: userID})
	if err != nil {
		return renderError(ctx, err)
	}

	resp := generated.GetUserResponse{
		UserId:               int(result.UserID),
		FullName:             result.FullName,
		PhoneNo:              result.PhoneNo,
		SuccessfulLoginCount: int(result.SuccessfulLoginCount),
	}
	return ctx.JSON(http.StatusOK, resp)
}

// Update logged-in user profile
// (PATCH /user)
func (s *Server) UpdateUser(ctx echo.Context) error {
	userID, err := s.getCurrentUser(ctx)
	if err != nil {
		return renderError(ctx, err)
	}

	var payload generated.UpdateUserRequest
	if err = json.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		err = JsonBodyInvalid
		return renderError(ctx, err)
	}
	_, err = s.userUsecase.UpdateUserProfile(ctx.Request().Context(), usecase.UpdateUserProfileInput{
		UserID:   userID,
		PhoneNo:  payload.PhoneNo,
		FullName: payload.FullName,
	})
	if err != nil {
		return renderError(ctx, err)
	}

	resp := generated.UpdateUserResponse{Message: "profile updated"}
	return ctx.JSON(http.StatusOK, resp)
}
