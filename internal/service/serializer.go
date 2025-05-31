// создание ответа для 'Response'
package service

import (
	"sort"
	"strconv"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
)

// цитата из базы
type QuoteSerializer struct {
	model.Quote
}

// перевод 'model.Quote' в формат для ответа
func (q *QuoteSerializer) Response() *QuoteResponse {
	return &QuoteResponse{
		ID:     strconv.FormatUint(uint64(q.ID), 10),
		Author: q.Author,
		Body:   q.Body,
	}
}

// шаблон ответа для 'model.Quote'
type QuoteResponse struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"quote"`
}

// список статей из базы
type QuoteListSerializer struct {
	Quotes []model.Quote
}

// перевод '[]model.Quote' в формат для ответа
func (ql *QuoteListSerializer) Response() []QuoteResponse {
	quotes := ql.Quotes // псевдоним - для удобства
	n := len(quotes)
	quoteResponse := make([]QuoteResponse, 0, n)

	// передаём в порядке возрастания ID цитаты
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].ID < quotes[j].ID
	})

	for _, quote := range quotes {
		serialize := QuoteSerializer{Quote: quote}
		quoteResponse = append(quoteResponse, *serialize.Response())
	}

	return quoteResponse
}
