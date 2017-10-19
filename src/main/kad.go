package main

import (
	"fmt"
	"os"
	"strconv"
	//	kademlia "d7024e"
	"encoding/hex"
	"io/ioutil"
	"messages"
	"net"
	"time"

	proto "github.com/golang/protobuf/proto"
)

func checkID(ID string) bool {
	if len(ID) != 40 {
		fmt.Println("Argument must be 40 charachters long")
		return false
	}
	_, err := hex.DecodeString(ID)
	if err != nil {
		fmt.Println("Argument not in correct hex format")
		return false
	}
	return true
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func getStuff(port int) (net.Conn, messages.Message) {
	var msg messages.Message = messages.Message{}
	ServerAddr, err := net.ResolveUDPAddr("udp", "localhost:"+fmt.Sprint(port))
	CheckError(err)
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	Conn.SetDeadline(time.Now().Add(5 * time.Second))
	retAddr := Conn.LocalAddr().String()
	var contact messages.Contact = messages.Contact{"", retAddr, ""}
	var request messages.Request = messages.Request{}
	msg.Sender = &contact
	msg.Request = &request
	msg.Type = messages.Message_REQUEST
	return Conn, msg
}
func pinMessage(pin bool, hash string, port int) {
	conn, msg := getStuff(port)
	defer conn.Close()
	var t messages.Request_Type
	if pin {
		t = messages.Request_CLIENT_PIN
	} else {
		t = messages.Request_CLIENT_UNPIN
	}
	msg.Request.Type = t
	msg.Request.ID = hash
	var buff []byte
	buff, err := proto.Marshal(&msg)
	CheckError(err)
	_, err = conn.Write(buff)
	CheckError(err)
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Must be in form\n\t kad <port> <command> <argument>")
		return
	}
	port, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Error converting", args[0], "to integer")
		return
	}

	switch args[1] {
	case "Store":
		fallthrough
	case "store":
		//fmt.Println("storing")
		//Find file and parse it
		dat, err := ioutil.ReadFile(args[2])
		CheckError(err)
		conn, msg := getStuff(port)
		retAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())

		msg.Request.Type = messages.Request_CLIENT_STORE
		msg.Request.Data = dat
		
		var buff []byte
		var returned []byte = make([]byte, 4096)
		buff, err = proto.Marshal(&msg)
		CheckError(err)
		
		n, err := conn.Write(buff)
		CheckError(err)
		conn.Close()
		
		ServerConn, err := net.ListenUDP("udp", retAddr)
		ServerConn.SetDeadline(time.Now().Add(5*time.Second))
		defer ServerConn.Close()
		n, err = ServerConn.Read(returned)
		CheckError(err)
		cstore := messages.CSTORE{}
		err = proto.Unmarshal(returned[:n], &cstore)
		CheckError(err)
		fmt.Println(cstore.Hash)

	case "Pin":
		fallthrough
	case "pin":
		//fmt.Println("pinning")
		if !checkID(args[2]) {
			return
		}
		pinMessage(true, args[2], port)

	case "Unpin":
		fallthrough
	case "unpin":
		//fmt.Println("unpinning")
		if !checkID(args[2]) {
			return
		}
		pinMessage(false, args[2], port)

	case "Cat":
		fallthrough
	case "cat":
		//fmt.Println("Kitty Cat!")
		if !checkID(args[2]) {
			return
		}
		conn, msg := getStuff(port)
		retAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
		
		msg.Request.Type = messages.Request_CLIENT_LOOKUP
		msg.Request.ID = args[2]
		
		var buff []byte
		var returned []byte = make([]byte, 4096)
		
		buff, err = proto.Marshal(&msg)
		CheckError(err)
		n, err := conn.Write(buff)
		CheckError(err)
		conn.Close()
		
		ServerConn, err := net.ListenUDP("udp", retAddr)
		ServerConn.SetDeadline(time.Now().Add(5*time.Second))
		n, err = ServerConn.Read(returned)
		CheckError(err)
		clook := messages.CLOOKUP{}
		err = proto.Unmarshal(returned[:n], &clook)
		CheckError(err)
		var out string = string(clook.Data)
		fmt.Println(out)
		
	case "Local":
		fallthrough
	case "local":
		if args[2] != "data" {
			fmt.Printf("Unknown argument %s. Did you mean 'data'?\n", args[2])
			return
		}
		//fmt.Println("Localing")
		//ask target for local data
		conn, msg := getStuff(port)
		retAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
		
		msg.Request.Type = messages.Request_CLIENT_LOCAL

		var returned []byte = make([]byte, 4096)		
		var buff []byte
		buff, err = proto.Marshal(&msg)
		CheckError(err)
		n, err := conn.Write(buff)
		CheckError(err)
		conn.Close()
		
		ServerConn, err := net.ListenUDP("udp", retAddr)
		ServerConn.SetDeadline(time.Now().Add(5*time.Second))
		n, err = ServerConn.Read(returned)
		CheckError(err)

		cloc := messages.CLOCAL{}
		err = proto.Unmarshal(returned[:n], &cloc)
		CheckError(err)
		fmt.Println("Pinned Data")
		for i := 0; i < len(cloc.Mine); i++ {
			fmt.Println(cloc.Mine[i])
		}
		fmt.Println("Unpinned Data")
		for i := 0; i < len(cloc.Other); i++ {
			fmt.Println(cloc.Other[i])
		}

	default:
		fmt.Println("Unknown command", args[1])
	}
}
