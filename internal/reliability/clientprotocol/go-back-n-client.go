package clientprotocol

import (
	"log"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability"
	"github.com/hannesi/go-back-n/internal/virtualsocket"
	"github.com/hannesi/go-back-n/pkg/utils"
)

type GoBackNProtocolClient struct {
	socket     *virtualsocket.VirtualSocket
	sequencer  utils.Sequencer
	windowSize uint8
	buffer     [][]byte
}

func NewGoBackNProtocolClient(socket *virtualsocket.VirtualSocket) (GoBackNProtocolClient, error) {
	client := GoBackNProtocolClient{
		socket:    socket,
		sequencer: utils.NewSequencer(config.DefaultConfig.GoBackNMaxSequence),
		buffer:    [][]byte{},
	}

	helloResponse, err := client.sendHello()
	if err != nil {
		return GoBackNProtocolClient{}, err
	}

	client.windowSize = helloResponse.WindowSize
	client.sequencer = utils.NewSequencer(helloResponse.MaxSequence)
	client.sequencer.SetCurrentValue(helloResponse.CurrentSequence)

	return client, nil
}

func (client GoBackNProtocolClient) sendHello() (reliability.HelloResponse, error) {
	// TODO: replace the mystery constant below
	log.Println("Sending HELLO")
	buffer := make([]byte, 3)

	client.socket.Send([]byte(config.DefaultConfig.HelloMessage))

	_, err := client.socket.Receive(buffer)
	if err != nil {
		return client.sendHello()
	}

	res, err := reliability.DeserializeHelloResponse(buffer)
	if err != nil {
		return client.sendHello()
	}

	log.Printf("Received response to HELLO: %+v\n", res)

	return res, err
}

func (client GoBackNProtocolClient) Send(data [][]byte) error {
	client.buffer = append(client.buffer, data...)
	return nil
}
