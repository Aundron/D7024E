package main

import (
	"fmt"
	"testing"
)

func TestInsertToList(t *testing.T) {
	NodeList := NodeList{}

	mockContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000")
	targetContact := NewContact(NewKademliaID("1111111100000000000000000000000000000005"), "localhost:8000")

	NodeList.InsertNode(mockContact, targetContact.ID)

	if len(NodeList.Nodes) != 1 {
		t.Error("Fail inserting to empty list")
	}

	NodeList.InsertNode(mockContact, targetContact.ID)

	if len(NodeList.Nodes) != 1 {
		t.Error("Fail duplicate insertion")
	}

	mockContact2 := NewContact(NewKademliaID("1111111100000000000000000000000000000004"), "localhost:8000")
	NodeList.InsertNode(mockContact2, targetContact.ID)

	if !NodeList.Nodes[0].Contact.ID.Equals(mockContact2.ID) {
		t.Error("Fail wrong order in list")
	}
}

func TestAlphaUnvisited(t *testing.T) {
	NodeList := NodeList{}

	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))

	AlphaUnvisited := NodeList.GetAlphaUnvisited()

	if len(AlphaUnvisited) > alpha {
		t.Error("AlphaUnvisited returned too many nodes")
	}

	for _, elem := range AlphaUnvisited {
		if elem.Visited != false || elem.Sent != false {
			t.Error("AlphaUnvisited returned visited or sent nodes")
		}

		elem.Sent = true
		elem.Visited = true
	}

	AlphaUnvisited = NodeList.GetAlphaUnvisited()

	if len(AlphaUnvisited) != 1 {
		t.Error("AlphaUnvisited returned too many nodes")
	}
}

func TestCheckIfDone(t *testing.T) {
	NodeList := NodeList{}

	// Insert 22 nodes
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000007"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000008"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000009"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000010"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000011"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000012"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000013"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000014"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000015"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000016"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000017"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000018"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000019"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000020"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))
	NodeList.InsertNode(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000021"), "localhost"), NewKademliaID("FFFFFFFF00000000000000000000000000000005"))

	// Set K closest nodes to visited
	for i := 0; i < bucketSize; i++ {
		NodeList.Nodes[i].Sent = true
		NodeList.Nodes[i].Visited = true
	}

	// CheckIfDone should now return true since K closest nodes has been visited
	if !NodeList.CheckIfDone() {
		t.Error("CheckIfDone returns false, expected true")
	}

	// Set fourth closest to a timeout
	NodeList.Nodes[3].Visited = false

	// CheckIfDone should now return false since only K-1 closest nodes have been visited
	if NodeList.CheckIfDone() {
		t.Error("CheckIfDone returns true, expected false")
	}

	// Set 21th element to visited
	NodeList.Nodes[20].Sent = true
	NodeList.Nodes[20].Visited = true

	// CheckIfDone should now return true again
	if !NodeList.CheckIfDone() {
		t.Error("CheckIfDone returns false, expected true")
	}

}

func TestGetIPAddress(t *testing.T) {
	ip := GetIPAddress()
	fmt.Println(ip)
}
