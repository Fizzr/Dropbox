package d7024e

import (
	"sync"
	"sort"
)

const k = 20
const alpha = 3

type Kademlia struct {
	rt *RoutingTable
	network Net
}

func NewKademlia(address string, network Net, base *Contact) *Kademlia{
	var c Contact = NewContact(NewKademliaID(randomHex(40)), "localghost")
	var rt *RoutingTable = NewRoutingTable(c)
	if(base != nil){
		rt.AddContact(*base)
	}
	return &Kademlia{rt, network}
}

func (kademlia *Kademlia) asyncLookup(target *Contact, potentials *CloseContacts, mutex sync.Mutex) {
	var a CloseContacts = kademlia.network.SendFindContactMessage(target)
	b :=(append(*potentials, a...))
	potentials = &b
	sort.Sort(potentials)	 
}

func (kademlia *Kademlia) LookupContact(target *Contact) *Contact{
	// Step 1. Get k closest to target
	var potentials CloseContacts = kademlia.rt.FindClosestContacts(target.ID, k)
	// Step 2. See if target exists. If so, return it
	for i := 0; i < len(potentials); i++ {
		if(potentials[i].contact.ID == target.ID){
			return &(potentials[i].contact)
		}
	}
	// Step 3. If not, send LookupContact to k closest contacts, including returned values, running alpha number of lookups in parallel
	var mutex sync.Mutex = sync.Mutex{}
	for i:= 0; i < alpha; i++ {
		go kademlia.asyncLookup(target, &potentials, mutex)
	}
	//wait for result here
	// Step 4. If found, return contact.
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// Step 1.Look for data in own hashtable. If found, return
	// Step 2. If not, Similar to lookupContact, Send lookupData request to k closest, running alpha number of lookups in parallel 
	// (Step 2 makes sense if we call LookupData through console, but not if someone call it on us... Same for LookupData)
	// Step 3.If file found, return it
}

func (kademlia *Kademlia) Store(data []byte) {
	// Hash data to get handle
	// Store data in own file (I think?)
	// Do lookup on data handle (I think?)
	// Store data in k closest nodes (I think?)
	// return handle
}
