// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

//go:generate mockgen -source=interfaces.go -destination=../mocks/repository.go -package mocks

// UserRepository is a simple CRUD Interface to interact with the user data model
type UserRepository interface {

	// CreateUser will create a new user as specified by the CreateUserInput input, and return the created record ID
	// Will return error on Database Error or Record Conflict
	CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)

	// GetUser will return a user data using either the User ID or PhoneNo, as specified by GetUserInput input
	// Will return error on Database Error or No Record Found
	GetUser(ctx context.Context, input GetUserInput) (output GetUserOutput, err error)

	// UpdateUser will update a user data with the specified ID on UpdateUserInput input
	// Will return error on Database Error or No Record Found
	UpdateUser(ctx context.Context, input UpdateUserInput) (output UpdateUserOutput, err error)
}
