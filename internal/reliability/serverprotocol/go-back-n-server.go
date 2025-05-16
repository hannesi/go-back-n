package serverprotocol

import (
	"log"
	"net"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability"
)

type GoBackNProtocolServer struct {
	Socket         *net.UDPConn
	lastOkSequence uint8
	maxSequence    uint8
	windowSize     uint8
}

func NewGoBackNProtocolServer(socket *net.UDPConn) GoBackNProtocolServer {
	return GoBackNProtocolServer{
		Socket:         socket,
		lastOkSequence: 0,
		maxSequence:    config.DefaultConfig.GoBackNMaxSequence,
		windowSize:     config.DefaultConfig.GoBackNWindowSize,
	}
}

func (server GoBackNProtocolServer) Receive() error {
	buffer := make([]byte, 1024)
	n, addr, err := server.Socket.ReadFromUDP(buffer)
	if err != nil {
		return server.Receive()
	}

	if reliability.IsHelloMessage(buffer[:n]) {
		server.sendHelloResponse(addr)
		return server.Receive()
	}

	packet, err := reliability.DeserializeReliableDataTransferPacket(buffer[:n])

	// if bit error is detected or a packet with unexpected sequence is received
	if !packet.IsChecksumValid() || server.lastOkSequence+1 != packet.Sequence {
		return server.Receive()
	}

	server.lastOkSequence = packet.Sequence

	ack := reliability.NewAckPacket("ACK", server.lastOkSequence)
	ackPacket, err := ack.Serialize()
	if err != nil {
		log.Fatal("Server failed to serialize an ack packet.")
	}
	server.Socket.WriteToUDP(ackPacket, addr)
	return nil
}

func (server GoBackNProtocolServer) sendHelloResponse(addr *net.UDPAddr) {
	helloResponse := reliability.NewHelloResponse(server.lastOkSequence, server.maxSequence, server.windowSize)
	res, _ := helloResponse.Serialize()
	server.Socket.WriteToUDP(res, addr)
}
