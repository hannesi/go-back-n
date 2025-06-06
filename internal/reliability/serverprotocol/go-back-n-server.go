package serverprotocol

import (
	"log"
	"net"
	"time"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability"
)

type GoBackNProtocolServer struct {
	Socket           *net.UDPConn
	expectedSequence uint8
}

func NewGoBackNProtocolServer(socket *net.UDPConn) GoBackNProtocolServer {
	return GoBackNProtocolServer{
		Socket:           socket,
		expectedSequence: 0,
	}
}

func (server GoBackNProtocolServer) Receive() error {
	buffer := make([]byte, 1024)
	n, addr, err := server.Socket.ReadFromUDP(buffer)
	if err != nil {
		return server.Receive()
	}

	if reliability.IsHelloMessage(buffer[:n]) {
		// reset sequencer on HELLO
		server.expectedSequence = 0
		server.sendHelloResponse(addr)
		return server.Receive()
	}

	packet, err := reliability.DeserializeReliableDataTransferPacket(buffer[:n])
	time.Sleep(100 * time.Millisecond)

	log.Println(string(packet.Payload))

	// if bit error is detected or a packet with unexpected sequence is received
	if !packet.IsChecksumValid() || server.expectedSequence != packet.Sequence {
		server.sendAck(addr)
		return server.Receive()
	}

	// if packet is ok, increase expectedSequence and send ack
	server.expectedSequence++
	server.sendAck(addr)

	return server.Receive()
}

func (server GoBackNProtocolServer) sendAck(dest *net.UDPAddr) {
	ack := reliability.NewAckPacket("ACK", server.expectedSequence)
	serializedAck, err := ack.Serialize()
	log.Printf("Sending %v", serializedAck)
	log.Printf("seq: %d", server.expectedSequence)
	if err != nil {
		log.Println("Server failed to serialize an ack packet. Trying again.")
		server.sendAck(dest)
	}
	server.Socket.WriteToUDP(serializedAck, dest)
}

func (server GoBackNProtocolServer) sendHelloResponse(addr *net.UDPAddr) {
	server.Socket.WriteToUDP([]byte(config.DefaultConfig.HelloMessage), addr)
}
