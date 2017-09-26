package d7024e

import (
	"sync"
//	"sync/atomic"
	"sort"
	//"fmt"
)

const k = 20
const alpha = 3

type Kademlia struct {
	rt      *RoutingTable
	network Net
}

type potentials struct {
	cc CloseContacts
	searched CloseContacts
	cond *sync.Cond
	index int
	activeThreads int
	wg *sync.WaitGroup
}

func NewKademlia(address string, network Net, base *Contact) *Kademlia {
	var c Contact = NewContact(NewRandomKademliaID(), "localghost")
	var rt *RoutingTable = NewRoutingTable(c)
	if base != nil {
		rt.AddContact(*base)
	}
	return &Kademlia{rt, network}
}

func (kademlia *Kademlia) asyncLookup(target *Contact, pot *potentials) {
	defer pot.wg.Done()
	for {
		pot.cond.L.Lock()
/*		for len(pot.cc) == 0 {
			if pot.activeThreads == 0{	//No hope of getting new items.
				pot.cond.L.Unlock()
				return
			}
			pot.cond.Wait()		//Wait until other threads return with new information.
		}*/
		var c CloseContact
		var checked bool = false
		for !checked {
			for len(pot.cc) == 0 {
				if pot.activeThreads == 0 {		//No hope of getting new items.
					pot.cond.L.Unlock()
					return						//No new item we haven't already checked. 
				}
				pot.cond.Wait()		//Wait until other threads return with new information
			}
			c = pot.cc[0]
			pot.cc = pot.cc[1:]
			for i:= 0; i < len(pot.searched); i++ {					//Make sure we're not seaching one we already searched.
				if c.contact.ID == pot.searched[i].contact.ID {
					continue
				}
			}
			checked = true			//Checked through all searched items, didn't find duplicates
		}
		pot.searched = append(pot.searched, c)
		var num int = pot.index
		pot.index ++
		if(num >= k) {
			pot.cond.Broadcast()
			pot.cond.L.Unlock()
			return 
		}
		
		pot.searched = append(pot.searched, c)
		pot.activeThreads ++
		pot.cond.L.Unlock()
			
		var a CloseContacts = kademlia.network.SendFindContactMessage(&c.contact, target.ID)
		//Go through all the results, and spawn routines to add them to RT
		
		pot.cond.L.Lock()
		pot.cc = (append(pot.cc, a...))
		sort.Sort(pot.cc)
		pot.activeThreads --
		pot.cond.Broadcast()
		pot.cond.L.Unlock()
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) *CloseContacts{
	// Step 1. Get k closest to target
	var cc CloseContacts = kademlia.rt.FindClosestContacts(target.ID, alpha)
	// Step 2. See if target exists. If so, return it
	/*for i := 0; i < len(cc); i++ {
		if(cc[i].contact.ID == target.ID){
			return &(cc[i].contact)
		}
	}*/
	// Step 3. If not, send LookupContact to k closest contacts, including returned values, running alpha number of lookups in parallel
	var wg sync.WaitGroup
	var pot *potentials = &potentials{cc, nil, &sync.Cond{L: &sync.Mutex{}}, 0, 0,&wg}
	wg.Add(alpha)
	for i:= 0; i < alpha; i++ {
		go kademlia.asyncLookup(target, pot)
	}
	//wait for result here
	wg.Wait()
	// Step 4. If found, return contact.
	
	var result CloseContacts = append(pot.cc, pot.searched...)
	sort.Sort(result)
	for i:= 0; i < len(result)-1; i ++ {
		if(result[i].contact.ID.Equals(result[i+1].contact.ID)) {
			result = append(result[:i], result[i+1:]...)
			i--
		}
	}
	
	return &result
}

func (kademlia *Kademlia) LookupData(hash string) {
	// Step 1.Look for data in own hashtable. If found, return
	// Step 2. If not, Similar to lookupContact, Send lookupData request to k closest, running alpha number of lookups in parallel
	// (Step 2 makes sense if we call LookupData through console, but not if someone call it on us... Same for LookupData)
	// Step 3.If file found, return it
}

func (kademlia *Kademlia) Store(data []byte) {
	/*
		// Hash data to get handle
		hasher := sha1.New()
		hasher.Write(data)
	*/
	// Store data in own file (I think?)

	// Do lookup on data handle (I think?)
	// Store data in k closest nodes (I think?)
	// return handle
}
