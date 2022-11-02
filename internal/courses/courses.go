package courses

import (
	"context"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/db"
	"golang-united-courses/internal/models"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	C *db.Handler
	api.UnimplementedCoursesServer
}

func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	var course models.Course
	course.Title = request.Title
	course.Description = request.Description
	course.CreatedBy = request.CreatedBy
	t := s.C.DB.Create(&course)
	if t.Error != nil {
		return nil, t.Error
	}
	return &api.CreateResponse{
		Id: uint32(course.ID),
	}, nil
}

func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.ID = uint(request.Id)
	t := s.C.DB.First(&course)
	if t.Error != nil {
		return nil, t.Error
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

func (s *Server) Update(ctx context.Context, request *api.UpdateRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = uint(request.Id)
	course.Title = request.Title
	course.Description = request.Description
	course.UpdatedBy = request.UpdatedBy
	t := s.C.DB.Updates(&course)
	if t.Error != nil {
		return nil, t.Error
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.ID = uint(request.Id)
	course.DeletedBy = request.DeletedBy
	t := s.C.DB.Updates(&course).Delete(&course)
	if t.Error != nil {
		return nil, t.Error
	}
	return &emptypb.Empty{}, nil
}
