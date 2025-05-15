package clientprotocol

import (
	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/virtualsocket"
	"github.com/hannesi/go-back-n/pkg/utils"
)

type GoBackNProtocolClient struct {
	socket         *virtualsocket.VirtualSocket
	sequencer      *utils.Sequencer
	gbnWindowSize  uint8
	gbnMaxSequence uint8
}

func NewGoBackNProtocolClient(socket *virtualsocket.VirtualSocket) *GoBackNProtocolClient {
	client := GoBackNProtocolClient{
		socket:    socket,
		sequencer: utils.NewSequencer(config.DefaultConfig.GoBackNWindowMaxValue),
	}

}

func (gbn *GoBackNProtocolClient) Send(data [][]byte) error {

}
