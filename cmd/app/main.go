package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ekvo/go-map-rwmu-mux/internal/app"
	"github.com/Ekvo/go-map-rwmu-mux/internal/config"
)

func main() {
	cfg, err := config.NewConfig("./init/.env")
	if err != nil {
		log.Fatalf("main: config error - {%v};", err)
	}

	qb := app.NewQuotationBook(cfg)

	qb.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sigChan

	qb.Stop()
}
