package service

import (
	"context"
	"golang-united-courses/pkg/api"
	"golang-united-courses/pkg/db"
	"golang-united-courses/pkg/models"

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
	err := s.Course.DB.Create(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	return &api.CreateResponse{
		Id: uint32(course.ID),
	}, nil
}

func (s *CourseServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.ID = uint(request.Id)
	err := s.Course.DB.Unscoped().First(&course).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	result := &api.GetResponse{
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
	err := s.Course.DB.First(&course, request.Id).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	course.ID = uint(request.Id)
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
	err := s.Course.DB.First(&course, request.Id).Error
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()
	}
	course.ID = uint(request.Id)
	course.DeletedBy = request.DeletedBy
	err = s.Course.DB.Updates(&course).Delete(&course).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &emptypb.Empty{}, nil
}

func (s *CourseServer) List(ctx context.Context, request *api.ListRequest) (*api.ListResponse, error) {
	var courses []models.Course
	err := s.Course.DB.Limit(int(request.Limit)).Offset(int(request.Offset)).Find(&courses, "user_id = ?", request.UserId).Error
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	result := &api.ListResponse{}
	result.Courses = make([]*api.GetResponse, 0, len(courses))
	for _, c := range courses {
		result.Courses = append(result.Courses, &api.GetResponse{
			Title:       c.Title,
			Description: c.Description,
			CreatedBy:   c.CreatedBy,
			UserId:      c.UserID,
			CreatedAt:   timestamppb.New(c.CreatedAt),
			UpdatedBy:   c.UpdatedBy,
			UpdatedAt:   timestamppb.New(c.UpdatedAt),
			DeletedBy:   c.DeletedBy,
			DeletedAt:   timestamppb.New(c.DeletedAt.Time),
		})
	}
	return result, nil
}
