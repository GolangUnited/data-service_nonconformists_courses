package internal

import (
	"errors"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrCourseWasDeleted     = errors.New("course was deleted")
	ErrUserCourseWasDeleted = errors.New("this course was deleted from user")
	ErrIncorrectArgument    = errors.New("incorrect argument")
	ErrCourseNotFound       = errors.New("course not found")
	ErrUserCourseNotFound   = errors.New("user didn't join this course")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrInvalidStatus        = errors.New("this status can't be set")
	ErrInvalidUUIDFormat    = errors.New("invalid UUID format")
	ErrOutOfRange           = errors.New("value out of range")
)
