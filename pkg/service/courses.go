package service

import (
	"context"
	"golang-united-courses/pkg/api"
	"golang-united-courses/pkg/db"
	"golang-united-courses/pkg/models"

	"github.com/google/uuid"
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
	course.CreatedBy = request.CreatedBy
	course.UpdatedBy = request.CreatedBy
	course.UserID = request.UserId
	course.ID = uuid.New().String()
	err := s.Course.DB.Create(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &api.CreateResponse{
		Id: course.ID,
	}, nil
}

func (s *CourseServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.ID = request.Id
	err := s.Course.DB.Unscoped().First(&course).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	result := &api.GetResponse{
		Id:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		CreatedBy:   course.CreatedBy,
		UserId:      course.UserID,
		CreatedAt:   timestamppb.New(course.CreatedAt),
		UpdatedBy:   course.UpdatedBy,
		UpdatedAt:   timestamppb.New(course.UpdatedAt),
		DeletedBy:   course.DeletedBy,
	}
	if !course.DeletedAt.Time.IsZero() {
		result.DeletedAt = timestamppb.New(course.DeletedAt.Time)
	}
	return result, nil
}

func (s *CourseServer) Update(ctx context.Context, request *api.UpdateRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = request.Id
	err := s.Course.DB.First(&course).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	if !course.DeletedAt.Time.IsZero() {
		return nil, status.New(codes.Aborted, "course was deleted").Err()
	}
	course.Title = request.Title
	course.Description = request.Description
	course.UpdatedBy = request.UpdatedBy
	course.UserID = request.UserId
	err = s.Course.DB.Updates(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = request.Id
	err := s.Course.DB.First(&course).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	if !course.DeletedAt.Time.IsZero() {
		return nil, status.New(codes.Aborted, "course was deleted").Err()
	}
	course.DeletedBy = request.DeletedBy
	err = s.Course.DB.Updates(&course).Delete(&course).Error
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
		q.Unscoped()
	}
	if request.UserId != "" {
		q.Where("user_id = ?", request.UserId)
	}
	err := q.Find(&courses).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	result := &api.ListResponse{}
	result.Courses = make([]*api.GetResponse, 0, len(courses))
	for _, c := range courses {
		courseResponse := &api.GetResponse{
			Id:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			CreatedBy:   c.CreatedBy,
			UserId:      c.UserID,
			CreatedAt:   timestamppb.New(c.CreatedAt),
			UpdatedBy:   c.UpdatedBy,
			UpdatedAt:   timestamppb.New(c.UpdatedAt),
			DeletedBy:   c.DeletedBy,
		}
		if !c.DeletedAt.Time.IsZero() {
			courseResponse.DeletedAt = timestamppb.New(c.DeletedAt.Time)
		}
		result.Courses = append(result.Courses, courseResponse)
	}
	return result, nil
}
