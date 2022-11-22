package courses

import (
	"golang-united-courses/internal/models"
)

func (p *PostgreSql) Create(title, desciption string) (string, error) {
	var course models.Course
	course.Title = title
	course.Description = desciption
	err := p.DB.Create(&course).Error
	if err != nil {
		return "", err
	}
	return course.ID.String(), nil
}

func (p *PostgreSql) Delete(id string) error {
	course, err := p.GetById(id)
	if err != nil {
		return err
	}
	if course.IsDeleted != 0 {
		return ErrCourseWasDeleted
	}
	err = p.DB.Model(&Course).Where("id = ?", id).Update("is_deleted", 1).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) Update(id, title, desciption string) error {
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
	err = p.DB.Updates(&course).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) GetById(id string) (models.Course, error) {
	var course models.Course
	err := p.DB.Model(&Course).Where("id = ?", id).First(&course).Error
	if err != nil {
		switch err.Error() {
		case ErrRecordNotFound.Error():
			return course, ErrCourseNotFound
		default:
			return course, err
		}
	}
	return course, nil
}

func (p *PostgreSql) List(showDeleted bool, limit, offset int32) ([]models.Course, error) {
	var courses []models.Course
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
	err := q.Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
