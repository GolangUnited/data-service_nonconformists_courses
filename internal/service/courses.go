package service

import (
	"context"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/db"
	"golang-united-courses/internal/models"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CourseServer struct {
	Course *db.Database
	api.UnimplementedCoursesServer
}

func (s *CourseServer) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	var course models.Course
	course.Title = request.Title
	course.Description = request.Description
	err := Create(&course, s.Course)
	if err != nil {
		return nil, err
	}
	return &api.CreateResponse{
		Id: course.ID,
	}, nil
}

func (s *CourseServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.ID = request.Id
	err := Find(&course, s.Course)
	if err != nil {
		return nil, err
	}
	result := &api.GetResponse{
		Id:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		CreatedAt:   timestamppb.New(course.CreatedAt),
		IsDeleted:   int32(course.IsDeleted),
	}
	return result, nil
}

func (s *CourseServer) Update(ctx context.Context, request *api.UpdateRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = request.Id
	err := Find(&course, s.Course)
	if err != nil {
		return nil, err
	}
	err = IsDel(course.IsDeleted)
	if err != nil {
		return nil, err
	}
	course.Title = request.Title
	course.Description = request.Description
	err = s.Course.DB.Updates(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = request.Id
	err := Find(&course, s.Course)
	if err != nil {
		return nil, err
	}
	err = IsDel(course.IsDeleted)
	if err != nil {
		return nil, err
	}
	course.IsDeleted = 1
	err = s.Course.DB.Updates(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) List(ctx context.Context, request *api.ListRequest) (*api.ListResponse, error) {
	var courses []models.Course
	var course models.Course
	q := s.Course.DB.Model(&course)
	if request.Limit > 0 {
		q.Limit(int(request.Limit))
	}
	if request.Offset > 0 {
		q.Offset(int(request.Offset))
	}
	if request.ShowDeleted {
		q.Where("is_deleted = ?", 1)
	}
	err := q.Find(&courses).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	result := &api.ListResponse{}
	result.Courses = make([]*api.GetResponse, 0, len(courses))
	for _, c := range courses {
		cResponse := &api.GetResponse{
			Id:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			CreatedAt:   timestamppb.New(c.CreatedAt),
			IsDeleted:   int32(c.IsDeleted),
		}
		result.Courses = append(result.Courses, cResponse)
	}
	return result, nil
}

func (s *CourseServer) JoinCourse(ctx context.Context, request *api.JoinCourseRequest) (*emptypb.Empty, error) {
	var userCourse models.UserCourse
	var course models.Course
	course.ID = request.CourseId
	err := Find(&course, s.Course)
	if err != nil {
		return nil, err
	}
	err = IsDel(course.IsDeleted)
	if err != nil {
		return nil, err
	}
	userCourse.CourseID = request.CourseId
	userCourse.UserID = request.UserId
	userCourse.StartDate = time.Now()
	err = s.Course.DB.Create(&userCourse).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) RenewCourse(ctx context.Context, request *api.RenewCourseRequest) (*emptypb.Empty, error) {
	var userCourse models.UserCourse
	var course models.Course
	course.ID = request.CourseId
	err := Find(&course, s.Course)
	if err != nil {
		return nil, err
	}
	err = IsDel(course.IsDeleted)
	if err != nil {
		return nil, err
	}
	userCourse.CourseID = request.CourseId
	userCourse.UserID = request.UserId
	err := Find(&userCourse, s.Course)
	if err != nil {
		return nil, err
	}
	err = IsDel(userCourse.IsDeleted)
	if err != nil {
		return nil, err
	}
	// also can check for >100
	if request.PersentFinished == 0 {
		return nil, status.New(codes.InvalidArgument, "course was deleted").Err()
	}
	userCourse.PercentFinished = request.PersentFinished
	err = s.Course.DB.Updates(&userCourse).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) DeclineCourse(ctx context.Context, request *api.DeclineCourseRequest) (*emptypb.Empty, error) {
	var userCourse models.UserCourse
	userCourse.CourseID = request.CourseId
	userCourse.UserID = request.UserId
	err := s.Course.DB.First(&userCourse).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	userCourse.IsDeleted = 1
	err = s.Course.DB.Updates(&userCourse).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) ListCourse(ctx context.Context, request *api.ListCourseRequest) (*api.ListCourseResponse, error) {
	var userCourses []models.UserCourse
	var userCourse models.UserCourse
	q := s.Course.DB.Model(&userCourse)
	if request.Limit > 0 {
		q.Limit(int(request.Limit))
	}
	if request.Offset > 0 {
		q.Offset(int(request.Offset))
	}
	if request.ShowDeleted {
		q.Where("is_deleted = ?", 1)
	}
	if request.UserId != "" {
		q.Where("user_id = ?", request.UserId)
	}
	if request.CourseId != "" {
		q.Where("course_id = ?", request.CourseId)
	}
	err := q.Find(&userCourses).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	result := &api.ListCourseResponse{}
	result.UserCourses = make([]*api.UserCourse, 0, len(userCourses))
	for _, uc := range userCourses {
		ucResponse := &api.UserCourse{
			UserId:          uc.UserID,
			CourseId:        uc.CourseID,
			IsDeleted:       uc.IsDeleted,
			StartDate:       timestamppb.New(uc.StartDate),
			PersentFinished: uc.PercentFinished,
		}
		if !uc.FinishDate.IsZero() {
			ucResponse.FinishDate = timestamppb.New(uc.FinishDate)
		}
		result.UserCourses = append(result.UserCourses, ucResponse)
	}
	return result, nil
}
