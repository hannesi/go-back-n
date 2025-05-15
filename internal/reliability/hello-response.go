package reliability

import (
	"bytes"
)

type HelloResponse struct {
	MaxSequence uint8
	WindowSize  uint8
}

func (res *HelloResponse) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := buffer.WriteByte(res.MaxSequence)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(res.WindowSize)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DeserializeHelloResponse(data []byte) (*HelloResponse, error) {
	buffer := bytes.NewReader(data)

	maxSequence, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}

	windowSize, err := buffer.ReadByte()
	if err != nil {
		return nil, err
	}

	return &HelloResponse{
		MaxSequence: maxSequence,
		WindowSize:  windowSize,
	}, nil
}
