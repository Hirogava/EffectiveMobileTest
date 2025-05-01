package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

type Manager struct {
	Conn *sql.DB
	WG   *sync.WaitGroup
	MU   *sync.RWMutex
}

func NewDBManager(driverName string, sourceName string) *Manager {
	log.Printf("Подключение к базе данных %s", driverName)
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		panic(fmt.Sprintf("Не удалось подключиться к базе данных: %v", err))
	}

	log.Printf("Проверка соединения с базой данных")
	if err = db.Ping(); err != nil {
		log.Printf("Ошибка проверки соединения: %v", err)
		panic(fmt.Sprintf("База данных не отвечает: %v", err))
	}

	log.Printf("Соединение с базой данных успешно установлено")
	return &Manager{
		Conn: db,
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
	}
}

func (manager *Manager) Close() {
	log.Printf("Закрытие соединения с базой данных")
	if manager.Conn != nil {
		manager.Conn.Close()
		manager.Conn = nil
		log.Printf("Соединение с базой данных закрыто")
	}
}