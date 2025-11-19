package database

import "context"

type Write[T any] interface {
	Create(ctx context.Context, data *T) error
	Update(ctx context.Context, data *T) error
}

type ReadOne[T any] interface {
	GetRepo(ctx context.Context, id string) (T, error)
}

type Delete[T any] interface {
	Delete(ctx context.Context, id string) error
}
