package assetimmudb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	assetdomain "github.com/xfrr/finantrack/internal/contexts/assets/domain"
)

// Asset represents the asset entity.
type eventDTO struct {
	ID            string    `db:"id"`
	AggregateID   string    `db:"aggregate_id"`
	AggregateName string    `db:"aggregate_name"`
	Timestamp     time.Time `db:"created_at"`
	Payload       []byte    `db:"payload"`
}

// Repository implements the Repository interface using ImmuDB.
type Repository struct {
	db *sql.DB
}

// NewImmuRepository creates a new ImmuRepository with the given ImmuDB client.
func NewImmuRepository(db *sql.DB) (*Repository, error) {
	repo := &Repository{
		db: db,
	}

	return repo, nil
}

// Save saves the given asset to ImmuDB.
func (r *Repository) Save(ctx context.Context, asset *assetdomain.Asset) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	events := asset.AggregateChanges()
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var eventDTOs []eventDTO
	for _, event := range events {
		var (
			payload []byte
			err     error
		)

		if event.Payload() != nil {
			payload, err = json.Marshal(event.Payload())
			if err != nil {
				return err
			}
		}

		eventID, ok := event.ID().(uuid.UUID)
		if !ok {
			return fmt.Errorf("event ID is not a UUID")
		}

		aggID, ok := event.ID().(uuid.UUID)
		if !ok {
			return fmt.Errorf("agggregate ID is not a UUID")
		}

		eventDTOs = append(eventDTOs, eventDTO{
			ID:            eventID.String(),
			AggregateID:   aggID.String(),
			AggregateName: event.Aggregate().Name,
			Timestamp:     event.Time(),
			Payload:       payload,
		})

		// Insert the event into the database
		dbInsertSqlQuery := fmt.Sprintf(`
			INSERT INTO %s (id, aggregate_id, aggregate_name, created_at, payload)
			VALUES (?, ?, ?, ?, ?);`,
			"events",
		)

		_, err = tx.ExecContext(ctx, dbInsertSqlQuery, eventDTOs[len(eventDTOs)-1].ID, eventDTOs[len(eventDTOs)-1].AggregateID, eventDTOs[len(eventDTOs)-1].AggregateName, eventDTOs[len(eventDTOs)-1].Timestamp, string(eventDTOs[len(eventDTOs)-1].Payload))
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an asset by its ID from ImmuDB.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*assetdomain.Asset, error) {
	// TODO: Implement GetByID method
	return nil, errors.New("not implemented")
}

// GetAll retrieves all assets from ImmuDB.
func (r *Repository) GetAll(ctx context.Context) ([]*assetdomain.Asset, error) {
	// TODO: Implement GetAll method
	return nil, errors.New("not implemented")
}

// Exists checks whether an asset exists in ImmuDB by its ID.
func (r *Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	stmt := `SELECT id FROM events WHERE id = ?`
	args := []any{id.String()}

	rows, err := r.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		switch {
		case sql.ErrNoRows == err:
			return false, nil
		default:
			return false, err
		}
	}

	defer rows.Close()
	return rows.Next(), nil
}
