package user_courses

import (
	"golang-united-courses/internal"
	"golang-united-courses/internal/models"
	"golang-united-courses/internal/repositories/db"
)

type UserCoursePGSQL struct {
	*db.PostgreSql
}

func (p *UserCoursePGSQL) Join(uc models.UserCourse) error {
	uc.PercentFinished = 0
	uc.Status = models.Joined
	if err := p.DB.Create(&uc).Error; err != nil {
		return err
	}
	return nil
}

func (p *UserCoursePGSQL) GetUserCourse(user_id, course_id string) (models.UserCourse, error) {
	var uc models.UserCourse
	if err := p.DB.Model(&models.UserCourse{}).Where("course_id = ? AND user_id = ?", course_id, user_id).First(&uc).Error; err != nil {
		switch err.Error() {
		case internal.ErrRecordNotFound.Error():
			return models.UserCourse{}, internal.ErrUserCourseNotFound
		default:
			return models.UserCourse{}, err
		}
	}
	return uc, nil
}

func (p *UserCoursePGSQL) UpdateUserCourse(uc models.UserCourse) error {
	if err := p.DB.Updates(&uc).Error; err != nil {
		return err
	}
	return nil
}

func (p *UserCoursePGSQL) ListUserCourse(user_id, course_id string, limit, offset int32, showDeleted bool) ([]models.UserCourse, error) {
	var userCourses []models.UserCourse
	q := p.DB.Model(&models.UserCourse{})
	if limit > 0 {
		q.Limit(int(limit))
	}
	if offset > 0 {
		q.Offset(int(offset))
	}
	if !showDeleted {
		q.Where("status != ?", models.Declined)
	}
	if user_id != "" {
		q.Where("user_id = ?", user_id)
	}
	if course_id != "" {
		q.Where("course_id = ?", course_id)
	}
	if err := q.Find(&userCourses).Error; err != nil {
		return nil, err
	}
	return userCourses, nil
}
