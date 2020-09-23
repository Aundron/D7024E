package main

import (
	"testing"
)

/*func TestRoutingTable(t *testing.T) {
	kademliaNode := NewKademlia()
	network := NewNetwork(kademliaNode)
	rt := kademliaNode.RoutingTable

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

	for i := range rt.buckets {
		fmt.Println(rt.buckets[i].Len())
	}
	fmt.Println(rt.getBucketIndex(NewKademliaID("2111111400000000000000000000000000000000")))
}*/

func TestFullKBucket(t *testing.T) {
	kademliaNode := NewKademlia()
	network := NewNetwork(kademliaNode)
	rt := kademliaNode.RoutingTable

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000007"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000008"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000009"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000A"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000B"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000C"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000D"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000E"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF0000000000000000000000000000000F"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000011"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000012"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000013"), "localhost"), network)
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000014"), "localhost"), network)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 21)
	/*for i := range contacts {
		fmt.Println(rt.getBucketIndex(contacts[i].ID))
	}*/

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000015"), "localhost"), network)

	contacts2 := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 21)

	if len(contacts) != len(contacts2) {
		t.Log("Bucket size exceeded")
		t.Fail()
	}
}
