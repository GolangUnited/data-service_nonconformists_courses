package courses

import (
	"golang-united-courses/internal/models"
	"time"

	"github.com/google/uuid"
)

func (p *PostgreSql) checkCourseById(course_id string) (uuid.UUID, error) {
	course, err := p.GetById(course_id)
	if err != nil {
		return uuid.Nil, err
	}
	if course.IsDeleted != 0 {
		return uuid.Nil, ErrorCourseWasDeleted
	}
	return course.ID, nil
}

func (p *PostgreSql) checkUserById(user_id string) (uuid.UUID, error) {
	id, err := uuid.Parse(user_id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (p *PostgreSql) checkId(course_id, user_id string) (models.UserCourse, error) {
	var userCourse models.UserCourse
	cid, err := p.checkCourseById(course_id)
	if err != nil {
		return userCourse, err
	}
	userCourse.CourseID = cid
	uid, err := p.checkUserById(user_id)
	if err != nil {
		return userCourse, err
	}
	userCourse.UserID = uid
	return userCourse, nil
}

func (p *PostgreSql) checkUserCourse(course_id, user_id string) (models.UserCourse, error) {
	userCourse, err := p.checkId(course_id, user_id)
	if err != nil {
		return userCourse, err
	}
	return userCourse, nil
}

func (p *PostgreSql) Join(user_id, course_id string) error {
	userCourse, err := p.checkId(course_id, user_id)
	if err != nil {
		return err
	}
	err = p.DB.Create(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) GetUserCourse(user_id, course_id string) error {
	userCourse, err := p.checkId(course_id, user_id)
	if err != nil {
		return err
	}
	err = p.DB.First(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) SetProgress(user_id, course_id string, percent uint32) error {
	userCourse, err := p.checkUserCourse(course_id, user_id)
	if err != nil {
		return err
	}
	switch {
	case percent <= 100:
		userCourse.PercentFinished = percent
	default:
		return ErrorIncorrectArgument
	}
	err = p.DB.Updates(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) SetStatus(user_id, course_id, status string) error {
	userCourse, err := p.checkUserCourse(course_id, user_id)
	if err != nil {
		return err
	}
	switch status {
	case "start":
		userCourse.StartDate = time.Now()
		userCourse.Status = status
	case "finish":
		userCourse.FinishDate = time.Now()
		userCourse.Status = status
	default:
		return ErrorIncorrectArgument
	}
	err = p.DB.Updates(&userCourse).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) ListUserCourse(user_id, course_id string, limit, offset int32, showDeleted bool) ([]models.UserCourse, error) {
	var userCourses []models.UserCourse
	q := p.DB.Model(&Course)
	if limit > 0 {
		q.Limit(int(limit))
	}
	if offset > 0 {
		q.Offset(int(offset))
	}
	if !showDeleted {
		q.Where("status != ?", "deleted")
	}
	if user_id != "" {
		q.Where("user_id = ?", user_id)
	}
	if course_id != "" {
		q.Where("course_id = ?", course_id)
	}
	err := q.Find(&userCourses).Error
	if err != nil {
		return nil, err
	}
	return userCourses, nil
}
