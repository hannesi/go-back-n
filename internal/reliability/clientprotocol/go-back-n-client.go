package clientprotocol

import (
	"log"
	"math"
	"time"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability"
	"github.com/hannesi/go-back-n/internal/virtualsocket"
	"github.com/hannesi/go-back-n/pkg/utils"
)

type GoBackNProtocolClient struct {
	socket                virtualsocket.VirtualSocket
	sequencer             utils.Sequencer
	packetQueue           []reliability.ReliableDataTransferPacket
	highestAckSeqReceived uint8
}

func NewGoBackNProtocolClient(socket virtualsocket.VirtualSocket) (GoBackNProtocolClient, error) {
	client := GoBackNProtocolClient{
		socket:      socket,
		sequencer:   utils.NewSequencer(config.DefaultConfig.GoBackNMaxSequence),
		packetQueue: []reliability.ReliableDataTransferPacket{},
	}

	helloAttempts := 0
	gotResponse := false
	for helloAttempts < config.DefaultConfig.HelloCountBeforeQuit && !gotResponse {
		gotResponse, _ = client.sendHello()
		helloAttempts++
	}

	if !gotResponse {
		return GoBackNProtocolClient{}, reliability.HelloError{}
	}

	return client, nil
}

func (client GoBackNProtocolClient) sendHello() (bool, error) {
	// TODO: replace the mystery constant below
	log.Println("Sending HELLO...")
	buffer := make([]byte, 5)

	client.socket.Send([]byte(config.DefaultConfig.HelloMessage))

	_, err := client.socket.Receive(buffer)
	if err != nil {
		return false, err
	}

	res := string(buffer[:])

	log.Printf("Received response to HELLO: %+v\n", res)

	return res == config.DefaultConfig.HelloMessage, err
}

func (client GoBackNProtocolClient) Send(data [][]byte) error {
	// form rdt packets from each byte array
	for _, payload := range data {
		rdtPacket := reliability.NewReliableDataTransferPacket(client.sequencer.Next(), payload)
		client.packetQueue = append(client.packetQueue, rdtPacket)
	}

	// send queued packets
	client.sendPacketQueue()

	return nil
}

func (client GoBackNProtocolClient) sendPacketQueue() {
	batch := client.makeBatch()
	ackChan := make(chan uint8)
	go client.listenForAcks(ackChan)
	client.sendBatch(batch)
	client.highestAckSeqReceived = <-ackChan
}

func (client GoBackNProtocolClient) listenForAcks(ackChannel chan uint8) {
	log.Print("Listening for acks...")
	highestAckSeqReceived := client.highestAckSeqReceived
	startTime := time.Now()
	for time.Now().Sub(startTime) < config.DefaultConfig.GoBackNAckCollectingTime {
		buffer := make([]byte, 4)
		_, err := client.socket.Receive(buffer)
		if err != nil {
			continue
		}
		ack, err := reliability.DeserializeAckBytes(buffer)
		if err != nil {
			continue
		}
		highestAckSeqReceived = ack.Sequence
	}
	log.Printf("Highest ack received: %d", highestAckSeqReceived)
	ackChannel <- highestAckSeqReceived
}

func (client GoBackNProtocolClient) makeBatch() [][]byte {
	n := int(math.Min(float64(config.DefaultConfig.GoBackNWindowSize), float64(len(client.packetQueue))))
	serializedPackets := make([][]byte, n)
	for i := range n {
		data, err := client.packetQueue[i].Serialize()
		if err != nil {
			log.Printf("%v", err)
		}
		serializedPackets[i] = data
	}
	return serializedPackets
}

func (client GoBackNProtocolClient) sendBatch(batch [][]byte) {
	for _, packet := range batch {
		client.socket.Send(packet)
	}
}
