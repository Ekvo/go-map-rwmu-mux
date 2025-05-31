// логика удаления цитаты
package service

import (
	"context"
	"log"
)

// содержит метод 'DeleteQuote' -> удаление цитаты
type RemoveQuote interface {
	DeleteQuote(ctx context.Context, id uint) error
}

// удаляем по ID
func (s *serviceQuote) DeleteQuote(ctx context.Context, id uint) error {
	if err := s.DBProvider.RemoveQuote(ctx, id); err != nil {
		log.Printf("service: DeleteQuote error - {%v};", err)
		return err
	}

	return nil
}
