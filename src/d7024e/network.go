package d7024e

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Net interface {
	SendPingMessage(contact *Contact)
	SendFindContactMessage(contact *Contact) CloseContacts
	SendFindDataMessage(hash string)
	SendStoreMessage(data []byte)
}

type Network struct {
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func Listen(ip string, port int) {
	// TODO
	portcov := strconv.Itoa(port)
	ServerAddr, err := net.ResolveUDPAddr("udp", ip+":"+portcov)
	CheckError(err)

	// Listen to Selected Port
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	fmt.Println("Ip: " + ip + " and Port " + portcov)
	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func SendPingMessage() {
	// func (network *Network) SendPingMessage(contact *Contact)
	ServerAddr, err := net.ResolveUDPAddr("udp", "localhost:8001")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func SendPingMessage1() {
	// func (network *Network) SendPingMessage(contact *Contact)
	ServerAddr, err := net.ResolveUDPAddr("udp", "localhost:8002")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// Query for k contacts closest to contact target
	// Should run synchronous (I guess)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
