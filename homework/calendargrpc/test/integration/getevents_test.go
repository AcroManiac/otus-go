package main

import (
	"context"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/pkg/api"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
)

type getEventsTest struct {
	searchResponse *api.GetEventsResponse
}

func ToTimePeriod(s string) (p api.TimePeriod) {
	switch s {
	case "day":
		p = api.TimePeriod_TIME_DAY
	case "week":
		p = api.TimePeriod_TIME_WEEK
	case "month":
		p = api.TimePeriod_TIME_MONTH
	default:
		p = api.TimePeriod_TIME_UNKNOWN
	}
	return
}

func (t *getEventsTest) iSendGetEventsRequestWithPeriodAndStartTime(arg1, arg2 string) error {

	// Create cancel context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get events from gRPC server
	parsed, err := time.Parse(time.RFC3339, arg2)
	if err != nil {
		return errors.Wrap(err, "couldn't parse incoming time")
	}
	searchTime, err := ptypes.TimestampProto(parsed)
	if err != nil {
		return errors.Wrap(err, "error converting timestamp")
	}
	t.searchResponse, err = grpcClient.GetEvents(ctx, &api.GetEventsRequest{
		Period:    ToTimePeriod(arg1),
		StartTime: searchTime,
	})
	if err != nil {
		return errors.Wrap(err, "failed sending GetEvents request")
	}

	return nil
}

func (t *getEventsTest) searchShouldReturnEvent(arg1 int) error {

	// Get events from response
	events := t.searchResponse.GetEvents()
	if events == nil || len(events) != arg1 {
		return errors.New("response returned not enough data")
	}

	return nil
}
