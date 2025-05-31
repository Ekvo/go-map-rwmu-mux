// логика получения полного списка цитат определенного автора
package service

import (
	"context"
	"log"
	"strings"
)

// содержит 'FindListQuoteByAuthor' -> -> return список цитат по переданному 'author'
type FindListQuoteByAuthor interface {
	ReadQuoteListByAuthor(ctx context.Context, author string) ([]QuoteResponse, error)
}

// перевод 'author' в нижний регистр
// получение цитат из базы по автору, и дальнейшая сериализация для ответа
func (s *serviceQuote) ReadQuoteListByAuthor(
	ctx context.Context,
	author string) ([]QuoteResponse, error) {
	// все данные в базе - в нижнем регистре
	author = strings.ToLower(author)

	quotes, err := s.DBProvider.QuoteListByAuthor(ctx, author)
	if err != nil {
		log.Printf("{service: ReadQuoteListByAuthor error - {%v};", err)
		return nil, err
	}

	serialize := QuoteListSerializer{Quotes: quotes}

	return serialize.Response(), nil
}
