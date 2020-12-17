package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/FaranushKarimov/crud/cmd/app"
	"github.com/FaranushKarimov/crud/pkg/customers"
	"github.com/FaranushKarimov/crud/pkg/customers/security"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@192.168.1.186:5433/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host, port, dsn string) (err error) {
	container := dig.New()
	container.Provide(app.NewServer)
	container.Provide(mux.NewRouter)
	container.Provide(func() (*pgxpool.Pool, error) {
		connCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
		return pgxpool.Connect(connCtx, dsn)
	})
	container.Provide(customers.NewService)
	container.Provide(func(server *app.Server) *http.Server {
		return &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: server,
		}
	})
	container.Provide(security.NewService)

	if err := container.Invoke(func(server *app.Server) {
		server.Init()
	}); err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
