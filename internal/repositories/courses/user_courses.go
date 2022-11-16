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

func (p *PostgreSql) SetProgress(percent uint32, user_id, course_id, status string) error {
	userCourse, err := p.checkUserCourse(course_id, user_id)
	if err != nil {
		return err
	}
	// FUTURE: to uint, check if % only > 100
	if percent > 100 {
		return ErrorIncorrectArgument
	}
	userCourse.PercentFinished = percent
	// FUTURE: statuses to INT -> 0 - not started is default; 1 - started; 2 - finished; 3 - deleted/paused; other - ErrIncorrectStatus
	if status != "" {
		userCourse.Status = status
	}
	if status == "start" && userCourse.StartDate.IsZero() {
		userCourse.StartDate = time.Now()
	}
	if status == "finish" && userCourse.FinishDate.IsZero() {
		userCourse.FinishDate = time.Now()
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
