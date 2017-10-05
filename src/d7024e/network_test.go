package d7024e

import (
	"fmt"
	"messages"
	"net"
	"time"
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func writeByte(address string, b []byte) {
	var laddr, raddr *net.UDPAddr

	laddr, err := net.ResolveUDPAddr("udp", "localhost")
	raddr, err = net.ResolveUDPAddr("udp", address)

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Printf("ERROR! \n %v\n", err)
		return
	}
	conn.Write(b)
	conn.Close()
}

func TestNetwork(t *testing.T) {
	var addr *net.UDPAddr
	addr, err := net.ResolveUDPAddr("udp", "localhost:8001")

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("ERROR! \n %v\n", err)
		t.Fail()
		return
	}
	var b []byte = make([]byte, 255)
	var w []byte = []byte{0x01, 0x02, 0x03, 0xFF}
	//RACE!
	go writeByte("localhost:8001", w)
	//RACE!
	num, err := conn.Read(b)
	if err != nil {
		fmt.Printf("READ ERROR! \n %v\n", err)
		t.Fail()
		return
	}
	for i := 0; i < num; i++ {
		if w[i] != b[i] {
			fmt.Printf("Byte &d is was %X, expected %X\n", i, b[i], w[i])
			t.Fail()
		}
	}
	fmt.Println("Success - Network communication")
	conn.Close()
}

func TestProtobufNetwork(t *testing.T) {
	var addr *net.UDPAddr
	addr, err := net.ResolveUDPAddr("udp", "localhost:8001")

	var a messages.Message = messages.Message{}
	sender := &messages.Contact{}
	sender.ID = "12344321"
	sender.Address = "localhost:8002"
	a.Sender = sender
	a.Type = 0
	aInner := &messages.Request{0,1, "4321"}
	a.Request = aInner

	p, _ := proto.Marshal(&a)

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("ERROR! \n %v\n", err)
		t.Fail()
		return
	}
	var b []byte = make([]byte, 255)
	//RACE!
	go writeByte("localhost:8001", p)
	//RACE!
	num, err := conn.Read(b)
	if err != nil {
		fmt.Printf("READ ERROR! \n %v\n", err)
		t.Fail()
		return
	}

	recieved := &messages.Message{}
	err = proto.Unmarshal(b[:num], recieved)
	if err != nil {
		fmt.Printf("UNMARSHAL ERROR! \n %v\n", err)
		t.Fail()
		return
	}
	for i := 0; i < num; i++ {
		if b[i] != p[i] {
			fmt.Printf("Byte %d differ. Expected %02X, got %02X", i, p[i], b[i])
			t.Fail()
		}
	}
	fmt.Println("Success - Protobuf communication")
	conn.Close()
}

func TestGetResponse(t *testing.T) {
	net := NewNetwork("localhost","8011", nil)
	ID := 2
	a:= func () {
		time.Sleep(100 * time.Millisecond)
		tmp := messages.Response{int64(ID), 1, nil}
		net.responseCond.L.Lock()
		*net.responseList = append(*net.responseList, tmp)
		*net.newResponse = true
		net.responseCond.Broadcast()
		net.responseCond.L.Unlock()
	}
	var bueno bool = true;
	go a()
	response := *net.getResponse(int64(ID))
	//fmt.Printf("%T\n",response.Type)
	//fmt.Println(response)
	bueno = bueno && response.MessageID == int64(ID)
	nilResponse := net.getResponse(int64(11))
	bueno = bueno && nilResponse == nil
	if(bueno){
		fmt.Println("Success - Network getResponse")
	}else {
		t.Fail()
	}
}

func TestComunnications2(t *testing.T) {
	//fmt.Println("222")
	var a1, a2 string = "localhost", "localhost"
	var p1, p2 string = "8025", "8026"
	var net1, net2 Network
	kad1 := NewKademlia(a1+":"+p1, &net1, nil)
	kad2 := NewKademlia(a2+":"+p2, &net2, nil)
	net1 = NewNetwork(a1, p1, kad1)
	net2 = NewNetwork(a2, p2, kad2)
	
	time.Sleep(1 * time.Second)
	var bueno bool = net1.SendPingMessage(&kad2.rt.me)
	if bueno {
		fmt.Println("Success - Network Ping")
	} else {
		t.Fail()
	}
	bueno = true
	var target *KademliaID = NewRandomKademliaID()
	
	var cc CloseContacts = net1.SendFindContactMessage(&kad2.rt.me, target)
	
	bueno = len(cc) == 1
	if(!bueno) { fmt.Printf("Expected length 1, found %d\n", len(cc))}
	bueno = bueno && cc[0].contact.ID.Equals(kad2.rt.me.ID)
	if(!bueno) { fmt.Println("Expected ID %v, found %v\n", kad2.rt.me.ID, cc[0].contact.ID)}
	bueno = bueno && cc[0].distance.Equals(kad2.rt.me.ID.CalcDistance(target))
		if(!bueno) { fmt.Printf("Expected distance %v, found %v\n",kad1.rt.me.ID.CalcDistance(kad2.rt.me.ID), cc[0].distance)}
	if bueno {
		fmt.Println("Success - Network FindContact")
	} else {
		t.Fail()
	}
}
