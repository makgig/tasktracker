package sqlite

import (
	"fmt"
	"tasktracker/internal/domain/task"
)

type Repository struct {
	db *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(t *task.Task) error {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (?, ?, ?, ?)
        RETURNING id`

	row := r.db.QueryRow(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err := row.Scan(&t.ID); err != nil {
		return fmt.Errorf("ошибка при создании задачи: %w", err)
	}

	return nil
}

func (r *Repository) GetTasks(query *task.ListQuery) ([]task.Task, error) {
	var tasks []task.Task
	var queryStr string
	var args []interface{}

	queryStr = "SELECT id, date, title, comment, repeat FROM scheduler"

	hasConditions := false

	if query.Date != "" {
		queryStr += " WHERE date = ?"
		args = append(args, query.Date)
		hasConditions = true
	}

	if query.Comment != "" {
		if hasConditions {
			queryStr += " AND"
		} else {
			queryStr += " WHERE"
		}
		queryStr += " title LIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", query.Comment))
	}

	queryStr += " ORDER BY date ASC LIMIT ?"
	args = append(args, query.Limit)

	err := r.db.Select(&tasks, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выборки задач: %w", err)
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(id int64) (*task.Task, error) {
	var task task.Task
	err := r.db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	return &task, nil
}

func (r *Repository) UpdateTask(t *task.Task) error {
	result, err := r.db.Exec(`
        UPDATE scheduler 
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`,
		t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
func (r *Repository) DeleteTask(id int64) error {
	result, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func (r *Repository) UpdateTaskDate(id int64, newDate string) error {
	result, err := r.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
