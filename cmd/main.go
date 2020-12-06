package main

import (
	"context"
	"time"

	"github.com/FaranushKarimov/crud/cmd/app"
	"github.com/FaranushKarimov/crud/pkg/customers"
	"github.com/gorilla/mux"
	"go.uber.org/dig"

	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	//это хост
	host := "0.0.0.0"
	//это порт
	port := "9999"
	//это строка подключения к бд
	dbConnectionString := "postgres://app:pass@localhost:5432/db"
	//запускаем функцию execute c проверкой на err
	if err := execute(host, port, dbConnectionString); err != nil {
		//если получили ошибку то закрываем приложения
		log.Print(err)
		os.Exit(1)
	}
}

//функция запуска сервера
func execute(host, port, dbConnectionString string) (err error) {

	//здес обявляем слайс с зависимостями то есть добавляем все сервисы и конструкторы
	dependencies := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			connCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(connCtx, dbConnectionString)
		},
		customers.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    host + ":" + port,
				Handler: server,
			}
		},
	}

	//обявляем новый контейнер
	container := dig.New()
	//в цикле регистрируем все зависимостив контейнер
	for _, v := range dependencies {
		err = container.Provide(v)
		if err != nil {
			return err
		}
	}

	/*вызываем метод Invoke позволяет вызвать в контейнере функцию*/
	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	//если получили ошибку то вернем его
	if err != nil {
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
