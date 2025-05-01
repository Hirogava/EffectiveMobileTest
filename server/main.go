package main

import (
	"effective/db"
	"effective/routes"
	"effective/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Запуск сервера...")

	log.Println("Загрузка переменных окружения")
	services.LoadEnvFile(".env")

	log.Println("Инициализация подключения к базе данных")
	manager := db.NewDBManager("postgres", os.Getenv("DB_CONNECTION_STRING"))

	log.Println("Выполнение миграций базы данных")
	db.Migrate(manager)
	log.Println("База данных успешно инициализирована и мигрирована")

	defer func() {
		log.Println("Закрытие соединения с базой данных")
		manager.Close()
	}()

	r := mux.NewRouter()
	log.Println("Инициализация маршрутов")
	routes.Init(r, manager)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	serverPort := os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: r,
	}

	log.Printf("Сервер запущен на порту %s", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}