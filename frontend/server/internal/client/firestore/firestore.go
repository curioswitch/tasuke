package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// Client is an interface wrapping the Firebase Admin SDK for our usage of firestore.
// It does not use a fluent API to make mocking easier.
type Client[T any] interface {
	// GetDocument retrieves a document at the given path.
	GetDocument(ctx context.Context, path string) (*T, error)

	// SetDocument creates a document in the given collection at the given path with the given data.
	SetDocument(ctx context.Context, path string, data *T) error
}

// NewClient returns a new Client wrapping the given firestore.Client for querying the given
// collection.
func NewClient[T any](store *firestore.Client, collection string) Client[T] {
	return &client[T]{
		store:      store,
		collection: collection,
	}
}

type client[T any] struct {
	store      *firestore.Client
	collection string
}

// GetDocument implements Client.
func (c *client[T]) GetDocument(ctx context.Context, path string) (*T, error) {
	doc, err := c.store.Collection(c.collection).Doc(path).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestore: failed to get document: %w", err)
	}

	var res T
	if err := doc.DataTo(&res); err != nil {
		return nil, fmt.Errorf("firestore: failed to parse document: %w", err)
	}

	return &res, nil
}

// SetDocument implements Client.
func (c *client[T]) SetDocument(ctx context.Context, path string, data *T) error {
	if _, err := c.store.Collection(c.collection).Doc(path).Set(ctx, data); err != nil {
		return fmt.Errorf("firestore: failed to create document: %w", err)
	}

	return nil
}
