package d7024e

import (
	"fmt"
	"messages"
	"net"
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

	a := messages.Message{}
	a.SenderID = "1234"
	a.SenderAddress = "localhost:8002"
	a.Type = 0
	aInner := &messages.Request{1, "4321"}
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

/*/Without Protobuf Implementation
func TestListenPing(t *testing.T) {
	Listen("localhost", 8001)
	SendPingMessage("localhost", 8001)
}*/
