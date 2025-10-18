package main

import (
	"consensus-kvstore/node"
	"fmt"
	"time"
)

func main() {
	fmt.Println("main")
	electionTimeout := 3 * time.Second
	heartbeatTimeout := 5 * time.Second
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	n := node.NewNode(1, 0, node.Leader, time.Second*2, electionTimer, time.Second*4, heartbeatTimer)
	fmt.Printf("Created Node: %+v\n", n)
}
