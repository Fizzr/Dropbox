package d7024e

import (
	//"math"
	"fmt"
	"sort"
)

var myID *KademliaID
var myBucketID int

type Node interface {
	getBucketFor(KademliaID) *bucket
	addContact(Contact) (bool, bool)
	getClosestContacts(target *KademliaID, count int) CloseContacts
}

type Branch struct {
	prefix   [20]byte
	exponent int
	left     Node
	right    Node
}
type Leaf struct {
	prefix   [20]byte
	exponent int
	ID       int
	buck     *bucket
}

//Returns true if the bit of the KademliaID at the relevant exponent for this branch is 1
// confusing function. Consider purging in favour of bitAt
func (branch *Branch) isOne(ID KademliaID) bool {
	return ID.bitAt(branch.exponent) == 1
}

func (branch *Branch) getBucketFor(ID KademliaID) *bucket {
	if branch.isOne(ID) {
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
	if(isLeft){
		ok, added = branch.left.addContact(contact)
	} else {
		ok, added = branch.right.addContact(contact)
	}
	
	if(ok){
		return ok, added
	} else { 
		//fmt.Print("Splitting Leaf ")
		var leaf *Leaf
		if isLeft {
			leaf = branch.left.(*Leaf)
		} else {
			leaf = branch.right.(*Leaf)
		}
		if(leaf.ID != myBucketID){
			fmt.Println("HAHAHAH")
			return true, false;	// We don't wanna split. Return true, cos everything is just fine!
		}

		var splitExponent int = leaf.exponent
		var buckets [2]bucket = leaf.buck.splitOn(splitExponent)
		var oldID int = leaf.ID
		var prefix [20]byte = leaf.prefix

		var myBitAtExponent byte = myID.bitAt(splitExponent)
		//fmt.Println(myID.toBinary())
		var leftID int = int(myBitAtExponent) + oldID
		var rightID int = ((int(myBitAtExponent) - 1) * -1) + oldID //if 1, becomes 0. If 0, becomes 1
		var leftPrefix, rightPrefix [20]byte = prefix, prefix
		var IDindex int = (IDLength - 1) - (splitExponent / 8)
		leftPrefix[IDindex] = leftPrefix[IDindex] | (1 << uint(splitExponent%8))
		var left Leaf = Leaf{leftPrefix, splitExponent-1, leftID, &(buckets[1])}
		var right Leaf = Leaf{rightPrefix, splitExponent-1, rightID, &(buckets[0])}

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
	var dirRes CloseContacts = direction.getClosestContacts(target, count)
	var diff int = count - len(dirRes)
	if(diff == 0){
		return dirRes
	}else{
		return append(dirRes, other.getClosestContacts(target, diff)...)
	}
}

func (branch *Branch) String() string {
	var tabs, info, openBrack, closeBrack string

	for i := 0; i < IDBits-1-branch.exponent; i++ {
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

func (leaf *Leaf) String() string {

	var tabs, info string
	for i := 0; i < IDBits-1-leaf.exponent; i++ {
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
		content += fmt.Sprintf("%s%s\n", tabs, e.Value.(Contact).ID)
	}
	return tabs + info + content
}

func (leaf *Leaf) getBucketFor(_ KademliaID) *bucket {
	return leaf.buck
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

type RoutingTable struct {
	me      Contact
	root    Node
	buckets [IDBits]*bucket
}

func (rt RoutingTable) String() string {
	return fmt.Sprint(rt.root)
}

func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDBits; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	myBucketID = 0
	myID = me.ID
	var prefix [20]byte
	routingTable.root = &Leaf{prefix, IDBits - 1, 0, newBucket()}
	routingTable.root.addContact(me)
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) (bool, bool) {
	
	var ok, added bool = routingTable.root.addContact(contact)
	if(!ok){
		var buckets [2]bucket
		var splitExponent int = IDBits - 1
		buckets = routingTable.root.(*Leaf).buck.splitOn(splitExponent)
		var left, right Leaf
		var myBitAtExponent byte = routingTable.me.ID.bitAt(splitExponent)
		var leftID int = int(myBitAtExponent)
		var rightID int = (int(myBitAtExponent) - 1) * -1 //if 1, becomes 0. If 0, becomes 1
		var leftPrefix, rightPrefix [20]byte
		leftPrefix[0] = 1 << uint(splitExponent%8)
		left = Leaf{leftPrefix, splitExponent - 1, leftID, &(buckets[1])}
		right = Leaf{rightPrefix, splitExponent - 1, rightID, &(buckets[0])}
		myBucketID = 1
		var newBranch Branch = Branch{routingTable.root.(*Leaf).prefix, splitExponent, &left, &right}
		routingTable.root = &newBranch
		return routingTable.root.addContact(contact)
	}else {return ok, added}
	
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) CloseContacts {
	return routingTable.root.getClosestContacts(target, count)
}

func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDBits - 1
}
