package server_grpc

import (
	"context"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	event "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/event-service/domain"
	gges "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/grpc/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	gges.UnimplementedEventServiceServer
	app *app.App
}

func NewService(app *app.App) *Service {
	return &Service{app: app}
}

func (s *Service) CreateEvent(ctx context.Context, req *gges.Event) (*gges.CreateEventResponse, error) {
	id, err := s.app.Event_service.CreateEvent(ctx, grpcEventToDomainEvent(req))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.CreateEventResponse{Id: id}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *gges.UpdateEventRequest) (*gges.UpdateEventResposnse, error) {
	err := s.app.Event_service.UpdateEvent(ctx, req.Id, grpcEventToDomainEvent(req.E))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.UpdateEventResposnse{}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *gges.DeleteEventRequest) (*gges.DeleteEventResponse, error) {
	err := s.app.Event_service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.DeleteEventResponse{}, nil
}

func (s *Service) GetEventListByDate(ctx context.Context, req *gges.EventListRequest) (*gges.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByDate(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByWeek(ctx context.Context, req *gges.EventListRequest) (*gges.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByMonth(ctx context.Context, req *gges.EventListRequest) (*gges.EventListResponse, error) {
	list, err := s.app.Event_service.GetEventListByMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventByID(ctx context.Context, req *gges.GetEventByIDRequest) (*gges.Event, error) {
	evt, err := s.app.Event_service.GetEventByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	event := domainEventToGrpcEvent(evt)
	return &event, nil
}

func grpcEventToDomainEvent(grpcEvent *gges.Event) event.Event {
	return event.Event{
		ID:               grpcEvent.Id,
		Title:            grpcEvent.Title,
		DateStart:        grpcEvent.DateStart.AsTime(),
		DateEnd:          grpcEvent.DateEnd.AsTime(),
		Description:      grpcEvent.Description,
		UserID:           grpcEvent.UserId,
		DateNotification: grpcEvent.DateNotification.AsTime(),
		Deleted:          grpcEvent.Deleted,
	}
}

func domainEventToGrpcEvent(domainEvent event.Event) gges.Event {
	return gges.Event{
		Id:               domainEvent.ID,
		Title:            domainEvent.Title,
		DateStart:        timestamppb.New(domainEvent.DateStart),
		DateEnd:          timestamppb.New(domainEvent.DateEnd),
		Description:      domainEvent.Description,
		UserId:           domainEvent.UserID,
		DateNotification: timestamppb.New(domainEvent.DateNotification),
		Deleted:          domainEvent.Deleted,
	}
}

func domainEventListToGrpcEventList(domainEvents []event.Event) []*gges.Event {
	result := make([]*gges.Event, 0, len(domainEvents))
	for _, event := range domainEvents {
		evt := domainEventToGrpcEvent(event)
		result = append(result, &evt)
	}
	return result
}
