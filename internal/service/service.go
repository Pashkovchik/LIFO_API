// Package service implements application business logic. Each logic group in own file.
package service

import (
	"context"
	"lifo-rest-api/internal/domain"
	"lifo-rest-api/internal/repository"
)

type Stack interface {
	Create(ctx context.Context, name string) (stackID uint64, err error)
	GetListOfStacks(ctx context.Context) (stacks []domain.Stack, err error)
	Delete(ctx context.Context, stackID uint64) (err error)
	Push(ctx context.Context, info domain.StackData) (err error)
	Get(ctx context.Context, typeOfOperation string, stackID uint64) (info domain.StackData, err error)
}

type Services struct {
	Stack Stack
}

type Dependencies struct {
	Repos *repository.Repositories
}

func NewServices(deps *Dependencies) *Services {
	return &Services{
		Stack: NewStackService(deps.Repos.Stack),
	}
}
