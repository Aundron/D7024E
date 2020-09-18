package d7024e

import (
	"testing"
)

func TestInsertToList(t *testing.T) {
	NodeList := NodeList{}

	mockContact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000")
	targetContact := NewContact(NewKademliaID("1111111100000000000000000000000000000005"), "localhost:8000")

	NodeList.InsertNode(&mockContact, &targetContact)

	if len(NodeList.Nodes) != 1 {
		t.Error("Fail inserting to empty list")
	}

	NodeList.InsertNode(&mockContact, &targetContact)

	if len(NodeList.Nodes) != 1 {
		t.Error("Fail duplicate insertion")
	}

	mockContact2 := NewContact(NewKademliaID("1111111100000000000000000000000000000004"), "localhost:8000")
	NodeList.InsertNode(&mockContact2, &targetContact)

	if !NodeList.Nodes[0].Contact.ID.Equals(mockContact2.ID) {
		t.Error("Fail wrong order in list")
	}

}
