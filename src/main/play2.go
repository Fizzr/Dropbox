package main

import(
	"messages"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"net"
	)

	
func main(){
	a := messages.Message{}
	a.SenderID = "1234"
	a.SenderAddress = "localhost:8002"
	a.Type = 0
	b := &messages.Request{1, "4321"}
	a.Request = b
	
	p, _ := proto.Marshal(&a)

	var laddr, raddr *net.UDPAddr

	laddr, err := net.ResolveUDPAddr("udp", "localhost")
	raddr, err = net.ResolveUDPAddr("udp", "localhost:8001")

	conn, err := net.DialUDP("udp", laddr, raddr)
	if(err != nil){
		fmt.Printf("ERROR! \n %v\n", err)
		return	
	}
//	var b []byte = []byte{0xAB}
	fmt.Println("Press enter to send")
	for {
		var a string = "1"
		fmt.Scanf(a)
		fmt.Println(a)
		num, err := conn.Write(p)
		if(err != nil) {
			fmt.Printf("WRITE ERROR! \n %v\n", err)
			return
		}
		fmt.Printf("Wrote %d bytes\n", num)
	}
	
	conn.Close()
}
