package main

import (
	kademlia "Dropbox1/src/d7024e"
	"fmt"
	"strconv"
)

func createNode(ip string, port int) kademlia.Contact {

	port2string := strconv.Itoa(port)
	newRandomID := kademlia.NewRandomKademliaID()
	newContact := kademlia.NewContact(newRandomID, "localhost"+":"+port2string)

	//  !rt.AddContact(NewContact(NewKademliaID(address), fmt.Sprintf("localhost:%s", 8000 + j+(i*bucketSize))))
	return newContact
}

func threads(numNodes int) {
	port := 8001
	baseContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:8000")
	newRT := kademlia.NewRoutingTable(baseContact)

	for i := 0; i < numNodes; i++ {
		nContact := createNode("localhost", port)
		fmt.Println(nContact)
		newRT.AddContact(nContact)
		port++
	}
	//fmt.Println(newRT)
}

func main() {
	threads(2)
	//c := kademlia.NewContact(kademlia.NewRandomKademliaID(), "localhost:8001")
	//kademlia.SendPingMessage()
	//kademlia.SendPingMessage1()
}
