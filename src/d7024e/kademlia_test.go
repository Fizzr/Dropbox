package d7024e

import (
	//	"messages"
	//	"net"
	"fmt"
	//	proto "github.com/golang/protobuf/proto"
	"sort"
	"testing"
	"time"
)

type MockNetwork struct {
	ip   string
	port int
	me   *Contact
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

func (mn *MockNetwork) SendPingMessage(contact *Contact) bool {
	return false
}
func (mn *MockNetwork) SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts {
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
	for i := 0; i < len(randRTs); i++ {
		if randRTs[i].me.ID.Equals(contact.ID) {
			//fmt.Println("yeh boi")
			cc := randRTs[i].FindClosestContacts(target, k)
			if mn.me != nil {
				//fmt.Println("lal")
				randRTs[i].AddContact(*mn.me)
			}
			return cc
		}
	}
	return nil
}
func (mn *MockNetwork) SendFindDataMessage(contact *Contact, hash string) (*CloseContacts, *[]byte) {
	if hash == "9cfef18a4799c191f79c9995dc2d7b9a49fcd213" {
		//returnByteData = []byte(hash)
		var returnDataString string = "apa"
		return nil, returnDataString
	}
	return nil, nil
}

var mapmap map[*KademliaID]*(map[string]*[]byte) = make(map[*KademliaID]*(map[string]*[]byte))
func (mn *MockNetwork) SendStoreMessage(contact *Contact, hash string, data []byte) {
	m, ok := mapmap[contact.ID]
	if(!ok){
		var newM map[string]*[]byte = make(map[string]*[]byte)
		mapmap[contact.ID] = &newM
		newM[hash] = &data
	} else {
		(*m)[hash] = &data
	}
	return
}

func getClosest(target KademliaID, conts []*Contact) CloseContacts {
	var cc CloseContacts
	for i := 0; i < len(conts); i++ {
		var dist *KademliaID = conts[i].CalcDistance(&target)
		cc = append(cc, CloseContact{*conts[i], dist})
	}
	sort.Sort(cc)
	return cc
}

var fib []int = []int{1, 2, 3, 5, 7, 12, 19, 31, 50, 81}

func TestBootstrap(t *testing.T) {
	var base *Kademlia
	var q, port int = 100, 8001
	randRTs = make([]RoutingTable, 0, q)
	base = newKademlia("localhost:8000", &MockNetwork{"localhost", 8000, nil}, nil)
	randRTs = append(randRTs, *base.rt)
	for i := 1; i < q; i++ {
		var tmp *Kademlia = newKademlia(fmt.Sprintf("localhost:%d", port), &MockNetwork{"localhost", port, nil}, &base.rt.me)
		port++
		base = tmp
		randRTs = append(randRTs, *base.rt)
	}
	//fmt.Println(base.rt.root)
}

func TestKademlia(t *testing.T) {
	var q, port int = 82, 8001
	var contacts []*Contact = make([]*Contact, 0, q)
	randRTs = make([]RoutingTable, 0, q)
	for i := 0; i < q; i++ {
		c := NewContact(NewRandomKademliaID(), fmt.Sprintf("localhost:%d", port))
		port++
		contacts = append(contacts, &c)
		randRTs = append(randRTs, *NewRoutingTable(*contacts[i]))
	}
	for i := 0; i < q; i++ {
		var cc CloseContacts = getClosest(*contacts[i].ID, contacts)
		for j := 0; j < len(fib); j++ {
			//fmt.Print(",")
			_, added := randRTs[i].AddContact(cc[fib[j]].contact)
			if !added {
				fmt.Println("Test contact rejected by bucket. Not a reliable test!")
				t.Fail()
			}
		}
	}

	testLookupContact(t)

	var find Contact = randRTs[10].me
	var closest CloseContacts = getClosest(*find.ID, contacts)
	testFindNode(find, closest, t)
}

func testLookupContact(t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8000, nil}
	var base Contact = randRTs[34].me
	//Don't want to run the bootstrap in test
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn, nil}
	rt.AddContact(base)
	var look Contact = randRTs[22].me
	var ret *Contact = kad.LookupContact(&look)

	var bueno bool

	bueno = ret.ID.Equals(randRTs[22].me.ID)
	if !bueno {
		fmt.Printf("LookupContact: Wrong ID. \n%v Expected\n%v found\n", randRTs[22].me.ID, ret.ID)
	}

	look = NewContact(NewRandomKademliaID(), "aa")
	ret = kad.LookupContact(&look)
	bueno = bueno && ret.ID == nil
	if !bueno {
		fmt.Printf("LookupContact: Wrong ID. Expected nil, found %v\n", ret.ID)
	}

	if bueno {
		fmt.Println("Success - Kademlia LookupContact")
	} else {
		t.Fail()
	}
}

func testFindData(t *testing.T) {
	var mn *MockNetwork = &MockNetowrk{"localhost", 8000, nil}
	var base Contact = randRTs[29].me
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn, nil} // Send Data In Here?
	rt.AddContact(base)

	var data string = "9cfef18a4799c191f79c9995dc2d7b9a49fcd213"
	var ret = kad.LookupData(data)

	var bueno bool
	bueno = ret.Equals([]byte("apa"))

	if !bueno {
		fmt.Printf("LookupData: Wrong DataID. \n%v Expected\n%v found\n", data, ret)
	}

	if bueno {
		fmt.Println("Sucess - Kademlia LookupData")
	} else {
		t.Fail()
	}

}

func testFindNode(look Contact, correct CloseContacts, t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8000, nil}
	var base Contact = randRTs[29].me
	//Don't want to run the bootstrap in test
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn, nil}
	rt.AddContact(base)
	//var kad *Kademlia = NewKademlia("localhost:8001", mn, &base)	var look Contact = randRTs[10].me
	var cc CloseContacts = kad.FindNode(&look)
	var bueno bool
	bueno = len(cc) == k
	if !bueno {
		fmt.Printf("FindContact: Wrong length. Expected %d, found %d\n", k, len(cc))
	}

	for i := 0; i < len(cc); i++ {
		var good bool = cc[i].contact.ID.Equals(correct[i].contact.ID)
		bueno = bueno && good
		if !good {
			fmt.Printf("FindContact: Wrong ID at index %d.\n %v Expected\n %v found\n", i, correct[i].contact.ID, cc[i].contact.ID)
		}
	}

	if bueno {
		fmt.Println("Success - Kademlia FindNode")
	} else {
		t.Fail()
	}
}

func TestStore(t *testing.T) {
	var mn *MockNetwork = &MockNetwork{"localhost", 8000, nil}
	var c Contact = NewContact(NewRandomKademliaID(), "localhost:8001")
	var rt *RoutingTable = NewRoutingTable(c)
	var kad *Kademlia = &Kademlia{rt, mn, nil}
	var encrypt string = "apa"
	var hash string = kad.Store([]byte(encrypt))
	time.Sleep(time.Second)

	m, ok := mapmap[c.ID]
	var bueno bool = true
	if(!ok) {
		fmt.Println("Store: Couldn't find contact map")
		bueno = false;
	} else {
		data, innerOK := (*m)[hash]
		if(!innerOK){
			fmt.Println("Store: Couldn't find data from hash")
		}else {
			var good bool = string(*data) == encrypt
			bueno = bueno && good
			if(!good) {
				fmt.Printf("Store: Incorrect data. Expected %s, found %s", encrypt, string(*data))
			}
		}
	}
	if(bueno) {
		fmt.Println("Success - Kademlia Store")
	} else {
		t.Fail()
	}
}
