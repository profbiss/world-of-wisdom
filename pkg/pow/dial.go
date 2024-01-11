package pow

import (
	"context"
	"encoding/gob"
	"net"
)

type Dialer struct {
	MaxIterations uint
	net.Dialer
}

type Conn struct {
	conn net.Conn
	enc  *gob.Encoder
	dec  *gob.Decoder
}

func (d Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	conn, err := d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	c := Conn{
		conn: conn,
		enc:  gob.NewEncoder(conn),
		dec:  gob.NewDecoder(conn),
	}

	challenge, err := c.requestChallenge()
	if err != nil {
		return nil, err
	}

	result, err := challenge.ComputeHash(d.MaxIterations)
	if err != nil {
		return nil, err
	}

	err = c.sendResult(result)
	if err != nil {
		return nil, err
	}

	err = c.checkSuccess()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c Conn) requestChallenge() (*ChallengePayload, error) {
	m := Message{Type: RequestChallenge}
	err := c.enc.Encode(m)
	if err != nil {
		return nil, err
	}

	err = c.dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	if m.Type != ResponseChallenge {
		return nil, ErrUnknownMessage
	}

	if challenge, ok := m.Payload.(ChallengePayload); ok {
		return &challenge, nil
	}

	return nil, ErrUnknownMessage
}

func (c Conn) sendResult(result ChallengePayload) error {
	err := c.enc.Encode(Message{
		Type:    RequestResource,
		Payload: result,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c Conn) checkSuccess() error {
	m := Message{}
	err := c.dec.Decode(&m)
	if err != nil {
		return err
	}

	switch m.Type {
	case ResponseResource:
		if resp, ok := m.Payload.(ResponseResourcePayload); ok {
			if resp.Success {
				return nil
			}

			return ErrUnknownError
		}

		return ErrUnknownPayloadType
	default:
		return ErrConnectionFailed
	}
}
