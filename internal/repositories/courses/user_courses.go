package courses

import (
	"fmt"
	"golang-united-courses/internal/models"

	"github.com/google/uuid"
)

func (p *PostgreSql) CheckID(course_id, user_id string) (models.UserCourse, error) {
	var userCourse models.UserCourse
	course, err := p.GetById(course_id)
	if err != nil {
		return userCourse, err
	}
	if course.IsDeleted != 0 {
		return userCourse, fmt.Errorf("%w", CourseWasDeletedError)
	}
	userCourse.CourseID = course.ID
	uid, err := uuid.Parse(user_id)
	if err != nil {
		return userCourse, err
	}
	userCourse.UserID = uid
	return userCourse, nil
}

func (p *PostgreSql) CheckUserCourse(user_id, course_id string) (models.UserCourse, error) {
	userCourse, err := p.CheckID(user_id, course_id)
	if err != nil {
		return userCourse, err
	}
	err = p.DB.First(&userCourse).Error
	if err != nil {
		return userCourse, err
	}
	if userCourse.IsDeleted != 0 {
		return userCourse, fmt.Errorf("%w", UserCourseWasDeletedError)
	}
	return userCourse, nil
}

func (p *PostgreSql) JoinCourse(user_id, course_id string) error {
	userCourse, err := p.CheckID(user_id, course_id)
	if err != nil {
		return err
	}
	err = p.DB.Create(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) DeclineCourse(user_id, course_id string) error {
	userCourse, err := p.CheckUserCourse(user_id, course_id)
	if err != nil {
		return err
	}
	userCourse.IsDeleted = 1
	err = p.DB.Updates(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) UpdateCourse(percent uint32, user_id, course_id string) error {
	userCourse, err := p.CheckUserCourse(user_id, course_id)
	if err != nil {
		return err
	}
	if percent == 0 || percent > 100 {
		return fmt.Errorf("%w", IncorrectArgumentError)
	}
	userCourse.PercentFinished = percent
	err = p.DB.Updates(&userCourse).Error
	if err != nil {
		return err
	}
}

func (p *PostgreSql) ListCourse(user_id, course_id string, limit, offset uint32, showDeleted bool) ([]models.UserCourse, error) {
	return nil, nil
}
