package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Client is an interface wrapping the Firebase Admin SDK for our usage of firestore.
// It does not use a fluent API to make mocking easier.
type Client[T any] interface {
	// Query retrieves all the documents matching the given EntityFilter. The return value
	// can be used as a range-over function.
	Query(ctx context.Context, ef firestore.EntityFilter) func(yield func(*T, error) bool)
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

func (c *client[T]) Query(ctx context.Context, ef firestore.EntityFilter) func(yield func(*T, error) bool) {
	return func(yield func(*T, error) bool) {
		iter := c.store.Collection(c.collection).WhereEntity(ef).Documents(ctx)
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}

			if err != nil {
				if !yield(nil, fmt.Errorf("firestore: failed to get document: %w", err)) {
					break
				}
				continue
			}

			var res T
			if err := doc.DataTo(&res); err != nil {
				if !yield(nil, fmt.Errorf("firestore: failed to parse document: %w", err)) {
					break
				}
				continue
			}

			if !yield(&res, nil) {
				break
			}
		}
	}
}
