// This file contains the usecase implementation layer.
package users

import (
	"crypto/rsa"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/usecase"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	allDigitRegex = regexp.MustCompile(`^(\+)*[\d]+$`) // match all number string with optional `+` prefix
	digitRegex    = regexp.MustCompile(`^.*\d.*$`)     // match string containing at least one number
	capitalRegex  = regexp.MustCompile(`^.*[A-Z].*$`)  // match string containing at least one Capital letter
	specialRegex  = regexp.MustCompile(`^.*[\W_].*$`)  // match string containing at least one special character (non alphanumeric)
)

// userUsecases is an implementation of usecase.UserUsecases
type userUsecases struct {
	userRepo  repository.UserRepository
	jwtSecret *rsa.PrivateKey
	jwtTtl    time.Duration
}

type NewUserUsecasesOptions struct {
	UserRepo  repository.UserRepository
	JwtSecret *rsa.PrivateKey
	JwtTtl    time.Duration
}

func NewUserUsecases(opts NewUserUsecasesOptions) usecase.UserUsecases {
	return &userUsecases{
		userRepo:  opts.UserRepo,
		jwtSecret: opts.JwtSecret,
		jwtTtl:    opts.JwtTtl,
	}
}

func validateUserPhoneNo(phoneNo string) []error {
	var res []error
	if len(phoneNo) < 10 || len(phoneNo) > 13 {
		res = append(res, fmt.Errorf(`phone_no must be between 10 and 13 characters long`))
	}
	if !strings.HasPrefix(phoneNo, "+62") {
		res = append(res, fmt.Errorf(`phone_no must start with indonesia country code ("+62")`))
	}
	if !allDigitRegex.MatchString(phoneNo) {
		res = append(res, fmt.Errorf(`besides the country code, phone_no must only contain numbers`))
	}
	return res
}

func validateUserFullName(fullName string) []error {
	var res []error
	if len(fullName) < 3 || len(fullName) > 60 {
		res = append(res, fmt.Errorf(`full_name must be between 3 and 60 characters long`))
	}
	return res
}

func validateUserPassword(password string) []error {
	var res []error
	if len(password) < 6 || len(password) > 64 {
		res = append(res, fmt.Errorf(`password must be between 6 and 64 characters long`))
	}
	if !digitRegex.MatchString(password) {
		res = append(res, fmt.Errorf(`password must contains at least one number [0-9]`))
	}
	if !capitalRegex.MatchString(password) {
		res = append(res, fmt.Errorf(`password must contains at least one capital letter [A-Z]`))
	}
	if !specialRegex.MatchString(password) {
		res = append(res, fmt.Errorf(`password must contains at least one non alphanumeric character`))
	}
	return res
}
