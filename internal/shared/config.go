package shared

import (
	"go-offline-test/internal/shared/dto/config"
	"log"
	"os"
)

func GetAddr() *config.AddrConfig {
	addr := os.Getenv("ADDR_CONFIG")
	port := os.Getenv("PORT_CONFIG")
	if addr == "" || port == "" {
		log.Fatal("addr || port не найден(ы) в файле .env")
	}

	return &config.AddrConfig{Addr: addr, Port: port}
}
