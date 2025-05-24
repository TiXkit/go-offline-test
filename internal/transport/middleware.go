package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"go-offline-test/internal/shared/dto"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const (
	quoteCtxKey   contextKey = "quote"
	quoteIDCtxKey contextKey = "quoteID"
	quoteMode     string     = "quote"
	authorMode    string     = "author"
)

func (c *Controller) MiddlewareValidate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Body == nil {
				c.error(w, r, fmt.Errorf("request body is required"), http.StatusBadRequest)
				return
			}

			var quote dto.Quote
			if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
				c.error(w, r, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
				return
			}
			log.Printf("INFO: данные запроса успешно получены:{\nID: %d\nQuote: %s\nAuthor: %s\n}", quote.ID, quote.Text, quote.AuthorName)

			if err := c.IQuoteService.ValidateData(quote.Text, quote.AuthorName, quoteMode); err != nil {
				c.error(w, r, err, 400)
				return
			}

			// Сохраняем данные в контекст
			ctx := context.WithValue(r.Context(), quoteCtxKey, &quote)
			r = r.WithContext(ctx)
		}
		// Валидация ID для DELETE и GET by ID
		if r.Method == http.MethodDelete || strings.HasPrefix(r.URL.Path, "/quotes/") {
			idStr := strings.TrimPrefix(r.URL.Path, "/quotes/")
			if idStr == "" {
				c.error(w, r, fmt.Errorf("quote ID is required"), http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				c.error(w, r, fmt.Errorf("invalid quote ID"), http.StatusBadRequest)
				return
			}

			log.Printf("INFO: данные запроса успешно получены:{\nID: %d\n}", id)

			// Сохраняем ID в контекст
			ctx := context.WithValue(r.Context(), quoteIDCtxKey, id)
			r = r.WithContext(ctx)
		}
		next(w, r)
	}

}
