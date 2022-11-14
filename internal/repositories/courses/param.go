package courses

import (
	"golang-united-courses/internal/models"

	"gorm.io/gorm"
)

type PostgreSql struct {
	DB *gorm.DB
}

type Error string

var (
	UserCourse models.UserCourse
	Course     models.Course
)

const (
	CourseWasDeletedError     = Error("course was deleted")
	UserCourseWasDeletedError = Error("this course was deleted from user")
	IncorrectArgumentError    = Error("incorrect argument")
)
