package main

import (
	"D7024E/d7024e"
	"os"
)

func main() {
	id := ""
	ip := "172.19.0.2"

	if len(os.Args) > 0 {
		startNode := d7024e.NewStartUpNode(id, ip)
		network := d7024e.NewNetwork(startNode)
		go network.Listen()
	} else {
		newNode := d7024e.NewKademlia()
		startNodeID := d7024e.NewKademliaID(id)
		network := d7024e.NewNetwork(newNode)
		d7024e.JoinNetwork(ip, startNodeID, newNode, network)
	}

	for {

	}
}
