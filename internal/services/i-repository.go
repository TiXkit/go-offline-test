package services

import (
	"context"
	"go-offline-test/internal/shared/dto"
)

type IQuoteRepository interface {
	AddQuote(ctx context.Context, quote *dto.Quote) error
	Quotes(ctx context.Context) ([]*dto.Quote, error)
	RandomQuote(ctx context.Context) (*dto.Quote, error)
	QuotesByAuthor(ctx context.Context, authorName string) ([]*dto.Quote, error)
	DeleteQuote(ctx context.Context, idQuote int) error
}
