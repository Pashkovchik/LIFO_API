package service

import (
	"context"
	"fmt"
	"lifo-rest-api/internal/domain"
	"lifo-rest-api/internal/domain/constant"
	"lifo-rest-api/internal/repository"
)

// stackService -.
type stackService struct {
	stackRepo repository.Stack
}

// NewStackService -.
func NewStackService(stackRepo repository.Stack) *stackService {
	return &stackService{
		stackRepo: stackRepo,
	}
}

func (s *stackService) Create(ctx context.Context, name string) (stackID uint64, err error) {
	stackID, err = s.stackRepo.Create(ctx, name)

	return stackID, err
}

func (s *stackService) GetListOfStacks(ctx context.Context) (stacks []domain.Stack, err error) {
	stacks, err = s.stackRepo.GetListOfStacks(ctx)

	return stacks, err
}

func (s *stackService) Delete(ctx context.Context, stackID uint64) (err error) {
	err = s.stackRepo.Delete(ctx, stackID)

	return err
}

func (s *stackService) Push(ctx context.Context, info domain.StackData) (err error) {
	err = s.stackRepo.Push(ctx, info)

	return err
}

func (s *stackService) Get(ctx context.Context, typeOfOperation string, stackID uint64) (info domain.StackData, err error) {
	if typeOfOperation == constant.TypePop {
		info, err = s.stackRepo.Pop(ctx, stackID)
	} else if typeOfOperation == constant.TypePeek {
		info, err = s.stackRepo.Peek(ctx, stackID)
	} else {
		err = fmt.Errorf("type argument must be one of them: %s, %s", constant.TypePop, constant.TypePeek)
	}

	return info, err
}
