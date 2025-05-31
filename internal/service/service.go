// бизнес логика сервиса
package service

import (
	"errors"

	"github.com/Ekvo/go-map-rwmu-mux/internal/db"
)

// ошибка для 'http.StatusBadRequest'
var ErrServiceInvalidData = errors.New("invalid data")

// содежит 'FindListQuote - read_quote_list.go; FindListQuoteByAuthor - read_quotes_by_author;'
type FindList interface {
	FindListQuote
	FindListQuoteByAuthor
}

// вся логика
type ServiceQuote interface {
	AddQuote
	FindRandomQuote
	FindList
	RemoveQuote
}

// содержит db.Provider
type serviceQuote struct {
	DBProvider db.Provider
}

// конструктор для serviceQuote
func NewService(dbProvider db.Provider) *serviceQuote {
	return &serviceQuote{DBProvider: dbProvider}
}
