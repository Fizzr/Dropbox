package d7024e

import (
	"container/list"
	//"fmt"
)

const bucketSize = 5

type bucket struct {
	list *list.List
}

func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

func (bucket *bucket) AddContact(contact Contact) bool{
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
			return true
		} else {
			//fmt.Println("Discarding contact")
			return false
		}
	} else {
		bucket.list.MoveToFront(element)
		//fmt.Println("Moved contact to top")
		return true
	}
}

func (buck *bucket) splitOn(exponent int) [2]bucket{
	var bucketList [2]bucket = [2]bucket{*newBucket(), *newBucket()}
	for e:= buck.list.Front(); e != nil; e = e.Next() {
		var c Contact = e.Value.(Contact)
		//If bit at exponent is 1, push to bucket at 1. If exponent is 0, push to bucket at 0
		bucketList[(c.ID[(IDLength-1)- (exponent/8)] >> uint((exponent%8)-1)) & 1].list.PushBack(c)
	}
	return bucketList
}

func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
