package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"lifo-rest-api/internal/domain"
	"lifo-rest-api/internal/domain/constant"
	"lifo-rest-api/pkg/database/postgres"

	"go.uber.org/zap"
)

type stackRepo struct {
	stackTable     string
	stackDataTable string
	db             *postgres.Postgres
}

func NewStackRepo(db *postgres.Postgres) *stackRepo {
	return &stackRepo{
		stackTable:     constant.PGTableStack,
		stackDataTable: constant.PGTableStackData,
		db:             db,
	}
}

func (r *stackRepo) Create(ctx context.Context, name string) (stackID uint64, err error) {
	qs, args, err := r.db.Builder.Insert(r.stackTable).Columns("name").
		Values(name).Suffix("RETURNING id").ToSql()
	if err != nil {
		return stackID, fmt.Errorf("stackRepo - CreateNewStack - r.InsertBuilder: %w", err)
	}

	zap.S().Debugf("CreateNewStack SQL: %s | args: %v", qs, args)

	row := r.db.Pool.QueryRow(ctx, qs, args...)

	if err = row.Scan(&stackID); err != nil {
		return stackID, fmt.Errorf("stackRepo - CreateNewStack - row.Exec: %w", err)
	}

	return
}

func (r *stackRepo) GetListOfStacks(ctx context.Context) (stacks []domain.Stack, err error) {
	stacks, err = r.getStacks(ctx)
	if err != nil {
		return stacks, err
	}

	for i := range stacks {
		stacks[i].StackData, err = r.getStackData(ctx, stacks[i].ID)
		if err != nil {
			return stacks, err
		}
	}

	return stacks, nil
}

func (r *stackRepo) getStacks(ctx context.Context) (stacks []domain.Stack, err error) {
	qs, args, err := r.db.Builder.Select("id", "name", "created_date").
		From(r.stackTable).ToSql()
	if err != nil {
		return stacks, fmt.Errorf("stackRepo - getStacks - r.SelectBuilder: %w", err)
	}

	zap.S().Debugf("getStacks SQL: %s | args: %v", qs, args)

	stackObject := new(domain.Stack)

	_, err = r.db.Pool.QueryFunc(ctx, qs, args, stackObject.Fields(), func(qfr pgx.QueryFuncRow) error {
		stack := *stackObject
		stacks = append(stacks, stack)

		return nil
	})
	if err != nil {
		return stacks, fmt.Errorf("stackRepo - getStacks - r.SelectBuilder: %w", err)
	}

	return stacks, nil
}

func (r *stackRepo) getStackData(ctx context.Context, stackID uint64) (stackDatas []domain.StackData, err error) {
	qs, args, err := r.db.Builder.Select("id", "stack_id", "info", "created_date").
		From(r.stackDataTable).Where(sq.Eq{"stack_id": stackID}).
		OrderBy("id desc").ToSql()
	if err != nil {
		return stackDatas, fmt.Errorf("stackRepo - getStackData - r.SelectBuilder: %w", err)
	}

	zap.S().Debugf("getStackData SQL: %s | args: %v", qs, args)

	stackDataObject := new(domain.StackData)

	_, err = r.db.Pool.QueryFunc(ctx, qs, args, stackDataObject.Fields(), func(qfr pgx.QueryFuncRow) error {
		stackData := *stackDataObject
		stackDatas = append(stackDatas, stackData)

		return nil
	})
	if err != nil {
		return stackDatas, fmt.Errorf("stackRepo - getStackData - r.SelectBuilder: %w", err)
	}

	return stackDatas, nil
}

func (r *stackRepo) Delete(ctx context.Context, stackID uint64) (err error) {
	qs, args, err := r.db.Builder.Delete(r.stackTable).Where(sq.Eq{"id": stackID}).ToSql()
	if err != nil {
		return fmt.Errorf("stackRepo - deleteStack - r.SelectBuilder: %w", err)
	}

	zap.S().Debugf("deleteStack SQL: %s | args: %v", qs, args)

	_, err = r.db.Pool.Exec(ctx, qs, args...)
	if err != nil {
		return fmt.Errorf("stackRepo - deleteStack - row.Scan: %w", err)
	}

	return
}

func (r *stackRepo) Push(ctx context.Context, info domain.StackData) (err error) {
	qs, args, err := r.db.Builder.Insert(r.stackDataTable).Columns("stack_id", "info").
		Values(info.StackID, info.Info).ToSql()
	if err != nil {
		return fmt.Errorf("stackRepo - PushInfoToStack - r.InsertBuilder: %w", err)
	}

	zap.S().Debugf("PushInfoToStack SQL: %s | args: %v", qs, args)

	_, err = r.db.Pool.Exec(ctx, qs, args...)
	if err != nil {
		return fmt.Errorf("stackRepo - PushInfoToStack - row.Exec: %w", err)
	}

	return
}

func (r *stackRepo) Pop(ctx context.Context, stackID uint64) (info domain.StackData, err error) {
	info, err = r.getLastInfo(ctx, stackID)
	if err != nil {
		return info, err
	}

	err = r.deleteLastInfo(ctx, info.ID)
	if err != nil {
		return info, err
	}

	return info, nil
}

func (r *stackRepo) Peek(ctx context.Context, stackID uint64) (info domain.StackData, err error) {
	info, err = r.getLastInfo(ctx, stackID)
	if err != nil {
		return info, err
	}

	return info, nil
}

func (r *stackRepo) getLastInfo(ctx context.Context, stackID uint64) (info domain.StackData, err error) {
	qs, args, err := r.db.Builder.Select("id", "stack_id", "info", "created_date").
		From(r.stackDataTable).Where(sq.Eq{"stack_id": stackID}).
		OrderBy("id desc").Limit(1).ToSql()
	if err != nil {
		return info, fmt.Errorf("stackRepo - getLastInfo - r.SelectBuilder: %w", err)
	}

	zap.S().Debugf("getLastInfo SQL: %s | args: %v", qs, args)

	row := r.db.Pool.QueryRow(ctx, qs, args...)

	err = row.Scan(&info.ID, &info.StackID, &info.Info, &info.CreatedDate)
	if err != nil {
		if err.Error() == constant.SqlNoRows {
			zap.S().Errorf("no rows in this stack: %v", err)
		}

		return info, fmt.Errorf("stackRepo - getLastInfo - row.Scan: %w", err)
	}

	return info, err
}

func (r *stackRepo) deleteLastInfo(ctx context.Context, infoID uint64) (err error) {
	qs, args, err := r.db.Builder.Delete(r.stackDataTable).Where(sq.Eq{"id": infoID}).ToSql()
	if err != nil {
		return fmt.Errorf("stackRepo - deleteLastInfo - r.SelectBuilder: %w", err)
	}

	zap.S().Debugf("deleteLastInfo SQL: %s | args: %v", qs, args)

	_, err = r.db.Pool.Exec(ctx, qs, args...)
	if err != nil {
		return fmt.Errorf("stackRepo - deleteLastInfo - row.Scan: %w", err)
	}

	return err
}
