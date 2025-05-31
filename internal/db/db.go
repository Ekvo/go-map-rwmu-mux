// ранилище для цитат, работает с оперативной памятью
package db

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Ekvo/go-map-rwmu-mux/internal/model"
)

var (
	// критическая ошибка работы базы
	ErrDBInternal = errors.New("internal provider error")

	ErrDBEmpty = errors.New("quote list is empty")

	ErrDBNotFound = errors.New("quote not found")

	ErrDBAlreadyExists = errors.New("quote already exists")
)

// логика взаимодейсвия с хранилищем
type Provider interface {
	NewQuote(ctx context.Context, quote model.Quote) error
	RandomQuote(ctx context.Context) (*model.Quote, error)
	QuoteList(ctx context.Context) ([]model.Quote, error)
	QuoteListByAuthor(ctx context.Context, author string) ([]model.Quote, error)
	RemoveQuote(ctx context.Context, id uint) error
}

// описание базы для цитат
type provider struct {
	// читаем -> RLock()
	// пишем -> Lock()
	rwMu sync.RWMutex

	// полные данные цитаты по индексу
	quoteByID map[uint]model.Quote

	// содержит только уникальные цитаты,
	// у цитаты может быть только один автор,
	// защита от дубликатов
	uniqQuote map[string]struct{}

	// ключ - автор, значение - индексы цитат,
	// для получения цитат по автору
	listOfQuoteIDByAuthor map[string][]uint

	// содержит все индексы текущих цитат,
	// для рандомного поиска
	validQuoteID []uint

	// хранит последний созданный индекс, стартовый 0
	curID uint

	// поиск случайного индекса, для 'validQuoteID'
	src rand.Source
}

// конструктор для 'provider'
func NewProvider() *provider {
	return &provider{
		rwMu:                  sync.RWMutex{},
		quoteByID:             make(map[uint]model.Quote),
		uniqQuote:             make(map[string]struct{}),
		listOfQuoteIDByAuthor: make(map[string][]uint),
		validQuoteID:          []uint{},
		curID:                 0,
		src:                   rand.NewSource(time.Now().Unix()),
	}
}

// перед записью цитаты в базу
func (p *provider) incrementID() {
	p.curID++
	log.Printf("db: incrementID curID - {%d};", p.curID)
}

// делаем случайный индекс и возвращаем ID цитаты
func (p *provider) randomID() (uint, error) {
	n := len(p.validQuoteID)
	if n == 0 {
		return 0, ErrDBEmpty
	}

	randID := uint(int(p.src.Int63()) % n)
	log.Printf("db: randomID created index - {%d} for find ID;", randID)

	return p.validQuoteID[randID], nil
}
