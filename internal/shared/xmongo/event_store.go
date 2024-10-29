package xmongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xfrr/finantrack/internal/shared/xevent"
	"github.com/xfrr/go-cqrsify/aggregate"
	"github.com/xfrr/go-cqrsify/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DefaultCollectionName is the default collection name for events.
const DefaultCollectionName = "events"

// Event represents an event that will be saved in the storage.
type Event = aggregate.Change

// EventStore defines the interface for saving and retrieving events.
type EventStore interface {
	Save(ctx context.Context, events ...Event) error
	Get(ctx context.Context, criteria Criteria) ([]Event, error)
	ExistsByAggregateID(ctx context.Context, aggregateID uuid.UUID) (bool, error)
}

// MongoEventStore is the MongoDB implementation of EventStore.
type MongoEventStore struct {
	client                 *Client
	payloadFactoryRegistry xevent.Registry
}

// NewMongoEventStore creates a new instance of MongoEventStore.
// It also creates the necessary indexes for the events collection.
// The registry is used to resolve the payload type for each event type.
func NewMongoEventStore(
	ctx context.Context,
	client *Client,
	registry xevent.Registry,
) (*MongoEventStore, error) {
	// Create indexes
	mes := &MongoEventStore{
		client:                 client,
		payloadFactoryRegistry: registry,
	}

	// Create indexes
	err := mes.createIndexes(ctx)
	if err != nil {
		return nil, err
	}

	return mes, nil
}

// Save saves the events in the storage.
func (s *MongoEventStore) Save(ctx context.Context, events ...Event) error {
	var dtos []interface{}

	for _, e := range events {
		dto, err := s.eventToDTO(e)
		if err != nil {
			return fmt.Errorf("failed to convert event to DTO: %w", err)
		}

		dtos = append(dtos, dto)
	}

	_, err := s.client.
		Collection(DefaultCollectionName).
		InsertMany(ctx, dtos)
	if err != nil {
		return fmt.Errorf("failed to insert events: %w", err)
	}

	return nil
}

// Get retrieves events from the storage that match the given criteria.
func (s *MongoEventStore) Get(ctx context.Context, criteria Criteria) ([]Event, error) {
	cursor, err := s.client.
		Collection(DefaultCollectionName).
		Find(ctx, criteria.ToBSON())
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}
	defer cursor.Close(ctx)

	var dtos []eventDTO
	if err = cursor.All(ctx, &dtos); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	var events []Event
	for _, dto := range dtos {
		var e Event
		e, err = s.createEventFromDTO(dto)
		if err != nil {
			return nil, fmt.Errorf("failed to convert DTO to event: %w", err)
		}
		events = append(events, e)
	}

	return events, nil
}

// ExistsByAggregateID checks if an event exists for the given aggregate ID.
func (s *MongoEventStore) ExistsByAggregateID(ctx context.Context, aggregateID uuid.UUID) (bool, error) {
	count, err := s.client.
		Collection(DefaultCollectionName).
		CountDocuments(ctx, bson.M{"aggregate_id": aggregateID.String()})
	if err != nil {
		return false, fmt.Errorf("failed to count documents: %w", err)
	}

	return count > 0, nil
}

// Registry returns the payload factory registry.
func (s *MongoEventStore) Registry() xevent.Registry {
	return s.payloadFactoryRegistry
}

// eventToDTO converts an Event into an eventDTO for storage.
func (s *MongoEventStore) eventToDTO(e Event) (*eventDTO, error) {
	ev, ok := event.Cast[uuid.UUID, any](e)
	if !ok {
		return nil, fmt.Errorf("event must have UUID IDs, got %v", e.ID())
	}

	if e.Aggregate() == nil {
		return nil, errors.New("event must have an aggregate reference")
	}

	aggregateID, ok := e.Aggregate().ID.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("aggregate ID must be a UUID, got %v", e.Aggregate().ID)
	}

	// Handle the event ID.
	eventID := ev.ID()
	if eventID == uuid.Nil {
		eventID = uuid.New()
	}

	// Convert UUIDs to strings.
	eventIDStr := eventID.String()
	aggregateIDStr := aggregateID.String()

	// Assign the payload.
	payload := ev.Payload()

	return &eventDTO{
		ID:               eventIDStr,
		Type:             e.Reason(),
		AggregateID:      aggregateIDStr,
		AggregateType:    e.Aggregate().Name,
		AggregateVersion: e.Aggregate().Version,
		Data:             payload,
		Timestamp:        e.Time(),
		Version:          1,
	}, nil
}

// createEventFromDTO converts an eventDTO back to an Event.
func (s *MongoEventStore) createEventFromDTO(dto eventDTO) (Event, error) {
	eventID, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event ID: %w", err)
	}

	aggregateID, err := uuid.Parse(dto.AggregateID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aggregate ID: %w", err)
	}

	// Get the payload factory function based on the event type.
	payloadFactory, err := s.payloadFactoryRegistry.GetFactory(dto.Type)
	if err != nil {
		return nil, err
	}

	// Create a new instance of the payload type.
	payload := payloadFactory()

	// Unmarshal the data into the payload instance.
	dataBytes, err := bson.Marshal(dto.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	if err = bson.Unmarshal(dataBytes, payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Reconstruct the event.
	return event.New[any](
		eventID,
		dto.Type,
		payload,
		event.WithAggregate(
			aggregateID,
			dto.AggregateType,
			dto.AggregateVersion,
		),
		event.WithTime(dto.Timestamp),
	), nil
}

// createIndexes creates the necessary indexes for the events collection.
func (s *MongoEventStore) createIndexes(ctx context.Context) error {
	// Create an index for the aggregate ID.
	err := s.createIndex(ctx, bson.D{{Key: "aggregate_id", Value: 1}}, false)
	if err != nil {
		return fmt.Errorf("failed to create index for aggregate ID: %w", err)
	}

	// Create a compound index for the aggregate ID and version.
	err = s.createIndex(ctx, bson.D{
		{Key: "aggregate_id", Value: 1},
		{Key: "aggregate_version", Value: 1},
	}, true)
	if err != nil {
		return fmt.Errorf("failed to create index for aggregate ID and version: %w", err)
	}

	// Create an index for the event type.
	err = s.createIndex(ctx, bson.D{{Key: "type", Value: 1}}, false)
	if err != nil {
		return fmt.Errorf("failed to create index for event type: %w", err)
	}

	// Create an index for the timestamp.
	err = s.createIndex(ctx, bson.D{{Key: "timestamp", Value: 1}}, false)
	if err != nil {
		return fmt.Errorf("failed to create index for timestamp: %w", err)
	}

	return nil
}

// createIndex creates an index in the events collection.
func (s *MongoEventStore) createIndex(ctx context.Context, keys bson.D, unique bool) error {
	index := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(unique),
	}

	_, err := s.client.
		Collection(DefaultCollectionName).
		Indexes().
		CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

// eventDTO represents the structure of an event stored in MongoDB.
type eventDTO struct {
	ID               string      `bson:"_id"`
	Type             string      `bson:"type"`
	AggregateID      string      `bson:"aggregate_id"`
	AggregateType    string      `bson:"aggregate_type"`
	AggregateVersion int         `bson:"aggregate_version"`
	Data             interface{} `bson:"data,omitempty"`
	Metadata         interface{} `bson:"metadata,omitempty"`
	Timestamp        time.Time   `bson:"timestamp"`
	Version          int         `bson:"version"`
}
