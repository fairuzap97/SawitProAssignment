// This file contains the interfaces for the usecase layer.
// The usecase layer is responsible for handling business / application logic.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package usecase

import "context"

//go:generate mockgen -source=interfaces.go -destination=../mocks/usecase.go -package mocks

// UserUsecases is an interface of business / application functions that's related to user actions
type UserUsecases interface {
	// RegisterUser will validate and register a new user based on the information on RegisterUserInput input
	RegisterUser(ctx context.Context, input RegisterUserInput) (output RegisterUserOutput, err error)

	// LoginUser will search for the target user, and validate the password with the records in the database
	// Will return a JWT Token if successful
	LoginUser(ctx context.Context, input LoginUserInput) (output LoginUserOutput, err error)

	// ValidateUserToken will validate a users JWT Token and return the UserID contained in the token
	ValidateUserToken(ctx context.Context, input ValidateUserTokenInput) (output ValidateUserTokenOutput, err error)

	GetUserProfile(ctx context.Context, input GetUserProfileInput) (output GetUserProfileOutput, err error)

	UpdateUserProfile(ctx context.Context, input UpdateUserProfileInput) (output UpdateUserProfileOutput, err error)
}
