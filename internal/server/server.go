// описание сервера
package server

import (
	"net"
	"net/http"
	"time"

	"github.com/Ekvo/go-map-rwmu-mux/internal/config"
)

// TimeoutShut - используется для 'Shutdown'
const TimeoutShut = 10 * time.Second

// обертка для 'http.Server'
type Srv struct {
	*http.Server
}

// создаем адрес для сервера, добавляем 'Handler'
func InitSRV(cfg *config.Config, router http.Handler) Srv {
	return Srv{
		Server: &http.Server{
			Addr:    net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
			Handler: router,
		},
	}
}
