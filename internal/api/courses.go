package api

import (
	"context"
	"golang-united-courses/internal/interfaces.go"
	"golang-united-courses/internal/models"
	"golang-united-courses/internal/repositories/courses"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CourseServer struct {
	DB interfaces.CourseManager
	UnimplementedCoursesServer
}

func checkIsValidUUID(id string) (uuid.UUID, error) {
	cid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return cid, nil
}

func getUserCourseUUID(cid, uid string) (models.UserCourse, error) {
	var uc models.UserCourse
	courseId, err := checkIsValidUUID(cid)
	if err != nil {
		return uc, err
	}
	uc.CourseID = courseId
	userId, err := checkIsValidUUID(uid)
	if err != nil {
		return uc, err
	}
	uc.UserID = userId
	return uc, nil
}

func (cs *CourseServer) checkUserCourse(cid, uid string) (models.UserCourse, error) {
	uc, err := getUserCourseUUID(cid, uid)
	if err != nil {
		return uc, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
	err = cs.DB.GetUserCourse(&uc)
	if err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return uc, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		default:
			return uc, status.Error(codes.Internal, err.Error())
		}
	}
	return uc, nil
}

func (cs *CourseServer) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	result, err := cs.DB.Create(request.GetTitle(), request.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &CreateResponse{Id: result}, nil
}

func (cs *CourseServer) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	_, err := checkIsValidUUID(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
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
	_, err := checkIsValidUUID(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
	if err := cs.DB.Update(request.GetId(), request.GetTitle(), request.GetDescription()); err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) Delete(ctx context.Context, request *DeleteRequest) (*emptypb.Empty, error) {
	_, err := checkIsValidUUID(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
	if err := cs.DB.Delete(request.GetId()); err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		case courses.ErrCourseWasDeleted.Error():
			return nil, status.Error(codes.Aborted, courses.ErrCourseWasDeleted.Error())
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
	uc, err := getUserCourseUUID(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
	result, err := cs.DB.GetById(request.GetCourseId())
	if err != nil {
		switch err.Error() {
		case courses.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, courses.ErrCourseNotFound.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if result.IsDeleted != 0 {
		return nil, status.Error(codes.Aborted, courses.ErrCourseWasDeleted.Error())
	}
	err = cs.DB.Join(uc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) GetUserCourse(ctx context.Context, request *GetUserCourseRequest) (*UserCourseResponse, error) {
	uc, err := cs.checkUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	result := &UserCourseResponse{
		UserId:          uc.UserID.String(),
		CourseId:        uc.CourseID.String(),
		PersentFinished: uc.PercentFinished,
		Status:          uc.Status,
		CreatedAt:       timestamppb.New(uc.CreatedAt),
	}
	if !uc.FinishDate.IsZero() {
		result.FinishDate = timestamppb.New(uc.FinishDate)
	}
	if !uc.StartDate.IsZero() {
		result.StartDate = timestamppb.New(uc.StartDate)
	}
	return result, nil
}

func (cs *CourseServer) SetProgress(ctx context.Context, request *SetProgressRequest) (*emptypb.Empty, error) {
	uc, err := cs.checkUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	if request.GetPersentFinished() > 100 {
		return nil, status.Error(codes.InvalidArgument, courses.ErrInvalidFormat.Error())
	}
	uc.PercentFinished = request.GetPersentFinished()
	err = cs.DB.UpdateUserCourse(uc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) SetStatus(ctx context.Context, request *SetStatusRequest) (*emptypb.Empty, error) {
	uc, err := cs.checkUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	//TODO: statutes -> INT and create MAP?
	switch request.GetStatus() {
	case "start":
		uc.StartDate = time.Now()
	case "finish":
		uc.FinishDate = time.Now()
	default:
		return nil, status.Error(codes.InvalidArgument, courses.ErrIncorrectArgument.Error())
	}
	uc.Status = request.GetStatus()
	err = cs.DB.UpdateUserCourse(uc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) ListUserCourse(ctx context.Context, request *ListUserCourseRequest) (*ListUserCourseResponse, error) {
	userCourses, err := cs.DB.ListUserCourse(request.GetUserId(), request.GetCourseId(), request.GetLimit(), request.GetOffset(), request.GetShowDeleted())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := &ListUserCourseResponse{}
	result.UserCourses = make([]*UserCourseResponse, 0, len(userCourses))
	for _, uc := range userCourses {
		ucResponse := &UserCourseResponse{
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
