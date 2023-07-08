package main

import (
	"context"
	//"database/sql"
	"fmt"
	"io"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "github.com/lib/pq"
	"github.com/edos10/ozon_test/handlers"
	"github.com/gorilla/mux"

)

const sizeUrl, sizeAlphabet = 10, 63

var numToLetter = make(map[int]byte)
var letterToNum = make(map[byte]int)


const (
	startSmallLts = 97
	endSmallLts = 122
	startBigLts = 65
	endBigLts = 90
	startNums = 48
	downcase = 95
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "default"
	dbname   = "guides_data"
)

/*
         Порядок следования всех символов в map:
	a-z = 0-25
	A-Z = 26-51
	0-9 = 52-61
	_ = 62
*/

func makeMaps() {
	startSym := byte('a')
	for i := 0; i < 26; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('A')
	for i := 26; i < 52; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('0')
	for i := 52; i < 62; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	numToLetter[62] = '_'

	startSym = byte('a')
	numForMap := 0
	for startSym <= 'z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('A')
	for startSym <= 'Z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('0')
	for startSym <= '9' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	letterToNum['_'] = 62
}

type Database interface {
	Save(shortUrl, originalUrl string) error
	Get(shortUrl string) error
	io.Closer
}

func main() {
	typeStorage := flag.String("storage", "redis", "Storage type: redis or postgres")

	flag.Parse()

	var db Database
	switch *typeStorage {

	case "redis":
		db = New
	case "postgres":
		_ = 6
	default:
		log.Fatal("Use postgres or redis, error...")

	}
	fmt.Println("Successfully connected!")
	makeMaps()
	r := mux.NewRouter()
	r.HandleFunc("/shorten", ShortenURLHandler(db)).Methods("POST")
	r.HandleFunc("/get", GetURLHandler(db)).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server gracefully stopped")
}
