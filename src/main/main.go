package main

import (
	"fmt"
	kademlia "d7024e"
	"time"
)

var bla []*kademlia.Kademlia //= make([]*kademlia.Kademlia, 99) 

func main()  {
	var port, num int = 8000, 9
	var base *kademlia.Kademlia = kademlia.NewKademlia("localhost", fmt.Sprintf("%d", port), nil)
	port++
	for i := 0; i < num; i++ { fmt.Print("'") }
	fmt.Println()
	for i := 0; i < num; i++ {
		bla = append(bla, kademlia.NewKademlia("localhost", fmt.Sprintf("%d", port), base.Me()))
		port++
		fmt.Print("'")
	}
	fmt.Println("\nStarted 100 nodes, starting at port", port-(num+1))
	go func () {
		for {
			time.Sleep(5* time.Second)
			fmt.Println("bip")
		}
	}()
	for {
		for i := 0; i <num; i++ {
			if(bla[i] == nil) {
				fmt.Println("pfft")
			}
		}
	}
}