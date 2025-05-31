// логика получения полного списка цитат
package service

import (
	"context"
	"log"
)

// содержит 'ReadQuoteList' -> return список цитат
type FindListQuote interface {
	ReadQuoteList(ctx context.Context) ([]QuoteResponse, error)
}

// получение цитат из базы, и дальнейшая сериализация для ответа
func (s *serviceQuote) ReadQuoteList(ctx context.Context) ([]QuoteResponse, error) {
	quotes, err := s.DBProvider.QuoteList(ctx)
	if err != nil {
		log.Printf("service: ReadQuoteList error - {%v};", err)
		return nil, err
	}

	serialize := QuoteListSerializer{Quotes: quotes}

	return serialize.Response(), nil
}
