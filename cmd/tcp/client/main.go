package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"world-of-wisdom/pkg/pow"
)

func main() {
	addr, exists := os.LookupEnv("ADDR")
	if !exists {
		addr = "localhost:2000"
	}
	d := &pow.Dialer{MaxIterations: 1000000}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for range time.Tick(time.Millisecond * 100) {
				connect(d, addr)
			}
		}()
	}

	wg.Wait()
}

func connect(d *pow.Dialer, addr string) {
	conn, err := d.DialContext(context.Background(), "tcp", addr)
	if err != nil {
		log.Printf("dial err: %v", err)
	}

	b, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("read resp err: %v", err)
	}

	fmt.Println(string(b))
}
