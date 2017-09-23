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
	rt *RoutingTable
	network Net
}

type asyncStruct struct {
	cc CloseContacts
	searched CloseContacts
	cond *sync.Cond
	index int
	activeThreads int
	wg *sync.WaitGroup
}

func NewAsyncStruct(base CloseContacts) *asyncStruct{
	var wg sync.WaitGroup
	return &asyncStruct{base, nil, &sync.Cond{L: &sync.Mutex{}}, 0, 0, &wg}
}

func (as *asyncStruct) getNext() *CloseContact {
	as.cond.L.Lock()
	var c CloseContact
	var checked bool = false
	for !checked {
		for len(as.cc) == 0 {
			if as.activeThreads == 0 {		//No hope of getting new items.
				as.cond.L.Unlock()
				return nil					//No new item we haven't already checked. 
			}
			as.cond.Wait()		//Wait until other threads return with new information
		}
		c = as.cc[0]
		as.cc = as.cc[1:]
		for i:= 0; i < len(as.searched); i++ {					//Make sure we're not seaching one we already searched.
			if c.contact.ID == as.searched[i].contact.ID {
				continue
			}
		}
		checked = true			//Checked through all searched items, didn't find duplicates
	}
	as.searched = append(as.searched, c)
	var num int = as.index
	as.index ++
	if(num >= k) {
		as.cond.Broadcast()
		as.cond.L.Unlock()
		return nil			//Already ran k times
	}
	
	as.searched = append(as.searched, c)
	as.activeThreads ++
	as.cond.L.Unlock()
	return &c
}

func (as *asyncStruct) addResult(res CloseContacts) {
	as.cond.L.Lock()
	as.cc = (append(as.cc, res...))
	sort.Sort(as.cc)
	as.activeThreads --
	as.cond.Broadcast()
	as.cond.L.Unlock()
}

func NewKademlia(address string, network Net, base *Contact) *Kademlia{
	var c Contact = NewContact(NewKademliaID(randomHex(40)), "localghost")
	var rt *RoutingTable = NewRoutingTable(c)
	if(base != nil){
		rt.AddContact(*base)
	}
	return &Kademlia{rt, network}
}

func (kademlia *Kademlia) asyncLookup(target *Contact, as *asyncStruct) {
	defer as.wg.Done()
	for {
		var c *CloseContact = as.getNext()
		if (c == nil) {return}
		
		var a CloseContacts = kademlia.network.SendFindContactMessage(&c.contact, target.ID)
		//Go through all the results, and spawn routines to add them to RT
		
		as.addResult(a)
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
	//var wg sync.WaitGroup
	//var as *asyncStruct = &asyncStruct{cc, nil, &sync.Cond{L: &sync.Mutex{}}, 0, 0,&wg}
	var as *asyncStruct = NewAsyncStruct(cc)
	as.wg.Add(alpha)
	for i:= 0; i < alpha; i++ {
		go kademlia.asyncLookup(target, as)
	}
	//wait for result here
	as.wg.Wait()
	// Step 4. If found, return contact.
	
	var result CloseContacts = append(as.cc, as.searched...)
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
	// Hash data to get handle
	// Store data in own file (I think?)
	// Do lookup on data handle (I think?)
	// Store data in k closest nodes (I think?)
	// return handle
}
