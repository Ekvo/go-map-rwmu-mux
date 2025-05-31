// маршрутизация
package transport

import (
	"github.com/Ekvo/go-map-rwmu-mux/internal/config"
	"github.com/Ekvo/go-map-rwmu-mux/internal/server"
	"github.com/Ekvo/go-map-rwmu-mux/internal/service"

	"net/http"
)

// оболочка для сервера, и 'ServeMux'
type Transport struct {
	*http.ServeMux
	server.Srv
}

// конструктор Transport
func NewTransport(cfg *config.Config) Transport {
	mux := http.NewServeMux()
	return Transport{
		ServeMux: mux,
		Srv:      server.InitSRV(cfg, mux),
	}
}

// создание маршрутов
func (r Transport) Routes(service service.ServiceQuote) {
	r.HandleFunc("POST /quotes", SaveOneQuote(service))
	r.HandleFunc("GET /quotes", RetrieveListOfQuote(service))
	r.HandleFunc("GET /quotes/random", RetrieveRandomQuote(service))
	r.HandleFunc("DELETE /quotes/{id}", ExpelQuote(service))
}
