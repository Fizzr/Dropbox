package d7024e

import (
	//"math"
	"fmt"
	"sort"
	"sync"
)

var myID *KademliaID
var myBucketID int

type RoutingTable struct {
	me Contact
	root Node
}

type Node interface {
	getBucketFor(KademliaID) *bucket
	addContact(Contact) (bool, bool)
	getClosestContacts(target *KademliaID, count int) CloseContacts
}

type Branch struct{
	prefix [20]byte
	exponent int
	left Node
	right Node
}
type Leaf struct{
	prefix [20]byte
	exponent int
	ID int
	buck *bucket
	rw sync.RWMutex
}


//Returns true if the bit of the KademliaID at the relevant exponent for this branch is 1
// confusing function. Consider purging in favour of bitAt
func (branch *Branch) isOne(ID KademliaID) bool{
	return ID.bitAt(branch.exponent) == 1
}
//Is this needed ever?!
func (branch *Branch) getBucketFor(ID KademliaID) *bucket{
	if(branch.isOne(ID)){
		return branch.left.getBucketFor(ID)
	} else {
		return branch.right.getBucketFor(ID)
	}
}

func (branch *Branch) addContact(contact Contact) (bool, bool) {
	//fmt.Printf("Branch add, exp %d\n", branch.exponent)
	var ok, added, isLeft bool
	//fmt.Println(contact.ID.bitAt(branch.exponent))
	isLeft = contact.ID.bitAt(branch.exponent) == 1
	var node Node
	if(isLeft){
		node = branch.left
		//ok, added = branch.left.addContact(contact)
	} else {
		node = branch.right
		//ok, added = branch.right.addContact(contact)
	}
	
	if(fmt.Sprintf("%T", node) == "*d7024e.Leaf"){
		node.(*Leaf).rw.Lock()	//Lock for writing
		defer node.(*Leaf).rw.Unlock()
	}
	ok, added = node.addContact(contact)
	if(ok){
		return ok, added
	} else { 
		// Need to split leaf
		// Only ever get here if just above a leaf.
	
	//fmt.Print("Splitting Leaf ")
		var leaf *Leaf = node.(*Leaf)
/*		if(isLeft){
			leaf = branch.left.(*Leaf)
		} else {
			leaf = branch.right.(*Leaf)
		}*/
		if(leaf.ID != myBucketID){
			//fmt.Println("HAHAHAH")
			return true, false;	// We don't wanna split. Return true, cos everything is just fine!
		}
		var splitExponent int = leaf.exponent
		var buckets [2]bucket = leaf.buck.splitOn(splitExponent)
		var oldID int = leaf.ID
		var prefix [20]byte = leaf.prefix
		
		var myBitAtExponent byte = myID.bitAt(splitExponent)
		//fmt.Println(myID.toBinary())
		var leftID int = int(myBitAtExponent) + oldID
		var rightID int = ((int(myBitAtExponent)-1) * -1)	+ oldID	//if 1, becomes 0. If 0, becomes 1
		var leftPrefix, rightPrefix [20]byte = prefix, prefix
		var IDindex int = (IDLength - 1) - (splitExponent/8)
		leftPrefix[IDindex] = leftPrefix[IDindex] | (1 << uint(splitExponent%8))
		var left Leaf = Leaf{leftPrefix, splitExponent-1, leftID, &(buckets[1]), sync.RWMutex{}}
		var right Leaf = Leaf{rightPrefix, splitExponent-1, rightID, &(buckets[0]), sync.RWMutex{}}

		myBucketID = oldID +1
		var newBranch Branch = Branch{prefix, splitExponent, &left, &right}
		if(isLeft){
			branch.left = &newBranch
			return branch.left.addContact(contact)
		} else {
			branch.right = &newBranch
			return branch.right.addContact(contact)
		}
	}
}

func (branch *Branch) getClosestContacts(target *KademliaID, count int) CloseContacts{
	var bitAtExponent byte = target.bitAt(branch.exponent)
	var direction, other Node 
	if(bitAtExponent == 1){
		direction = branch.left
		other = branch.right
	} else {
		direction = branch.right
		other = branch.left
	}
	if(fmt.Sprintf("%T", direction) == "*d7024e.Leaf") {
		direction.(*Leaf).rw.RLock();	//read Lock
		defer direction.(*Leaf).rw.RUnlock();
	}
	var dirRes CloseContacts = direction.getClosestContacts(target, count)
	var diff int = count - len(dirRes)
	if(diff == 0){
		return dirRes
	}else{
		if(fmt.Sprintf("%T", other) == "*d7024e.Leaf") {
			other.(*Leaf).rw.RLock();	//read Lock
			defer other.(*Leaf).rw.RUnlock();
		}
		return append(dirRes, other.getClosestContacts(target, diff)...)
	}
}

func (branch *Branch) String() string{
	var tabs, info, openBrack, closeBrack string
	
	for i := 0; i < IDBits-1 - branch.exponent; i++ {
		tabs += "\t"
	}
	info = tabs + fmt.Sprintf("Branch - Exponent %v, Prefix: ", branch.exponent)
	var to int
	
	to = IDLength-1 - (branch.exponent / 8)
	for i := 0; i < to /*- edge*/; i++ {
		info += fmt.Sprintf("%08b", branch.prefix[i])
	}
	for j := 7; j > branch.exponent % 8; j-- {
		info += fmt.Sprintf("%b", (branch.prefix[to] >> uint(j)) & 1) //Apparently can't use width to limit length
	}
	info += "\n"
	openBrack = tabs + "{\n"
	closeBrack = tabs + "}\n"
	return info + openBrack + fmt.Sprint(branch.right) + closeBrack + openBrack + fmt.Sprint(branch.left) + closeBrack
}


func (leaf *Leaf) String() string{
	
	var tabs, info string
	for i := 0; i < IDBits-1 - leaf.exponent; i++ {
		tabs += "\t"
	}
	info = fmt.Sprintf("Leaf - ID: %v, Number of entries: %v Exponent: %v\nPrefix:\n", leaf.ID, leaf.buck.Len(), leaf.exponent )
	var to int
	
	to = IDLength-1 - (leaf.exponent / 8)
	for i := 0; i < to/* - edge*/; i++ {
		info += fmt.Sprintf("%08b", leaf.prefix[i])
	}
	for j := 7; j > leaf.exponent % 8; j-- {
		info += fmt.Sprintf("%b", (leaf.prefix[to] >> uint(j)) & 1) //Apparently can't use width to limit length
	}
	info += "\n-------\n"
	var content string
	for e := leaf.buck.list.Front(); e != nil; e = e.Next() {
		content += fmt.Sprintf("%s%s\n",tabs, e.Value.(Contact).ID.toBinary())
	}
	return tabs + info + content
}

func (leaf *Leaf) getBucketFor(_ KademliaID) *bucket{
	return leaf.buck;
}

func (leaf *Leaf) addContact(contact Contact) (bool, bool) {
	//fmt.Printf("Leaf add exp %d ID %d\n", leaf.exponent, leaf.ID)
	var ok, added bool = leaf.buck.AddContact(contact)
	return ok, added
}

func (leaf *Leaf) getClosestContacts(target *KademliaID, count int) CloseContacts{
	var res CloseContacts = leaf.buck.GetContactAndCalcDistance(target)
	sort.Sort(res)
	if(len(res) > count){
		res = res[:count]
	}
	return res
}


func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	routingTable.me = me
	myBucketID = 1
	myID = me.ID
	var prefix [20]byte
	var bit int = int(myID.bitAt(IDBits-1))
	var left, right *Leaf 
	var leftPrefix [20]byte = prefix
	leftPrefix[0] = leftPrefix[0] | 1<<7
	left = &Leaf{leftPrefix, IDBits-2, bit, newBucket(), sync.RWMutex{}}
	right = &Leaf{prefix, IDBits-2, (bit-1)*-1, newBucket(), sync.RWMutex{}}
	routingTable.root = &Branch{prefix, IDBits-1, left, right}
	routingTable.root.addContact(me)
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) (bool, bool) {
	//TODO: Not threadsafe. Maybe just make root into a branch to start off with....
	return routingTable.root.addContact(contact)
	/*var ok, added bool = routingTable.root.addContact(contact)
	if(!ok){
		var buckets [2]bucket
		var splitExponent int = IDBits-1
		buckets = routingTable.root.(*Leaf).buck.splitOn(splitExponent)
		var left, right Leaf
		var myBitAtExponent byte = routingTable.me.ID.bitAt(splitExponent)
		var leftID int = int(myBitAtExponent)
		var rightID int = (int(myBitAtExponent)-1) * -1		//if 1, becomes 0. If 0, becomes 1
		var leftPrefix, rightPrefix [20]byte
		leftPrefix[0] = 1 << uint(splitExponent%8)
		left = Leaf{leftPrefix, splitExponent-1, leftID, &(buckets[1]), sync.RWMutex{}}
		right = Leaf{rightPrefix, splitExponent-1, rightID, &(buckets[0]), sync.RWMutex{}}
		myBucketID = 1
		var newBranch Branch = Branch{routingTable.root.(*Leaf).prefix, splitExponent, &left, &right}
		routingTable.root = &newBranch
		return routingTable.root.addContact(contact)
	}else {return ok, added}*/
	
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) CloseContacts {
	return routingTable.root.getClosestContacts(target, count)
}
