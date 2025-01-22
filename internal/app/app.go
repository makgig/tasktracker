package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tasktracker/internal/domain/task"
	"tasktracker/internal/storage/sqlite"
	"tasktracker/internal/transport"
	"tasktracker/tests" // Используем напрямую настройки из tests
)

// App представляет собой основное приложение
type App struct {
	server  *http.Server
	db      *sqlite.DB
	service *task.Service
	handler *transport.Handler
}

// New создает новый экземпляр приложения
func New() (*App, error) {
	dbFile := tests.DBFile
	if envFile := os.Getenv("TODO_DBFILE"); strings.TrimSpace(envFile) != "" {
		dbFile = envFile
	}
	database, err := sqlite.New(dbFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации БД: %w", err)
	}
	repository := sqlite.NewRepository(database)
	service := task.NewService(repository)
	handler := transport.NewHandler(service)

	return &App{
		db:      database,
		service: service,
		handler: handler,
	}, nil
}

// Start запускает веб-сервер и начинает обработку запросов
func (a *App) Start() error {
	// Получаем порт из конфигурации или переменной окружения
	port := tests.Port
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if eport, err := strconv.ParseInt(portStr, 10, 32); err == nil {
			port = int(eport)
		}
	}

	// Регистрируем маршруты
	a.handler.RegisterRoutes()

	// Создаем HTTP-сервер с нужными настройками
	a.server = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	fmt.Printf("Сервер запущен на http://localhost:%d\n", port)

	// Запускаем сервер
	return a.server.ListenAndServe()
}
