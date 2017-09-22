package main

import(
	"messages"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"net"
	)

const REQUEST = 0
const RESPONSE= 1
const PING = 0
const FINDNODE = 1
const FINDDATA = 2
const STORE = 3
	
func main(){

	var addr *net.UDPAddr
	addr,err := net.ResolveUDPAddr("udp", "localhost:8001")

	conn, err := net.ListenUDP("udp", addr)
	if(err != nil){
		fmt.Printf("ERROR! \n %v\n", err)
		return	
	}
	var b []byte = make([]byte, 255)
	
	for {
		num, err := conn.Read(b)
		if(err != nil) {
			fmt.Printf("READ ERROR! \n %v\n", err)
			return
		}
		mess := &messages.Message{}
		err = proto.Unmarshal(b[:num], mess)
		if(err != nil){
			fmt.Printf("UNMARSHAL ERROR! \n %v\n", err)
			return	
		}
		if(mess.Type == REQUEST) {
			req := mess.Request
			switch req.Type {
				case PING: fmt.Println("ping")
				case FINDNODE: fmt.Printf("findNode %s\n", req.ID)
				case FINDDATA: fmt.Printf("fintData %s\n", req.ID)
				case STORE: fmt.Println("store %s\n", req.ID)
				default: fmt.Println("sumtin wrong")
			}
		} else if(mess.Type == RESPONSE) {
			fmt.Println("Response")
		}
		fmt.Printf("Bytes: %X\n", b[:num])
	}
	
	conn.Close()
}
