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
	net := NewNetwork("localhost","8011")
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
	response := net.getResponse(int64(ID))
	//fmt.Printf("%T\n",response.Type)
	//fmt.Println(response)
	bueno = bueno && response.MessageID == int64(ID)
	response = net.getResponse(int64(11))
	bueno = bueno && response.MessageID == 0
	if(bueno){
		fmt.Println("Success - Network getResponse")
	}else {
		t.Fail()
	}
}

func TestComunnications2(t *testing.T) {
	//fmt.Println("222")
	net1 := NewNetwork("localhost","8025")
	net2 := NewNetwork("localhost","8026")
	c1 := NewContact(NewRandomKademliaID(), "localhost:8025")
	c2 := NewContact(NewKademliaID("53FAFFFBB0230099001E200000000C0000000000"), "localhost:8026")
	net1.me = *c1.ID
	net2.me = *c2.ID
	time.Sleep(1000000)
	fmt.Printf("Ping: %v\n", net1.SendPingMessage(&c2))
}
