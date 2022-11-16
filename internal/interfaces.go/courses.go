package interfaces

import (
	"golang-united-courses/internal/models"
)

type CourseManager interface {
	Init(url string) error
	Close() error
	CourseStoreManager
	UserCourseStoreManager
}

type CourseStoreManager interface {
	Create(title, desciption string) (string, error)
	Delete(id string) error
	Update(id, title, desciption string) error
	GetById(id string) (models.Course, error)
	List(showDeleted bool, limit, offset int32) ([]models.Course, error)
}

type UserCourseStoreManager interface {
	Join(user_id, course_id string) error
	SetProgress(percent int32, user_id, course_id, status string) error
	ListUserCourse(user_id, course_id string, limit, offset int32, showDeleted bool) ([]models.UserCourse, error)
}
