package xmongo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/finantrack/internal/shared/xos"
	"github.com/xfrr/go-cqrsify/event"

	. "github.com/xfrr/finantrack/internal/shared/xmongo"
)

const databaseName = "test-finantrack-event-store"

type mockEventPayload struct {
	Key string `json:"key"`
}

func TestEventStore_Save(t *testing.T) {
	var (
		uri                = xos.GetEnvWithDefault("FINANTRACK_TEST_MONGO_URI", "mongodb://localhost:27017")
		mockUUID           = uuid.New().String()
		mockEventTimestamp = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := NewClient(ctx, uri, databaseName)
	if err != nil {
		t.Fatal(err)
	}

	registry := xevent.NewPayloadRegistry()
	registry.Register("event-type", func() interface{} {
		return &mockEventPayload{}
	})

	sut, err := NewMongoEventStore(ctx, client, registry)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("save event", func(t *testing.T) {
		// Clean up the database after the test.
		defer cleanUp(ctx, t, client)

		inputEvent := &EventMock{
			IDFunc: func() any {
				return uuid.MustParse(mockUUID)
			},
			ReasonFunc: func() string {
				return "event-type"
			},
			PayloadFunc: func() any {
				return mockEventPayload{Key: "value"}
			},
			TimeFunc: func() time.Time {
				return mockEventTimestamp
			},
			AggregateFunc: func() *event.AggregateRef[any] {
				return &event.AggregateRef[any]{
					ID:      uuid.MustParse(mockUUID),
					Name:    "aggregate-type",
					Version: 1,
				}
			},
		}

		err = sut.Save(ctx, inputEvent)
		if err != nil {
			t.Fatal(err)
		}

		var events []Event

		// Check if the event was saved in the database.
		events, err = sut.Get(ctx, WithEventIDCriteria(mockUUID)())
		require.NoError(t, err)
		require.Len(t, events, 1)

		resultEvent, ok := event.Cast[uuid.UUID, any](events[0])
		require.True(t, ok)

		assert.Equal(t, mockUUID, resultEvent.ID().String())
		assert.Equal(t, "event-type", inputEvent.Reason())
		assert.Equal(t, "aggregate-type", inputEvent.Aggregate().Name)
		assert.Equal(t, 1, inputEvent.Aggregate().Version)
		assert.Equal(t, "value", inputEvent.Payload().(mockEventPayload).Key)
		assert.Equal(t, mockEventTimestamp, inputEvent.Time())
	})

	t.Run("save event with same id should return error", func(t *testing.T) {
		// Clean up the database after the test.
		defer cleanUp(ctx, t, client)

		event := &EventMock{
			IDFunc: func() any {
				return uuid.MustParse(mockUUID)
			},
			ReasonFunc: func() string {
				return "event-type"
			},
			PayloadFunc: func() any {
				return mockEventPayload{Key: "value"}
			},
			TimeFunc: func() time.Time {
				return mockEventTimestamp
			},
			AggregateFunc: func() *event.AggregateRef[any] {
				return &event.AggregateRef[any]{
					ID:      uuid.MustParse(mockUUID),
					Name:    "aggregate-type",
					Version: 1,
				}
			},
		}

		err = sut.Save(ctx, event)
		require.NoError(t, err)

		err = sut.Save(ctx, event)
		require.Error(t, err)
	})
}

func TestEventStore_Get(t *testing.T) {
	var (
		uri = xos.GetEnvWithDefault("FINANTRACK_TEST_MONGO_URI", "mongodb://localhost:27017")
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := NewClient(ctx, uri, databaseName)
	if err != nil {
		t.Fatal(err)
	}

	cleanUp(ctx, t, client)
	registry := xevent.NewPayloadRegistry()
	sut, err := NewMongoEventStore(ctx, client, registry)
	if err != nil {
		t.Fatal(err)
	}

	mockEvents := generateMockEvents(ctx, t, sut)
	defer cleanUp(ctx, t, client)

	var specs = []struct {
		name     string
		criteria Criteria
		expected []*EventMock
	}{
		{
			name:     "get all events by aggregate id",
			criteria: WithAggregateIDCriteria("00000000-0000-0000-0000-000000000001")(),
			expected: []*EventMock{mockEvents[1]},
		},
		{
			name:     "get all events by aggregate type",
			criteria: WithAggregateTypeCriteria("aggregate-type-2")(),
			expected: []*EventMock{mockEvents[2]},
		},
		{
			name:     "get all events by aggregate version",
			criteria: WithAggregateVersionCriteria(3)(),
			expected: []*EventMock{mockEvents[2]},
		},
		{
			name:     "get all events by event id",
			criteria: WithEventIDCriteria("00000000-0000-0000-0000-000000000004")(),
			expected: []*EventMock{mockEvents[4]},
		},
		{
			name:     "get all events by event type",
			criteria: WithEventTypeCriteria("event-type-4")(),
			expected: []*EventMock{mockEvents[4]},
		},
		{
			name: "get all events by aggregate id and type",
			criteria: And(
				WithAggregateIDCriteria("00000000-0000-0000-0000-000000000003")(),
				WithEventTypeCriteria("event-type-3")(),
			)(),
			expected: []*EventMock{mockEvents[3]},
		},
		{
			name: "get all events by aggregate id, type and aggregate version",
			criteria: And(
				WithAggregateIDCriteria("00000000-0000-0000-0000-000000000001")(),
				WithEventTypeCriteria("event-type-1")(),
				WithAggregateVersionCriteria(2)(),
			)(),
			expected: []*EventMock{mockEvents[1]},
		},
		{
			name: "get all events by aggregate id or type",
			criteria: Or(
				WithAggregateIDCriteria("00000000-0000-0000-0000-000000000001")(),
				WithEventTypeCriteria("event-type-2")(),
			)(),
			expected: []*EventMock{
				mockEvents[1],
				mockEvents[2],
			},
		},
		{
			name: "get all events by aggregate id and type or aggregate version",
			criteria: And(
				WithAggregateIDCriteria("00000000-0000-0000-0000-000000000002")(),
				Or(
					WithEventTypeCriteria("event-type-2")(),
					WithAggregateVersionCriteria(2)(),
				)(),
			)(),
			expected: []*EventMock{mockEvents[2]},
		},
		{
			name: "get all events by aggregate id or type and not version",
			criteria: And(
				Or(
					WithAggregateIDCriteria("00000000-0000-0000-0000-000000000004")(),
					WithEventTypeCriteria("event-type-4")(),
				)(),
				Not(WithAggregateIDCriteria("00000000-0000-0000-0000-000000000004")())(),
			)(),
			expected: []*EventMock{},
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			var events []Event
			events, err = sut.Get(ctx, spec.criteria)
			require.NoError(t, err)

			assert.Equal(t, len(spec.expected), len(events))
			for i, e := range events {
				ev, ok := event.Cast[uuid.UUID, *mockEventPayload](e)
				require.True(t, ok)

				expected, ok := event.Cast[uuid.UUID, *mockEventPayload](spec.expected[i])
				require.True(t, ok)

				assert.Equal(t, expected.ID().String(), ev.ID().String())
				assert.Equal(t, expected.Reason(), ev.Reason())
				assert.Equal(t, expected.Aggregate().ID, ev.Aggregate().ID)
				assert.Equal(t, expected.Aggregate().Name, ev.Aggregate().Name)
				assert.Equal(t, expected.Aggregate().Version, ev.Aggregate().Version)
				assert.Equal(t, expected.Payload(), ev.Payload())
				assert.Equal(t, expected.Time(), ev.Time())
			}
		})
	}
}

func cleanUp(ctx context.Context, t *testing.T, client *Client) {
	err := client.Drop(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

// generateMockEvents generates a list of events to be used in the tests.
// Must create an event for each property of the EventMock struct to test the Get criteria.
func generateMockEvents(ctx context.Context, t *testing.T, sut *MongoEventStore) []*EventMock {
	var events []*EventMock
	for i := range [5]int{} {
		event := &EventMock{
			IDFunc: func() any {
				return uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-00000000000%d", i))
			},
			ReasonFunc: func() string {
				return fmt.Sprintf("event-type-%d", i)
			},
			PayloadFunc: func() any {
				return &mockEventPayload{Key: fmt.Sprintf("value-%d", i)}
			},
			TimeFunc: func() time.Time {
				return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			},
			AggregateFunc: func() *event.AggregateRef[any] {
				return &event.AggregateRef[any]{
					ID:      uuid.MustParse(fmt.Sprintf("00000000-0000-0000-0000-00000000000%d", i)),
					Name:    fmt.Sprintf("aggregate-type-%d", i),
					Version: i + 1,
				}
			},
		}

		// Register the payload factory for the event type.
		sut.Registry().Register(fmt.Sprintf("event-type-%d", i+1), func() interface{} {
			return &mockEventPayload{}
		})

		err := sut.Save(ctx, event)
		if err != nil {
			t.Fatal(err)
		}

		events = append(events, event)
	}

	return events
}
