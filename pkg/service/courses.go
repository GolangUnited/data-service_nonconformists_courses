package service

import (
	"context"
	"errors"
	"fmt"
	"golang-united-courses/pkg/api"
	"golang-united-courses/pkg/db"
	"golang-united-courses/pkg/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
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
	course.UserID = request.UserId
	t := s.Course.DB.Create(&course)
	if t.Error != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Can't create item. Reason: %s", t.Error))
	}
	return &api.CreateResponse{
		Id: uint32(course.ID),
	}, nil
}

func (s *CourseServer) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.ID = uint(request.Id)
	t := s.Course.DB.First(&course)
	if t.Error != nil {
		switch {
		case errors.Is(t.Error, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Can't get item with id %d.", course.ID))
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("Can't perform Get request. Reason: %s", t.Error))
		}
	}
	return &api.GetResponse{
		Title:       course.Title,
		Description: course.Description,
		CreatedBy:   course.CreatedBy,
		UserId:      course.UserID,
		CreatedAt:   timestamppb.New(course.CreatedAt),
		UpdatedBy:   course.UpdatedBy,
		UpdatedAt:   timestamppb.New(course.UpdatedAt),
		DeletedBy:   course.DeletedBy,
		DeletedAt:   timestamppb.New(course.DeletedAt.Time),
	}, nil
}

func (s *CourseServer) Update(ctx context.Context, request *api.UpdateRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = uint(request.Id)
	course.Title = request.Title
	course.Description = request.Description
	course.UpdatedBy = request.UpdatedBy
	course.UserID = request.UserId
	t := s.Course.DB.First(&course, course.ID).Updates(&course)
	if t.Error != nil {
		switch {
		case errors.Is(t.Error, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Can't update item with id %d.", course.ID))
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("Can't perform Update request. Reason: %s", t.Error))
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = uint(request.Id)
	course.DeletedBy = request.DeletedBy
	t := s.Course.DB.First(&course, course.ID).Updates(&course).Delete(&course)
	if t.Error != nil {
		switch {
		case errors.Is(t.Error, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Can't delete item with id %d.", course.ID))
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("Can't perform Delete request. Reason: %s", t.Error))
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) List(ctx context.Context, request *api.ListRequest) (*api.ListResponse, error) {
	var courses []models.Course
	t := s.Course.DB.Limit(int(request.Limit)).Offset(int(request.Offset)).Find(&courses, "user_id = ?", request.UserId)
	if t.Error != nil {
		switch {
		case errors.Is(t.Error, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Can't list items with id %d.", request.UserId))
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("Can't perform List request. Reason: %s", t.Error))
		}
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
