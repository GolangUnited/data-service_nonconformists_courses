package courses

import (
	"context"
	"licheropew/golang-united-courses/pkg/api"
	"licheropew/golang-united-courses/pkg/db"
	"licheropew/golang-united-courses/pkg/models"
	"log"

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
	t := s.C.DB.Create(&course)
	if t.Error != nil {
		log.Fatal(t.Error)
	}
	return &api.CreateResponse{
		Id: course.Id,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	var course models.Course
	course.Id = request.Id
	t := s.C.DB.First(&course)
	if t.Error != nil {
		log.Fatal(t.Error)
	}
	return &api.GetResponse{Title: course.Title, Description: course.Description, CreatedBy: course.CreatedBy, CreatedAt: timestamppb.New(course.CreatedAt)}, nil
}

func (s *Server) Update(ctx context.Context, request *api.UpdateRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.Id = request.Id
	course.Title = request.Title
	course.Description = request.Description
	s.C.DB.Updates(&course)
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	var course models.Course
	course.Id = request.Id
	s.C.DB.Delete(&course)
	return &emptypb.Empty{}, nil
}
