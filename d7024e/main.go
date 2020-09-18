package main

import (
	"os"
)

func main() {
	id := "FFFFFFFF00000000000000000000000000000000"
	ip := "172.19.0.2"

	if len(os.Args) > 1 {
		startNode := NewStartUpNode(id, ip)
		network := NewNetwork(startNode)
		go network.Listen()
	} else {
		newNode := NewKademlia()
		startNodeID := NewKademliaID(id)
		network := NewNetwork(newNode)
		JoinNetwork(ip, startNodeID, newNode, network)
	}

	for {

	}
}
