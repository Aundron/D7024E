package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Packet struct {
	Type    string
	Address string
	//KademliaID KademliaID
	KademliaID string
	Data       string
}

type Network struct {
}

func NewNetwork() *Network {
	newNetwork := &Network{}

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

		recPacket := Packet{}

		json.Unmarshal(buffer[:size], &recPacket)

		fmt.Println(recPacket)
		//HandleRequest()
	}
}

func HandleRequest() {
	// 	if ping
	// 		HandlePing()
	// 	if store
	// 		HandleStore()
	// 	...
	//
	// SendResponse()
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO

	// Encode new ping packet
	packet := Packet{
		Type:       "Ping",
		Address:    contact.Address,
		KademliaID: "jdsaihuasgudhsauiud",
		Data:       "Test",
	}

	b, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("error")
	}

	// Open connection to contact ip via DialUDP (!!, you specify a remote in DialUDP so you only receive packets from that specific address (in this case the node we ping))
	udpAddr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, udpAddr)

	// Send encoded packet to contact

	if err != nil {
		fmt.Println("error sending packet")
	}
	fmt.Println("Sending PING packet to " + contact.Address)
	conn.Write(b)

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
