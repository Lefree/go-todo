package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lefree/go-todo/internal/models"
	"log/slog"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewTaskRepository(db *sql.DB, log *slog.Logger) *TaskRepository {
	return &TaskRepository{
		db:     db,
		logger: log,
	}
}

func (r *TaskRepository) Save(dto *models.Task) (string, error) {
	id := uuid.New()
	sqlStatement := "insert into task (id, name, description, is_finished) values($1, $2, $3, $4)"
	_, err := r.db.Query(sqlStatement, id, dto.Name, dto.Description, dto.IsFinished)
	if err != nil {
		r.logger.With("error", err.Error()).Error("failed to save data into database")
		return "", err
	}
	return id.String(), nil
}

func (r *TaskRepository) GetAll() ([]*models.Task, error) {
	sqlStatement := "select id, name, description, is_finished from task order by id"
	rows, err := r.db.Query(sqlStatement)
	result := make([]*models.Task, 0)
	if err != nil {
		r.logger.With("error", err.Error()).Error("failed fetch from database")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &models.Task{}
		scErr := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished)
		if scErr != nil {
			r.logger.With("error", scErr.Error()).Error("failed parse database result to dto")
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (r *TaskRepository) GetById(id string) (*models.Task, error) {
	sqlStatement := "select id, name, description, is_finished from task where id = $1"
	rows, err := r.db.Query(sqlStatement, id)
	if err != nil {
		r.logger.With("error", err.Error()).Error(fmt.Sprintf("failed fetch record with id = %s from database", id))
		return nil, err
	}
	defer rows.Close()
	task := &models.Task{}
	for rows.Next() {
		scErr := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished)
		if scErr != nil {
			r.logger.With("error", scErr.Error()).Error("failed parse database result to dto")
			return nil, scErr
		}
	}
	return task, nil
}

func (r *TaskRepository) Patch(dto *models.Task) (*models.Task, error) {
	task, err := r.GetById(dto.ID)
	if err != nil {
		r.logger.With("error", err.Error()).Error(fmt.Sprintf("failed to fetch task with id = %s before update", dto.ID))
		return nil, err
	}
	if task == nil {
		r.logger.Warn(fmt.Sprintf("task with id = %s not found, update skipped", dto.ID))
		return nil, errors.New(fmt.Sprintf("task with id = %s not found", dto.ID))
	}
	if dto.Name != "" {
		task.Name = dto.Name
	}
	if dto.Description != "" {
		task.Description = dto.Description
	}
	task.IsFinished = dto.IsFinished
	_, updErr := r.Save(task)
	if updErr != nil {
		r.logger.With("error", updErr).Error(fmt.Sprintf("error on update task with id = %s", task.ID))
		return nil, updErr
	}
	return task, nil

}

func (r *TaskRepository) Delete(id string) error {
	sqlStatement := "delete from task where id = $1"
	_, err := r.db.Query(sqlStatement, id)
	if err != nil {
		r.logger.With("error", err.Error()).Error(fmt.Sprintf("failed to delete task with id = %s from database.", id))
		return err
	}
	return nil
}
