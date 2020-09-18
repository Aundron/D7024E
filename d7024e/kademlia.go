package main

import "encoding/json"

const alpha = 3

type Kademlia struct {
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

func NewKademlia(routingTable *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.RoutingTable = routingTable
	return kademlia
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
