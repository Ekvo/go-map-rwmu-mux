// реализация для 'CreateQuote'
package service

// логика создания новой цитаты
import (
	"context"
	"log"
	"strings"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
)

// содержит метод 'CreateQuote' -> создания цитаты
type AddQuote interface {
	CreateQuote(ctx context.Context, quote model.Quote) error
}

// перед сохранением, переводим все данные 'quote' в нижний регистр
func (s *serviceQuote) CreateQuote(ctx context.Context, quote model.Quote) error {
	// предотвращаем повторение данных при разных регистрах
	quote.Author = strings.ToLower(quote.Author)
	quote.Body = strings.ToLower(quote.Body)

	if err := s.DBProvider.NewQuote(ctx, quote); err != nil {
		log.Printf("service: CreateQuote error - {%v};", err)
		return err
	}

	return nil
}
