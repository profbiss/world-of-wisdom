package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"world-of-wisdom/internal/repository/quotes"
	"world-of-wisdom/pkg/pow"
)

func main() {
	addr, exists := os.LookupEnv("ADDR")
	if !exists {
		addr = ":2000"
	}

	complexity, exists := os.LookupEnv("COMPLEXITY")
	if !exists {
		complexity = "2"
	}
	iComplexity, err := strconv.Atoi(complexity)
	if err != nil {
		log.Fatal(err)
	}

	l, err := pow.Listen(func() uint8 { return uint8(iComplexity) }, "tcp", addr)
	if err != nil {
		log.Printf("listen err: %v", err)

		return
	}
	defer l.Close()

	repo := quotes.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)

		msg, err := repo.GetQuote()
		if err != nil {
			log.Printf("error write message: %v", err)

			return
		}

		n, err := writer.Write(msg)
		if err != nil {
			log.Printf("error write message: %v", err)

			return
		}
		if n < len(msg) {
			log.Printf("error write message: %v", io.ErrShortWrite)

			return
		}
	})

	err = http.Serve(l, mux)
	if err != nil {
		log.Printf("server err: %v", err)

		return
	}
}
