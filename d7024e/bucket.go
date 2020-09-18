package main

import (
	"container/list"
	"fmt"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact, network *Network) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {

			// Get least recent contact

			leastRecent := bucket.list.Back().Value.(Contact)

			// Ping least recently
			response := network.SendPingMessage(&leastRecent)

			// If we get no answer from our ping, we remove the least recent contact and add the new contact to the front of the list
			if response.Type == "TIMEOUT" {
				fmt.Println("PING timeout, replacing contact")
				bucket.list.Remove(bucket.list.Back())
				bucket.list.PushFront(contact)
			} else if response.Type == "PONG" {
				fmt.Println("PONG received, discarding new contact")
				// If least recent node responds, move it to the front and discard new contact
				bucket.list.Remove(bucket.list.Back())
				bucket.list.PushFront(leastRecent)
			} else {
				fmt.Println("ERROR: Unexpected response from PING")
			}

		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
