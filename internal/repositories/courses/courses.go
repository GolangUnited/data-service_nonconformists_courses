package courses

import (
	"golang-united-courses/internal"
	"golang-united-courses/internal/models"
	"golang-united-courses/internal/repositories/db"
)

type CoursePGSQL struct {
	*db.PostgreSql
}

func (p *CoursePGSQL) Create(title, desciption string) (string, error) {
	var course models.Course
	course.Title = title
	course.Description = desciption
	if err := p.DB.Create(&course).Error; err != nil {
		return "", err
	}
	return course.ID.String(), nil
}

func (p *CoursePGSQL) Delete(id string) error {
	var Course models.Course
	course, err := p.GetById(id)
	if err != nil {
		return err
	}
	if course.IsDeleted != 0 {
		return internal.ErrCourseWasDeleted
	}
	if err = p.DB.Model(&Course).Where("id = ?", id).Update("is_deleted", 1).Error; err != nil {
		return err
	}
	return nil
}

func (p *CoursePGSQL) Update(id, title, desciption string) error {
	course, err := p.GetById(id)
	if err != nil {
		return err
	}
	if title != "" {
		course.Title = title
	}
	if desciption != "" {
		course.Description = desciption
	}
	if err = p.DB.Updates(&course).Error; err != nil {
		return err
	}
	return nil
}

func (p *CoursePGSQL) GetById(id string) (models.Course, error) {
	var course models.Course
	if err := p.DB.Model(&course).Where("id = ?", id).First(&course).Error; err != nil {
		switch err.Error() {
		case internal.ErrRecordNotFound.Error():
			return models.Course{}, internal.ErrCourseNotFound
		default:
			return models.Course{}, err
		}
	}
	return course, nil
}

func (p *CoursePGSQL) List(showDeleted bool, limit, offset int32) ([]models.Course, error) {
	var courses []models.Course
	var Course models.Course
	q := p.DB.Model(&Course)
	if limit > 0 {
		q.Limit(int(limit))
	}
	if offset > 0 {
		q.Offset(int(offset))
	}
	if !showDeleted {
		q.Where("is_deleted = ?", 0)
	}
	if err := q.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}
