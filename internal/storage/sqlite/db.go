package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

type DB struct {
	*sqlx.DB
}

func New(dbFile string) (*DB, error) {
	// Получаем путь к исполняемому файлу приложения
	appPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить путь к исполняемому файлу: %w", err)
	}
	// Получаем директорию приложения
	appDir := filepath.Dir(appPath)
	var absPath string
	if filepath.IsAbs(dbFile) {
		absPath = dbFile
	} else {
		// Если путь относительный, объединяем его с директорией приложения
		// Так как используется один и тот же конфиг с тестами, то для получения относительного пути используем strings.Replace()
		absPath = filepath.Join(appDir, strings.Replace(dbFile, "../", "./", 1))
	}

	// Проверяем существование файла БД
	_, err = os.Stat(absPath)
	var install bool
	if err != nil {
		install = true // Устанавливаем флаг при любой ошибке
		fmt.Println("База данных не найдена, будет создана новая")

		// Создаем директории для БД
		dbDir := filepath.Dir(absPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию %s: %w", dbDir, err)
		}

		// Создаем файл БД
		file, err := os.Create(absPath)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать файл БД: %w", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("не удалось закрыть файл БД: %w", err)
		}
		fmt.Printf("Файл БД создан: %s\n", absPath)
	} else {
		fmt.Printf("Найден существующий файл БД: %s\n", absPath)
	}
	// Подключаемся к БД
	db, err := sqlx.Connect("sqlite3", absPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	database := &DB{DB: db}

	if install {
		if err := database.createSchema(); err != nil {
			// При ошибке создания схемы закрываем соединение
			if closeErr := db.Close(); closeErr != nil {
				return nil, fmt.Errorf("ошибка создания схемы БД: %w; ошибка закрытия соединения: %v", err, closeErr)
			}
			return nil, fmt.Errorf("не удалось создать схему БД: %w", err)
		}
		fmt.Println("База данных успешно инициализирована")
	}
	return database, nil
}

func (db *DB) createSchema() error {
	// Создаем таблицу scheduler
	_, err := db.Exec(`
        CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat VARCHAR(128)
        )
    `)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	// Создаем индекс по дате для быстрой сортировки
	_, err = db.Exec(`
        CREATE INDEX idx_scheduler_date ON scheduler(date)
    `)
	if err != nil {
		return fmt.Errorf("не удалось создать индекс: %w", err)
	}

	return nil
}
