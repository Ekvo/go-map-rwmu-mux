// инициализация, запуск и остановка сервиса
package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Ekvo/go-map-rwmu-mux/internal/config"
	"github.com/Ekvo/go-map-rwmu-mux/internal/db"
	"github.com/Ekvo/go-map-rwmu-mux/internal/server"
	"github.com/Ekvo/go-map-rwmu-mux/internal/service"
	"github.com/Ekvo/go-map-rwmu-mux/internal/transport"
)

// ключевые узлы приложения
type QuotationBook struct {
	repository db.Provider
	service    service.ServiceQuote
	transport  transport.Transport
}

// конструктор для QuotationBook
func NewQuotationBook(cfg *config.Config) *QuotationBook {
	qb := &QuotationBook{}

	qb.repository = db.NewProvider()
	qb.service = service.NewService(qb.repository)
	qb.transport = transport.NewTransport(cfg)

	log.Print("app: NewQuotationBook is created")

	return qb
}

// вызываем 'transport.Routes' для создания маршрутов, и запускаем сервер в горутине
func (qb *QuotationBook) Run() {
	log.Print("app: Run Quotation Book")

	qb.transport.Routes(qb.service)

	go func() {
		log.Print("app: listen and serve - start")
		if err := qb.transport.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("go app: Run error - {%v};", err)
		}
		log.Print("app: listen and serve - end")
	}()
}

// запуск 'Shutdown' при помощи 'context'
func (qb *QuotationBook) Stop() {
	log.Print("app: Stop Quotation Book")

	ctx, cancel := context.WithTimeout(context.Background(), server.TimeoutShut)
	defer cancel()

	if err := qb.transport.Shutdown(ctx); err != nil {
		log.Fatalf("app: Stop Shutdown error - {%v};", err)
	}

	log.Print("app: shutdown complete")
}
