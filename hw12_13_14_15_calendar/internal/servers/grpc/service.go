package server_grpc

import (
	"context"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	esgg "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/grpc-gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	esgg.UnimplementedEventServiceServer
	app *app.App
}

func NewService(app *app.App) *Service {
	return &Service{app: app}
}

func (s *Service) CreateEvent(ctx context.Context, req *esgg.Event) (*esgg.CreateEventResponse, error) {
	id, err := s.app.Event_service.CreateEvent(ctx, grpcEventToDomainEvent(req))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.CreateEventResponse{Id: id}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *esgg.UpdateEventRequest) (*esgg.UpdateEventResposnse, error) {
	err := s.app.Event_service.UpdateEvent(ctx, req.Id, grpcEventToDomainEvent(req.E))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.UpdateEventResposnse{}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *esgg.DeleteEventRequest) (*esgg.DeleteEventResponse, error) {
	err := s.app.Event_service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.DeleteEventResponse{}, nil
}

func (s *Service) GetEventListByDate(ctx context.Context, req *esgg.EventListRequest) (*esgg.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByDate(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByWeek(ctx context.Context, req *esgg.EventListRequest) (*esgg.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByMonth(ctx context.Context, req *esgg.EventListRequest) (*esgg.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &esgg.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func grpcEventToDomainEvent(grpcEvent *esgg.Event) event.Event {
	return event.Event{
		ID:               grpcEvent.Id,
		Title:            grpcEvent.Title,
		DateStart:        grpcEvent.DateStart.AsTime(),
		DateEnd:          grpcEvent.DateEnd.AsTime(),
		Descripion:       grpcEvent.Description,
		UserID:           grpcEvent.UserId,
		DateNotification: grpcEvent.DateNotification.AsTime(),
		Deleted:          grpcEvent.Deleted,
	}
}

func domainEventToGrpcEvent(domainEvent event.Event) esgg.Event {
	return esgg.Event{
		Id:               domainEvent.ID,
		Title:            domainEvent.Title,
		DateStart:        timestamppb.New(domainEvent.DateStart),
		DateEnd:          timestamppb.New(domainEvent.DateEnd),
		Description:      domainEvent.Descripion,
		UserId:           domainEvent.UserID,
		DateNotification: timestamppb.New(domainEvent.DateNotification),
		Deleted:          domainEvent.Deleted,
	}
}

func domainEventListToGrpcEventList(domainEvents []event.Event) []*esgg.Event {
	result := make([]*esgg.Event, 0, len(domainEvents))
	for _, event := range domainEvents {
		evt := domainEventToGrpcEvent(event)
		result = append(result, &evt)
	}
	return result
}

// type Event struct {
// 	ID               int64     `db:"id"`
// 	Title            string    `db:"title"`
// 	DateStart        time.Time `db:"date_start"`
// 	DateEnd          time.Time `db:"date_end"`
// 	Descripion       string    `db:"description"`
// 	UserID           string    `db:"user_id"`
// 	DateNotification time.Time `db:"date_notification"`
// 	Deleted          bool      `db:"deleted"`
// }

// type Event struct {
// 	state         protoimpl.MessageState
// 	sizeCache     protoimpl.SizeCache
// 	unknownFields protoimpl.UnknownFields

// 	Id               int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
// 	Title            string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
// 	DateStart        *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=date_start,json=dateStart,proto3" json:"date_start,omitempty"`
// 	DateEnd          *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date_end,json=dateEnd,proto3" json:"date_end,omitempty"`
// 	Description      string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
// 	UserId           string                 `protobuf:"bytes,6,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
// 	DateNotification *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=date_notification,json=dateNotification,proto3" json:"date_notification,omitempty"`
// 	Deleted          bool                   `protobuf:"varint,8,opt,name=deleted,proto3" json:"deleted,omitempty"`
// }
