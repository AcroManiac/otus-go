package main

import (
	"os"
	"testing"

	"github.com/cucumber/messages-go/v10"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/database"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:      "progress",
		Paths:       []string{"features"},
		Randomize:   0,
		Concurrency: 0,
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {

	// Initialize connection to Calendar API
	s.Step(`^Connection to Calendar API on "([^"]*)"$`, connectionToCalendarAPIOn)

	// Make AddEvent test
	add := &addEventTest{
		eventData:      nil,
		createResponse: nil,
	}
	s.Step(`^There is the event:$`, add.thereIsTheEvent)
	s.Step(`^I send AddEvent request$`, add.iSendAddEventRequest)
	s.Step(`^response should have event$`, add.responseShouldHaveEvent)

	// Make GetEvents test
	get := &getEventsTest{}
	s.Step(`^I send GetEvents request with period "([^"]*)" and start time "([^"]*)"$`,
		get.iSendGetEventsRequestWithPeriodAndStartTime)
	s.Step(`^search should return (\d+) event$`, get.searchShouldReturnEvent)

	// Close connection to Calendar API
	s.AfterScenario(closeClient)

	// Make SendNotification test
	send := &sendNotificationTest{}
	s.Step(`^Connection to PostgreSQL service with DSN "([^"]*)"$`,
		send.connectionToPostgreSQLServiceWithDSN)
	s.Step(`^I wait for (\d+) minute$`, send.iWaitForMinute)
	s.Step(`^selection from table "([^"]*)" should return (\d+) notice$`,
		send.selectionFromTableShouldReturnNotice)
}

// Global variables for multiple tests execution
var (
	clientConn *grpc.ClientConn
	grpcClient api.CalendarApiClient
	db         *database.Connection
	eventID    string
)

func connectionToCalendarAPIOn(arg1 string) error {

	// Start gRPC client
	var err error
	clientConn, err = grpc.Dial(arg1, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "could not connect gRPC server")
	}

	grpcClient = api.NewCalendarApiClient(clientConn)
	if grpcClient == nil {
		return errors.New("failed creating calendar API client")
	}

	return nil
}

func closeClient(*messages.Pickle, error) {
	if clientConn != nil {
		_ = clientConn.Close()
	}
}
