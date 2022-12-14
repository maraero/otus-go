package servergrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/config"
	er "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-repository"
	es "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/events-service"
	l "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/logger"
	gges "github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/servers/grpc/generated"
	"github.com/maraero/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SuiteTest struct {
	suite.Suite
	client gges.EventServiceClient
	closer func()
}

func buildTestConfig() config.Config {
	loggerConfig := config.Logger{
		Level:            "info",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	serverConfig := config.Server{
		Host:     "localhost",
		HTTPPort: "8000",
		GrpcPort: "8001",
	}
	storageConfig := config.Storage{
		Type:   "in-memory",
		Driver: "postgres",
		DSN:    "postgresql://admin:admin@localhost:5432/calendar?sslmode=disable",
	}
	return config.Config{
		Logger:  loggerConfig,
		Server:  serverConfig,
		Storage: storageConfig,
	}
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()
	buffer := 101024 * 1024
	lsnr := bufconn.Listen(buffer)

	config := buildTestConfig()

	logger, err := l.New(config.Logger)
	s.Require().NoError(err)

	storage := storage.Storage{Source: storage.StorageInMemory, Connection: nil}
	eventsRepository := er.New(&storage)
	eventsService := es.New(eventsRepository)
	calendar := app.New(eventsService, logger)

	baseServer := grpc.NewServer()
	gges.RegisterEventServiceServer(baseServer, NewService(calendar))

	go func() {
		if err := baseServer.Serve(lsnr); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lsnr.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	s.closer = func() {
		err := lsnr.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	s.client = gges.NewEventServiceClient(conn)
}

func (s *SuiteTest) TeardownTest() {
	s.closer()
}

func TestGrpcServer(t *testing.T) {
	suite.Run(t, new(SuiteTest))
}

func (s *SuiteTest) TestCreateEvent() {
	s.Run("successful", func() {
		ctx := context.Background()
		in := &gges.Event{
			Id:               0,
			Title:            "test",
			DateStart:        timestamppb.New(time.Now()),
			DateEnd:          timestamppb.New(time.Now().Add(1 * time.Hour)),
			Description:      "test description",
			UserId:           "test user id",
			DateNotification: timestamppb.New(time.Time{}),
			Deleted:          false,
		}
		out, err := s.client.CreateEvent(ctx, in)
		s.Require().NoError(err)
		s.Require().Equal(int64(1), out.Id)
	})

	s.Run("invalid argument", func() {
		ctx := context.Background()
		in := &gges.Event{
			Id:               0,
			Title:            "test",
			DateStart:        timestamppb.New(time.Time{}), // <= zero time
			DateEnd:          timestamppb.New(time.Now().Add(1 * time.Hour)),
			Description:      "test description",
			UserId:           "test user id",
			DateNotification: timestamppb.New(time.Time{}),
			Deleted:          false,
		}
		out, err := s.client.CreateEvent(ctx, in)
		s.Require().Error(err)
		s.Require().Nil(out)
	})
}

func (s *SuiteTest) TestUpdateEvent() {
	s.Run("successful", func() {
		ctx := context.Background()
		createdIn := &gges.Event{
			Id:               0,
			Title:            "created title",
			DateStart:        timestamppb.New(time.Now()),
			DateEnd:          timestamppb.New(time.Now().Add(1 * time.Hour)),
			Description:      "test description",
			UserId:           "test user id",
			DateNotification: timestamppb.New(time.Time{}),
			Deleted:          false,
		}
		createdOut, err := s.client.CreateEvent(ctx, createdIn)
		s.Require().NoError(err)
		updatedEvent := createdIn
		updatedEvent.Title = "updated title"
		updatedIn := &gges.UpdateEventRequest{
			Id:    createdOut.Id,
			Event: updatedEvent,
		}
		out, err := s.client.UpdateEvent(ctx, updatedIn)
		s.Require().NoError(err)
		s.Require().NotNil(out)
	})

	s.Run("invalid argument", func() {
		ctx := context.Background()
		event := &gges.Event{
			Id:               0,
			Title:            "created title",
			DateStart:        timestamppb.New(time.Now()),
			DateEnd:          timestamppb.New(time.Now().Add(1 * time.Hour)),
			Description:      "test description",
			UserId:           "test user id",
			DateNotification: timestamppb.New(time.Time{}),
			Deleted:          false,
		}
		in := &gges.UpdateEventRequest{
			Id:    0,
			Event: event,
		}
		out, err := s.client.UpdateEvent(ctx, in)
		s.Require().Error(err)
		s.Require().Nil(out)
	})
}

func (s *SuiteTest) TestDeleteEvent() {
	s.Run("successful", func() {
		ctx := context.Background()
		createdIn := &gges.Event{
			Id:               0,
			Title:            "event to delete",
			DateStart:        timestamppb.New(time.Now()),
			DateEnd:          timestamppb.New(time.Now().Add(1 * time.Hour)),
			Description:      "test description",
			UserId:           "test user id",
			DateNotification: timestamppb.New(time.Time{}),
			Deleted:          false,
		}
		createdOut, err := s.client.CreateEvent(ctx, createdIn)
		s.Require().NoError(err)
		deletedIn := &gges.DeleteEventRequest{
			Id: createdOut.Id,
		}
		out, err := s.client.DeleteEvent(ctx, deletedIn)
		s.Require().NoError(err)
		s.Require().NotNil(out)

		getEventIn := &gges.GetEventByIDRequest{Id: createdOut.Id}
		getEventOut, err := s.client.GetEventByID(ctx, getEventIn)
		s.Require().NoError(err)
		s.Require().NotNil(getEventOut)
		s.Require().Equal(true, getEventOut.Deleted)
	})

	s.Run("invalid argument", func() {
		ctx := context.Background()
		in := &gges.DeleteEventRequest{
			Id: 0,
		}
		out, err := s.client.DeleteEvent(ctx, in)
		s.Require().Error(err)
		s.Require().Nil(out)
	})
}

func (s *SuiteTest) TestGetEventListByDate() {
	s.Run("empty list", func() {
		ctx := context.Background()
		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(0, len(out.Events))
	})

	s.Run("list by date", func() {
		ctx := context.Background()

		week := 7 * 24 * time.Hour
		multiplier := 3

		for i := 0; i < 2; i++ {
			start := time.Now().Add(time.Duration(i*multiplier) * week)
			evt := &gges.Event{
				Id:               0,
				Title:            "test_" + fmt.Sprint(i),
				DateStart:        timestamppb.New(start),
				DateEnd:          timestamppb.New(start.Add(1 * time.Hour)),
				Description:      "test description",
				UserId:           "test user id",
				DateNotification: timestamppb.New(time.Time{}),
				Deleted:          false,
			}
			_, err := s.client.CreateEvent(ctx, evt)
			s.Require().NoError(err)
		}

		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(1, len(out.Events))
		s.Require().Equal("test_0", out.Events[0].Title)
	})
}

func (s *SuiteTest) TestGetEventListByWeek() {
	s.Run("empty list", func() {
		ctx := context.Background()
		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(0, len(out.Events))
	})

	s.Run("list by week", func() {
		ctx := context.Background()

		week := 7 * 24 * time.Hour
		multiplier := 3

		for i := 0; i < 2; i++ {
			start := time.Now().Add(time.Duration(i*multiplier) * week)
			evt := &gges.Event{
				Id:               0,
				Title:            "test_" + fmt.Sprint(i),
				DateStart:        timestamppb.New(start),
				DateEnd:          timestamppb.New(start.Add(1 * time.Hour)),
				Description:      "test description",
				UserId:           "test user id",
				DateNotification: timestamppb.New(time.Time{}),
				Deleted:          false,
			}
			_, err := s.client.CreateEvent(ctx, evt)
			s.Require().NoError(err)
		}

		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(1, len(out.Events))
		s.Require().Equal("test_0", out.Events[0].Title)
	})
}

func (s *SuiteTest) TestGetEventListByMonth() {
	s.Run("empty list", func() {
		ctx := context.Background()
		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(0, len(out.Events))
	})

	s.Run("list by week", func() {
		ctx := context.Background()

		week := 7 * 24 * time.Hour
		multiplier := 5

		for i := 0; i < 2; i++ {
			start := time.Now().Add(time.Duration(i*multiplier) * week)
			evt := &gges.Event{
				Id:               0,
				Title:            "test_" + fmt.Sprint(i),
				DateStart:        timestamppb.New(start),
				DateEnd:          timestamppb.New(start.Add(1 * time.Hour)),
				Description:      "test description",
				UserId:           "test user id",
				DateNotification: timestamppb.New(time.Time{}),
				Deleted:          false,
			}
			_, err := s.client.CreateEvent(ctx, evt)
			s.Require().NoError(err)
		}

		in := &gges.EventListRequest{
			Date: timestamppb.New(time.Now()),
		}
		out, err := s.client.GetEventListByDate(ctx, in)
		s.Require().NoError(err)
		s.Require().NotNil(out)
		s.Require().Equal(1, len(out.Events))
		s.Require().Equal("test_0", out.Events[0].Title)
	})
}
