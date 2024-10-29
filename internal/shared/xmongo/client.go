package xmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// MongoConnectDefaultTimeout is the default timeout for connecting to MongoDB.
	MongoConnectDefaultTimeout = 5 * time.Second
)

// Client is a MongoDB client wrapper that embeds a MongoDB database.
type Client struct {
	*mongo.Database
}

// NewClient returns a new MongoDB client with a connection to the specified URI and database.
func NewClient(ctx context.Context, uri string, database string) (*Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	return &Client{
		Database: client.Database(database),
	}, nil
}

// Close closes the MongoDB client connection.
func (c *Client) Close(ctx context.Context) error {
	err := c.Client().Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}

	return nil
}
