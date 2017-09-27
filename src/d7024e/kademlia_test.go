package d7024e

import (
//	"messages"
//	"net"
	"fmt"
//	proto "github.com/golang/protobuf/proto"
	"testing"
	"sort"
	"time"
)

type MockNetwork struct {
	ip string
	port int
	me *Contact
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
	
var randRTs []RoutingTable

func (mn *MockNetwork) SendPingMessage(contact *Contact){
	return
}
func (mn *MockNetwork) SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts{
	time.Sleep(10000)
	/*for i := 0; i < len(lookList); i++ {
		if contact.ID.String() == lookList[i] {
			var res CloseContacts
			for j := i; j < i+3 && j < len(lookList); j++ {
				var c Contact = NewContact(NewKademliaID(lookList[j]), "localhost")
				res = append(res, CloseContact{c, c.ID.CalcDistance(target)})
			}
			sort.Sort(res)
			return res
		}
	}*/
	//fmt.Println("fml")
	for i:=0; i < len(randRTs); i++ {
		if(randRTs[i].me.ID.Equals(contact.ID)){
			//fmt.Println("yeh boi")
			cc := randRTs[i].FindClosestContacts(target, k)
			if(mn.me != nil) {
				//fmt.Println("lal")
				randRTs[i].AddContact(*mn.me)
			}
			return cc
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

func getClosest(target KademliaID, conts []*Contact) CloseContacts{
	var cc CloseContacts
	for i := 0; i < len(conts); i++ {
		var dist *KademliaID = conts[i].CalcDistance(&target)
		cc = append(cc, CloseContact{*conts[i], dist})
	}
	sort.Sort(cc)
	return cc
}

var fib []int = []int{1,2,3,5,7,12,19,31}

func TestBootstrap(t *testing.T) {
	var base *Kademlia
	var q, port int = 100, 8001
	randRTs = make([]RoutingTable, 0, q)
	base = NewKademlia( "localhost:8000", &MockNetwork{"localhost", 8000, nil}, nil)
	randRTs = append(randRTs, *base.rt)
	for i := 1; i < q; i++ {
		var tmp *Kademlia = NewKademlia(fmt.Sprintf("localhost:%d", port), &MockNetwork{"localhost", port, nil}, &base.rt.me)
		port ++
		base = tmp
		randRTs = append(randRTs, *base.rt)
	}
	fmt.Println(base.rt.root)
}

func TestKademlia(t *testing.T) {
	var q, port int = 100, 8001
	var contacts []*Contact = make([]*Contact, 0, q)
	randRTs = make([]RoutingTable, 0, q)
	for i:= 0; i < q; i++ {
		c := NewContact(NewRandomKademliaID(), fmt.Sprintf("localhost:%d", port))
		port ++
		contacts = append(contacts, &c)
		randRTs = append(randRTs, *NewRoutingTable(*contacts[i]))
	}
	for i:= 0; i < q; i++ {
		var cc CloseContacts = getClosest(*contacts[i].ID, contacts)
		for j := 0; j < len(fib); j++ {
			//fmt.Print(",")
			randRTs[i].AddContact(cc[fib[j]].contact)
		}
	}
	
	
	testLookupContact(t)
	testFindNode(t)
}

func testLookupContact(t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8000, nil}
	var base Contact = randRTs[34].me
	//Don't want to run the bootstrap in test
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn}
	rt.AddContact(base)
	//var kad *Kademlia = NewKademlia("localhost:8001", mn, &base)
	var look Contact = randRTs[22].me
	var ret *Contact = kad.LookupContact(&look)
	/*for i := 0; i < len(cc); i++ {
		fmt.Println(cc[i].contact.ID)
	}*/
	fmt.Println(ret)
	
	look = NewContact(NewRandomKademliaID(), "aa")
	ret = kad.LookupContact(&look)
	fmt.Println(ret)
}

func testFindNode(t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8000, nil}
	var base Contact = randRTs[29].me
	//Don't want to run the bootstrap in test
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn}
	rt.AddContact(base)
	//var kad *Kademlia = NewKademlia("localhost:8001", mn, &base)	var look Contact = randRTs[10].me
	var look Contact = randRTs[10].me
	var cc CloseContacts = kad.FindNode(&look)
	for i := 0; i < len(cc); i++ {
		fmt.Println(cc[i])
	}
	look = NewContact(NewRandomKademliaID(), "")
	cc = kad.FindNode(&look)
	fmt.Println("------------")
	for i := 0; i < len(cc); i++ {
		fmt.Println(cc[i])
	}
}