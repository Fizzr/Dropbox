package d7024e

import (
	"fmt"
	"net"
	"os"
	proto "github.com/golang/protobuf/proto"
	"time"
	"sync"
	"messages"
)

type responses []messages.Response

type Net interface {
	SendPingMessage(contact *Contact)
	SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts
	SendFindDataMessage(hash string)
	SendStoreMessage(data []byte)
}

type Network struct {
	address string
	port string
	me KademliaID
	mID int64
	IDLock *sync.Mutex
	responseList *responses
	newResponse *bool
	responseCond *sync.Cond
}

//Wakes the threads waiting for a response once every now and then to account for timeout checks
func (net *Network) timeoutCheck(){
	for {
		time.Sleep(1 * time.Second)
		net.responseCond.Broadcast()
	}
}

func NewNetwork (address string, port string) Network {
	b := false
	c := make(responses, 0)
	a:= Network{address: address, port: port, mID: 1, IDLock: &sync.Mutex{}, responseList: &c, newResponse: &b, responseCond: &sync.Cond{L: &sync.Mutex{}}}
	go a.timeoutCheck()
	go a.Listen()
	return a
}

func (net *Network) getMessageID() int64 {
	net.IDLock.Lock()
	ID := net.mID
	net.mID ++
	net.IDLock.Unlock()
	return ID
}

const timeout = 5 * time.Second

func (net *Network) getResponse (ID int64) messages.Response{
	// IMPORTANT! Response might not ever arrive (UDP)! Add some robust timeout options
	var start time.Time = time.Now()
	for {
		net.responseCond.L.Lock()
		defer net.responseCond.L.Unlock()
		for(!*net.newResponse) {
			if(time.Since(start) > timeout) {
				//fmt.Println("i ded")
				return messages.Response{}
			}
			net.responseCond.Wait()
		}
		for i := 0; i < len(*net.responseList); i ++ {
			if((*net.responseList)[i].MessageID == ID){
				a := (*net.responseList)[i]
				*net.responseList = append((*net.responseList)[:i], (*net.responseList)[i+1:]...)
				if len(*net.responseList) == 0 {
					*net.newResponse = false
				}
				return a
			}
		}
	}
	return messages.Response{}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func (network *Network) addResponse (response messages.Response) {
	network.responseCond.L.Lock()
	*network.responseList = append(*network.responseList, response)
	*network.newResponse = true
	network.responseCond.Broadcast()
	network.responseCond.L.Unlock()
}

func (network *Network) Listen() {
	ServerAddr, err := net.ResolveUDPAddr("udp", network.address+":"+network.port)
	CheckError(err)

	// Listen to Selected Port
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()
	//fmt.Println("Ip: " + network.address + " and Port " + network.port)
	buff := make([]byte, 1024)

	for {
		n, err := ServerConn.Read(buff)
		/*fmt.Printf("Received %d bytes. Buff len %d \n", n, len(buff))
		for i:= 0; i < n; i++ {
			fmt.Print(buff[i])
		}
		fmt.Println()
		*/
		if err != nil {
			fmt.Println("Error: ", err)
		}
		
		var received *messages.Message = &messages.Message{}
		err = proto.Unmarshal(buff[:n], received)
		CheckError(err)
		if(received.Type == messages.Message_RESPONSE) {
			go network.addResponse(*received.Response)
		} else if (received.Type == messages.Message_REQUEST) {
			//TODO: Request code!
			request := received.Request
			switch request.Type {
				case messages.Request_PING:
						network.respondPingMessage(*received)
					break
				default:
					fmt.Println("Error: Unknown request type")
			}
		} else {
			fmt.Println("Error: Not valid message type!")
		}
	}
}

func (network *Network) newRequestMessage() messages.Message{
	var msg messages.Message = messages.Message{}
	msg.Type = messages.Message_REQUEST
	var me messages.Contact = messages.Contact{fmt.Sprint(network.me), network.address + ":" + network.port}
	msg.Sender = &me
	return msg
}

func (network *Network) respondPingMessage(received messages.Message) {
	//fmt.Println("respondPing")
	//fmt.Println("messege received ", received)
	ServerAddr, err := net.ResolveUDPAddr("udp", received.Sender.Address)
	CheckError(err)

//	LocalAddr, err := net.ResolveUDPAddr("udp", network.address + ":" + network.port)
	CheckError(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	var msg messages.Message = messages.Message{}
	msg.Type = messages.Message_RESPONSE
	var me messages.Contact = messages.Contact{fmt.Sprint(network.me), network.address + ":" + network.port}
	//fmt.Println(network.me)
	msg.Sender = &me

	var ping messages.Response = messages.Response{received.Request.MessageID, messages.Response_PING, nil}
	msg.Response = &ping
	//fmt.Println("messege to send ",msg)
	var buff []byte
	buff, err = proto.Marshal(&msg)
	CheckError(err)

	_, err = Conn.Write(buff)
	if err != nil {
		fmt.Println(msg, err)
	}	
}

func (network *Network) SendPingMessage(contact *Contact) bool{
	//portconv := strconv.Itoa(port)
	//fmt.Println("sendPing")
	ServerAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)

//	LocalAddr, err := net.ResolveUDPAddr("udp", network.address + ":" + network.port)
	CheckError(err)

	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)

	defer Conn.Close()
	var msg messages.Message = network.newRequestMessage()
	var mID int64 = network.getMessageID()
	var ping messages.Request = messages.Request{mID, messages.Request_PING, ""}
	msg.Request = &ping
	
	var buff []byte
	buff, err = proto.Marshal(&msg)
	CheckError(err)
	_, err = Conn.Write(buff)
	if err != nil {
		fmt.Println(msg, err)
	}
	var response messages.Response = network.getResponse(mID)
	//fmt.Printf("ID should be %v, is %v, type %v\n", mID, response.MessageID, response.Type)
	//fmt.Println(response)
	if(response.MessageID == mID && response.Type == messages.Response_PING) {
		return true
	}
	return false;
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
