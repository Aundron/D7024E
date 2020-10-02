package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Packet struct {
	Type       string
	Address    string
	KademliaID *KademliaID
	Data       []byte
}

type Data struct {
	Key   *KademliaID
	Value []byte
}

type Message struct {
	Packet        *Packet
	IPAddress     string
	ReturnChannel chan *Packet
}

type Network struct {
	KademliaNode *Kademlia
	Channel      *(chan *Message)
}

func NewNetwork(kademliaNode *Kademlia) *Network {
	newNetwork := &Network{}
	newNetwork.KademliaNode = kademliaNode
	channel := make(chan *Message)
	newNetwork.Channel = &channel

	return newNetwork
}

func (network *Network) Listen() {
	// TODO
	fmt.Println("Listening")

	udpAddr, err := net.ResolveUDPAddr("udp", GetIPAddress()+":67")
	if err != nil {
		fmt.Println("HandleFindNode: ResolveUDPAddr error, failed to resolve address: " + GetIPAddress())
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 8192)
		size, rAddr, _ := conn.ReadFromUDP(buffer)

		recPacket := DecodePacket(buffer[:size])

		go network.HandleRequest(&recPacket, conn, rAddr)
	}
}

func (network *Network) HandleRequest(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	fmt.Println("Received " + recPacket.Type + " packet from ID: " + recPacket.KademliaID.String() + ", IP: " + recPacket.Address)

	// Update routingtable
	newContact := NewContact(recPacket.KademliaID, recPacket.Address)
	network.KademliaNode.RoutingTable.AddContact(newContact, network)

	if recPacket.Type == "PING" {
		network.HandlePing(recPacket, conn, rAddr)
	}
	if recPacket.Type == "FIND_NODE" {
		network.HandleFindNode(recPacket, conn, rAddr)
	}
	if recPacket.Type == "STORE" {
		network.HandleStore(recPacket, conn, rAddr)
	}
	if recPacket.Type == "FIND_DATA" {
		network.HandleFindData(recPacket, conn, rAddr)
	}
}

func (network *Network) SendResponse(packetType string, data []byte, conn *net.UDPConn, rAddr *net.UDPAddr) {
	b := EncodePacket(packetType, network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, data)

	fmt.Println("Sending " + packetType + " packet to " + rAddr.IP.String())
	_, err := conn.WriteToUDP(b, rAddr)
	if err != nil {
		fmt.Println("Error writing FOUND_K_NODES packet")
		log.Fatal(err.Error())
	}
}

func (network *Network) HandleFindNode(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	// Get K closest nodes to target node
	kademliaID := ""
	json.Unmarshal(recPacket.Data, &kademliaID)

	targetNode := NewKademliaID(kademliaID)
	closestContacts := network.KademliaNode.RoutingTable.FindClosestContacts(targetNode, bucketSize)

	byte, err := json.Marshal(closestContacts)
	if err != nil {
		fmt.Println("error marshaling closestContacts")
		log.Fatal(err)
	}

	network.SendResponse("FOUND_K_NODES", byte, conn, rAddr)
}

func (network *Network) HandlePing(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	network.SendResponse("PONG", nil, conn, rAddr)
}

func (network *Network) HandleStore(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	data := &Data{}
	json.Unmarshal(recPacket.Data, data)

	// Insert the data in our storage
	network.KademliaNode.Storage = append(network.KademliaNode.Storage, data)
	fmt.Println("Saved string: '" + string(data.Value) + "' in storage")

	network.SendResponse("STORED", nil, conn, rAddr)
}

func (network *Network) HandleFindData(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	key := ""
	json.Unmarshal(recPacket.Data, &key)
	fmt.Println("DATA KEY: " + key)

	data := network.KademliaNode.SearchStorage(key)

	if data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			fmt.Println("HandleStore: Error marshaling data")
		}

		network.SendResponse("FOUND_DATA", dataByte, conn, rAddr)
	} else {
		// Get K closest nodes to target node
		kademliaID := ""
		json.Unmarshal(recPacket.Data, &kademliaID)

		keyNode := NewKademliaID(kademliaID)
		closestContacts := network.KademliaNode.RoutingTable.FindClosestContacts(keyNode, bucketSize)

		byte, err := json.Marshal(closestContacts)
		if err != nil {
			fmt.Println("HandleFindData: error marshaling closestContacts")
			log.Fatal(err)
		}

		network.SendResponse("FOUND_K_NODES_DATA", byte, conn, rAddr)
	}
}

func (network *Network) SendMessage() {

	for {
		message := <-*network.Channel

		go func(message *Message) {
			b := EncodePacket(message.Packet.Type, network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, message.Packet.Data)

			// Open connection to contact ip via DialUDP (!!, you specify a remote in DialUDP so you only receive packets from that specific address (in this case the node we ping))
			udpAddr, err := net.ResolveUDPAddr("udp", message.IPAddress+":67")
			if err != nil {
				fmt.Println("SendMessage: ResolveUDPAddr error, failed to resolve remote address: " + message.IPAddress)
			}

			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				fmt.Println("SendMessage: DialUDP error, failed dialing " + udpAddr.IP.String() + ":" + strconv.Itoa(udpAddr.Port))
			}
			conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
			defer conn.Close()
			// Send encoded packet to contact

			fmt.Println("Sending " + message.Packet.Type + " packet to " + message.IPAddress)
			_, err = conn.Write(b)
			if err != nil {
				fmt.Println("Error writing " + message.Packet.Type + " packet")
				fmt.Println(err)
			}

			// Wait for response
			buffer := make([]byte, 8192)
			recPacket := Packet{}
			size, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("TIMEOUT!")
				fmt.Println(message.Packet.Type + " message timed out!")
				recPacket.Type = "TIMEOUT"
			} else {
				json.Unmarshal(buffer[:size], &recPacket)
				fmt.Println("Received " + recPacket.Type + " response from " + message.IPAddress)
			}

			message.ReturnChannel <- &recPacket
		}(message)
	}

}

func (network *Network) SendPingMessage(contact *Contact) *Packet {
	// TODO
	pingPacket := &Packet{
		Type:       "PING",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       nil,
	}
	returnChannel := make(chan *Packet)
	pingMessage := &Message{
		Packet:        pingPacket,
		IPAddress:     contact.Address,
		ReturnChannel: returnChannel,
	}

	*network.Channel <- pingMessage
	return <-returnChannel
}

func (network *Network) SendFindContactMessage(node *Node, target *KademliaID) *Packet {
	if node.Sent == true {
		fmt.Println("Already sent")
		recPacket := Packet{}
		recPacket.Type = "TIMEOUT"
		return &recPacket
	}
	node.Sent = true
	targetByte, err := json.Marshal(target.String())
	if err != nil {
		fmt.Println("Error marshaling data")
	}

	findContactPacket := &Packet{
		Type:       "FIND_NODE",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       targetByte,
	}
	returnChannel := make(chan *Packet)
	findContactMessage := &Message{
		Packet:        findContactPacket,
		IPAddress:     node.Contact.Address,
		ReturnChannel: returnChannel,
	}
	*network.Channel <- findContactMessage
	responsePacket := <-returnChannel

	if responsePacket.Type == "FOUND_K_NODES" {
		node.Visited = true
	}

	return responsePacket
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

/*func EncodePacket2(packet *Packet) []byte {
	b, err := json.Marshal(packet)
	if err != nil {
		fmt.Println("error")
	}
	return b
}*/

func (network *Network) SendFindDataMessage(node *Node, key *KademliaID) *Packet {
	if node.Sent == true {
		fmt.Println("Already sent")
		recPacket := Packet{}
		recPacket.Type = "TIMEOUT"
		return &recPacket
	}
	node.Sent = true
	keyByte, err := json.Marshal(key.String())
	if err != nil {
		fmt.Println("SendFindDataMessage: Error marshaling data")
	}

	findDataPacket := &Packet{
		Type:       "FIND_DATA",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       keyByte,
	}
	returnChannel := make(chan *Packet)
	findDataMessage := &Message{
		Packet:        findDataPacket,
		IPAddress:     node.Contact.Address,
		ReturnChannel: returnChannel,
	}
	*network.Channel <- findDataMessage
	responsePacket := <-returnChannel

	if responsePacket.Type != "TIMEOUT" {
		node.Visited = true
	}

	return responsePacket
}

func (network *Network) SendStoreMessage(contact *Contact, data *Data) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Println("SendStoreMessage: Error marshaling data")
	}

	findStorePacket := &Packet{
		Type:       "STORE",
		Address:    network.KademliaNode.RoutingTable.me.Address,
		KademliaID: network.KademliaNode.RoutingTable.me.ID,
		Data:       dataByte,
	}
	returnChannel := make(chan *Packet)
	findStoreMessage := &Message{
		Packet:        findStorePacket,
		IPAddress:     contact.Address,
		ReturnChannel: returnChannel,
	}
	*network.Channel <- findStoreMessage
	<-returnChannel
}
