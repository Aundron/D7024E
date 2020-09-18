package main

import (
	"encoding/json"
	"net"
)

const alpha = 3

type Kademlia struct {
	id           *KademliaID
	ip           string
	RoutingTable *RoutingTable
}

type Node struct {
	Contact          *Contact
	DistanceToTarget *KademliaID
	Visited          bool
	TimedOut         bool
}

type NodeList struct {
	Nodes []*Node
}

func NewKademlia() *Kademlia {
	kademliaNode := &Kademlia{}
	kademliaNode.id = NewRandomKademliaID()
	//felhantering?
	kademliaNode.ip = GetIPAddress()
	kademliaNode.RoutingTable = NewRoutingTable(NewContact(kademliaNode.id, kademliaNode.ip))
	return kademliaNode
}

func (nodeList *NodeList) InsertNode(newContact *Contact, target *KademliaID) *Node {
	distance := newContact.ID.CalcDistance(target)
	NewNode := &Node{
		Contact:          newContact,
		DistanceToTarget: distance,
		Visited:          false,
		TimedOut:         false,
	}

	for i, elem := range (*nodeList).Nodes {
		if newContact.ID.Equals(elem.Contact.ID) {
			return nil
		}
		if NewNode.DistanceToTarget.Less(elem.DistanceToTarget) {
			fst := (*nodeList).Nodes[:i]
			lst := (*nodeList).Nodes[i:]
			mid := []*Node{NewNode}
			(*nodeList).Nodes = append(append(fst, mid...), lst...)
			return NewNode
		}
	}
	(*nodeList).Nodes = []*Node{NewNode}
	return NewNode
}

func (kademlia *Kademlia) LookupContact(target *Contact, network *Network) {
	NodeList := NodeList{}

	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, bucketSize)

	channel := make(chan *Packet, alpha)

	for i := range closestContacts {
		NodeList.InsertNode(&closestContacts[i], target.ID)
	}

	kademlia.LookupContactRec(target, network, &NodeList, &channel)

}

func (kademlia *Kademlia) LookupContactRec(target *Contact, network *Network, nodeList *NodeList, channel *(chan *Packet)) {
	AlphaUnvisited := nodeList.GetAlphaUnvisited()
	for _, elem := range AlphaUnvisited {
		go func() {
			*channel <- network.SendFindContactMessage(elem, target.ID, nodeList)
		}()
	}

	for !nodeList.CheckIfDone() {
		packet := <-*channel

		if packet.Type == "FOUND_K_NODES" {
			contactList := []Contact{}
			json.Unmarshal(packet.Data, &contactList)
			for j, elem := range contactList {
				newContact := NewContact(contactList[j].ID, contactList[j].Address)
				network.KademliaNode.RoutingTable.AddContact(newContact, network)
				nodeList.InsertNode(&elem, target.ID)
			}
			go kademlia.LookupContactRec(target, network, nodeList, channel)
		}

	}
}

func (nodeList *NodeList) GetAlphaUnvisited() []*Node {
	UnvisitedNodes := []*Node{}
	for _, elem := range nodeList.Nodes {
		if elem.Visited == false && elem.TimedOut == false {
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
		} else if elem.Visited == false && elem.TimedOut == false {
			return false
		} else if elem.Visited == true && elem.TimedOut == false {
			count++
		}
	}
	return true
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func JoinNetwork(startNodeIP string, startNodeID *KademliaID, newNode *Kademlia, network *Network) {
	startNode := NewContact(startNodeID, startNodeIP)
	go network.Listen()
	//new node insert start node into one of its k-buckets
	newNode.RoutingTable.AddContact(startNode, network)

	//TODO: newNode.LookupContact(newNode.id, )
	//new node performs a node lookup of its own ID against the start node (the only other node it knows)

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
