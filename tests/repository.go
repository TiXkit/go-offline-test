package repository_test

import (
	"context"
	"errors"
	"go-offline-test/internal/repository"
	"go-offline-test/internal/shared/dto"
	"sync"
	"testing"
)

func TestAddQuote(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()

	t.Run("Успешное добавление цитаты", func(t *testing.T) {
		quote := &dto.Quote{AuthorName: "Test Author", Text: "Test Quote"}
		err := qr.AddQuote(ctx, quote)
		if err != nil {
			t.Fatalf("AddQuote() error = %v, want nil", err)
		}
		if quote.ID != 1 {
			t.Errorf("quote.ID = %d, want 1", quote.ID)
		}
	})

	t.Run("Добавление с отменённым контекстом", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()
		err := qr.AddQuote(canceledCtx, &dto.Quote{})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("AddQuote() error = %v, want %v", err, context.Canceled)
		}
	})
}

func TestQuotes(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()

	t.Run("Пустой репозиторий", func(t *testing.T) {
		_, err := qr.Quotes(ctx)
		if !errors.Is(err, repository.ErrQuotesNotFound) {
			t.Errorf("Quotes() error = %v, want %v", err, repository.ErrQuotesNotFound)
		}
	})

	t.Run("Успешное получение", func(t *testing.T) {
		qr.AddQuote(ctx, &dto.Quote{AuthorName: "Author", Text: "Quote"})
		quotes, err := qr.Quotes(ctx)
		if err != nil {
			t.Fatalf("Quotes() error = %v, want nil", err)
		}
		if len(quotes) != 1 {
			t.Errorf("len(quotes) = %d, want 1", len(quotes))
		}
	})
}

func TestRandomQuote(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()

	t.Run("Пустой репозиторий", func(t *testing.T) {
		_, err := qr.RandomQuote(ctx)
		if !errors.Is(err, repository.ErrQuotesNotFound) {
			t.Errorf("RandomQuote() error = %v, want %v", err, repository.ErrQuotesNotFound)
		}
	})

	t.Run("Успешное получение", func(t *testing.T) {
		qr.AddQuote(ctx, &dto.Quote{AuthorName: "Author", Text: "Quote 1"})
		qr.AddQuote(ctx, &dto.Quote{AuthorName: "Author", Text: "Quote 2"})
		quote, err := qr.RandomQuote(ctx)
		if err != nil {
			t.Fatalf("RandomQuote() error = %v, want nil", err)
		}
		if quote == nil {
			t.Error("RandomQuote() returned nil, want quote")
		}
	})
}

func TestQuotesByAuthor(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()

	t.Run("Несуществующий автор", func(t *testing.T) {
		_, err := qr.QuotesByAuthor(ctx, "Unknown")
		if !errors.Is(err, repository.ErrAuthorNotFound) {
			t.Errorf("QuotesByAuthor() error = %v, want %v", err, repository.ErrAuthorNotFound)
		}
	})

	t.Run("Успешное получение", func(t *testing.T) {
		qr.AddQuote(ctx, &dto.Quote{AuthorName: "Author", Text: "Quote"})
		quotes, err := qr.QuotesByAuthor(ctx, "Author")
		if err != nil {
			t.Fatalf("QuotesByAuthor() error = %v, want nil", err)
		}
		if len(quotes) != 1 {
			t.Errorf("len(quotes) = %d, want 1", len(quotes))
		}
	})
}

func TestDeleteQuote(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()

	t.Run("Несуществующая цитата", func(t *testing.T) {
		err := qr.DeleteQuote(ctx, 999)
		if !errors.Is(err, repository.ErrQuoteNotFound) {
			t.Errorf("DeleteQuote() error = %v, want %v", err, repository.ErrQuoteNotFound)
		}
	})

	t.Run("Успешное удаление", func(t *testing.T) {
		qr.AddQuote(ctx, &dto.Quote{AuthorName: "Author", Text: "Quote"})
		err := qr.DeleteQuote(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteQuote() error = %v, want nil", err)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	qr := repository.NewQuoteRepository()
	ctx := context.Background()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = qr.AddQuote(ctx, &dto.Quote{AuthorName: "Concurrent", Text: "Quote"})
		}()
	}

	wg.Wait()

	quotes, err := qr.Quotes(ctx)
	if err != nil {
		t.Fatalf("Quotes() error = %v", err)
	}
	if len(quotes) != goroutines {
		t.Errorf("len(quotes) = %d, want %d", len(quotes), goroutines)
	}
}
