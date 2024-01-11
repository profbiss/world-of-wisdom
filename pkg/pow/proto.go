package pow

import "encoding/gob"

var (
	_ = func() bool {
		gob.Register(ChallengePayload{})
		gob.Register(ResponseResourcePayload{})

		return true
	}()
)

type MessageType uint

const (
	_                 MessageType = iota
	RequestChallenge              // from client to server - request new challenge from server
	ResponseChallenge             // from server to client - message with challenge for client
	RequestResource               // from client to server - message with solved challenge
	ResponseResource              // from server to client - message with useful info is solution is correct, or with error if not
)

type Message struct {
	Type    MessageType
	Payload any
}

type ResponseResourcePayload struct {
	Success bool
}
