package main

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
	Data       []byte
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

		recPacket := DecodePacket(buffer[:size])

		fmt.Println(recPacket)

		network.HandleRequest(&recPacket)
	}
}

func (network *Network) HandleRequest(recPacket *Packet) {

	if recPacket.Type == "PING" {
		network.HandlePing(recPacket)
	}
	if recPacket.Type == "FIND_NODE" {
		network.HandleFindNode(recPacket)
	}
	// 	if ping
	// 		HandlePing()
	// 	if store
	// 		HandleStore()
	// 	...
	//
	// SendResponse()
}

func (network *Network) HandleFindNode(recPacket *Packet) {
	newContact := NewContact(recPacket.KademliaID, recPacket.Address)
	network.KademliaNode.RoutingTable.AddContact(newContact, network)

	// Get K closest nodes to target node
	kademliaID := ""
	json.Unmarshal(recPacket.Data, &kademliaID)
	targetNode := NewKademliaID(kademliaID)
	closestContacts := network.KademliaNode.RoutingTable.FindClosestContacts(targetNode, bucketSize)

	byte, err := json.Marshal(closestContacts)
	if err != nil {
		fmt.Println("error")
	}

	b := EncodePacket("FOUND_K_NODES", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, byte)

	udpAddr, _ := net.ResolveUDPAddr("udp", recPacket.Address)
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Println("error sending packet")
	}

	fmt.Println("Sending FOUND_K_NODES packet to " + recPacket.Address)
	conn.Write(b)
}

func (network *Network) HandlePing(recPacket *Packet) {

	// TODO: Update routing tables
	newContact := NewContact(recPacket.KademliaID, recPacket.Address)
	network.KademliaNode.RoutingTable.AddContact(newContact, network)

	pongPacket := Packet{
		Type:       "PONG",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       nil,
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
		Data:       nil,
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

func (network *Network) SendFindContactMessage(node *Node, target *KademliaID, nodeList *NodeList) *Packet {

	targetByte, err := json.Marshal(target)
	if err != nil {
		fmt.Println("error")
	}

	b := EncodePacket("NODE_LOOKUP", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, targetByte)

	udpAddr, _ := net.ResolveUDPAddr("udp", node.Contact.Address)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	// Send encoded packet to contact

	if err != nil {
		fmt.Println("error sending packet")
	}
	fmt.Println("Sending NODE_LOOKUP packet to " + node.Contact.Address)
	conn.Write(b)

	// Wait for response
	buffer := make([]byte, 8192)
	recPacket := Packet{}
	size, err := conn.Read(buffer)
	if err != nil {
		recPacket.Type = "TIMEOUT"
		node.TimedOut = true
	} else {
		node.Visited = true
		json.Unmarshal(buffer[:size], &recPacket)
	}

	// Return response
	return &recPacket

}

func DecodePacket(recByte []byte) Packet {
	recPacket := Packet{}
	json.Unmarshal(recByte, &recPacket)
	return recPacket
}

func EncodePacket(packetType string, address string, kademliaID *KademliaID, data []byte) []byte {
	packet := Packet{
		Type:       packetType,
		Address:    address,
		KademliaID: kademliaID,
		Data:       data,
	}
	b, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("error")
	}
	return b
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
