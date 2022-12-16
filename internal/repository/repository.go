// Package application implements application db processing.
package repository

import (
	"context"
	"lifo-rest-api/internal/domain"
	"lifo-rest-api/pkg/database/postgres"
)

type Stack interface {
	Create(ctx context.Context, name string) (stackID uint64, err error)
	GetListOfStacks(ctx context.Context) (stacks []domain.Stack, err error)
	Delete(ctx context.Context, stackID uint64) (err error)
	Push(ctx context.Context, info domain.StackData) (err error)
	Pop(ctx context.Context, stackID uint64) (info domain.StackData, err error)
	Peek(ctx context.Context, stackID uint64) (info domain.StackData, err error)
}

type Repositories struct {
	Stack Stack
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	stackRepo := NewStackRepo(db)

	return &Repositories{
		Stack: stackRepo,
	}
}
