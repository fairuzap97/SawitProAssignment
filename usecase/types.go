// This file contains types that are used in the repository layer.
package usecase

type RegisterUserInput struct {
	PhoneNo  string
	FullName string
	Password string
}

type RegisterUserOutput struct {
	UserID uint64
}

type LoginUserInput struct {
	PhoneNo  string
	Password string
}

type LoginUserOutput struct {
	JwtToken string
}

type ValidateUserTokenInput struct {
	JwtToken string
}

type ValidateUserTokenOutput struct {
	UserID uint64
}

type GetUserProfileInput struct {
	UserID uint64
}

type GetUserProfileOutput struct {
	UserID               uint64
	PhoneNo              string
	FullName             string
	SuccessfulLoginCount uint64
}

type UpdateUserProfileInput struct {
	UserID   uint64
	PhoneNo  *string
	FullName *string
}

type UpdateUserProfileOutput struct{}
