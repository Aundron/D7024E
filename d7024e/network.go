package d7024e

import (
	"fmt"
	"log"
	"net"
)

type Packet struct {
	Type       string
	Address    string
	KademliaID KademliaID
	Data       []byte
}

type Network struct {
	Channel chan Packet
}

func NewNetwork(channel chan Packet) *Network {
	newNetwork := &Network{}
	newNetwork.Channel = channel

	return newNetwork
}

func Listen() {
	// TODO
	fmt.Println("Listening")

	// TODO: Bygg en getIP() funktion, annars lyssnar vi på alla IP adresser på systemet
	pc, err := net.ListenPacket("udp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buffer := make([]byte, 8192)
		size, _, _ := pc.ReadFrom(buffer)
		fmt.Println(size)
		// HandleRequest()
	}
}

// func HandleRequest() {
// 	if ping
// 		HandlePing()
// 	if store
// 		HandleStore()
// 	...
//
// SendResponse()
// }

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO

	// Encode new ping packet

	// Open connection to contact ip via DialUDP (!!, you specify a remote in DialUDP so you only receive packets from that specific address (in this case the node we ping))

	// Send encoded packet to contact

	// Set deadline for response

	// Wait for response

	// Decode response

	// HandleResponse(response)

	// Return response
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
