package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Packet struct {
	Type       string
	Address    string
	KademliaID *KademliaID
	Data       string
}

type Network struct {
	KademliaNode *Kademlia
}

func NewNetwork(kademliaNode *Kademlia) *Network {
	newNetwork := &Network{}
	newNetwork.KademliaNode = kademliaNode

	return newNetwork
}

func (network *Network) Listen() {
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

		network.HandleRequest(&recPacket)
	}
}

func (network *Network) HandleRequest(recPacket *Packet) {

	if recPacket.Type == "PING" {
		network.HandlePing(recPacket)
	}
	// 	if ping
	// 		HandlePing()
	// 	if store
	// 		HandleStore()
	// 	...
	//
	// SendResponse()
}

func (network *Network) HandlePing(recPacket *Packet) {

	// TODO: Update routing tables
	newContact := NewContact(recPacket.KademliaID, recPacket.Address)
	network.KademliaNode.RoutingTable.AddContact(newContact, network)

	pongPacket := Packet{
		Type:       "PONG",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       "Test",
	}

	b, err := json.Marshal(pongPacket)
	if err != nil {
		fmt.Println("error")
	}

	udpAddr, _ := net.ResolveUDPAddr("udp", recPacket.Address)
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Println("error sending packet")
	}
	fmt.Println("Sending PONG packet to " + recPacket.Address)
	conn.Write(b)
}

func (network *Network) SendPingMessage(contact *Contact) *Packet {
	// TODO

	// Encode new ping packet
	packet := Packet{
		Type:       "PING",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       "Test",
	}

	b, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("error")
	}

	// Open connection to contact ip via DialUDP (!!, you specify a remote in DialUDP so you only receive packets from that specific address (in this case the node we ping))
	udpAddr, _ := net.ResolveUDPAddr("udp", contact.Address)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	// Send encoded packet to contact

	if err != nil {
		fmt.Println("error sending packet")
	}
	fmt.Println("Sending PING packet to " + contact.Address)
	conn.Write(b)

	// Wait for response
	buffer := make([]byte, 8192)
	recPacket := Packet{}
	size, err := conn.Read(buffer)
	if err != nil {
		recPacket.Type = "TIMEOUT"
	} else {
		json.Unmarshal(buffer[:size], &recPacket)
	}

	// Return response
	return &recPacket
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
