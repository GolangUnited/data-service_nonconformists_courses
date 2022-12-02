package utils

import (
	"errors"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrCourseWasDeleted     = errors.New("course was deleted")
	ErrUserCourseWasDeleted = errors.New("this course was deleted from user")
	ErrIncorrectArgument    = errors.New("incorrect argument")
	ErrCourseNotFound       = errors.New("course not found")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrInvalidStatus        = errors.New("this status can't be set")
)
