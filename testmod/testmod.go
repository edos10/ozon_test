package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	storageType := flag.String("storage", "redis", "Storage type: redis or postgres")
	flag.Parse()

	var db Database

	switch *storageType {
	case "redis":
		db = NewRedisDatabase()
	case "postgres":
		db = NewPostgresDatabase()
	default:
		log.Fatal("Invalid storage type")
	}

	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenURLHandler(db)).Methods("POST")
	r.HandleFunc("/get", GetURLHandler(db)).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Создаем канал для получения сигналов ОС
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Ожидаем получения сигнала
	<-stop

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Закрываем соединение с базой данных
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

	// Останавливаем HTTP-сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}