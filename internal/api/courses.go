package api

import (
	"context"
	"golang-united-courses/internal/interfaces.go"
	"golang-united-courses/internal/repositories/courses"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CourseServer struct {
	DB interfaces.CourseManager
	UnimplementedCoursesServer
}

func (cs *CourseServer) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	result, err := cs.DB.Create(request.GetTitle(), request.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &CreateResponse{Id: result}, nil
}

func (cs *CourseServer) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	result, err := cs.DB.GetById(request.GetId())
	if err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &GetResponse{
		Id:          result.ID.String(),
		Title:       result.Title,
		Description: result.Description,
		CreatedAt:   timestamppb.New(result.CreatedAt),
		IsDeleted:   result.IsDeleted,
	}, nil
}

func (cs *CourseServer) Update(ctx context.Context, request *UpdateRequest) (*emptypb.Empty, error) {
	if err := cs.DB.Update(request.GetId(), request.GetTitle(), request.GetDescription()); err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		case courses.ErrorCourseWasDeleted.Error():
			return nil, status.Error(codes.Aborted, courses.ErrorCourseWasDeleted.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) Delete(ctx context.Context, request *DeleteRequest) (*emptypb.Empty, error) {
	if err := cs.DB.Delete(request.GetId()); err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		case courses.ErrorCourseWasDeleted.Error():
			return nil, status.Error(codes.Aborted, courses.ErrorCourseWasDeleted.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) List(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	courses, err := cs.DB.List(request.ShowDeleted, request.Limit, request.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := &ListResponse{}
	result.Courses = make([]*GetResponse, 0, len(courses))
	for _, c := range courses {
		cResponse := &GetResponse{
			Id:          c.ID.String(),
			Title:       c.Title,
			Description: c.Description,
			CreatedAt:   timestamppb.New(c.CreatedAt),
			IsDeleted:   c.IsDeleted,
		}
		result.Courses = append(result.Courses, cResponse)
	}
	return result, nil
}

func (cs *CourseServer) JoinCourse(ctx context.Context, request *JoinCourseRequest) (*emptypb.Empty, error) {
	err := cs.DB.Join(request.GetUserId(), request.GetCourseId())
	if err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return &emptypb.Empty{}, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		case courses.ErrorCourseWasDeleted.Error():
			return &emptypb.Empty{}, status.Error(codes.Aborted, courses.ErrorCourseWasDeleted.Error())
		default:
			return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) GetUserCourse(ctx context.Context, request *GetUserCourseRequest) (*UserCourse, error) {
	result, err := cs.DB.GetUserCourse(request.GetUserId(), request.GetCourseId())
	if err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		case courses.ErrorCourseWasDeleted.Error():
			return nil, status.Error(codes.Aborted, courses.ErrorCourseWasDeleted.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	uc := &UserCourse{
		UserId:          result.UserID.String(),
		CourseId:        result.CourseID.String(),
		PersentFinished: result.PercentFinished,
		Status:          result.Status,
		CreatedAt:       timestamppb.New(result.CreatedAt),
	}
	if !result.FinishDate.IsZero() {
		uc.FinishDate = timestamppb.New(result.FinishDate)
	}
	if !result.StartDate.IsZero() {
		uc.StartDate = timestamppb.New(result.StartDate)
	}
	return uc, nil
}

func (cs *CourseServer) SetProgress(ctx context.Context, request *SetProgressRequest) (*emptypb.Empty, error) {
	err := cs.DB.SetProgress(request.GetUserId(), request.GetCourseId(), request.GetPersentFinished())
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) SetStatus(ctx context.Context, request *SetStatusRequest) (*emptypb.Empty, error) {
	err := cs.DB.SetStatus(request.GetUserId(), request.GetCourseId(), request.GetStatus())
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) ListUserCourse(ctx context.Context, request *ListUserCourseRequest) (*ListUserCourseResponse, error) {
	userCourses, err := cs.DB.ListUserCourse(request.GetUserId(), request.GetCourseId(), request.GetLimit(), request.GetOffset(), request.GetShowDeleted())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := &ListUserCourseResponse{}
	result.UserCourses = make([]*UserCourse, 0, len(userCourses))
	for _, uc := range userCourses {
		ucResponse := &UserCourse{
			UserId:          uc.UserID.String(),
			CourseId:        uc.CourseID.String(),
			PersentFinished: uc.PercentFinished,
			Status:          uc.Status,
			CreatedAt:       timestamppb.New(uc.CreatedAt),
		}
		if !uc.FinishDate.IsZero() {
			ucResponse.FinishDate = timestamppb.New(uc.FinishDate)
		}
		if !uc.StartDate.IsZero() {
			ucResponse.StartDate = timestamppb.New(uc.StartDate)
		}
		result.UserCourses = append(result.UserCourses, ucResponse)
	}
	return result, nil
}
