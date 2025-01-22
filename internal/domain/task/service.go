package task

import (
	"fmt"
	"time"
)

type Repository interface {
	Create(*Task) error
	GetTasks(*ListQuery) ([]Task, error)
	GetTaskByID(int64) (*Task, error)
	UpdateTask(*Task) error
	DeleteTask(int64) error
	UpdateTaskDate(int64, string) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("заголовок задачи не может быть пустым")
	}

	now := time.Now()
	today := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	if task.Date == "today" {
		task.Date = today
	}

	if err := ValidateDate(task.Date); err != nil {
		return fmt.Errorf("некорректная дата: %w", err)
	}

	if task.Date <= today && task.Repeat != "" {
		if task.Repeat == "d 1" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
			}
			task.Date = nextDate
		}
	} else if task.Date < today {
		task.Date = today
	}

	return s.repository.Create(task)
}

func (s *Service) GetNearestTasks(dateFilter, commentFilter string) ([]Task, error) {
	query := &ListQuery{
		Date:    dateFilter,
		Comment: commentFilter,
		Limit:   50,
	}

	tasks, err := s.repository.GetTasks(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка задач: %w", err)
	}

	if tasks == nil {
		return []Task{}, nil
	}

	return tasks, nil
}

func (s *Service) GetTask(id int64) (*Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("некорректный идентификатор задачи")
	}
	return s.repository.GetTaskByID(id)
}

func (s *Service) UpdateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("заголовок задачи не может быть пустым")
	}

	now := time.Now()
	today := now.Format(DateFormat)

	if err := ValidateDate(task.Date); err != nil {
		return fmt.Errorf("некорректная дата: %w", err)
	}

	if task.Date <= today && task.Repeat != "" {
		if task.Repeat == "d 1" {
			task.Date = today
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
			}
			task.Date = nextDate
		}
	} else if task.Date < today {
		task.Date = today
	}

	return s.repository.UpdateTask(task)
}
func (s *Service) MarkTaskDone(id int64, now time.Time) error {
	task, err := s.repository.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.repository.DeleteTask(id)
	}

	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
	}

	return s.repository.UpdateTaskDate(id, nextDate)
}

func (s *Service) DeleteTask(id int64) error {
	// Проверяем существование задачи перед удалением
	if _, err := s.repository.GetTaskByID(id); err != nil {
		return err
	}
	return s.repository.DeleteTask(id)
}
