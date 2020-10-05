package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Set logging output
	log.SetOutput(ioutil.Discard)
	//log.SetOutput(os.Stderr)

	scanner := bufio.NewReader(os.Stdin)
	id := "FFFFFFFF00000000000000000000000000000000"
	ip := "172.19.0.2"

	var kademliaNode *Kademlia
	var network *Network

	if len(os.Args) > 1 {
		kademliaNode = NewStartUpNode(id, ip)
		fmt.Println("KademliaID: " + kademliaNode.ID.String())
		fmt.Print("IP: ")
		fmt.Println(kademliaNode.IP)
		network = NewNetwork(kademliaNode)
		go network.SendMessage()
		go network.Listen()
	} else {
		kademliaNode = NewKademlia()
		fmt.Println("KademliaID: " + kademliaNode.ID.String())
		fmt.Print("IP: ")
		fmt.Println(kademliaNode.IP)
		startNodeID := NewKademliaID(id)
		network = NewNetwork(kademliaNode)
		go network.SendMessage()
		go network.Listen()
		JoinNetwork(ip, startNodeID, kademliaNode, network)
	}

	for {
		fmt.Print("$ : ")

		command, _ := scanner.ReadString('\n')
		commandArr := strings.SplitN(strings.TrimSpace(command), " ", 2)

		switch commandArr[0] {
		case "exit":
			os.Exit(0)
		case "put":
			fmt.Println(kademliaNode.Store(commandArr[1], network))
		case "get":
			fmt.Println(kademliaNode.FindValue(commandArr[1], network))
		case "help":
			fmt.Println("Supported commands:")
			fmt.Println("exit - Terminates the node")
			fmt.Println("put <string> - Stores a string in the system")
			fmt.Println("get <hash> - Retrieves a string in the system")
		default:
			fmt.Println("Unknown command " + commandArr[0] + ", the supported commands are:")
			fmt.Println("exit - Terminates the node")
			fmt.Println("put <string> - Stores a string in the system")
			fmt.Println("get <hash> - Retrieves a string in the system")
		}

	}
}
