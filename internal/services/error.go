package services

import "errors"

var (
	ErrNoQuotesAvailable    = errors.New("в хранилище нет доступных цитат")
	ErrAuthorNotFound       = errors.New("автор не найден")
	ErrQuoteNotFound        = errors.New("цитата не найдена")
	ErrNoQuotesByThisAuthor = errors.New("цитаты этого автора не найдены")
	ErrQuoteAlreadyExist    = errors.New("цитата уже существует")
	ErrAddQuote             = errors.New("ошибка создания цитаты")
	ErrGetQuotes            = errors.New("ошибка получения списка цитат")
	ErrGetQuote             = errors.New("ошибка получения цитаты")
	ErrGetQuoteByAuthor     = errors.New("ошибка получения цитаты по автору")
)

type ErrInvalidName struct {
	Code    int
	Message string
}

func (ri *ErrInvalidName) Error() string {
	return ri.Message
}

func NewErrInvalidData(code int, message string) *ErrInvalidName {
	return &ErrInvalidName{Code: code, Message: message}
}
