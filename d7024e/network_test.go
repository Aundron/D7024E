package d7024e

import (
	"fmt"
	"testing"
)

func TestSendPingMessage(t *testing.T) {
	channel := make(chan Packet)

	mockNetwork := NewNetwork(channel)

	mockContact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")

	go mockNetwork.SendPingMessage(&mockContact)

	val := <-channel
	fmt.Println(val.Address)
	fmt.Println(val.Type)
}

func TestListen(t *testing.T) {
	go Listen()
}
