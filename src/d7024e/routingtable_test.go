package d7024e

import (
	"fmt"
	"testing"
)

var testList []string = []string{
	"0111111400000000000000000000000000000000",
	"0011111400000000000000000000000000000000",
	"F001111400000000000000000000000000000000",
	"0111111400000000000000000000000000000000",
	"FFFFFFFF00000000000000000000000000000000",
	"F111111100000000000000000000000000000000",
	"F111111200000000000000000000000000000000",
	"F111111300000000000000000000000000000000",
	"1111111400000000000000000000000000000000",
	"2111111400000000000000000000000000000000"}

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	for i := 0; i < len(testList); i++{
		contact := NewContact(NewKademliaID(testList[i]), fmt.Sprintf("localhost:8%03d", i))
		var bits string
		for j := 0; j < IDLength; j++ {
			bits += fmt.Sprintf("%08b", contact.ID[j])
		}
		fmt.Println(bits)
		rt.AddContact(contact)
	}
	
/*	rt.AddContact(NewContact(NewKademliaID("0111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0011111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0001111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))
*/
	//contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	//for i := range contacts {
	//	fmt.Println(contacts[i].String())
	//}
	fmt.Println("")
	fmt.Println(rt.root)
	fmt.Println("")
	
	fmt.Printf("%T, %T, %T \n", rt.root.(*Branch).left, rt.root, rt.root.(*Branch).right)
}
