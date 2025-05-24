package repository

import (
	"context"
	"errors"
	"go-offline-test/internal/shared/dto"
	"math/rand"
	"sync"
)

var (
	ErrQuotesNotFound       = errors.New("нет доступных цитат в памяти")
	ErrAuthorQuotesNotFound = errors.New("нет доступных цитат автора в памяти")
	ErrAuthorNotFound       = errors.New("автор не найден в памяти")
	ErrQuoteNotFound        = errors.New("цитата не найдена в памяти")
	ErrQuoteAlreadyExist    = errors.New("цитата уже существует")
)

type QuoteRepository struct {
	quotes        map[int]*dto.Quote
	authors       map[string]*dto.Author
	quoteCounter  int
	authorCounter int
	freeIDs       map[int]bool
	mu            sync.RWMutex
}

func NewQuoteRepository() *QuoteRepository {
	return &QuoteRepository{
		quotes:  make(map[int]*dto.Quote),
		authors: make(map[string]*dto.Author),
		freeIDs: make(map[int]bool),
	}
}

func (qr *QuoteRepository) AddQuote(ctx context.Context, quote *dto.Quote) error {
	qr.mu.Lock()
	defer qr.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	// Проверяем есть ли свободные id в списке для ключа
	if len(qr.freeIDs) > 0 {
		// Если есть - записываем по ключу этого id данные
		for id := range qr.freeIDs {
			quote.ID = id
			delete(qr.freeIDs, id)
			break
		}
		// В ином случае, пользуемся обычной логикой, прибавляем к счётчику единицу и записываем по этому ключу данные.
		// Есть счётчик, который регистрирует свободные ключи в qr.freeIDs, если свободных ключей нет - высчитывается и потом сохраняется последний int счётчика
		// и данные записываются по нему, как по ключу.
	} else {
		qr.quoteCounter++
		quote.ID = qr.quoteCounter
	}

	// Записываем данные
	qr.quotes[quote.ID] = quote

	// Проверяем, существует ли указанный автор, если нет - создаём. Логика со счётчиками такая же, как и с цитатами
	if author, exists := qr.authors[quote.AuthorName]; !exists {
		qr.authorCounter++
		qr.authors[quote.AuthorName] = &dto.Author{
			ID:         qr.authorCounter,
			AuthorName: quote.AuthorName,
			Quotes:     []*dto.Quote{quote},
		}
	} else {
		// Проверяем, нет ли такой же цитаты у автора
		for _, q := range author.Quotes {
			if q.Text == quote.Text {
				return ErrQuoteAlreadyExist
			}
		}
		author.Quotes = append(author.Quotes, quote)
	}

	return nil
}

func (qr *QuoteRepository) Quotes(ctx context.Context) ([]*dto.Quote, error) {
	qr.mu.RLock()
	defer qr.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Проверяем есть ли цитаты в памяти
	if len(qr.quotes) == 0 {
		return nil, ErrQuotesNotFound
	}

	// Переписываем из мапы в слайс.
	quotes := make([]*dto.Quote, 0, len(qr.quotes))
	for _, quote := range qr.quotes {
		if quote != nil {
			quotes = append(quotes, quote)
		}
	}

	if len(quotes) == 0 {
		return nil, ErrQuotesNotFound
	}

	return quotes, nil
}

func (qr *QuoteRepository) RandomQuote(ctx context.Context) (*dto.Quote, error) {
	qr.mu.RLock()
	defer qr.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Юзаем раннее описанную фанку и получаем слайс.
	quotes, err := qr.Quotes(ctx)
	if err != nil {
		return nil, err
	}

	// Проверяем ещё раз слайс на пустоту.
	if len(quotes) == 0 {
		return nil, ErrQuotesNotFound
	}

	// Задаёт рандомную генерацию размером до длинны слайса.
	maxQuotes := len(quotes)
	random := rand.Intn(maxQuotes)

	// возвращаем.
	return quotes[random], nil
}

func (qr *QuoteRepository) QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error) {
	qr.mu.RLock()
	defer qr.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Проверяем, существует ли автор.
	author, exists := qr.authors[authorName]
	if !exists {
		return nil, ErrAuthorNotFound
	}

	// Проверяем есть ли цитаты у автора.
	if len(author.Quotes) == 0 {
		return nil, ErrAuthorQuotesNotFound
	}

	return author.Quotes, nil
}

func (qr *QuoteRepository) DeleteQuote(ctx context.Context, idQuote int) error {
	qr.mu.Lock()
	defer qr.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	// Проверяем существование цитаты с данным ID.
	if _, exists := qr.quotes[idQuote]; !exists {
		return ErrQuoteNotFound
	}

	// Удаляем из автора.
	quote := qr.quotes[idQuote]
	if author, exists := qr.authors[quote.AuthorName]; exists {
		for i, q := range author.Quotes {
			if q.ID == idQuote {
				author.Quotes = append(author.Quotes[:i], author.Quotes[i+1:]...)
				break
			}
		}
	}

	// Удаляем из цитат всех.
	delete(qr.quotes, idQuote)
	qr.freeIDs[idQuote] = true

	// Декрементируем счётчик, если id был максимальным для счётчика
	if idQuote == qr.quoteCounter {
		qr.quoteCounter--
		for i := qr.quoteCounter + 1; ; i++ {
			if _, exists := qr.freeIDs[i]; exists {
				delete(qr.freeIDs, i)
			} else {
				break
			}
		}
	}

	return nil
}
