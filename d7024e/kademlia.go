package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

const alpha = 3

type Kademlia struct {
	ID           *KademliaID
	IP           string
	RoutingTable *RoutingTable
	Storage      []*Data
}

type Node struct {
	Contact          *Contact
	DistanceToTarget *KademliaID
	Visited          bool
	Sent             bool
}

type NodeList struct {
	Nodes []*Node
	mux   sync.Mutex
}

func NewStartUpNode(id string, ip string) *Kademlia {
	startUpNode := &Kademlia{}
	startUpNode.ID = NewKademliaID(id)
	startUpNode.IP = ip
	startUpNode.RoutingTable = NewRoutingTable(NewContact(startUpNode.ID, startUpNode.IP))
	return startUpNode
}

func NewKademlia() *Kademlia {
	kademliaNode := &Kademlia{}
	kademliaNode.ID = NewRandomKademliaID()
	//felhantering?
	kademliaNode.IP = GetIPAddress()
	kademliaNode.RoutingTable = NewRoutingTable(NewContact(kademliaNode.ID, kademliaNode.IP))
	return kademliaNode
}

func (nodeList *NodeList) InsertNode(newContact Contact, target *KademliaID) *Node {
	// Never insert yourself into the nodeList
	nodeList.mux.Lock()
	defer nodeList.mux.Unlock()
	if newContact.Address == GetIPAddress() {
		return nil
	}

	distance := newContact.ID.CalcDistance(target)
	NewNode := &Node{
		Contact:          &newContact,
		DistanceToTarget: distance,
		Visited:          false,
		Sent:             false,
	}

	for i, elem := range (*nodeList).Nodes {
		if newContact.ID.Equals(elem.Contact.ID) {
			//fmt.Println("Found same object in nodeList, skipping")
			return nil
		}
		if NewNode.DistanceToTarget.Less(elem.DistanceToTarget) {
			// "Simple" and "logical" from: https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
			nodeList.Nodes = append(nodeList.Nodes, &Node{})
			copy(nodeList.Nodes[i+1:], nodeList.Nodes[i:])
			nodeList.Nodes[i] = NewNode
			return NewNode
		}
	}
	(*nodeList).Nodes = append((*nodeList).Nodes, NewNode)
	return NewNode
}

func (kademlia *Kademlia) LookupContact(target *KademliaID, network *Network) *NodeList {
	NodeList := NodeList{}

	closestContacts := kademlia.RoutingTable.FindClosestContacts(target, bucketSize)
	//fmt.Print("LookupContact: Closest contacts: ")
	//fmt.Println(closestContacts)

	channel := make(chan *Packet, alpha)

	for i := range closestContacts {
		NodeList.InsertNode(closestContacts[i], target)
	}

	nl := kademlia.LookupContactRec(target, network, &NodeList, &channel)
	fmt.Println("LookupContact DONE!")
	return nl

}

func (kademlia *Kademlia) LookupContactRec(target *KademliaID, network *Network, nodeList *NodeList, channel *(chan *Packet)) *NodeList {
	/*fmt.Println("Nodelist before sending: ")
	for _, elem := range nodeList.Nodes {
		fmt.Println(elem)
	}*/
	AlphaUnvisited := nodeList.GetAlphaUnvisited()
	//fmt.Println("AlphaUnvisited:")
	for i := range AlphaUnvisited {
		//fmt.Println(AlphaUnvisited[i])
		go func(i int) {
			*channel <- network.SendFindContactMessage(AlphaUnvisited[i], target)
		}(i)
	}

	for range AlphaUnvisited {
		packet := <-*channel

		if packet.Type == "FOUND_K_NODES" {
			contactList := []Contact{}
			json.Unmarshal(packet.Data, &contactList)
			//fmt.Println(contactList)
			for _, elem := range contactList {
				newContact := NewContact(elem.ID, elem.Address)
				network.KademliaNode.RoutingTable.AddContact(newContact, network)
				nodeList.InsertNode(elem, target)
			}
		}

		if !nodeList.CheckIfDone() {
			kademlia.LookupContactRec(target, network, nodeList, channel)
		}

	}
	return nodeList
}

func (kademlia *Kademlia) LookupData(key *KademliaID, network *Network) *Data {
	NodeList := NodeList{}

	closestContacts := kademlia.RoutingTable.FindClosestContacts(key, bucketSize)
	//fmt.Print("LookupContact: Closest contacts: ")
	//fmt.Println(closestContacts)

	channel := make(chan *Packet, alpha)

	for i := range closestContacts {
		NodeList.InsertNode(closestContacts[i], key)
	}

	return kademlia.LookupDataRec(key, network, &NodeList, &channel)

}

func (kademlia *Kademlia) LookupDataRec(target *KademliaID, network *Network, nodeList *NodeList, channel *(chan *Packet)) *Data {
	/*fmt.Println("Nodelist before sending: ")
	for _, elem := range nodeList.Nodes {
		fmt.Println(elem)
	}*/
	AlphaUnvisited := nodeList.GetAlphaUnvisited()
	//fmt.Println("AlphaUnvisited:")
	for i := range AlphaUnvisited {
		//fmt.Println(AlphaUnvisited[i])
		go func(i int) {
			*channel <- network.SendFindDataMessage(AlphaUnvisited[i], target)
		}(i)
	}

	for range AlphaUnvisited {
		packet := <-*channel

		if packet.Type == "FOUND_K_NODES_DATA" {
			contactList := []Contact{}
			json.Unmarshal(packet.Data, &contactList)
			//fmt.Println(contactList)
			for _, elem := range contactList {
				newContact := NewContact(elem.ID, elem.Address)
				network.KademliaNode.RoutingTable.AddContact(newContact, network)
				nodeList.InsertNode(elem, target)
			}
		} else if packet.Type == "FOUND_DATA" {
			data := &Data{}
			json.Unmarshal(packet.Data, data)
			return data
		}

		if !nodeList.CheckIfDone() {
			kademlia.LookupDataRec(target, network, nodeList, channel)
		}

	}
	return nil
}

func (nodeList *NodeList) GetAlphaUnvisited() []*Node {
	UnvisitedNodes := []*Node{}
	for _, elem := range nodeList.Nodes {
		if elem.Visited == false && elem.Sent == false {
			UnvisitedNodes = append(UnvisitedNodes, elem)
			if len(UnvisitedNodes) == alpha {
				break
			}
		}
	}
	return UnvisitedNodes
}

func (nodeList *NodeList) CheckIfDone() bool {
	count := 0
	for _, elem := range nodeList.Nodes {
		if count >= bucketSize {
			return true
		} else if elem.Visited == false && elem.Sent == false {
			return false
		} else if elem.Visited == true {
			count++
		}
	}
	return true
}

func (kademlia *Kademlia) FindValue(hash string, network *Network) string {

	data := kademlia.SearchStorage(hash)
	if data != nil {
		return string(data.Value)
	}

	key := NewKademliaID(hash)
	data = kademlia.LookupData(key, network)
	if data == nil {
		return "Data object not found."
	}

	return string(data.Value)
}

func (kademlia *Kademlia) Store(data string, network *Network) string {
	// Hash the data to get a key
	// https://gobyexample.com/sha1-hashes
	hash := sha1.New()
	hash.Write([]byte(data))
	key := hex.EncodeToString(hash.Sum(nil))
	newData := &Data{
		Key:   NewKademliaID(key),
		Value: []byte(data),
	}

	// Find K closest known nodes to Key
	nodeList := kademlia.LookupContact(newData.Key, network)
	sentCount := 0
	for _, elem := range nodeList.Nodes {
		if elem.Visited == true {
			go network.SendStoreMessage(elem.Contact, newData)
			sentCount++
		}
		if sentCount == bucketSize {
			break
		}
	}
	return key
}

func JoinNetwork(startNodeIP string, startNodeID *KademliaID, newNode *Kademlia, network *Network) {
	startNode := NewContact(startNodeID, startNodeIP)
	//new node insert start node into one of its k-buckets
	newNode.RoutingTable.AddContact(startNode, network)
	newNodeContact := NewContact(newNode.ID, newNode.IP)
	//new node performs a node lookup of its own ID against the start node (the only other node it knows)
	newNode.LookupContact(newNodeContact.ID, network)
	//new node refreshes all k-buckets further away than the k-bucket the start node falls in
	//(This refresh is just a lookup of a random key that is within that k-bucket range)
}

/*func (kademlia *Kademlia) refresh() {

}*/

//example from: https://play.golang.org/p/BDt3qEQ_2H
func GetIPAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	return ""
}

func (kademlia *Kademlia) SearchStorage(hash string) *Data {
	key := NewKademliaID(hash)

	for _, elem := range kademlia.Storage {
		if elem.Key.Equals(key) {
			fmt.Println("SearchStorage: Found data '" + string(elem.Value) + "'")
			return elem
		}
	}
	return nil
}
