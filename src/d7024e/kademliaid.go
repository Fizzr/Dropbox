package d7024e

import (
	"encoding/hex"
	"math/rand"
	"fmt"
)

const IDLength = 20
const IDBits = IDLength * 8

type KademliaID [IDLength]byte

func NewKademliaID(data string) *KademliaID {
	var decoded []byte
	decoded, _ = hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

func (kademliaID *KademliaID) String() string {
	var ret string
	for i:=0; i < IDLength; i++{
		ret += fmt.Sprintf("%02X",kademliaID[i])
	}
	return ret
	//return hex.EncodeToString(kademliaID[0:IDLength])
}

func (kademliaID *KademliaID) toBinary() string {
	var result string
	for i := 0; i < IDLength; i++ {
		result += fmt.Sprintf("%08b", kademliaID[i])
	}
	return result
}

func (ID KademliaID) bitAt(exponent int) byte {
	//fmt.Println(ID)
	//exponent = exponent/8 //from bit to byte
	var IDIndex, expBy8, expMod8 int
	expBy8 = exponent/8
	expMod8 = exponent%8
	if(expMod8 == 0 && exponent != 0){
		expBy8 -= 1
	}
	IDIndex = (IDLength-1) - expBy8
	defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered. exponent: %v expBy8: %v IDIndex: %v expMod8: %v\n ID: %v\n", exponent, expBy8, IDIndex, expMod8, ID)
        }
    }()

	
	return (ID[IDIndex] >> uint(expMod8)) & 1;

}
