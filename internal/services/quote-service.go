package services

import (
	"context"
	"errors"
	"fmt"
	"go-offline-test/internal/repository"
	"go-offline-test/internal/shared/dto"
	"log"
	"strings"
	"unicode"
)

type IQuoteService interface {
	// AddQuote - Добавляет новую цитату.
	AddQuote(ctx context.Context, quote *dto.Quote) error
	// ListQuotes - Получает все существующие цитаты.
	ListQuotes(ctx context.Context) ([]*dto.Quote, error)
	// RandomQuote - Получает рандомную цитату.
	RandomQuote(ctx context.Context) (*dto.Quote, error)
	// QuotesByAuthor - Получает все цитаты автора.
	QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error)
	// DeleteQuote - Удаляет цитату.
	DeleteQuote(ctx context.Context, quoteID int) error
	// ValidateData - Валидирует данные.
	ValidateData(text, authorName, mode string) error
}

type QuoteService struct {
	repo *repository.QuoteRepository
}

func NewQuoteService(repo *repository.QuoteRepository) *QuoteService {
	return &QuoteService{repo: repo}
}

func (qs *QuoteService) AddQuote(ctx context.Context, quote *dto.Quote) error {
	if err := qs.repo.AddQuote(ctx, quote); err != nil {
		if errors.Is(err, repository.ErrQuoteAlreadyExist) {
			log.Printf("WARN: не удалось создать цитату. Ошибка: %v", err)
			return ErrQuoteAlreadyExist
		}
		log.Printf("ERROR: не удалось создать цитату. Ошибка: %v", err)
		return fmt.Errorf(ErrAddQuote.Error(), ": %s", err)
	}

	return nil
}

func (qs *QuoteService) ListQuotes(ctx context.Context) ([]*dto.Quote, error) {
	quotes, err := qs.repo.Quotes(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrQuotesNotFound) {
			log.Printf("WARN: не удалось получить список цитат из памяти. Ошибка: %v", err)
			return nil, ErrNoQuotesAvailable
		}
		log.Printf("ERROR: не удалось получить список цитат из памяти. Ошибка: %v", err)
		return nil, fmt.Errorf(ErrGetQuotes.Error(), ": %s", err)
	}

	return quotes, nil
}

func (qs *QuoteService) RandomQuote(ctx context.Context) (*dto.Quote, error) {
	quote, err := qs.repo.RandomQuote(ctx)
	if err != nil {

		if errors.Is(err, repository.ErrQuoteNotFound) {
			log.Printf("WARN: не удалось получить рандомную цитату из памяти. Ошибка: %v", err)
			return nil, ErrNoQuotesAvailable
		}
		log.Printf("ERROR: не удалось получить рандомную цитату из памяти. Ошибка: %v", err)
		return nil, fmt.Errorf(ErrGetQuote.Error(), ": %s", err)
	}

	return quote, nil
}

func (qs *QuoteService) QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error) {
	quote, err := qs.repo.QuotesByAuthor(ctx, authorName)
	if err != nil {
		switch {

		case errors.Is(err, repository.ErrAuthorNotFound):
			log.Printf("WARN: не удалось получить список цитат по автору из памяти. Ошибка: %v", err)
			return nil, ErrAuthorNotFound
		case errors.Is(err, repository.ErrAuthorQuotesNotFound):
			log.Printf("WARN: не удалось получить список цитат по автору из памяти. Ошибка: %v", err)
			return nil, ErrNoQuotesByThisAuthor
		default:
			log.Printf("ERROR: не удалось получить список цитат по автору из памяти. Ошибка: %v", err)
			return nil, fmt.Errorf(ErrGetQuoteByAuthor.Error(), ": %s", err)
		}
	}

	return quote, nil
}

func (qs *QuoteService) DeleteQuote(ctx context.Context, quoteID int) error {
	if err := qs.repo.DeleteQuote(ctx, quoteID); err != nil {
		if errors.Is(err, repository.ErrQuotesNotFound) {
			log.Printf("WARN: не удалось удалить цитату по id=%d из памяти. Ошибка: %v", err, quoteID)
			return ErrQuoteNotFound
		}
		log.Printf("WARN: не удалось удалить цитату по id=%d из памяти. Ошибка: %v", err, quoteID)
		return err
	}
	return nil
}

func (qs *QuoteService) ValidateData(text, authorName, mode string) error {
	switch mode {
	case "quote":
		text = strings.TrimSpace(text)

		if text == "" {
			err := NewErrInvalidData(400, "цитата не может быть пустой")
			log.Printf("WARN: ошибка валидации цитаты: %v", err)
			return err
		}

		minLength, maxLength := 1, 500
		if len(text) < minLength {
			err := NewErrInvalidData(400, fmt.Sprintf("цитата слишком короткая (минимум %d символов)", minLength))
			log.Printf("WARN: ошибка валидации цитаты: %v", err)
			return err
		}
		if len(text) > maxLength {
			err := NewErrInvalidData(400, fmt.Sprintf("цитата слишком длинная (максимум %d символов)", minLength))
			log.Printf("WARN: ошибка валидации цитаты: %v", err)
			return err
		}

		if err := validateAuthor(authorName); err != nil {
			return err
		}
	case "author":
		if err := validateAuthor(authorName); err != nil {
			return err
		}
	default:
		log.Printf("ERROR: не удалось валидировать переданные данные, text: %s, authorName: %s. Ошибка: указан не существующий метод валидации данных,", text, authorName)
		return fmt.Errorf("не существующий метод проверки")
	}
	return nil
}

func validateAuthor(authorName string) error {
	authorName = strings.TrimSpace(authorName)

	if authorName == "" {
		err := NewErrInvalidData(400, "имя автора не может быть пустым")
		log.Printf("WARN: ошибка валидации имени автора: %v", err)
		return err
	}

	minLength, maxLength := 2, 100
	if len(authorName) < minLength {
		err := fmt.Sprintf("имя автора слишком короткое (минимум %d символов)", minLength)
		log.Printf("WARN: ошибка валидации имени автора: %v", err)
		return NewErrInvalidData(400, err)
	}
	if len(authorName) > maxLength {
		err := fmt.Sprintf("имя автора слишком длинное (максимум %d символов)", maxLength)
		log.Printf("WARN: ошибка валидации имени автора: %v", err)
		return NewErrInvalidData(400, err)
	}

	if strings.HasPrefix(authorName, "-") || strings.HasSuffix(authorName, "-") {
		err := "имя автора не может начинаться или заканчиваться дефисом"
		log.Printf("WARN: ошибка валидации имени автора: %v", err)
		return NewErrInvalidData(400, err)
	}

	for _, r := range authorName {
		if !(unicode.IsLetter(r) || unicode.IsSpace(r) || r == '-') {
			err := fmt.Sprintf("имя автора содержит недопустимые символы: '%s' (символ '%c')", authorName, r)
			log.Printf("WARN: ошибка валидации имени автора: %v", err)
			return NewErrInvalidData(400, err)
		}
	}

	return nil
}
