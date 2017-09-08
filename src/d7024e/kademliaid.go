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
	decoded, _ := hex.DecodeString(data)

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
	return hex.EncodeToString(kademliaID[0:IDLength])
}

func (ID KademliaID) bitAt(exponent int) byte {
	//fmt.Println(ID)
	exponent = exponent/8
	var IDIndex, expBy8, expMod8 int
	expBy8 = exponent/8
	IDIndex = (IDLength-1) - expBy8
	expMod8 = exponent%8
	defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered. exponent: %v expBy8: %v IDIndex: %v expMod8: %v\n ID: %v\n", exponent, expBy8, IDIndex, expMod8, ID)
        }
    }()

	
	return (ID[IDIndex] >> uint(expMod8)) & 1;

}
