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

const timeout = 1 * time.Second

type responses []messages.Response

type Net interface {
	SendPingMessage(contact *Contact) bool
	SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts
	SendFindDataMessage(contact *Contact, hash string) (*CloseContacts, *[]byte)
	SendStoreMessage(contact *Contact, data []byte)
}

type Network struct {
	address string
	port string
	kad *Kademlia
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

func NewNetwork (address string, port string, kad *Kademlia) Network {
	b := false
	c := make(responses, 0)
	a:= Network{address: address, port: port, mID: 1, kad: kad, IDLock: &sync.Mutex{}, responseList: &c, newResponse: &b, responseCond: &sync.Cond{L: &sync.Mutex{}}}
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


func (net *Network) getResponse (ID int64) *messages.Response{
	// IMPORTANT! Response might not ever arrive (UDP)!
	// Also, add timeout for messages in queue, to save memory and prevent unecessary looping
	var start time.Time = time.Now()
	for {
		//fmt.Print("brap")
		if(time.Since(start) > timeout) {
			return nil
		}
		net.responseCond.L.Lock()
		for(!*net.newResponse) {
			if(time.Since(start) > timeout) {
				net.responseCond.L.Unlock()
				return nil
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
				net.responseCond.L.Unlock()
				return &a
			}
		}
		net.responseCond.L.Unlock()
	}
	return nil
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func (network *Network) addResponse (response messages.Response) {
	network.responseCond.L.Lock()
	//fmt.Println("Twerk. Adding ID ", response.MessageID)
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
	buff := make([]byte, 4096)	//This number is pretty arbitrary. But it fits 20 contacts being returned! 4kb might not fit data returns tho...

	for {
		n, err := ServerConn.Read(buff)
		//fmt.Printf("Received %d bytes. Buff len %d \n", n, len(buff))
		/*for i:= 0; i < n; i++ {
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
		
		var sender Contact = NewContact(NewKademliaID(received.Sender.ID), received.Sender.Address)
		go network.kad.rt.AddContact(sender)
		
		if(received.Type == messages.Message_RESPONSE) {
			//fmt.Println("Listened response")
			go network.addResponse(*received.Response)
		} else if (received.Type == messages.Message_REQUEST) {
			//fmt.Println("Listened request")
			switch received.Request.Type {
				case messages.Request_PING:
					go network.respondPingMessage(*received)
					break
				case messages.Request_FINDNODE:
					go network.respondFindNodeMessage(*received)					
					break
				case messages.Request_FINDDATA:
					break
				case messages.Request_STORE: 
					break
				default:
					fmt.Println("Error: Unknown request type")
			}
		} else {
			fmt.Println("Error: Not valid message type!")
		}
	}
}

func (network *Network) newMessage(typ messages.Message_Type) messages.Message{
	var msg messages.Message = messages.Message{}
	msg.Type = typ
	var me messages.Contact = messages.Contact{fmt.Sprint(network.kad.rt.me.ID), network.address + ":" + network.port, ""}
	msg.Sender = &me
	return msg
}

func (network *Network) newResponseMessage() messages.Message{
	return network.newMessage(messages.Message_RESPONSE)
}

func (network *Network) newRequestMessage() messages.Message{
	return network.newMessage(messages.Message_REQUEST)
}

func (network *Network) respondPingMessage(received messages.Message) {
	ServerAddr, err := net.ResolveUDPAddr("udp", received.Sender.Address)
	CheckError(err)
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	var msg messages.Message = network.newResponseMessage()

	var ping messages.Response = messages.Response{received.Request.MessageID, messages.Response_PING, nil, nil}
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

	ServerAddr, err := net.ResolveUDPAddr("udp", contact.Address)
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
	//fmt.Printf("wrote %d bytes\n", n)
	if err != nil {
		fmt.Println(msg, err)
	}
	var response *messages.Response = network.getResponse(mID)
	if(response == nil){
		return false
	}
	//fmt.Printf("ID should be %v, is %v, type %v\n", mID, response.MessageID, response.Type)
	//fmt.Println(response)
	if(response.MessageID == mID && response.Type == messages.Response_PING) {
		return true
	}
	return false;
}

func (network *Network) respondFindNodeMessage(received messages.Message) {
	//fmt.Println("respond Find")
	var cc CloseContacts = network.kad.rt.FindClosestContacts(NewKademliaID(received.Request.ID), k)
	var msg messages.Message = network.newResponseMessage()
	var response messages.Response = messages.Response{}
	response.Type = messages.Response_FINDNODE
	for i := 0; i < len(cc); i++ {
		var cont messages.Contact = messages.Contact{fmt.Sprint(cc[i].contact.ID), cc[i].contact.Address, fmt.Sprint(cc[i].distance)}
		response.Contacts = append(response.Contacts, &cont)
	}
	response.MessageID = received.Request.MessageID
	msg.Response = &response
	var buff []byte
	buff, err := proto.Marshal(&msg)
	CheckError(err)
	
	ServerAddr, err := net.ResolveUDPAddr("udp", received.Sender.Address)
	CheckError(err)
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	_,err = Conn.Write(buff)
	//fmt.Printf("Responded with %d bytes\n", n)
	CheckError(err)
}

func (network *Network) SendFindContactMessage(contact *Contact, target *KademliaID) CloseContacts{
	ServerAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)	
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	var msg messages.Message = network.newRequestMessage()
	var mID int64 = network.getMessageID()
	msg.Request = &messages.Request{mID, messages.Request_FINDNODE, fmt.Sprint(target)}
	
	var buff []byte
	buff, err = proto.Marshal(&msg)
	CheckError(err)
	_, err = Conn.Write(buff)
	//fmt.Printf("wrote %d bytes\n", n)
	if err != nil {
		fmt.Println(msg, err)
	}
	//fmt.Println("sent find. Waiting for ID ", mID)
	var response *messages.Response = network.getResponse(mID)
	//fmt.Println("response ", response)
	if(response == nil){
		return nil
	}
	var res CloseContacts
	for i := 0; i < len(response.Contacts); i ++ {
		res = append(res, CloseContact{Contact{NewKademliaID(response.Contacts[i].ID), response.Contacts[i].Address}, NewKademliaID(response.Contacts[i].Distance)})
	}
	return res
}

func (network *Network) respondFindDataMessage(received messages.Message) {
	var msg messages.Message = network.newResponseMessage()
	var response messages.Response = messages.Response{}
	response.MessageID = received.Request.MessageID
	
	//FIND DATA IN FILE 
	var data []byte
	var dataFound bool = false
	
	if(dataFound) {
		response.Data = data
		response.Type = messages.Response_FINDDATA_FOUND
	} else {
		response.Type = messages.Response_FINDDATA_NODES
		var cc CloseContacts = network.kad.rt.FindClosestContacts(NewKademliaID(received.Request.ID), k)
		for i := 0; i < len(cc); i++ {
			var cont messages.Contact = messages.Contact{fmt.Sprint(cc[i].contact.ID), cc[i].contact.Address, fmt.Sprint(cc[i].distance)}
			response.Contacts = append(response.Contacts, &cont)
		}
	}
	msg.Response = &response
	var buff []byte
	buff, err := proto.Marshal(&msg)
	CheckError(err)
	
	ServerAddr, err := net.ResolveUDPAddr("udp", received.Sender.Address)
	CheckError(err)
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	_,err = Conn.Write(buff)
	//fmt.Printf("Responded with %d bytes\n", n)
	CheckError(err)
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) (*CloseContacts, *[]byte) {
	ServerAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	CheckError(err)	
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	var msg messages.Message = network.newRequestMessage()
	var mID int64 = network.getMessageID()
	msg.Request = &messages.Request{mID, messages.Request_FINDDATA, hash}
	
	var buff []byte
	buff, err = proto.Marshal(&msg)
	CheckError(err)
	_, err = Conn.Write(buff)
	//fmt.Printf("wrote %d bytes\n", n)
	if err != nil {
		fmt.Println(msg, err)
	}
	var response *messages.Response = network.getResponse(mID)
	//fmt.Println("response ", response)
	if(response == nil){
		return nil, nil
	}
	if(response.Type == messages.Response_FINDDATA_FOUND) {
		return nil, &response.Data
	} else if (response.Type == messages.Response_FINDDATA_NODES) {
		var res CloseContacts
		for i := 0; i < len(response.Contacts); i ++ {
			res = append(res, CloseContact{Contact{NewKademliaID(response.Contacts[i].ID), response.Contacts[i].Address}, NewKademliaID(response.Contacts[i].Distance)})
		}
		return &res, nil
	} else {
		fmt.Println("Error: Mismatched response type in FINDDATA! Received ", response.Type)
		return nil, nil
	}
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) {
	// TODO
}
