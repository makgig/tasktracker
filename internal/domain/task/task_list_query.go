package task

// ListQuery содержит параметры для фильтрации списка задач
type ListQuery struct {
	Date    string // Фильтр по дате
	Comment string // Фильтр по комментарию
	Limit   int    // Ограничение количества возвращаемых задач
}
