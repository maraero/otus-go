package servergrpc

import (
	"context"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events"
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
	id, err := s.app.EventsService.CreateEvent(ctx, grpcEventToDomainEvent(req))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.CreateEventResponse{Id: id}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *gges.UpdateEventRequest) (*gges.UpdateEventResponse, error) {
	err := s.app.EventsService.UpdateEvent(ctx, req.Id, grpcEventToDomainEvent(req.Event))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.UpdateEventResponse{}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *gges.DeleteEventRequest) (*gges.DeleteEventResponse, error) {
	err := s.app.EventsService.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.DeleteEventResponse{}, nil
}

func (s *Service) GetEventListByDate(ctx context.Context, req *gges.EventListRequest) (*gges.EventListResponse, error) {
	list, err := s.app.EventsService.GetEventListByDate(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByWeek(ctx context.Context, req *gges.EventListRequest) (*gges.EventListResponse, error) {
	list, err := s.app.EventsService.GetEventListByWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventListByMonth(
	ctx context.Context,
	req *gges.EventListRequest,
) (*gges.EventListResponse, error) {
	list, err := s.app.EventsService.GetEventListByMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &gges.EventListResponse{Events: domainEventListToGrpcEventList(list)}, nil
}

func (s *Service) GetEventByID(ctx context.Context, req *gges.GetEventByIDRequest) (*gges.Event, error) {
	evt, err := s.app.EventsService.GetEventByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	event := domainEventToGrpcEvent(evt)
	return &event, nil
}

func grpcEventToDomainEvent(grpcEvent *gges.Event) events.Event {
	return events.Event{
		ID:               grpcEvent.Id,
		Title:            grpcEvent.Title,
		DateStart:        grpcEvent.DateStart.AsTime(),
		DateEnd:          grpcEvent.DateEnd.AsTime(),
		Description:      grpcEvent.Description,
		UserID:           grpcEvent.UserId,
		DateNotification: grpcEvent.DateNotification.AsTime(),
	}
}

func domainEventToGrpcEvent(domainEvent events.Event) gges.Event {
	return gges.Event{
		Id:               domainEvent.ID,
		Title:            domainEvent.Title,
		DateStart:        timestamppb.New(domainEvent.DateStart),
		DateEnd:          timestamppb.New(domainEvent.DateEnd),
		Description:      domainEvent.Description,
		UserId:           domainEvent.UserID,
		DateNotification: timestamppb.New(domainEvent.DateNotification),
	}
}

func domainEventListToGrpcEventList(domainEvents []events.Event) []*gges.Event {
	result := make([]*gges.Event, 0, len(domainEvents))
	for _, event := range domainEvents {
		evt := domainEventToGrpcEvent(event)
		result = append(result, &evt)
	}
	return result
}
