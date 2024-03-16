package repository

import "errors"

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorRecordConflict = errors.New("record conflict error")
)
