package courses

import (
	"errors"
	"golang-united-courses/internal/models"

	"gorm.io/gorm"
)

type PostgreSql struct {
	DB *gorm.DB
}

var (
	UserCourse models.UserCourse
	Course     models.Course
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrCourseWasDeleted     = errors.New("course was deleted")
	ErrUserCourseWasDeleted = errors.New("this course was deleted from user")
	ErrIncorrectArgument    = errors.New("incorrect argument")
	ErrCourseNotFound       = errors.New("course not found")
	ErrInvalidFormat        = errors.New("invalid format")
)
