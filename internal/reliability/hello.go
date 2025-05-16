package reliability

import (
	"bytes"

	"github.com/hannesi/go-back-n/internal/config"
)

type HelloResponse struct {
	CurrentSequence uint8
	MaxSequence     uint8
	WindowSize      uint8
}

func NewHelloResponse(currentSequence uint8, maxSequence uint8, windowSize uint8) HelloResponse {
	return HelloResponse{
		CurrentSequence: currentSequence,
		MaxSequence: maxSequence,
		WindowSize: windowSize,
	}
}

func (res HelloResponse) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := buffer.WriteByte(res.CurrentSequence)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(res.MaxSequence)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(res.WindowSize)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DeserializeHelloResponse(data []byte) (HelloResponse, error) {
	buffer := bytes.NewReader(data)

	currentSequence, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	maxSequence, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	windowSize, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	return HelloResponse{
		CurrentSequence: currentSequence,
		MaxSequence:     maxSequence,
		WindowSize:      windowSize,
	}, nil
}

func IsHelloMessage(data []byte) bool {
	return string(data) == config.DefaultConfig.HelloMessage
}
