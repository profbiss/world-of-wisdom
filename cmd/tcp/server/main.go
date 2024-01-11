package main

import (
	"io"
	"log"
	"net"
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
		log.Printf("parse complexity err: %v", err)
	}

	repo := quotes.New()

	l, err := pow.Listen(func() uint8 {
		return uint8(iComplexity)
	}, "tcp", addr)
	if err != nil {
		log.Printf("listen err: %v", err)
	}
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil || conn == nil {
			log.Printf("error accept connection: %v", err)

			continue
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn, repo)
	}
}

func handleConnection(conn net.Conn, repo *quotes.QuoteRepo) {
	defer conn.Close()

	msg, err := repo.GetQuote()
	if err != nil {
		log.Printf("error write message: %v", err)

		return
	}

	n, err := conn.Write(msg)
	if err != nil {
		log.Printf("error write message: %v", err)

		return
	}
	if n < len(msg) {
		log.Printf("error write message: %v", io.ErrShortWrite)

		return
	}
}
