package d7024e

import (
	"container/list"
	//"fmt"
)

const bucketSize = k

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

func (bucket *bucket) AddContact(contact Contact) (bool, bool){
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
			//fmt.Println("Added contact to bucket")
			return true, true
		} else {
			//fmt.Println("Discarding contact")
			return false, false
		}
	} else {
		bucket.list.MoveToFront(element)
		//fmt.Println("Moved contact to top")
		return true, false
	}
}

func (buck *bucket) splitOn(exponent int) [2]bucket{
	var bucketList [2]bucket = [2]bucket{*newBucket(), *newBucket()}
	for e:= buck.list.Front(); e != nil; e = e.Next() {
		var c Contact = e.Value.(Contact)
		//If bit at exponent is 1, push to bucket at 1. If exponent is 0, push to bucket at 0
		bucketList[c.ID.bitAt(exponent)].list.PushBack(c)
	}
	return bucketList
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) CloseContacts {
	var contacts CloseContacts

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		var contact Contact = elt.Value.(Contact)
		var dist *KademliaID = contact.CalcDistance(target)
		contacts = append(contacts, CloseContact{contact, dist})
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}