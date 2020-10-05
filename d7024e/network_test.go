package main

import (
	"testing"
)

func TestSendPingMessage(t *testing.T) {
	mockKademlia := NewKademlia()
	mockNetwork := NewNetwork(mockKademlia)
	mockKademliaID := NewRandomKademliaID()
	mockContact := NewContact(mockKademliaID, "localhost")
	channel := make(chan *Packet)
	go func(channel chan *Packet) {
		channel <- mockNetwork.SendPingMessage(&mockContact)
	}(channel)
	pingMessage := <-*mockNetwork.Channel
	if pingMessage.Packet.Type != "PING" {
		t.Error("SendPingMessage sent wrong packet type")
	}

	pingMessage.ReturnChannel <- &Packet{
		Type:       "PONG",
		Address:    mockContact.Address,
		KademliaID: mockContact.ID,
		Data:       nil,
	}
	packet := <-channel
	if packet.Type != "PONG" {
		t.Error("SendPingMessage returned wrong packet type")
	}
}

func TestSendFindContactMessage(t *testing.T) {
	mockKademlia := NewKademlia()
	mockNetwork := NewNetwork(mockKademlia)
	mockContact := NewContact(NewRandomKademliaID(), "localhost")
	target := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	distance := mockContact.ID.CalcDistance(target)
	node := &Node{
		Contact:          &mockContact,
		DistanceToTarget: distance,
		Visited:          false,
		Sent:             false,
	}
	channel := make(chan *Packet)
	go func(channel chan *Packet) {
		channel <- mockNetwork.SendFindContactMessage(node, target)
	}(channel)
	message := <-*mockNetwork.Channel
	if message.Packet.Type != "FIND_NODE" {
		t.Error("SendFindContactMessage sent wrong packet type")
	}

	message.ReturnChannel <- &Packet{
		Type:       "FOUND_K_NODES",
		Address:    mockContact.Address,
		KademliaID: mockContact.ID,
		Data:       nil,
	}

	packet := <-channel
	if packet.Type != "FOUND_K_NODES" {
		t.Error("SendFindContactMessage returned wrong packet type")
	}
}
func TestSendFindDataMessage(t *testing.T) {
	mockKademlia := NewKademlia()
	mockNetwork := NewNetwork(mockKademlia)
	mockContact := NewContact(NewRandomKademliaID(), "localhost")
	target := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	distance := mockContact.ID.CalcDistance(target)
	node := &Node{
		Contact:          &mockContact,
		DistanceToTarget: distance,
		Visited:          false,
		Sent:             false,
	}
	channel := make(chan *Packet)
	go func(channel chan *Packet) {
		channel <- mockNetwork.SendFindDataMessage(node, target)
	}(channel)
	message := <-*mockNetwork.Channel
	if message.Packet.Type != "FIND_DATA" {
		t.Error("SendFindDataMessage sent wrong packet type")
	}

	message.ReturnChannel <- &Packet{
		Type:       "FOUND_DATA",
		Address:    mockContact.Address,
		KademliaID: mockContact.ID,
		Data:       nil,
	}

	packet := <-channel
	if packet.Type != "FOUND_DATA" {
		t.Error("SendFindDataMessage returned wrong packet type")
	}
}

func TestEncodeDecodePacket(t *testing.T) {
	packetType := "HEJ"
	address := "localhost"
	id := NewRandomKademliaID()
	encoded := EncodePacket(packetType, address, id, nil)
	decoded := DecodePacket(encoded)
	if decoded.Type != packetType {
		t.Error("Encode/Decode failed.")
	}
}
