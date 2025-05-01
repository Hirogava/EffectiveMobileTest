package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(manager *Manager) {
	log.Printf("Начало миграции базы данных")
	driver, err := postgres.WithInstance(manager.Conn, &postgres.Config{})
	if err != nil {
		log.Printf("Ошибка создания драйвера миграции: %v", err)
		panic(fmt.Sprintf("Не удалось создать драйвер миграции: %v", err))
	}

	log.Printf("Создание мигратора")
	m, err := migrate.NewWithDatabaseInstance(
		"file://db//migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Printf("Ошибка создания мигратора: %v", err)
		panic(fmt.Sprintf("Не удалось создать мигратор: %v", err))
	}

	log.Printf("Применение миграций")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("Ошибка применения миграций: %v", err)
		panic(fmt.Sprintf("Не удалось применить миграции: %v", err))
	}

	log.Printf("Миграции успешно применены")
}