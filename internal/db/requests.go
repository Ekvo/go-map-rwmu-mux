// реализация запросов в базу
package db

import (
	"context"
	"log"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
	"github.com/Ekvo/go-map-rwmu-mux/pkg/utils"
)

// добавление цитаты
func (p *provider) NewQuote(_ context.Context, quote model.Quote) error {
	p.rwMu.Lock()
	defer p.rwMu.Unlock()

	// проверка на уникальность
	if _, ex := p.uniqQuote[quote.Body]; ex {
		log.Printf("db: NewQuote quote with body  - {%s} is exists;", quote.Body)
		return ErrDBAlreadyExists
	}

	// создаем ID для цитаты
	p.incrementID()

	// запись данных
	quote.ID = p.curID
	p.quoteByID[quote.ID] = quote
	p.uniqQuote[quote.Body] = struct{}{}
	p.listOfQuoteIDByAuthor[quote.Author] = append(p.listOfQuoteIDByAuthor[quote.Author], quote.ID)
	p.validQuoteID = append(p.validQuoteID, quote.ID)

	log.Printf("db: NewQuote with ID - {%d};", quote.ID)

	return nil
}

// получение случайной цитаты
func (p *provider) RandomQuote(_ context.Context) (*model.Quote, error) {
	p.rwMu.RLock()
	defer p.rwMu.RUnlock()

	randID, err := p.randomID()
	if err != nil {
		return nil, err
	}

	if quote, ex := p.quoteByID[randID]; ex {
		log.Printf("db: RandomQuote by ID - {%d};", p.curID)
		return &quote, nil
	}

	log.Print("db: RandomQuote - internal error;")

	return nil, ErrDBInternal
}

// получение всех цитат
func (p *provider) QuoteList(_ context.Context) ([]model.Quote, error) {
	p.rwMu.RLock()
	defer p.rwMu.RUnlock()

	n := len(p.quoteByID)

	if n == 0 {
		return nil, ErrDBEmpty
	}

	arrQuote := make([]model.Quote, 0, n)

	for _, quote := range p.quoteByID {
		arrQuote = append(arrQuote, quote)
	}

	return arrQuote, nil
}

// список всех цитат по автору
func (p *provider) QuoteListByAuthor(_ context.Context, author string) ([]model.Quote, error) {
	p.rwMu.RLock()
	defer p.rwMu.RUnlock()

	// данный автор отсутствует
	if _, ex := p.listOfQuoteIDByAuthor[author]; !ex {
		log.Printf("db: QuoteListByAuthor autor - {%s} not found;", author)
		return nil, ErrDBNotFound
	}

	arrQuote := make([]model.Quote, 0, len(p.listOfQuoteIDByAuthor[author]))

	for _, quoteID := range p.listOfQuoteIDByAuthor[author] {
		quote, ex := p.quoteByID[quoteID]
		if !ex {
			log.Printf("db: QuoteListByAuthor - internal - not exist quoteID - {%d};", quoteID)
			return nil, ErrDBInternal
		}
		arrQuote = append(arrQuote, quote)
	}

	return arrQuote, nil
}

// удаление цитаты по ID
// проверяем наличие данных о цитате в (quoteByID, uniqQuote, istOfQuoteIDByAuthor, validQuoteID)
// все хорошо -> удаляем
func (p *provider) RemoveQuote(_ context.Context, id uint) error {
	p.rwMu.Lock()
	defer p.rwMu.Unlock()

	// цитата по ID
	quote, ex := p.quoteByID[id]
	if !ex {
		return ErrDBNotFound
	}

	// проверка в uniqQuote
	if _, ex := p.uniqQuote[quote.Body]; !ex {
		log.Printf("db: RemoveQuote - internal - not exist key - {%s} in uniqQuote;", quote.Body)
		return ErrDBInternal
	}

	// список ID цитат по автору
	quotesID, ex := p.listOfQuoteIDByAuthor[quote.Author]
	if !ex {
		log.Printf("db: RemoveQuote - internal - not exist key - {%s} in listOfQuoteIDByAuthor;", quote.Author)
		return ErrDBInternal
	}

	// индекс для ID цитаты из quotesID
	indexFromAuhtor, ex := utils.IndexByValue(quotesID, id)
	if !ex {
		log.Printf(
			"db: RemoveQuote - internal - not exist id - {%d} in listOfQuoteIDByAuthor by author - {%s};",
			indexFromAuhtor, quote.Author)
		return ErrDBInternal
	}

	// индекс для ID цитаты из всех текущих статей
	indexFromQuotes, ex := utils.IndexByValue(p.validQuoteID, id)
	if !ex {
		log.Printf("db: RemoveQuote - internal - not exist id - {%d} in validQuoteID", id)
		return ErrDBInternal
	}

	delete(p.quoteByID, id)
	delete(p.uniqQuote, quote.Body)

	// удаляем цитату из списка автора
	quotesID = append(quotesID[:indexFromAuhtor], quotesID[indexFromAuhtor+1:]...)
	if len(quotesID) == 0 {
		// нет цитат у автора -> удаляем
		delete(p.listOfQuoteIDByAuthor, quote.Author)
	} else {
		// сохраняем новый список
		p.listOfQuoteIDByAuthor[quote.Author] = quotesID
	}

	// удаляем из всех текущих индексов цитат
	p.validQuoteID = append(p.validQuoteID[:indexFromQuotes], p.validQuoteID[indexFromQuotes+1:]...)

	log.Printf("db: RemoveQuote by ID - {%d} is deleted;", id)

	return nil
}
