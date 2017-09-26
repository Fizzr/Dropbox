package d7024e

import (
//	"messages"
//	"net"
	"fmt"
//	proto "github.com/golang/protobuf/proto"
	"testing"
	"sort"
)

type MockNetwork struct {
	ip string
	port int
}
var lookList []string = []string{
	"0000000000000000000000000000000000000000",
	"1111111111111111111111111111111111111111",
	"2222222222222222222222222222222222222222",
	"3333333333333333333333333333333333333333",
	"4444444444444444444444444444444444444444",
	"5555555555555555555555555555555555555555",
	"6666666666666666666666666666666666666666",
	"7777777777777777777777777777777777777777",
	"8888888888888888888888888888888888888888",
	"9999999999999999999999999999999999999999",
	"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
	"CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC",
	"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD",
	"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"}
	
func (mn *MockNetwork) SendPingMessage(contact *Contact){
	return
}
func (mn *MockNetwork) SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts{
	for i := 0; i < len(lookList); i++ {
		if contact.ID.String() == lookList[i] {
			var res CloseContacts
			for j := i; j < i+3 && j < len(lookList); j++ {
				var c Contact = NewContact(NewKademliaID(lookList[j]), "localhost")
				res = append(res, CloseContact{c, c.ID.CalcDistance(target)})
			}
			sort.Sort(res)
			return res
		}
	}
	return nil
}
func (mn *MockNetwork) SendFindDataMessage(hash string) {
	return
}
func (mn *MockNetwork) SendStoreMessage(data []byte) {
	return
}

func TestLookupContact(t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8001}
	var base Contact = NewContact(NewKademliaID(lookList[0]), "localhost:8001")
	var kad *Kademlia = NewKademlia("localhost:8001", mn, &base)
	var look Contact = NewContact(NewKademliaID(lookList[15]), "")
	var cc CloseContacts = *kad.LookupContact(&look)
	for i := 0; i < len(cc); i++ {
		fmt.Println(cc[i].contact.ID)
	}	
}