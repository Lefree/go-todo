package service

import (
	"database/sql"
	"lefree/go-todo/internal/models"
	"log/slog"

	"github.com/google/uuid"
)

type TaskService struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewTaskService(log *slog.Logger, db *sql.DB) *TaskService {
	return &TaskService{
		logger: log,
		db:     db,
	}
}

func (s *TaskService) AddTask(dto *models.Task) (string, error) {
	id := uuid.New()
	sqlStatement := "insert into task (id, name, description, is_finished) values($1, $2, $3, $4)"
	_, err := s.db.Query(sqlStatement, id, dto.Name, dto.Description, dto.IsFinished)
	if err != nil {
		s.logger.Error("failed to save data into database, %s", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s *TaskService) GetAll() ([]*models.Task, error) {
	sqlStatement := "select id, name, description, is_finished from task order by id"
	rows, err := s.db.Query(sqlStatement)
	result := make([]*models.Task, 0)
	if err != nil {
		s.logger.Error("failed fetch from database %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &models.Task{}
		scErr := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished)
		if scErr != nil {
			s.logger.Error("failed parse database result to dto %s", scErr.Error())
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (s *TaskService) GetById(id string) (*models.Task, error) {
	sqlStatement := "select id, name, description, is_finished from task where id = $1"
	rows, err := s.db.Query(sqlStatement, id)
	if err != nil {
		s.logger.Error("failed fetch record with id = %s from database %s", id, err.Error())
		return nil, err
	}
	defer rows.Close()
	task := &models.Task{}
	for rows.Next() {
		scErr := rows.Scan(&task.ID, &task.Name, &task.Description, &task.IsFinished)
		if scErr != nil {
			s.logger.Error("failed parse database result to dto %s", scErr.Error())
			return nil, scErr
		}
	}
	return task, nil
}

func (s *TaskService) Delete(id string) error {
	sqlStatement := "delete from task where id = $1"
	_, err := s.db.Query(sqlStatement, id)
	if err != nil {
		s.logger.Error("failed to delete task with id = %s from database. %s", id, err.Error())
		return err
	}
	return nil
}
