package pow

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func NewChallenge(res string, zerosCount uint8) ChallengePayload {
	return ChallengePayload{
		Version:    1,
		ZerosCount: zerosCount,
		Date:       time.Now().Unix(),
		Resource:   res,
		Rand:       base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Int63()))),
		Counter:    0,
	}
}

type ChallengePayload struct {
	Version    uint
	ZerosCount uint8
	Date       int64
	Resource   string
	Rand       string
	Counter    uint
}

func (h ChallengePayload) Bytes() []byte {
	return []byte(fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter))
}

func (h ChallengePayload) ComputeHash(maxIterations uint) (ChallengePayload, error) {
	s := sha1.New()
	if s.Size() < int(h.ZerosCount) {
		return h, ErrZerosCountBig
	}

	want := string(make([]byte, h.ZerosCount))

	for h.Counter <= maxIterations || maxIterations <= 0 {
		b := h.Bytes()
		n, err := s.Write(b)
		if err != nil {
			return h, err
		}
		if len(b) != n {
			return h, io.ErrShortWrite
		}

		hash := s.Sum(nil)
		s.Reset()

		if string(hash[:h.ZerosCount]) == want {
			return h, nil
		}

		h.Counter++
	}

	return h, ErrMaxIterationsExceeded
}

func (h ChallengePayload) CheckHash(res string) error {
	s := sha1.New()
	if s.Size() < int(h.ZerosCount) {
		return ErrZerosCountBig
	}

	if h.Resource != res {
		return ErrResourceMissMatch
	}

	want := string(make([]byte, h.ZerosCount))

	s.Write(h.Bytes())
	hash := s.Sum(nil)

	if string(hash[:h.ZerosCount]) == want {
		return nil
	}

	return ErrHashIncorrect
}
