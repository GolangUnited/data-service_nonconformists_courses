package api

import (
	"context"
	"golang-united-courses/internal"
	"golang-united-courses/internal/interfaces"
	"golang-united-courses/internal/models"
	"log"
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

func New(icm interfaces.CourseManager) *CourseServer {
	return &CourseServer{
		DB: icm,
	}
}

func checkUserCourseUUID(cid, uid string) bool {
	if _, err := uuid.Parse(cid); err != nil {
		return false
	}
	if _, err := uuid.Parse(uid); err != nil {
		return false
	}
	return true
}

func (cs *CourseServer) getUserCourse(cid, uid string) (models.UserCourse, error) {
	if uuidCheck := checkUserCourseUUID(cid, uid); !uuidCheck {
		return models.UserCourse{}, status.Error(codes.InvalidArgument, internal.ErrInvalidUUIDFormat.Error())
	}
	uc, err := cs.DB.GetUserCourse(uid, cid)
	if err != nil {
		switch err.Error() {
		case internal.ErrUserCourseNotFound.Error():
			return models.UserCourse{}, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("internal error: %s", err.Error())
			return models.UserCourse{}, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	}
	return uc, nil
}

func (cs *CourseServer) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	result, err := cs.DB.Create(request.GetTitle(), request.GetDescription())
	if err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	return &CreateResponse{Id: result}, nil
}

func (cs *CourseServer) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	if _, err := uuid.Parse(request.GetId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, internal.ErrInvalidUUIDFormat.Error())
	}
	result, err := cs.DB.GetById(request.GetId())
	if err != nil {
		switch err.Error() {
		case internal.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
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
	if _, err := uuid.Parse(request.GetId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, internal.ErrInvalidUUIDFormat.Error())
	}
	if err := cs.DB.Update(request.GetId(), request.GetTitle(), request.GetDescription()); err != nil {
		switch err.Error() {
		case internal.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) Delete(ctx context.Context, request *DeleteRequest) (*emptypb.Empty, error) {
	if _, err := uuid.Parse(request.GetId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, internal.ErrInvalidUUIDFormat.Error())
	}
	if err := cs.DB.Delete(request.GetId()); err != nil {
		switch err.Error() {
		case internal.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, err.Error())
		case internal.ErrCourseWasDeleted.Error():
			log.Println("attempt to delete already deleted value")
		default:
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) List(ctx context.Context, request *ListRequest) (*ListResponse, error) {
	courses, err := cs.DB.List(request.GetShowDeleted(), request.GetLimit(), request.GetOffset())
	if err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	// second request to get and return number of values without Limit and Offset (send it as zero values) -> goes to Total value as len()
	coursesTotal, err := cs.DB.List(request.GetShowDeleted(), 0, 0)
	if err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	result := &ListResponse{}
	result.Total = int32(len(coursesTotal))
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
	if uuidCheck := checkUserCourseUUID(request.GetCourseId(), request.GetUserId()); !uuidCheck {
		return nil, status.Error(codes.InvalidArgument, internal.ErrInvalidUUIDFormat.Error())
	}
	course, err := cs.DB.GetById(request.GetCourseId())
	if err != nil {
		switch err.Error() {
		case internal.ErrCourseNotFound.Error():
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	}
	if course.IsDeleted != 0 {
		return nil, status.Error(codes.Aborted, internal.ErrCourseWasDeleted.Error())
	}
	uc, err := cs.DB.GetUserCourse(request.GetUserId(), request.GetCourseId())
	if err != nil {
		switch err.Error() {
		case internal.ErrUserCourseNotFound.Error():
			uc.CourseID = uuid.MustParse(request.GetCourseId())
			uc.UserID = uuid.MustParse(request.GetUserId())
			if err = cs.DB.Join(uc); err != nil {
				log.Printf("internal error: %s", err.Error())
				return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
			}
		default:
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	}
	if uc.Status == models.Declined {
		uc.Status = models.Joined
		if err = cs.DB.UpdateUserCourse(uc); err != nil {
			log.Printf("internal error: %s", err.Error())
			return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
		}
	} else {
		log.Println("Attempt to Join already accepted course.")
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) GetUserCourse(ctx context.Context, request *GetUserCourseRequest) (*UserCourseResponse, error) {
	uc, err := cs.getUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	result := &UserCourseResponse{
		UserId:          uc.UserID.String(),
		CourseId:        uc.CourseID.String(),
		Status:          selectStatus(uc.Status),
		CreatedAt:       timestamppb.New(uc.CreatedAt),
		PercentFinished: uc.PercentFinished,
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
	uc, err := cs.getUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	if request.GetPercentFinished() > 100 {
		return nil, status.Error(codes.OutOfRange, internal.ErrOutOfRange.Error())
	}
	uc.PercentFinished = request.GetPercentFinished()
	if err = cs.DB.UpdateUserCourse(uc); err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) SetStatus(ctx context.Context, request *SetStatusRequest) (*emptypb.Empty, error) {
	uc, err := cs.getUserCourse(request.GetCourseId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	switch request.GetStatus() {
	case Statuses_STATUS_STARTED:
		uc.StartDate = time.Now()
		uc.Status = models.Started
	case Statuses_STATUS_FINISHED:
		uc.FinishDate = time.Now()
		uc.Status = models.Finished
	case Statuses_STATUS_DECLINED:
		uc.Status = models.Declined
	case Statuses_STATUS_JOINED:
		uc.Status = models.Joined
	default:
		return nil, status.Error(codes.InvalidArgument, internal.ErrInvalidStatus.Error())
	}
	if err = cs.DB.UpdateUserCourse(uc); err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	return &emptypb.Empty{}, nil
}

func (cs *CourseServer) ListUserCourse(ctx context.Context, request *ListUserCourseRequest) (*ListUserCourseResponse, error) {
	userCourses, err := cs.DB.ListUserCourse(request.GetUserId(), request.GetCourseId(), request.GetLimit(), request.GetOffset(), request.GetShowDeleted())
	if err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	// second request to get and return number of values without Limit and Offset (send it as zero values) -> goes to Total value as len()
	userCoursesTotal, err := cs.DB.ListUserCourse(request.GetUserId(), request.GetCourseId(), 0, 0, request.GetShowDeleted())
	if err != nil {
		log.Printf("internal error: %s", err.Error())
		return nil, status.Error(codes.Internal, internal.ErrInternal.Error())
	}
	result := &ListUserCourseResponse{}
	result.Total = int32(len(userCoursesTotal))
	result.UserCourses = make([]*UserCourseResponse, 0, len(userCourses))
	for _, uc := range userCourses {
		ucResponse := &UserCourseResponse{
			UserId:          uc.UserID.String(),
			CourseId:        uc.CourseID.String(),
			CreatedAt:       timestamppb.New(uc.CreatedAt),
			Status:          selectStatus(uc.Status),
			PercentFinished: uc.PercentFinished,
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

func selectStatus(status models.Statuses) (s Statuses) {
	switch status {
	case models.Joined:
		s = Statuses_STATUS_JOINED
	case models.Started:
		s = Statuses_STATUS_STARTED
	case models.Finished:
		s = Statuses_STATUS_FINISHED
	case models.Declined:
		s = Statuses_STATUS_DECLINED
	default:
		s = Statuses_STATUS_UNKNOWN
	}
	return
}
