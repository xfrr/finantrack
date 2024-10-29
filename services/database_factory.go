package services

import (
	"context"
	"fmt"
)

type DatabaseEngineType string

const (
	InMemoryDatabaseEngine DatabaseEngineType = "inmemory"
	MongoDatabaseEngine    DatabaseEngineType = "mongodb"
	ImmuDBDatabaseEngine   DatabaseEngineType = "immudb"
)

// RepositoryNotFoundError is the error returned when a repository is not found.
type RepositoryNotFoundError struct {
	engine DatabaseEngineType
}

// Error returns the error message.
func (r *RepositoryNotFoundError) Error() string {
	return fmt.Sprintf("repository not found for database engine %s", r.engine)
}

// ErrRepositoryNotFound returns a new instance of RepositoryNotFoundError.
func ErrRepositoryNotFound(engine DatabaseEngineType) error {
	return &RepositoryNotFoundError{engine: engine}
}

// RepositoryFactory is the interface that wraps the basic methods for creating a repository.
type RepositoryFactory[R any] interface {
	// RegisterRepository registers a repository for a specific database engine.
	RegisterRepository(databaseEngine DatabaseEngineType, fn RepositoryFactoryFunc[R]) error

	// CreateRepository creates a repository for a specific database engine.
	CreateRepository(ctx context.Context, databaseEngine DatabaseEngineType) (R, func() error, error)
}

// repositoryFactory is the default implementation of the RepositoryFactory interface.
type repositoryFactory[R any] struct {
	repositories map[DatabaseEngineType]RepositoryFactoryFunc[R]
}

type RepositoryFactoryFunc[R any] func(ctx context.Context) (R, func() error, error)

// NewRepositoryFactory creates a new instance of RepositoryFactory.
func NewRepositoryFactory[R any]() RepositoryFactory[R] {
	return &repositoryFactory[R]{
		repositories: make(map[DatabaseEngineType]RepositoryFactoryFunc[R]),
	}
}

// RegisterRepository registers a repository for a specific database engine.
func (r *repositoryFactory[R]) RegisterRepository(databaseEngine DatabaseEngineType, fn RepositoryFactoryFunc[R]) error {
	if _, ok := r.repositories[databaseEngine]; ok {
		return fmt.Errorf("repository already registered for database engine %s", databaseEngine)
	}

	r.repositories[databaseEngine] = fn

	return nil
}

// CreateRepository creates a repository for a specific database engine.
func (r *repositoryFactory[R]) CreateRepository(
	ctx context.Context,
	databaseEngine DatabaseEngineType,
) (R, func() error, error) {
	var (
		repository R
		closer     func() error
	)

	fn, ok := r.repositories[databaseEngine]
	if !ok {
		return repository, closer, ErrRepositoryNotFound(databaseEngine)
	}

	repository, closer, err := fn(ctx)
	if err != nil {
		return repository, closer, err
	}

	return repository, closer, nil
}
