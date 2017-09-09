package d7024e

import (
	"fmt"
)

const bucketSize = 5
var myID KademliaID
var myBucketID int

type Node interface {
	getBucketFor(KademliaID) *bucket
	addContact(Contact) bool
}

type Branch struct{
	prefix uint//[20]byte
	exponent int
	left Node
	right Node
}
type Leaf struct{
	prefix uint//[20]byte
	exponent int
	ID int
	buck *bucket
}


//Returns true if the bit of the KademliaID at the relevant exponent for this branch is 1, 
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
	isLeft = contact.ID.bitAt(branch.exponent) == 1//branch.isOne(*contact.ID)
	if(isLeft){
		ok = branch.left.addContact(contact)
	} else {
		ok = branch.right.addContact(contact)
	}
	
	if(ok){
		return ok
	} else { //>>>>>>>>>>>>>>>>> TODO: Check if we can split! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		//fmt.Print("Splitting Leaf ")
		var leaf *Leaf
		if(isLeft){
			leaf = branch.left.(*Leaf)
		} else {
			leaf = branch.right.(*Leaf)
		}
		if(leaf.ID != myBucketID){
			return true;	// We don't wanna split. Return true, cos everything is just fine!
		}
		
		var splitExponent int = branch.exponent-2
		var buckets [2]bucket = leaf.buck.splitOn(splitExponent)
		var oldID int = leaf.ID
		var prefix uint = leaf.prefix
		
		var myBitAtExponent byte = myID.bitAt(splitExponent)
		var left, right Leaf
		if (myBitAtExponent == 1){
			left = Leaf{prefix ^ (1 << uint(splitExponent)),splitExponent, oldID + 1, &(buckets[1])}
			right = Leaf{prefix, splitExponent, oldID, &(buckets[0])}
		} else {
			left = Leaf{prefix ^ (0x1 << uint(splitExponent)), splitExponent, oldID, &(buckets[1])}
			right = Leaf{prefix, splitExponent, oldID + 1, &(buckets[0])}
		}
		myBucketID = oldID +1
		var newBranch Branch = Branch{prefix, splitExponent+1, &left, &right}
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

func (branch *Branch) String() string{
	var tabs, info, openBrack, closeBrack string
	for i := 0; i < IDBits - branch.exponent; i++ {
		tabs += "\t"
	}
	info = tabs + fmt.Sprintf("Exponent %v, Prefix: %0*b \n", branch.exponent, IDBits - branch.exponent, branch.prefix >> uint(branch.exponent))
	openBrack = tabs + "{\n"
	closeBrack = tabs + "}\n"
	return info + openBrack + fmt.Sprint(branch.right) + closeBrack + openBrack + fmt.Sprint(branch.left) + closeBrack
}


func (leaf *Leaf) String() string{
	
	var tabs string
	for i := 0; i < IDBits - leaf.exponent; i++ {
		tabs += "\t"
	}
	return tabs + fmt.Sprintf("ID: %v, Number of entries: %v Exponent: %v Prefix: %0*b\n", leaf.ID, leaf.buck.Len(), leaf.exponent ,IDBits -leaf.exponent , leaf.prefix) //>> uint(leaf.exponent))
}

func (leaf *Leaf) getBucketFor(_ KademliaID) *bucket{
	return leaf.buck;
}

func (leaf *Leaf) addContact(contact Contact) bool {
	//fmt.Printf("Leaf add exp %d ID %d\n", leaf.exponent, leaf.ID)
	return leaf.buck.AddContact(contact)
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
	routingTable.root = &Leaf{0,IDBits, 0, newBucket()}
	routingTable.root.addContact(me)
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) {
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
		//Can be made neater!
		if(myBitAtExponent == 1){
			left = Leaf{1 << uint(splitExponent), splitExponent, 1, &(buckets[1])}
			right = Leaf{0, splitExponent, 0, &(buckets[0])}
		}else{
			left = Leaf{1 << uint(splitExponent), splitExponent , 0, &(buckets[1])}
			right = Leaf{0, splitExponent, 1, &(buckets[0])}
		}
		var newBranch Branch = Branch{0, IDBits, &left, &right}
		routingTable.root = &newBranch
		routingTable.root.addContact(contact)
	}
	
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
