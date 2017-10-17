package main

import (
	"fmt"
	kademlia "d7024e"
	"time"
)

var bla []*kademlia.Kademlia //= make([]*kademlia.Kademlia, 99) 

func main()  {
	var port int = 8000
	var base *kademlia.Kademlia = kademlia.NewKademlia("localhost", fmt.Sprintf("%d", port), nil)
	port++
	for i := 1; i < 100; i++ { fmt.Print("'") }
	fmt.Println()
	for i := 1; i < 100; i++ {
		bla = append(bla, kademlia.NewKademlia("localhost", fmt.Sprintf("%d", port), base.Me()))
		port++
		fmt.Print("'")
	}
	fmt.Println("\nStarted 100 nodes, starting at port", port-100)
	go func () {
		for {
			time.Sleep(5* time.Second)
			fmt.Println("bip")
		}
	}()
	for {
		for i := 0; i <99; i++ {
			if(bla[i] == nil) {
				fmt.Println("pfft")
			}
		}
	}
}