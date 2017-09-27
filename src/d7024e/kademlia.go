package d7024e

import (
	"sync"
//	"sync/atomic"
	"sort"
//	"fmt"
)

const k = 20
const alpha = 3

type Kademlia struct {
	rt      *RoutingTable
	network Net
}

type asyncStruct struct {
	cc CloseContacts
	searched CloseContacts
	cond *sync.Cond
	index int
	activeThreads int
	wg *sync.WaitGroup
	run bool
}

func NewKademlia(address string, network Net, base *Contact) *Kademlia{
	var c Contact = NewContact(NewRandomKademliaID(), "localghost")
	var rt *RoutingTable = NewRoutingTable(c)
	if base != nil {
		rt.AddContact(*base)
	}
	return &Kademlia{rt, network}
}

func NewAsyncStruct(base CloseContacts) *asyncStruct{
	var wg sync.WaitGroup
	return &asyncStruct{base, nil, &sync.Cond{L: &sync.Mutex{}}, 0, 0, &wg, true}
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
		checked = true
		for i:= 0; checked && i < len(as.searched); i++ {					//Make sure we're not seaching one we already searched.
			if c.contact.ID.Equals(as.searched[i].contact.ID) {
				checked = false
			}
		}
					//Checked through all searched items, didn't find duplicates
	}
	as.searched = append(as.searched, c)
//	fmt.Println("Searched ------ " + fmt.Sprint(as.searched))
	var num int = as.index
	as.index ++
	if(num >= k) {
		as.run = false
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
	//Both res and cc are sorted with respect to the same distance
	//So we can step throught them together, and see if they are equals or not
	//Insert elements from res into index of cc if res is less than cc element
	//If they are the same, discard res element
	/*for i:= 0; i < len(as.cc); i++ {
		fmt.Printf("Input cc - %v\n", as.cc[i]) 
	}
	for i:= 0; i < len(res); i++ {
		fmt.Printf("Input res - %v\n", res[i]) 
	}*/

	var newCC CloseContacts = make([]CloseContact, 0, len(res) + len(as.cc))
	//fmt.Printf("resLen %d ccLen %d newCCLen %d\n",len(res), len(as.cc), len(newCC))
	for i, j, k := 0,0,0; j+i < len(as.cc) + len(res); {
		//fmt.Printf("i %d j %d k %d\n",i,j,k)
		if j == len(res) {
			newCC = append(newCC, as.cc[i])
			k++
			i++
			continue
		}
		if i == len(as.cc) {
			newCC = append(newCC, res[j])
			k++
			j++
			continue
		}
		if as.cc[i].contact.ID.Equals(res[j].contact.ID) {
			j++
			continue	//Skip this result element
		}
		if as.cc[i].distance.Less(res[j].distance) {
			newCC = append(newCC, as.cc[i])
			k++
			i++
			continue 		//go to next element
		} else {
			newCC = append(newCC, res[j])
			k++
			j++
			continue
		}
	}
	/*for i := 0; i < len(newCC); i++ {
		fmt.Println(newCC[i])
	}*/
//	as.cc = (append(as.cc, res...))
//	sort.Sort(as.cc)
	as.cc = newCC
	as.activeThreads --
	as.cond.Broadcast()
	as.cond.L.Unlock()
}


func (kademlia *Kademlia) asyncLookup(target *Contact, as *asyncStruct, result *Contact, num int) {
	defer as.wg.Done()
	for as.run{
	var c *CloseContact = as.getNext()
		if (c == nil) {return}
		
		//fmt.Printf("Thread %v - Searching %s\n", num, c)
		
		var a CloseContacts = kademlia.network.SendFindContactMessage(&c.contact, target.ID)
		//Go through all the results, and spawn routines to add them to RT. Also check for target
		//fmt.Printf("Thread %v - len %v\n", num, len(a))
		for i := 0; i < len(a); i++ {
			//fmt.Printf("Thread %v - result %s\n",num, a[i])
			go kademlia.rt.AddContact(a[i].contact)
			if(a[i].contact.ID.Equals(target.ID)){
				*result =  a[i].contact
				//fmt.Printf("AAAAAAAAAAA %v\n", a[i].contact)
				as.run = false
			}
		}
	
		as.addResult(a)
		for i := 0; i < len(as.cc); i++{
			//fmt.Printf("Thread %v - cc %v\n", num, as.cc[i])
		}
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) *Contact{
	// Step 1. Get alpha closest to target
	var cc CloseContacts = kademlia.rt.FindClosestContacts(target.ID, alpha)
	// Step 2. See if target exists. If so, return it
	for i := 0; i < len(cc); i++ {
		if(cc[i].contact.ID == target.ID){
			return &(cc[i].contact)
		}
	}
	// Step 3. If not, send LookupContact to k closest contacts, including returned values, running alpha number of lookups in parallel
	var as *asyncStruct = NewAsyncStruct(cc)
	//alpha = 1
	as.wg.Add(alpha)
	var result *Contact = &Contact{}
	for i:= 0; i < alpha; i++ {
		go kademlia.asyncLookup(target, as, result, i)
	}
	//wait for result here
	as.wg.Wait()
	
	return result
}

func (kademlia *Kademlia) asyncFindNode(target *Contact, as *asyncStruct) {
	defer as.wg.Done()
	for as.run {
		var c *CloseContact = as.getNext()
		if (c == nil) {return}
		var a CloseContacts = kademlia.network.SendFindContactMessage(&c.contact, target.ID)
		for i:=0; i < len(a); i++ {
			go kademlia.rt.AddContact(a[i].contact)
		}
		as.addResult(a)
	}
}

func (kademlia *Kademlia) FindNode(target *Contact) CloseContacts{
	var cc CloseContacts = kademlia.rt.FindClosestContacts(target.ID, alpha)
	var as *asyncStruct = NewAsyncStruct(cc)
	as.wg.Add(alpha)
	for i:= 0; i < alpha; i++ {
		go kademlia.asyncFindNode(target, as)
	}
	as.wg.Wait()
	
	var result CloseContacts = append(as.cc, as.searched...)
	sort.Sort(result)
	for i:= 0; i < len(result)-1; i ++ {
		if(result[i].contact.ID.Equals(result[i+1].contact.ID)) {
			result = append(result[:i], result[i+1:]...)
			i--
		}
	}
	return result
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
