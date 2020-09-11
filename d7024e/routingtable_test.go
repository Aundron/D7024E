package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {

	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	kademliaNode := NewKademlia(rt)
	network := NewNetwork(kademliaNode)

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"), network)
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"), network)
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"), network)
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"), network)
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"), network)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}

func TestFullKBucket(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8000"))
	kademliaNode := NewKademlia(rt)
	network := NewNetwork(kademliaNode)

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000007"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000008"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000009"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000A"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000B"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000C"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000D"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000E"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000F"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000011"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000012"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000013"), "localhost:8001"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000014"), "localhost:8001"), network)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 21)
	/*for i := range contacts {
		fmt.Println(rt.getBucketIndex(contacts[i].ID))
	}*/

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000015"), "localhost:8001"), network)

	contacts2 := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 21)

	if len(contacts) != len(contacts2) {
		t.Log("Bucket size exceeded")
		t.Fail()
	}
}
