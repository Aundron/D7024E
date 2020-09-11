package d7024e

type Kademlia struct {
	RoutingTable *RoutingTable
}

func NewKademlia(routingTable *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.RoutingTable = routingTable
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
