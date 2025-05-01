package main

import (
	"effective/db"
	"effective/routes"
	"effective/services"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "effective/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Effective Mobile API
// @version 1.0
// @description API для управления данными о людях
// @host localhost:8080
// @BasePath /api
func main() {
	log.Println("Запуск сервера...")

	log.Println("Загрузка переменных окружения")
	if err := services.LoadEnvFile(".env"); err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
	}

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
	routes.Init(r.PathPrefix("/api").Subrouter(), manager)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: r,
	}

	log.Printf("Сервер запущен на порту %s", serverPort)
	log.Printf("Swagger UI доступен по адресу: http://localhost:%s/swagger/index.html", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
