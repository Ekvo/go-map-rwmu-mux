// логика получения случайной цитаты
package service

import (
	"context"
	"log"
)

// содержит 'ReadRandomQuote' -> -> return цитату
type FindRandomQuote interface {
	ReadRandomQuote(ctx context.Context) (*QuoteResponse, error)
}

// идем в базу, при нахождении сериализуем цитату в ответ
func (s *serviceQuote) ReadRandomQuote(ctx context.Context) (*QuoteResponse, error) {
	quote, err := s.DBProvider.RandomQuote(ctx)
	if err != nil {
		log.Printf("service: ReadRandomQuote error - {%v};", err)
		return nil, err
	}

	serialize := QuoteSerializer{Quote: *quote}

	return serialize.Response(), nil
}
