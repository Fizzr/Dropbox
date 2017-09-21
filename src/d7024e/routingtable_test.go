package d7024e

import (
	"testing"
	"math/rand"
	"math"
	//"encoding/hex"
)

// var testList []string = []string{
// 	"A82FBAAAAAAAAAAAAAA555555555555555555551",
// 	"0011111400000000000000000000000000000000",
// 	"F001111400000000000000000000000000000000",
// 	"0111111400000000000000000000000000000000",
// 	"FFFFFFFF00000000000000000000000000000000",
// 	"F111111100000000000000000000000000000000",
// 	"F111111200000000000000000000000000000000",
// 	"F111111300000000000000000000000000000000",
// 	"2111111400000000000000000000000000000000",
// 	"3111111400000000000000000000000000000000",
// 	"4111111400000000000000000000000000000000",
// 	"5FFFFFFF00000000000000000000000000000000",
// 	"6111111100000000000000000000000000000000",
// 	"7111111200000000000000000000000000000000",
// 	"8111111300000000000000000000000000000000",
// 	"9111111400000000000000000000000000000000",
// 	"AFFFFFFF00000000000000000000000000000000",
// 	"B111111100000000000000000000000000000000",
// 	"C111111200000000000000000000000000000000",
// 	"D111111300000000000000000000000000000000",
// 	"E111111300000000000000000000000000000000"}
//
// //Test functions for branches and leafes for use in TestRoutingTable
// func (branch *Branch) verifyFullTree(i int) int {
// 	var left, right int
// 	//to get around the fact that I can't define functions for interfaces here
// 	switch a := branch.left.(type) {
// 	case *Leaf:
// 		left = a.verifyFullTree(i + 1)
// 	case *Branch:
// 		left = a.verifyFullTree(i + 1)
// 	}
// 	switch a := branch.right.(type) {
// 	case *Leaf:
// 		right = a.verifyFullTree(i + 1)
// 	case *Branch:
// 		right = a.verifyFullTree(i + 1)
// 	}
// 	if left == -1 {
// 		return -1
// 	}
// 	if right == -1 {
// 		return -1
// 	}
// 	var big, small int
// 	if left > right {
// 		big, small = left, right
// 	} else {
// 		big, small = right, left
// 	}
// 	if small-i != 1 {
// 		fmt.Printf("Exp %v, big %v, small %v i %v\n", branch.exponent, big, small, i)
// 		fmt.Println("Incorrect tree structure!")
// 		return -1
// 	} else {
// 		return big
// 	}
// }
//
// func (leaf *Leaf) verifyFullTree(i int) int {
// 	if leaf.buck.Len() != bucketSize {
// 		if leaf.buck.Len() == 1 {
// 			if fmt.Sprint(leaf.buck.list.Front().Value.(Contact).ID) == testList[0] {
// 				return i
// 			}
// 		}
// 		fmt.Printf("Not right length! Expected 1 or 5, got %v\n", leaf.buck.Len())
// 		fmt.Printf("ID %v, Exp %v\n %v\n ", leaf.ID, leaf.exponent, leaf)
// 		return -1
// 	}
// 	for e := leaf.buck.list.Front(); e != nil; e = e.Next() {
// 		var ID *KademliaID = e.Value.(Contact).ID
// 		for i := IDBits - 1; i > leaf.exponent; i-- {
// 			//fmt.Println(i)
// 			var preIndex int = IDLength - 1 - i/8
// 			//fmt.Printf("bitAt: %v prefix: %v preIndex %v i %v shift %v\n", ID.bitAt(i), ((leaf.prefix[preIndex/8] >> uint(i%8)) & 1), preIndex, i, preIndex%8)
// 			if (ID.bitAt(i)) != ((leaf.prefix[preIndex] >> uint(i%8)) & 1) {
// 				fmt.Printf("start of ID isn't identical to prefix. Exponent: %v bitAt: %v prefix bit: %v\nPrefix at %v\n", i, ID.bitAt(i), ((leaf.prefix[preIndex/8] >> uint(i%8)) & 1), ID.toBinary())
// 				return -1
// 			}
// 		}
// 	}
// 	return i
// }
//
// func TestKadmeliaIDbitAt(t *testing.T) {
//
// 	//a := NewKademliaID(testList[0])
// 	//for i := 88; i < 95; i++{
// 	//	fmt.Printf("%b", a.bitAt(i))
// 	//}
// 	//fmt.Printf("\n%s\n", a.toBinary()[88:95])
//
// 	try := func(s string) bool {
// 		var ID *KademliaID = NewKademliaID(s)
// 		var address string
// 		for i := IDBits - 1; i >= 0; i-- {
// 			address = address + fmt.Sprintf("%b", ID.bitAt(i))
// 		}
// 		/*fmt.Println(ID.toBinary())
// 		var address string
// 		var bits byte
// 		for i := IDLength*8 -1; i >= 0; i--{
// 			fmt.Printf("%b", ID.bitAt(i))
// 			var num uint = uint(i % 8)
// 			if(num == 0){
// 				address = fmt.Sprintf("%s%02X", address, bits)
// 				bits = 0x00
// 			}
// 			bits = bits | ID.bitAt(i) << num
//
// 		}
// 		fmt.Printf("\n%s - calculated\n%s - should be\n", address, s)*/
// 		if address == ID.toBinary() {
// 			return true
// 		} else {
// 			return false
// 		}
// 	}
//
// 	var pass = true
// 	pass = pass && try(testList[0])
// 	pass = try(randomHex(40)) && pass
// 	pass = try(randomHex(40)) && pass
// 	pass = try(randomHex(40)) && pass
// 	if pass {
// 		fmt.Println("Success - KadmeliaID bitAt")
// 	} else {
// 		t.Fail()
// 	}
// }
//
// const letters = "0123456789ABCDEF"
//
// func randomHex(n int) string {
// 	b := make([]byte, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }
//
// func strToByte(str string) byte {
// 	switch str {
// 	case "0":
// 		return 0x0
// 	case "1":
// 		return 0x1
// 	case "2":
// 		return 0x2
// 	case "3":
// 		return 0x3
// 	case "4":
// 		return 0x4
// 	case "5":
// 		return 0x5
// 	case "6":
// 		return 0x6
// 	case "7":
// 		return 0x7
// 	case "8":
// 		return 0x8
// 	case "9":
// 		return 0x9
// 	case "A":
// 		return 0xA
// 	case "B":
// 		return 0xB
// 	case "C":
// 		return 0xC
// 	case "D":
// 		return 0xD
// 	case "E":
// 		return 0xE
// 	case "F":
// 		return 0xF
// 	}
// 	return 0x0
// }
//
// func TestSplitOn(t *testing.T) {
// 	var b *bucket = newBucket()
// 	for i := 0; i < bucketSize/2; i++ {
// 		b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("A%039X", i)), ""))
// 	}
// 	for i := 0; i < bucketSize/2; i++ {
// 		b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("0%039X", i)), ""))
// 	}
// 	var pass bool = true
// 	var bucks [2]bucket = b.splitOn(159)
// 	for e := bucks[0].list.Front(); e != nil; e = e.Next() {
// 		pass = pass && e.Value.(Contact).ID.bitAt(159) == 0
// 	}
// 	for e := bucks[1].list.Front(); e != nil; e = e.Next() {
// 		pass = pass && e.Value.(Contact).ID.bitAt(159) == 1
// 	}
// 	if pass {
// 		fmt.Println("Success - bucket splitOn")
// 	} else {
// 		t.Fail()
// 	}
//
// }
//
// func TestRoutingTable(t *testing.T) {
// 	var c Contact = NewContact(NewKademliaID(testList[0]), "localhost:8000")
//
// 	rt := NewRoutingTable(c)
// 	levels := 153 //number of levels down we'll go
// 	//var start string
// 	for i := 0; i < levels; i++ { // i = level we're at. i.e. what exponent we're differating on!
// 		/*if(i%4 == 0){
// 			tmp := i/4
// 			start += testList[0][tmp:tmp+1]
// 			fmt.Println(start)
// 		}*/
//
// 		var current, after, before int
// 		before = i / 4
// 		after = 39 - before //hexes after active hex
// 		current = before    //index of hex that's active
//
// 		var active byte = strToByte(testList[0][current : current+1])
// 		active = active ^ 1<<uint(3-(i%4))
//
// 		var start string = testList[0][:current]
//
// 		for j := 0; j < bucketSize; j++ {
// 			var tail string = randomHex(after)
// 			var address string = start + fmt.Sprintf("%01X", active) + tail
// 			//id := NewKademliaID(address)
// 			//fmt.Println(id.toBinary())
// 			if !rt.AddContact(NewContact(NewKademliaID(address), fmt.Sprintf("localhost:%s", 8000+j+(i*bucketSize)))) {
// 				fmt.Printf("Failed to add level %v, number %v\n", i, j)
// 			}
// 		}
// 	}
//
// 	//	fmt.Println(rt.root)
// 	var v int = rt.root.(*Branch).verifyFullTree(1)
// 	if v != levels+1 {
// 		if v != -1 {
// 			fmt.Printf("Expected height %v, got %v\n", levels+1, v)
// 		}
// 		t.Fail()
// 	} else {
// 		fmt.Println("Success - RoutingTable")
// 	}
// }
//
// // Test Function to check Distance between 2 Nodes.
// func TestDistFunc(t *testing.T) {
//
// 	DistT1 := NewKademliaID(testList[0])
// 	DistT2 := NewKademliaID(testList[5])
// 	DistBetween := DistT2.CalcDistance(DistT1)
// 	fmt.Println(DistBetween)
// }

func (leaf *Leaf) verifyFullTree(i int) int {
	if(leaf.buck.Len() != bucketSize){ 
		if(leaf.buck.Len() == 1){
			if fmt.Sprint(leaf.buck.list.Front().Value.(Contact).ID) == testList[0] {
					return i
			}
		}
		fmt.Printf("Not right length! Expected 1 or %d, got %v\n", bucketSize, leaf.buck.Len())
		fmt.Printf("ID %v, Exp %v\n %v\n ", leaf.ID, leaf.exponent, leaf)
		return -1
		}
	for e:= leaf.buck.list.Front(); e != nil; e = e.Next() {
		var ID *KademliaID = e.Value.(Contact).ID
		for i := IDBits-1; i > leaf.exponent; i--{
			//fmt.Println(i)
			var preIndex int = IDLength-1 - i/8
			//fmt.Printf("bitAt: %v prefix: %v preIndex %v i %v shift %v\n", ID.bitAt(i), ((leaf.prefix[preIndex/8] >> uint(i%8)) & 1), preIndex, i, preIndex%8)
			if (ID.bitAt(i)) != ((leaf.prefix[preIndex] >> uint(i%8)) & 1){
				fmt.Printf("start of ID isn't identical to prefix. Exponent: %v bitAt: %v prefix bit: %v\nPrefix at %v\n",i, ID.bitAt(i), ((leaf.prefix[preIndex/8] >> uint(i%8)) & 1), ID.toBinary())
				return -1
			}
		}  
	}
	return i
}


const letters = "0123456789ABCDEF"

func randomHex(n int) string{
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
	
func TestKadmeliaIDbitAt(t *testing.T){
	
	try := func (s string) bool {
		var ID *KademliaID = NewKademliaID(s)
		var address string
		for i := IDBits-1; i >= 0; i-- {
			address = address + fmt.Sprintf("%b", ID.bitAt(i))
		}
		/*fmt.Println(ID.toBinary())
		var address string
		var bits byte
		for i := IDLength*8 -1; i >= 0; i--{
			fmt.Printf("%b", ID.bitAt(i))
			var num uint = uint(i % 8)
			if(num == 0){
				address = fmt.Sprintf("%s%02X", address, bits)
				bits = 0x00
			}
			bits = bits | ID.bitAt(i) << num
			
		}
		fmt.Printf("\n%s - calculated\n%s - should be\n", address, s)*/
		if(address == ID.toBinary()){
			return true
		}else {
			return false
		}
	}
	
	var pass = true
	pass = pass && try(testList[0])
	
	for i := 0; i < 160; i++ {
		pass = pass && try(randomHex(40))
	}
	
	if(pass){
		fmt.Println("Success - KadmeliaID bitAt")
	}else {
		t.Fail()
	}
}


func TestSplitOn(t *testing.T){
	try := func(exponent int) bool {
		var b *bucket = newBucket()
		for i := 0; i < bucketSize/2; i ++{
			//b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("%s%d%s", randomHex((IDLength*2)- 1 - exponent/4), int(math.Pow(2,float64(exponent%4))), randomHex(exponent/4))), ""))
			b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("%0s%d%0s", randomHex((IDLength*2)- 1 - exponent/4), int(math.Pow(2,float64(exponent%4))), randomHex(exponent/4))), ""))
		}
		for i := 0; i < bucketSize/2; i ++{
			b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("%s0%s", randomHex((IDLength*2)- 1 - exponent/4), randomHex(exponent/4))), ""))
			//b.AddContact(NewContact(NewKademliaID(fmt.Sprintf("%040X", i)), ""))
		}
		var pass bool = true
		var	 bucks [2]bucket = b.splitOn(exponent)
		var num int = 0
		for e := bucks[0].list.Front(); e != nil; e = e.Next() {
			pass = pass && e.Value.(Contact).ID.bitAt(exponent) == 0
			num ++
		}
		if(num != bucketSize/2){ pass = false }
		num = 0
		for e := bucks[1].list.Front(); e != nil; e = e.Next() {
			pass = pass && e.Value.(Contact).ID.bitAt(exponent) == 1
			num ++
		}
		if(num != bucketSize/2){ pass = false }
		return pass
	}
	
	pass := true
	for i := 0; i < 160; i++ {
		pass = pass && try(i)
	}
	
	if(pass){
		fmt.Println("Success - bucket splitOn")
	}else{
		t.Fail()
	}

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
	levels := 130 					//number of levels down we'll go. Can't go too low,
									// or we will get address space collisions, and won't be able to fill the tree.
	
	for i:= 0; i < levels; i ++ {	// i = level we're at. i.e. what exponent we're differating on!
		
		var current, after, before int
		before = i/4
		after = (IDBits/4)-1 - before		//hexes after active hex
		current = before		//index of hex that's active
		
		var active byte = strToByte(testList[0][current:current+1])
		active = active ^ 1 << uint(3-(i%4))
		
		var start string = testList[0][:current]
		
		for j:= 0; j < bucketSize; j++{
			var tail string = randomHex(after)
			var address string = start + fmt.Sprintf("%01X", active) + tail
			//id := NewKademliaID(address)
			//fmt.Println(id.toBinary())
			var ok, added = rt.AddContact(NewContact(NewKademliaID(address), fmt.Sprintf("localhost:%s", 8000 + j+(i*bucketSize))))
			if ok {
				if !added {
					fmt.Printf("Failed to add number %d on level %d\n", j, i)
				}
			} else {
				t.Error("AddContact returned false. Not supposed to happen!")
			}
		}
	}
	
	//	fmt.Println(rt.root)
	var a CloseContacts = rt.FindClosestContacts(NewKademliaID(testList[0]), 5)
	fmt.Println(testList[0])
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i].contact.ID)
	}
	fmt.Println(len(a))
	var v int = rt.root.(*Branch).verifyFullTree(1)
	if(v != levels+1){
		if (v != -1){fmt.Printf("Expected height %v, got %v\n", levels+1, v)}
		t.Fail()
	}else {
		fmt.Println("Success - RoutingTable")
	}
}

func TestEquals(t *testing.T){
	var address string = randomHex(40)
	var ID1, ID2 *KademliaID
	ID1 = NewKademliaID(address)
	ID2 = NewKademliaID(address)
	if(!ID1.Equals(ID2)){
		t.Fail()
	}else {
		fmt.Println("Success - KademliaID Equals")
	}
}

// Test Function to check Distance between 2 Nodes.
func TestDistFunc(t *testing.T) {

	DistT1 := NewKademliaID(testList[0])
	DistT2 := NewKademliaID(testList[5])
	DistBetween := DistT2.CalcDistance(DistT1)
	fmt.Printf("%v\n",DistBetween)
}
