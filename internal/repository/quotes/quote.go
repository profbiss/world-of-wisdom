package quotes

import (
	"bufio"
	"bytes"
	_ "embed"
	"math/rand"
)

//go:embed quotes.txt
var quotes []byte

type QuoteRepo struct {
	quotes [][]byte
}

func New() *QuoteRepo {
	qr := new(QuoteRepo)

	reader := bytes.NewReader(quotes)
	s := bufio.NewScanner(reader)

	for s.Scan() {
		if q := bytes.TrimSpace(s.Bytes()); len(q) > 0 {
			qr.quotes = append(qr.quotes, q)
		}
	}

	return qr
}

func (q *QuoteRepo) GetQuote() ([]byte, error) {
	return q.quotes[rand.Intn(len(q.quotes))], nil
}
