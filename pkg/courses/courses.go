package courses

import (
	"context"
	"fmt"
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
	t := s.Course.DB.Create(&course)
	if t.Error != nil {
		err := status.Error(codes.Internal, fmt.Sprintf("Can't create item. Reason: %s", t.Error))
		return nil, err
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
		err := status.Error(codes.NotFound, fmt.Sprintf("item with id %d was not found", course.ID))
		return nil, err
	}
	return &api.GetResponse{
		Title:       course.Title,
		Description: course.Description,
		CreatedBy:   course.CreatedBy,
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
	t := s.Course.DB.Updates(&course)
	if t.Error != nil {
		err := status.Error(codes.Internal, fmt.Sprintf("Can't update item. Reason: %s", t.Error))
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CourseServer) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = uint(request.Id)
	course.DeletedBy = request.DeletedBy
	t := s.Course.DB.Updates(&course).Delete(&course)
	if t.Error != nil {
		err := status.Error(codes.Internal, fmt.Sprintf("Can't delete item. Reason: %s", t.Error))
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
