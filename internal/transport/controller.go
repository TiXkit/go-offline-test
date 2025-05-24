package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-offline-test/internal/services"
	"go-offline-test/internal/shared/dto"
	"log"
	"net/http"
)

type Controller struct {
	services.IQuoteService
}

func NewController(service services.IQuoteService) *Controller {
	return &Controller{service}
}

func (c *Controller) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
	log.Printf("INFO [%d] ответ успешно отправлен\n%s", status, data)
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, err error, status int) {
	if status == 0 {
		switch {
		case errors.Is(err, services.ErrNoQuotesAvailable) ||
			errors.Is(err, services.ErrAuthorNotFound) ||
			errors.Is(err, services.ErrQuoteNotFound) ||
			errors.Is(err, services.ErrNoQuotesByThisAuthor):
			status = 404
		default:
			status = 500
		}
	}
	log.Printf("ERROR [%d]: %v", status, err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
		"code":  status,
	})
}

func (c *Controller) AddQuote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quote, ok := r.Context().Value(quoteCtxKey).(*dto.Quote)
		if !ok {
			c.error(w, r, fmt.Errorf("quote data missing"), http.StatusBadRequest)
			return
		}

		if err := c.IQuoteService.AddQuote(r.Context(), quote); err != nil {
			c.error(w, r, err, 0)
			return
		}

		c.respond(w, r, quote, http.StatusCreated)

	}
}

func (c *Controller) DeleteQuote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value(quoteIDCtxKey).(int)
		if !ok {
			c.error(w, r, fmt.Errorf("quote data missing"), http.StatusBadRequest)
			return
		}

		if err := c.IQuoteService.DeleteQuote(r.Context(), id); err != nil {
			c.error(w, r, err, 0)
			return
		}

		c.respond(w, r, nil, http.StatusNoContent)
	}
}

func (c *Controller) RandomQuote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := c.IQuoteService.RandomQuote(r.Context())
		if err != nil {
			c.error(w, r, err, 0)
			return
		}
		c.respond(w, r, quote, http.StatusOK)
	}
}

func (c *Controller) GetQuotesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var quotes []*dto.Quote
		var err error

		authorHeader := r.Header.Get("author")
		log.Printf("INFO: данные запроса успешно получены:{\nAuthor: %s\n}", authorHeader)
		if authorHeader != "" {
			quotes, err = c.IQuoteService.QuotesByAuthor(r.Context(), authorHeader)
			if err := c.IQuoteService.ValidateData("", authorHeader, authorMode); err != nil {
				c.error(w, r, err, 400)
				return
			}
		} else {
			quotes, err = c.IQuoteService.ListQuotes(r.Context())
		}

		if err != nil {
			c.error(w, r, err, 0)
			return
		}

		c.respond(w, r, quotes, http.StatusOK)
	}
}
