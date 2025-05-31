package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
)

// набор цитат для записи в базу
// изменять лишь в случае полного понимания работы тестов в данном пакете
var quotesData = []model.Quote{
	{
		Author: `William James`,
		Body:   `The greatest weapon against stress is our ability to choose one thought over another`,
	},
	{
		Author: `Napoleon Bonaparte`,
		Body:   `My dictionary does not contain the word 'impossible'`,
	},
	{
		Author: `Steve Jobs`,
		Body:   `Your time is limited, so don’t waste it living someone else’s life`,
	},
}

func TestProvider_NewQuote(t *testing.T) {
	ctx := context.TODO()
	pr := NewProvider()

	testData := []struct {
		title           string
		quotesDataIndex uint
		err             error
	}{
		{
			title:           `add William James`,
			quotesDataIndex: 0,
			err:             nil,
		},
		{
			title:           `add Napoleon Bonaparte`,
			quotesDataIndex: 1,
			err:             nil,
		},
		{
			title:           `again add William James`,
			quotesDataIndex: 0,
			err:             ErrDBAlreadyExists,
		},
		{
			title:           `add Steve Jobs`,
			quotesDataIndex: 2,
			err:             nil,
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			if test.quotesDataIndex >= uint(len(quotesData)) {
				t.Fatal("Test package is broken")
			}

			err := pr.NewQuote(ctx, quotesData[test.quotesDataIndex])

			if !errors.Is(err, test.err) {
				t.Errorf("errors not equal {got}:{want} {%v}:{%v};", err, test.err)
			}
		})
	}
}

func TestProvider_QuoteList(t *testing.T) {
	ctx := context.TODO()
	pr := NewProvider()

	testData := []struct {
		title  string
		quotes []model.Quote
	}{
		{
			title: `William James`,
			quotes: []model.Quote{
				{
					ID:     1,
					Author: `William James`,
					Body:   `The greatest weapon against stress is our ability to choose one thought over another`,
				},
			},
		},
		{
			title: `William James, Napoleon Bonaparte`,
			quotes: []model.Quote{
				{
					ID:     1,
					Author: `William James`,
					Body:   `The greatest weapon against stress is our ability to choose one thought over another`,
				},
				{
					ID:     2,
					Author: `Napoleon Bonaparte`,
					Body:   `My dictionary does not contain the word 'impossible'`,
				},
			},
		},
		{
			title: `William James, Napoleon Bonaparte, Steve Jobs`,
			quotes: []model.Quote{
				{
					ID:     1,
					Author: `William James`,
					Body:   `The greatest weapon against stress is our ability to choose one thought over another`,
				},
				{
					ID:     2,
					Author: `Napoleon Bonaparte`,
					Body:   `My dictionary does not contain the word 'impossible'`,
				},
				{
					ID:     3,
					Author: `Steve Jobs`,
					Body:   `Your time is limited, so don’t waste it living someone else’s life`,
				},
			},
		},
	}

	// проверка при пустой базе данных
	_, err := pr.QuoteList(ctx)
	if !errors.Is(err, ErrDBEmpty) {
		t.Fatalf("the errors must be equal - {got}:{want} {%v}:{%v}", err, ErrDBEmpty)
	}

	if len(testData) != len(quotesData) {
		t.Fatal("Test package is broken")
	}

	for i, quote := range quotesData {
		t.Run(testData[i].title, func(t *testing.T) {
			// добавляем в базу по одному
			err := pr.NewQuote(ctx, quote)
			if err != nil {
				t.Fatalf("NewQuote: error should be nil - {%v}", err)
			}

			// получаем текуший список
			quoteList, err := pr.QuoteList(ctx)
			if err != nil {
				t.Fatalf("QuoteList: error should be nil - {%v}", err)
			}

			// упорядочиваем ответ, так как он пришёл из map
			sort.Slice(quoteList, func(i, j int) bool {
				return quoteList[i].ID < quoteList[j].ID
			})

			// сравниваем с ожидаемым
			if !reflect.DeepEqual(quoteList, testData[i].quotes) {
				t.Errorf("quotes not equal - {got}:{want} {%v}:{%v}", quoteList, testData[i].quotes)
			}
		})
	}
}

func TestProvider_QuoteListByAuthor(t *testing.T) {
	ctx := context.TODO()
	pr := NewProvider()

	tmpQuotesData := append(quotesData, model.Quote{
		Author: `William James`,
		Body:   `Action may not always bring happiness, but there is no happiness without action`,
	})

	// записываем цитаты в базу
	for _, quote := range tmpQuotesData {
		t.Run("add quote", func(t *testing.T) {
			err := pr.NewQuote(ctx, quote)
			if err != nil {
				t.Fatalf("NewQuote: error should be nil - {%v}", err)
			}
		})
	}

	testData := []struct {
		title  string
		author string
		quotes []model.Quote
		err    error
	}{
		{
			title:  `Napoleon Bonaparte - one quote`,
			author: `Napoleon Bonaparte`,
			quotes: []model.Quote{
				{
					ID:     2,
					Author: `Napoleon Bonaparte`,
					Body:   `My dictionary does not contain the word 'impossible'`,
				},
			},
			err: nil,
		},
		{
			title:  `Steve Jobs - one quote`,
			author: `Steve Jobs`,
			quotes: []model.Quote{
				{
					ID:     3,
					Author: `Steve Jobs`,
					Body:   `Your time is limited, so don’t waste it living someone else’s life`,
				},
			},
			err: nil,
		},
		{
			title:  `William James - two quotes`,
			author: `William James`,
			quotes: []model.Quote{
				{
					ID:     1,
					Author: `William James`,
					Body:   `The greatest weapon against stress is our ability to choose one thought over another`,
				},
				{
					ID:     4,
					Author: `William James`,
					Body:   `Action may not always bring happiness, but there is no happiness without action`,
				},
			},
			err: nil,
		},
		{
			title:  `Alien - zero quote`,
			author: `Unknown`,
			quotes: nil,
			err:    ErrDBNotFound,
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			// получаем текуший список по автору
			quoteList, err := pr.QuoteListByAuthor(ctx, test.author)
			if !errors.Is(err, test.err) {
				t.Errorf("errors not equal {got}:{want} {%v}{%v}", err, test.err)
			}

			// сравниваем с ожидаемым
			if !reflect.DeepEqual(quoteList, test.quotes) {
				t.Errorf("quotes not equal - {got}:{want} {%v}:{%v}", quoteList, test.quotes)
			}
		})
	}
}

func TestProvider_RandomQuote(t *testing.T) {
	ctx := context.TODO()
	pr := NewProvider()

	// Создаем для проверки получение случайно выбранной цитаты из хранилища
	tmpQuoteMap := map[uint]model.Quote{}

	for i, quote := range quotesData {
		quote.ID = uint(i) + 1
		tmpQuoteMap[quote.ID] = quote
	}

	if len(tmpQuoteMap) != len(quotesData) {
		t.Fatal("Test package is broken")
	}

	for i, quote := range quotesData {
		t.Run(fmt.Sprintf("find random quote: part - {%d}", i+1), func(t *testing.T) {
			err := pr.NewQuote(ctx, quote)
			if err != nil {
				t.Fatalf("NewQuote: error should be nil - {%v}", err)
			}

			// получаем случайную цитату
			randomQuote, err := pr.RandomQuote(ctx)
			if err != nil {
				t.Fatalf("RandomQuote: error should be nil - {%v}", err)
			}

			// получаем цитату по ID для сравнения
			quote, ex := tmpQuoteMap[randomQuote.ID]
			if !ex {
				t.Fatalf("RandomQuote: random quote by ID - {%d} not exists;", randomQuote.ID)
			}

			// сравниваем с ожидаемым
			if !reflect.DeepEqual(*randomQuote, quote) {
				t.Errorf("quote not equal {got}:{want} {%v}{%v}", *randomQuote, quote)
			}
		})
	}
}

func TestProvider_RemoveQuote(t *testing.T) {
	ctx := context.TODO()
	pr := NewProvider()

	for _, quote := range quotesData {
		t.Run("add quote", func(t *testing.T) {
			err := pr.NewQuote(ctx, quote)
			if err != nil {
				t.Fatalf("NewQuote: error should be nil - {%v}", err)
			}
		})
	}

	testData := []struct {
		title   string
		idQuote uint
		err     error
	}{
		{
			title:   `remove quote with ID = 1`,
			idQuote: 1,
			err:     nil,
		},
		{
			title:   `remove quote with ID = 3`,
			idQuote: 3,
			err:     nil,
		},
		{
			title:   `remove quote with ID = 2`,
			idQuote: 2,
			err:     nil,
		},
		{
			title:   `remove quote with ID = 1`,
			idQuote: 1,
			err:     ErrDBNotFound,
		},
	}

	for _, test := range testData {
		t.Run(test.title, func(t *testing.T) {
			err := pr.RemoveQuote(ctx, test.idQuote)
			if !errors.Is(err, test.err) {
				t.Errorf("errors not equal {got}:{want} {%v}{%v}", err, test.err)
			}
		})
	}
}
