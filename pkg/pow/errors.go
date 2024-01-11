package pow

import "errors"

var (
	ErrUnknownError          = errors.New("unknown error")
	ErrUnknownMessage        = errors.New("unknown message")
	ErrMaxIterationsExceeded = errors.New("max iterations exceeded")
	ErrZerosCountBig         = errors.New("zeros count big")
	ErrHashIncorrect         = errors.New("hash incorrect")
	ErrResourceMissMatch     = errors.New("resource missMatch")
	ErrUnknownPayloadType    = errors.New("unknown payload type")
	ErrUnknownMessageType    = errors.New("unknown message type")
	ErrConnectionFailed      = errors.New("connection failed")
)
