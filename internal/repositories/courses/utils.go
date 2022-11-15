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
	ErrRecordNotFound         = errors.New("record not found")
	ErrorCourseWasDeleted     = errors.New("course was deleted")
	ErrorUserCourseWasDeleted = errors.New("this course was deleted from user")
	ErrorIncorrectArgument    = errors.New("incorrect argument")
	ErrCourseNotFound         = errors.New("course not found")
)
