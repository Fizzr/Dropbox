package d7024e

import (
	//"math"
	"fmt"
)

var myID *KademliaID
var myBucketID int

type Node interface {
	getBucketFor(KademliaID) *bucket
	addContact(Contact) bool
	getClosestContact(target KademliaID, count int) []CloseContact
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
}


//Returns true if the bit of the KademliaID at the relevant exponent for this branch is 1
// confusing function. Consider purging in favour of bitAt
func (branch *Branch) isOne(ID KademliaID) bool{
	return ID.bitAt(branch.exponent) == 1
}

func (branch *Branch) getBucketFor(ID KademliaID) *bucket{
	if(branch.isOne(ID)){
		return branch.left.getBucketFor(ID)
	} else {
		return branch.right.getBucketFor(ID)
	}
}

func (branch *Branch) addContact(contact Contact) bool {
	//fmt.Printf("Branch add, exp %d\n", branch.exponent)
	var ok, isLeft bool
	//fmt.Println(contact.ID.bitAt(branch.exponent))
	isLeft = contact.ID.bitAt(branch.exponent) == 1
	if(isLeft){
		ok = branch.left.addContact(contact)
	} else {
		ok = branch.right.addContact(contact)
	}
	
	if(ok){
		return ok
	} else { 
		//fmt.Print("Splitting Leaf ")
		var leaf *Leaf
		if(isLeft){
			leaf = branch.left.(*Leaf)
		} else {
			leaf = branch.right.(*Leaf)
		}
		if(leaf.ID != myBucketID){
		
			//fmt.Printf("Not my buckedent! Leaf %d, Mine %d\n", leaf.ID, myBucketID)
			return true;	// We don't wanna split. Return true, cos everything is just fine!
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
		var left Leaf = Leaf{leftPrefix, splitExponent-1, leftID, &(buckets[1])}
		var right Leaf = Leaf{rightPrefix, splitExponent-1, rightID, &(buckets[0])}
		//fmt.Printf("My bit at %d: %d, leftID %d, rightID %d\n",splitExponent, int(myBitAtExponent), leftID, rightID)
		//fmt.Printf("Branch exponent: %d, splitExponent: %d, newLeafExponent %d\n", branch.exponent, splitExponent, splitExponent-1)

		
		/*
		if (myBitAtExponent == 1){
			left = Leaf{prefix ^ (1 << uint(splitExponent)),splitExponent, oldID + 1, &(buckets[1])}
			right = Leaf{prefix, splitExponent, oldID, &(buckets[0])}
		} else {
			left = Leaf{prefix ^ (1 << uint(splitExponent)), splitExponent, oldID, &(buckets[1])}
			right = Leaf{prefix, splitExponent, oldID + 1, &(buckets[0])}
		}*/
		myBucketID = oldID +1
		var newBranch Branch = Branch{prefix, splitExponent, &left, &right}
		if(isLeft){
			//fmt.Printf("what used to be left? %T\n", branch.left)
			branch.left = &newBranch
			//fmt.Printf("what is left? %T\n", branch.left)
			return branch.left.addContact(contact)
		} else {
			//fmt.Printf("what used to be right? %T\n", branch.right)
			branch.right = &newBranch
			//fmt.Printf("what is right? %T\n", branch.right)
			return branch.right.addContact(contact)
		}
	}
}

func (branch *Branch) getClosestContact(target KademliaID, count int) []CloseContact{
	return nil
}

func (branch *Branch) String() string{
	var tabs, info, openBrack, closeBrack string
	
	for i := 0; i < IDBits-1 - branch.exponent; i++ {
		tabs += "\t"
	}
	info = tabs + fmt.Sprintf("Branch - Exponent %v, Prefix: ", branch.exponent)
	var /*i, edge*/ to int
	/*if(branch.exponent%8 != 0){
		edge = 1
	}*/
	to = IDLength-1 - (branch.exponent / 8)
	for i := 0; i < to /*- edge*/; i++ {
		info += fmt.Sprintf("%08b", branch.prefix[i])
	}
	//if(i != 0){
	//i++}
	//if(branch.exponent % 8 != 0){
		for j := 7; j > branch.exponent % 8; j-- {
			info += fmt.Sprintf("%b", (branch.prefix[to] >> uint(j)) & 1) //Apparently can't use width to limit length
		}
	//}
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
	info = fmt.Sprintf("Leaf - ID: %v, Number of entries: %v Exponent: %v Prefix: ", leaf.ID, leaf.buck.Len(), leaf.exponent )
	var /*i, edge*/ to int
	/*if(leaf.exponent%8 != 0){
		edge = 1
	}*/
	to = IDLength-1 - (leaf.exponent / 8)
	for i := 0; i < to/* - edge*/; i++ {
		info += fmt.Sprintf("%08b", leaf.prefix[i])
	}
/*	if(i != 0){
	i++}*/
	//if(leaf.exponent % 8 != 0){
		for j := 7; j > leaf.exponent % 8; j-- {
			info += fmt.Sprintf("%b", (leaf.prefix[to] >> uint(j)) & 1) //Apparently can't use width to limit length
		}
	//}
	info += "\n"
	var content string
	for e := leaf.buck.list.Front(); e != nil; e = e.Next() {
		content += fmt.Sprintf("%s%s\n",tabs, e.Value.(Contact).ID.toBinary())
	}
	return tabs + info + content
}

func (leaf *Leaf) getBucketFor(_ KademliaID) *bucket{
	return leaf.buck;
}

func (leaf *Leaf) addContact(contact Contact) bool {
	//fmt.Printf("Leaf add exp %d ID %d\n", leaf.exponent, leaf.ID)
	return leaf.buck.AddContact(contact)
}

func (leaf *Leaf) getClosestContact(target KademliaID, count int) []CloseContact{
	return nil
}

type RoutingTable struct {
	me Contact
	root Node
	buckets [IDBits]*bucket
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
	routingTable.root = &Leaf{prefix,IDBits-1, 0, newBucket()}
	routingTable.root.addContact(me)
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) bool {
	//fmt.Println("rt add")
	//bucketIndex := routingTable.getBucketIndex(contact.ID)
	//buck := routingTable.buckets[bucketIndex]
	//buck.AddContact(contact)
	
	var ok bool = routingTable.root.addContact(contact)
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
		left = Leaf{leftPrefix, splitExponent-1, leftID, &(buckets[1])}
		right = Leaf{rightPrefix, splitExponent-1, rightID, &(buckets[0])}
		myBucketID = 1
		//fmt.Printf("Root split left prefix! %b\n", left.prefix)
		//fmt.Printf("My bit at %d: %d, leftID %d, rightID %d\n", splitExponent, myBitAtExponent, leftID, rightID)
		
		/*/Can be made neater!
		if(myBitAtExponent == 1){
			left = Leaf{1 << uint(splitExponent), splitExponent, 1, &(buckets[1])}
			right = Leaf{0, splitExponent, 0, &(buckets[0])}
		}else{
			left = Leaf{1 << uint(splitExponent), splitExponent , 0, &(buckets[1])}
			right = Leaf{0, splitExponent, 1, &(buckets[0])}
		}*/
		var newBranch Branch = Branch{routingTable.root.(*Leaf).prefix, splitExponent, &left, &right}
		routingTable.root = &newBranch
		return routingTable.root.addContact(contact)
	}else {return true}
	
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDBits) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDBits {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
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
