package main

import (
	"log"
	"tasktracker/internal/app"
)

func main() {
	// Создаем приложение
	application, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	// Запускаем сервер
	if err := application.Start(); err != nil {
		log.Fatal(err)
	}
}
