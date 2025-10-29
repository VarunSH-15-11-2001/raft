package main

import (
	"consensus-kvstore/communication"
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
	leader := node.NewNode(0, 0, node.Leader, time.Second*2, electionTimer, time.Second*4, heartbeatTimer)
	leader.State = node.Leader
	fmt.Printf("Created Node: %+v\n", leader)

	follower1 := GetNewNode(1)
	follower2 := GetNewNode(2)
	follower3 := GetNewNode(3)

	var nodeMap map[int]*node.Node = make(map[int]*node.Node)
	nodeMap[leader.ID] = leader
	nodeMap[follower1.ID] = follower1
	nodeMap[follower2.ID] = follower2
	nodeMap[follower3.ID] = follower3

	communication.Init(nodeMap)
}

func GetNewNode(i int) *node.Node {
	IDs := []int{0, 1, 2, 3, 4}
	var Term int = 0

	ElectionTimeouts := []time.Duration{2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	electionTimer := time.NewTimer(ElectionTimeouts[i])

	return node.NewNode(IDs[i], Term, node.Follower, ElectionTimeouts[i], electionTimer, 1*time.Second, time.NewTimer(1*time.Second))

}
