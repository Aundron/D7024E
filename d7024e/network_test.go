package d7024e

/*func TestSendPingMessage(t *testing.T) {
	mockNetwork := NewNetwork()
	mockContact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:6001")
	channel := make(chan []byte)

	go func() {
		fmt.Println("Listening")
		conn, err := net.ListenPacket("udp", mockContact.Address)
		if err != nil {
			t.Error("Error")
		}
		defer conn.Close()

		buffer := make([]byte, 8192)
		size, _, err := conn.ReadFrom(buffer)
		if err != nil {
			t.Error("Error")
		}

		channel <- buffer[:size]
	}()
	go mockNetwork.SendPingMessage(&mockContact)

	buffer := <-channel

	recPacket := Packet{}

	json.Unmarshal(buffer, &recPacket)

	fmt.Println(recPacket)
}*/
