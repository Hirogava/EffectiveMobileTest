package services

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvFile(filename string) error {
	log.Printf("Загрузка переменных окружения из файла %s", filename)
	err := godotenv.Load(filename)
	if err != nil {
		log.Printf("Ошибка загрузки переменных окружения: %v", err)
	} else {
		log.Printf("Переменные окружения успешно загружены")
	}
	return err
}