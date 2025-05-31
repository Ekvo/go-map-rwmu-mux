// описывает получение цитаты из 'Request'
package service

import (
	"net/http"
	"strings"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
	"github.com/Ekvo/go-map-rwmu-mux/pkg/utils"
)

// поля для unmarshal 'json'
// 'model' - для создания 'model.Quote' из полученных данных
type QuoteDeserializer struct {
	Author string `json:"author"`
	Body   string `json:"quote"`

	model model.Quote `json:"-"`
}

// уонструктор
func NewQuoteDeserializer() *QuoteDeserializer {
	return &QuoteDeserializer{}
}

// получение 'model.Quote' после 'Decode'
func (q *QuoteDeserializer) Model() model.Quote {
	return q.model
}

// вызывает 'utils.DecodeJSON'
// удаляем пробелы, передаем данные а 'model'
func (q *QuoteDeserializer) Decode(req *http.Request) error {
	if err := utils.DecodeJSON(req, q); err != nil {
		return err
	}

	if q.Author = strings.TrimSpace(q.Author); q.Author == "" {
		return ErrServiceInvalidData
	}

	if q.Body = strings.TrimSpace(q.Body); q.Body == "" {
		return ErrServiceInvalidData
	}

	q.model.Author = q.Author
	q.model.Body = q.Body

	return nil
}
