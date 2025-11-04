package main

import (
	"consensus-kvstore/communication_impl"
	"consensus-kvstore/node"
	"fmt"
	"time"
)

func main() {
	fmt.Println("main")
	// main_done := make(chan bool)
	electionTimeout := 6 * time.Second
	heartbeatTimeout := 5 * time.Second
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	leader := node.NewNode(0, 0, node.Leader, time.Second*2, electionTimer, time.Second*4, heartbeatTimer)
	leader.State = node.Leader
	fmt.Printf("\nCreated Node: %+v\n", leader)

	follower1 := GetNewNode(1)
	// follower1_timer := time.NewTimer(1 * time.Second)
	// follower1.electionTimeout = 1 * time.Second

	follower2 := GetNewNode(2)
	follower3 := GetNewNode(3)

	var nodeMap map[int]*node.Node = make(map[int]*node.Node)
	nodeMap[leader.ID] = leader
	nodeMap[follower1.ID] = follower1
	nodeMap[follower2.ID] = follower2
	nodeMap[follower3.ID] = follower3

	network := communication_impl.RealNetwork{}
	network.Init(nodeMap)
	leader_done := make(chan bool)
	node.StartElectionTimeout(leader, leader_done, &network)
	follower1_done := make(chan bool)
	follower2_done := make(chan bool)
	follower3_done := make(chan bool)

	node.StartElectionTimeout(follower1, follower1_done, &network)
	node.StartElectionTimeout(follower2, follower2_done, &network)
	node.StartElectionTimeout(follower3, follower3_done, &network)
	<-follower1_done
	<-follower2_done
	<-follower3_done
	// <-main_done
	fmt.Print("Main completed")
}

func GetNewNode(i int) *node.Node {
	IDs := []int{0, 1, 2, 3, 4}
	var Term int = 0

	ElectionTimeouts := []time.Duration{2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	electionTimer := time.NewTimer(ElectionTimeouts[i])

	fmt.Printf("\nCreated Node: %+v with electionTimeout: %+v\n", IDs[i], ElectionTimeouts[i])

	return node.NewNode(IDs[i], Term, node.Follower, ElectionTimeouts[i], electionTimer, 1*time.Second, time.NewTimer(1*time.Second))

}
