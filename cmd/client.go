package main

import (
	"fmt"
	"log"

	"github.com/hannesi/go-back-n/internal/virtualsocket"
)

var predefinedMessages = []string{"Alekhine", "Botvinnik", "Capablanca", "Ding", "Euwe", "Finegold", "Giri", "Houska", "Ivanchuk", "Jaenisch", "Karpov", "Löwenthal", "Muzychuk", "Naroditsky", "Ojanen", "Polugaevsky", "Qin", "Réti", "Shirov", "Tal", "Ushenina", "Vachier-Lagrave", "Williams", "Xie", "Yusupov", "Zaitsev"}

func main() {
	fmt.Println("")

	socket, err := virtualsocket.NewVirtualSocket()

	if err != nil {
		log.Fatal("Failed to create virtual socket: %v", err)
	}

	defer socket.Close()
}
