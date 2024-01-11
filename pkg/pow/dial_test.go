package pow

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialer_Listener(t *testing.T) {
	l, err := Listen(func() uint8 {
		return 1
	}, "tcp", ":0")
	if err != nil {
		t.Error(err)
	}

	results := make([][]byte, 0, 100000)
	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil || conn == nil {
				log.Printf("error accept connection: %v", err)

				continue
			}

			go handleConnection(t, conn, &results)
		}
	}()

	d := Dialer{MaxIterations: 10000000}
	addr := l.Addr().String()
	for i := 0; i < 10; i++ {
		wg := sync.WaitGroup{}
		wg.Add(100)
		for y := 0; y < 100; y++ {
			go func(i, y int) {
				defer wg.Done()
				conn, err := d.DialContext(context.Background(), "tcp", addr)
				if err != nil {
					t.Error(err)
					return
				}

				n, err := conn.Write([]byte{byte(i), byte(y)})
				if err != nil {
					t.Error(err)
				}
				if n != 2 {
					t.Error(io.ErrShortWrite)
					return
				}
			}(i, y)
		}
		wg.Wait()
	}
	time.Sleep(time.Second)
	if len(results) != 1000 {
		t.Error("missed connections")
	}
}

func handleConnection(t *testing.T, conn net.Conn, results *[][]byte) {
	defer conn.Close()

	data := make([]byte, 10)
	n, err := conn.Read(data)
	if err != nil {
		t.Error(err)
		return
	}
	if n != 2 {
		t.Error(io.EOF)
		return
	}

	*results = append(*results, data[:n])
}
