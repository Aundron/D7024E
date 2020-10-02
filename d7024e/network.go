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
		network.HandleStore(recPacket)
	}
	if recPacket.Type == "FIND_DATA" {
		network.HandleFindData(recPacket, conn, rAddr)
	}
}

func (network *Network) HandleFindNode(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	fmt.Println("Received FIND_NODE packet from ID: " + recPacket.KademliaID.String() + ", IP: " + recPacket.Address)

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

	b := EncodePacket("FOUND_K_NODES", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, byte)

	fmt.Println("Sending FOUND_K_NODES packet to " + recPacket.Address)
	_, err = conn.WriteToUDP(b, rAddr)
	if err != nil {
		fmt.Println("Error writing FOUND_K_NODES packet")
		log.Fatal(err.Error())
	}
}

func (network *Network) HandlePing(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	fmt.Println("Received PING packet from ID: " + recPacket.KademliaID.String() + ", IP: " + recPacket.Address)

	b := EncodePacket("PONG", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, nil)

	fmt.Println("Sending PONG packet to " + recPacket.Address)
	_, err := conn.WriteToUDP(b, rAddr)
	if err != nil {
		fmt.Println("Error writing PONG packet")
		log.Fatal(err)
	}
}

func (network *Network) HandleStore(recPacket *Packet) {
	fmt.Println("Received STORE packet from ID: " + recPacket.KademliaID.String() + ", IP: " + recPacket.Address)

	data := &Data{}
	json.Unmarshal(recPacket.Data, data)

	// Insert the data in our storage
	network.KademliaNode.Storage = append(network.KademliaNode.Storage, data)
	fmt.Println("Saved string: '" + string(data.Value) + "' in storage")
}

func (network *Network) HandleFindData(recPacket *Packet, conn *net.UDPConn, rAddr *net.UDPAddr) {
	fmt.Println("Received FIND_DATA packet from ID: " + recPacket.KademliaID.String() + ", IP: " + recPacket.Address)

	key := ""
	json.Unmarshal(recPacket.Data, &key)
	fmt.Println("DATA KEY: " + key)

	data := network.KademliaNode.SearchStorage(key)

	if data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			fmt.Println("HandleStore: Error marshaling data")
		}

		b := EncodePacket("FOUND_DATA", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, dataByte)

		fmt.Println("Sending FOUND_DATA packet to " + recPacket.Address)
		_, err = conn.WriteToUDP(b, rAddr)
		if err != nil {
			fmt.Println("Error writing FOUND_DATA packet")
			log.Fatal(err.Error())
		}
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

		b := EncodePacket("FOUND_K_NODES_DATA", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, byte)

		fmt.Println("Sending FOUND_K_NODES_DATA packet to " + recPacket.Address)
		_, err = conn.WriteToUDP(b, rAddr)
		if err != nil {
			fmt.Println("Error writing FOUND_K_NODES_DATA packet")
			log.Fatal(err.Error())
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) *Packet {
	// TODO

	// Encode new ping packet
	b := EncodePacket("PING", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, nil)

	// Open connection to contact ip via DialUDP (!!, you specify a remote in DialUDP so you only receive packets from that specific address (in this case the node we ping))
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address+":67")
	if err != nil {
		fmt.Println("SendPingMessage: ResolveUDPAddr error, failed to resolve remote address: " + contact.Address)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("SendPingMessage: DialUDP error, failed dialing " + udpAddr.IP.String() + ":" + strconv.Itoa(udpAddr.Port))
	}
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	defer conn.Close()
	// Send encoded packet to contact

	fmt.Println("Sending PING packet to " + contact.Address)
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("Error writing PING packet")
		fmt.Println(err)
	}

	// Wait for response
	buffer := make([]byte, 8192)
	recPacket := Packet{}
	size, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ping message timed out!")
		recPacket.Type = "TIMEOUT"
	} else {
		json.Unmarshal(buffer[:size], &recPacket)
		fmt.Println("Received PONG response from " + contact.Address)
	}

	// Return response
	return &recPacket
}

func (network *Network) SendFindContactMessage(node *Node, target *KademliaID) *Packet {
	node.Sent = true
	targetByte, err := json.Marshal(target.String())
	if err != nil {
		fmt.Println("Error marshaling data")
	}

	b := EncodePacket("FIND_NODE", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, targetByte)

	udpAddr, err := net.ResolveUDPAddr("udp", node.Contact.Address+":67")
	if err != nil {
		fmt.Println("SendFindContactMessage: ResolveUDPAddr error, failed to resolve address: " + node.Contact.Address)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("SendFindContactMessage: DialUDP error, failed dialing " + udpAddr.IP.String() + ":" + strconv.Itoa(udpAddr.Port))
	}
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	defer conn.Close()
	// Send encoded packet to contact

	if err != nil {
		fmt.Println("error sending packet")
	}
	fmt.Println("Sending FIND_NODE packet to " + node.Contact.Address + " with target " + target.String())
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("Error writing FIND_NODE packet")
	}

	// Wait for response
	buffer := make([]byte, 8192)
	recPacket := Packet{}
	size, err := conn.Read(buffer)
	if err != nil {
		recPacket.Type = "TIMEOUT"
		fmt.Println("FIND_NODE message to " + node.Contact.Address + " timed out!")
	} else {
		node.Visited = true
		json.Unmarshal(buffer[:size], &recPacket)
		fmt.Println("Received FOUND_K_NODES response from " + node.Contact.Address)
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

func (network *Network) SendFindDataMessage(node *Node, key *KademliaID) *Packet {
	node.Sent = true
	keyByte, err := json.Marshal(key.String())
	if err != nil {
		fmt.Println("SendFindDataMessage: Error marshaling data")
	}

	b := EncodePacket("FIND_DATA", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, keyByte)

	udpAddr, err := net.ResolveUDPAddr("udp", node.Contact.Address+":67")
	if err != nil {
		fmt.Println("SendFindDataMessage: ResolveUDPAddr error, failed to resolve address: " + node.Contact.Address)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("SendFindDataMessage: DialUDP error, failed dialing " + udpAddr.IP.String() + ":" + strconv.Itoa(udpAddr.Port))
	}
	conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	defer conn.Close()
	// Send encoded packet to contact

	if err != nil {
		fmt.Println("SendFindDataMessage: error sending packet")
	}
	fmt.Println("Sending FIND_DATA packet to " + node.Contact.Address + " with key " + key.String())
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("Error writing FIND_DATA packet")
	}

	// Wait for response
	buffer := make([]byte, 8192)
	recPacket := Packet{}
	size, err := conn.Read(buffer)
	if err != nil {
		recPacket.Type = "TIMEOUT"
		fmt.Println("FIND_DATA message to " + node.Contact.Address + " timed out!")
	} else {
		node.Visited = true
		json.Unmarshal(buffer[:size], &recPacket)
		if recPacket.Type == "FOUND_DATA" {
			fmt.Println("Received FOUND_DATA response from " + node.Contact.Address)
		} else {
			fmt.Println("Received FOUND_K_NODES_DATA response from " + node.Contact.Address)
		}
	}

	// Return response
	return &recPacket

}

func (network *Network) SendStoreMessage(contact *Contact, data *Data) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Println("SendStoreMessage: Error marshaling data")
	}

	b := EncodePacket("STORE", network.KademliaNode.RoutingTable.me.Address, network.KademliaNode.RoutingTable.me.ID, dataByte)

	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address+":67")
	if err != nil {
		fmt.Println("SendStoreMessage: ResolveUDPAddr error, failed to resolve remote address: " + contact.Address)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("SendStoreMessage: DialUDP error, failed dialing " + udpAddr.IP.String() + ":" + strconv.Itoa(udpAddr.Port))
	}

	defer conn.Close()
	// Send encoded packet to contact

	fmt.Println("Sending STORE packet to " + contact.Address)
	_, err = conn.Write(b)
	if err != nil {
		fmt.Println("Error writing STORE packet")
		fmt.Println(err)
	}
}
