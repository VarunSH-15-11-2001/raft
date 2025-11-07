package main

import (
	"consensus-kvstore/communication_impl"
	"consensus-kvstore/node"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("main")
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

	network := communication_impl.InMemoryNetwork{}
	network.Init(nodeMap)
	leader_done := make(chan bool)
	follower1_done := make(chan bool)
	follower2_done := make(chan bool)
	follower3_done := make(chan bool)
	follower1.Listener(&network)

	var wg sync.WaitGroup

	start := make(chan struct{})

	wg.Add(4)

	go func() {
		<-start
		node.StartElectionTimeout(leader, leader_done, &network)
		<-leader_done
		wg.Done()
	}()

	go func() {
		<-start
		node.StartElectionTimeout(follower1, follower1_done, &network)
		<-follower1_done
		wg.Done()
	}()

	go func() {
		<-start
		node.StartElectionTimeout(follower2, follower2_done, &network)
		<-follower2_done
		wg.Done()
	}()

	go func() {
		<-start
		node.StartElectionTimeout(follower3, follower3_done, &network)
		<-follower3_done
		wg.Done()
	}()

	close(start)
	wg.Wait()

	if len(network.History) >= 2 {
		fmt.Printf("\nMessage added from nodeID:%d to nodeID:%d\n", network.History[1].FromID, network.History[1].ToID)
	}

	time.Sleep(8 * time.Second)
	fmt.Print("\nMain completed\n")
}

func GetNewNode(i int) *node.Node {
	IDs := []int{0, 1, 2, 3, 4}
	var Term int = 0

	ElectionTimeouts := []time.Duration{2 * time.Second, 3 * time.Second, 4 * time.Second, 5 * time.Second}
	electionTimer := time.NewTimer(ElectionTimeouts[i])

	fmt.Printf("\nCreated Node: %+v with electionTimeout: %+v\n", IDs[i], ElectionTimeouts[i])

	return node.NewNode(IDs[i], Term, node.Follower, ElectionTimeouts[i], electionTimer, 1*time.Second, time.NewTimer(1*time.Second))

}
