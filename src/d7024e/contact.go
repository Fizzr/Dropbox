package d7024e

import (
	"fmt"
)

type Contact struct {
	ID       *KademliaID
	Address  string
}

type CloseContact struct{
	contact Contact
	distance *KademliaID
}

type CloseContacts []CloseContact

func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address}
}

func (contact *Contact) CalcDistance(target *KademliaID) *KademliaID{
	return contact.ID.CalcDistance(target)
}

func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

func (c CloseContacts) Len() int {
	return len(c)
}
func (c CloseContacts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CloseContacts) Less(i, j int) bool {
	return c[i].distance.Less(c[j].distance)
}
func (c CloseContact) String() string {
	return fmt.Sprintf("%s - Distance %s", c.contact, c.distance)
}