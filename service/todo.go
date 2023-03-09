package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read              = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
		readWithSize      = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithIDAndSize = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos := []*model.TODO{}
	if prevID != 0 && size != 0 {
		stmt, err := s.db.PrepareContext(ctx, readWithIDAndSize)
		if err != nil {
			return nil, err
		}

		rows, err := stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var todo model.TODO
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	} else if size != 0 {
		stmt, err := s.db.PrepareContext(ctx, readWithSize)
		if err != nil {
			return nil, err
		}

		rows, err := stmt.QueryContext(ctx, size)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var todo model.TODO
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	} else {
		stmt, err := s.db.PrepareContext(ctx, read)
		if err != nil {
			return nil, err
		}

		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var todo model.TODO
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, &model.ErrNotFound{
			When: time.Now(),
			What: "error not found",
		}
	}

	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	query := fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	var args []interface{}

	for _, id := range ids {
		args = append(args, id)
	}

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return &model.ErrNotFound{
			When: time.Now(),
			What: "error not found",
		}
	}

	return nil
}
