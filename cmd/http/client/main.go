package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&pow.Dialer{MaxIterations: 100000000}).DialContext,
		},
	}

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for range time.Tick(time.Millisecond * 100) {
				request(client, addr)
			}
		}()
	}

	wg.Wait()
}

func request(client http.Client, addr string) {
	resp, err := client.Get("http://" + addr + "/")
	if err != nil {
		log.Printf("request err: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read resp err: %v", err)
	}

	fmt.Println(string(b))
}
