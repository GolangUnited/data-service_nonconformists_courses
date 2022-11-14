package interfaces

import "golang-united-courses/internal/models"

type CourseManager interface {
	Init(url string) error
	Close() error
	ICourseStore
	IUserCourseStore
}

type ICourseStore interface {
	Create(title, desciption string) (string, error)
	Delete(id string) error
	Update(id, title, desciption string) error
	GetById(id string) (models.Course, error)
	List(showDeleted bool, limit, offset uint32) ([]models.Course, error)
}

type IUserCourseStore interface {
	JoinCourse(user_id, course_id string) error
	DeclineCourse(user_id, course_id string) error
	UpdateCourse(percent uint32, user_id, course_id string) error
	ListCourse(user_id, course_id string, limit, offset uint32, showDeleted bool) ([]models.UserCourse, error)
}
