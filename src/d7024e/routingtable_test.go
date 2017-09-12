package d7024e

import (
	"fmt"
	"testing"
	"math/rand"
	//"encoding/hex"
)

var testList []string = []string{
	"FA10BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1",
	"0011111400000000000000000000000000000000",
	"F001111400000000000000000000000000000000",
	"0111111400000000000000000000000000000000",
	"FFFFFFFF00000000000000000000000000000000",
	"F111111100000000000000000000000000000000",
	"F111111200000000000000000000000000000000",
	"F111111300000000000000000000000000000000",
	"2111111400000000000000000000000000000000",
	"3111111400000000000000000000000000000000",
	"4111111400000000000000000000000000000000",
	"5FFFFFFF00000000000000000000000000000000",
	"6111111100000000000000000000000000000000",
	"7111111200000000000000000000000000000000",
	"8111111300000000000000000000000000000000",
	"9111111400000000000000000000000000000000",
	"AFFFFFFF00000000000000000000000000000000",
	"B111111100000000000000000000000000000000",
	"C111111200000000000000000000000000000000",
	"D111111300000000000000000000000000000000",
	"E111111300000000000000000000000000000000"}

func TestKadmeliaIDbitAt(t *testing.T){
	var ID *KademliaID = NewKademliaID(testList[0])
	
	var address string
	var bits byte
	for i := IDLength*8 -1; i >= 0; i--{
		//fmt.Printf("%01b", c.ID.bitAt(i))
		var num uint = uint(i % 8)
		if(num == 0){
			address = fmt.Sprintf("%s%01X", address, bits)
			bits = 0x0
		}
		bits = bits | ID.bitAt(i) << num
		
	}
//	fmt.Printf("\n%s\n", address)
	if(address == testList[0]){
		fmt.Println("Success - KadmeliaID bitAt")
	}else {
		t.Fail()
	}
}

const letters = "0123456789ABCDEF"

func randomHex(n int) string{
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func strToByte(str string) byte{
	switch str{
		case "0":
			return 0x0
		case "1":
			return 0x1
		case "2":
			return 0x2
		case "3":
			return 0x3
		case "4":
			return 0x4
		case "5":
			return 0x5
		case "6":
			return 0x6
		case "7":
			return 0x7
		case "8":
			return 0x8
		case "9":
			return 0x9
		case "A":
			return 0xA
		case "B":
			return 0xB
		case "C":
			return 0xC
		case "D":
			return 0xD
		case "E":
			return 0xE
		case "F":
			return 0xF
	}
	return 0x0
}
	
func TestRoutingTable(t *testing.T) {
	var c Contact = NewContact(NewKademliaID(testList[0]), "localhost:8000")
	
	rt := NewRoutingTable(c)
	levels := 5 					//number of levels down we'll go
	//var start string
	for i:= 0; i < levels; i ++ {	// i = level we're at. i.e. what exponent we're differating on!
		/*if(i%4 == 0){
			tmp := i/4
			start += testList[0][tmp:tmp+1]
			fmt.Println(start)
		}*/
		
		var current, after, before int
		before = i/4
		after = 39 - before		//hexes after active hex
		current = before		//index of hex that's active
		
		var active byte = strToByte(testList[0][current:current+1])
		active = active ^ 1 << uint(3-(i%4))
		
		var start string = testList[0][:current]
		
		for j:= 0; j < bucketSize; j++{
			var tail string = randomHex(after)
			var address string = start + fmt.Sprintf("%01X", active) + tail
			id := NewKademliaID(address)
			fmt.Println(id.toBinary())
			rt.AddContact(NewContact(NewKademliaID(address), fmt.Sprintf("localhost:8%03d", j+(levels*bucketSize))))
		}
	}
	
	
	/*for i := 0; i < len(testList); i++{
		contact := NewContact(NewKademliaID(testList[i]), fmt.Sprintf("localhost:8%03d", i))
		//var bits string
		//for j := 0; j < IDLength; j++ {
		//	bits += fmt.Sprintf("%08b", contact.ID[j])
		//}
		//fmt.Println(bits)
		rt.AddContact(contact)
	}*/
	//contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	//for i := range contacts {
	//	fmt.Println(contacts[i].String())
	//}
	fmt.Println("")
	fmt.Println(rt.root)
	fmt.Println("")
	
	fmt.Printf("%T, %T, %T \n", rt.root.(*Branch).left, rt.root, rt.root.(*Branch).right)
}
