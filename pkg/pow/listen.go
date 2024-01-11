package pow

import (
	"encoding/gob"
	"net"
)

type Listener struct {
	zerosCount func() uint8
	net.Listener
}

func Listen(zerosCount func() uint8, network, address string) (net.Listener, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	return Listener{
		zerosCount: zerosCount,
		Listener:   l,
	}, nil
}

func (s Listener) Accept() (net.Conn, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}

	if err := s.verify(conn); err != nil {
		if err := conn.Close(); err != nil {
			return nil, err
		}

		return nil, err
	}

	return conn, nil
}

func (s Listener) verify(conn net.Conn) error {
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

loop:
	m := &Message{}
	err := dec.Decode(m)
	if err != nil {
		return err
	}

	switch m.Type {
	case RequestChallenge:
		m.Type = ResponseChallenge
		m.Payload = NewChallenge(conn.RemoteAddr().String(), s.zerosCount())
		err = enc.Encode(m)
		if err != nil {
			return err
		}

		goto loop
	case RequestResource:
		if resp, ok := m.Payload.(ChallengePayload); ok {
			err := resp.CheckHash(conn.RemoteAddr().String())
			if err != nil {
				err := enc.Encode(Message{
					Type: ResponseResource,
					Payload: ResponseResourcePayload{
						Success: false,
					},
				})

				return err
			}
			err = enc.Encode(Message{
				Type: ResponseResource,
				Payload: ResponseResourcePayload{
					Success: true,
				},
			})
			if err != nil {
				return err
			}

			return nil
		}

		return ErrUnknownPayloadType
	default:
		return ErrUnknownMessageType
	}
}
