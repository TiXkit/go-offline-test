package transport

import (
	"go-offline-test/internal/shared"
	"log"
	"net/http"
)

func RunRouter(c *Controller) {
	router := http.NewServeMux()

	router.HandleFunc("POST /quotes", c.MiddlewareValidate(c.AddQuote()))
	router.HandleFunc("DELETE /quotes/{id}", c.MiddlewareValidate(c.DeleteQuote()))
	router.HandleFunc("GET /quotes", c.GetQuotesHandler())
	router.HandleFunc("GET /quotes/random", c.RandomQuote())

	addrConf := shared.GetAddr()

	log.Printf("INFO: сервер на порту%s запущен.", addrConf.Port)
	if err := http.ListenAndServe(addrConf.Port, router); err != nil {
		log.Fatal("не удалось запустить сервер,ошибка:", err)
	}
}
