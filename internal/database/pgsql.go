package database

import (
	"golang-united-courses/internal/repositories/courses"
	"golang-united-courses/internal/repositories/db"
	"golang-united-courses/internal/repositories/user_courses"
)

type PgSql struct {
	*db.PostgreSql
	courses.CoursePGSQL
	user_courses.UserCoursePGSQL
}

func NewPgSql() *PgSql {
	PgSqlDb := new(db.PostgreSql)
	return &PgSql{
		PgSqlDb,
		courses.CoursePGSQL{PostgreSql: PgSqlDb},
		user_courses.UserCoursePGSQL{PostgreSql: PgSqlDb},
	}
}
