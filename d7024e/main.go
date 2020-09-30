package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	id := "FFFFFFFF00000000000000000000000000000000"
	ip := "172.19.0.2"

	if len(os.Args) > 1 {
		startNode := NewStartUpNode(id, ip)
		fmt.Println("KademliaID: " + startNode.ID.String())
		fmt.Print("IP: ")
		fmt.Println(startNode.IP)
		network := NewNetwork(startNode)
		go network.Listen()
		time.Sleep(25 * time.Second)
		key := startNode.Store("hej", network)
		time.Sleep(5 * time.Second)
		fmt.Println("Found Data: " + startNode.FindValue(key, network))
	} else {
		newNode := NewKademlia()
		fmt.Println("KademliaID: " + newNode.ID.String())
		fmt.Print("IP: ")
		fmt.Println(newNode.IP)
		startNodeID := NewKademliaID(id)
		network := NewNetwork(newNode)
		go network.Listen()
		JoinNetwork(ip, startNodeID, newNode, network)
	}

	for {

	}
}
