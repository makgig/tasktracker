package task

import (
	"fmt"
	"time"
)

// Task представляет собой основную бизнес-сущность задачи в планировщике.
// Структура соответствует таблице scheduler в базе данных.
type Task struct {
	// ID задачи, автоинкрементное поле в базе данных
	ID int64 `db:"id" json:"id,string"`

	// Date хранит дату задачи в формате YYYYMMDD (например, 20240126)
	// Используется строковый тип для совместимости с API и базой данных
	Date string `db:"date" json:"date"`

	// Title содержит заголовок задачи
	// Поле не может быть пустым (NOT NULL в БД)
	Title string `db:"title" json:"title"`

	// Comment хранит дополнительное описание задачи
	// Может быть пустым
	Comment string `db:"comment" json:"comment"`

	// Repeat определяет правило повторения задачи
	// Ограничено 128 символами в базе данных
	// Поддерживает форматы:
	// - "d N" - повтор каждые N дней (1-400)
	// - "y" - ежегодный повтор
	// - "w N,M,..." - повтор в указанные дни недели (1-7)
	// - "m N[,M,...] [X,Y,...]" - повтор в указанные дни и месяцы
	Repeat string `db:"repeat" json:"repeat"`
}

// DateFormat определяет формат даты, используемый во всем приложении
const DateFormat = "20060102"

// ValidateDate проверяет корректность даты в формате YYYYMMDD
func ValidateDate(date string) error {
	if len(date) != 8 {
		return fmt.Errorf("некорректная длина даты, ожидается 8 символов")
	}

	_, err := time.Parse(DateFormat, date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %w", err)
	}

	return nil
}

// ParseDate преобразует строковое представление даты в time.Time
func ParseDate(date string) (time.Time, error) {
	return time.Parse(DateFormat, date)
}

// FormatDate преобразует time.Time в строковое представление
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}
