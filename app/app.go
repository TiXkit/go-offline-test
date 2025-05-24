package main

import (
	"go-offline-test/internal/repository"
	"go-offline-test/internal/services"
	"go-offline-test/internal/transport"
	"log"
)

func main() {
	// Для запуска без докер
	// projectRoot, err := filepath.Abs("../")
	// if err != nil {
	// 	log.Fatal("Не удалось получить projectRoot:", err)
	// }
	//
	// envPath := filepath.Join(projectRoot, ".env")
	// if err := godotenv.Load(envPath); err != nil {
	//	log.Fatal("Не удалось загрузить .env файл")
	// }
	// log.Printf("INFO: .env файл успешно загружен")

	repo := repository.NewQuoteRepository()
	log.Printf("INFO: слой репозитория успешно создан")
	service := services.NewQuoteService(repo)
	log.Printf("INFO: сервисный слой успешно создан")
	controller := transport.NewController(service)
	log.Printf("INFO: транспортный слой успешно создан")
	transport.RunRouter(controller)
}
