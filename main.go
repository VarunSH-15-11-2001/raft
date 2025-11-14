package main

import (
	"consensus-kvstore/communication_impl"
	"consensus-kvstore/node"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

func main() {

	file, err := os.OpenFile(
		"logs/main.log",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0666,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	logger_main := slog.New(slog.NewJSONHandler(file, nil))

	logger_main.Info("Main started")
	leader_logger, leader_file := GetLogger("logs/leader.log")
	electionTimeout := 20 * time.Second
	heartbeatTimeout := 5 * time.Second
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	leader := node.NewNode(0, 0, node.Leader, time.Second*20, electionTimer, time.Second*4, heartbeatTimer, leader_logger)
	leader.State = node.Leader

	logger_main.Info("Created Leader Node with ID:", "id", leader.ID)

	follower1_logger, follower1_file := GetLogger("logs/follower1.log")
	follower1 := GetNewNode(1, follower1_logger)
	// follower1_timer := time.NewTimer(1 * time.Second)
	// follower1.electionTimeout = 1 * time.Second

	follower2_logger, follower2_file := GetLogger("logs/follower2.log")
	follower2 := GetNewNode(2, follower2_logger)

	follower3_logger, follower3_file := GetLogger("logs/follower3.log")
	follower3 := GetNewNode(3, follower3_logger)

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
	leader.Listener(&network)
	follower1.Listener(&network)
	follower2.Listener(&network)
	follower3.Listener(&network)

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
		logger_main.Debug("Message added", "fromID", network.History[1].FromID, "toID", network.History[1].ToID)
	}

	time.Sleep(30 * time.Second)
	defer file.Close()
	defer leader_file.Close()
	defer follower1_file.Close()
	defer follower2_file.Close()
	defer follower3_file.Close()

	defer fmt.Print("Main completed")
}

func GetNewNode(i int, syslog *slog.Logger) *node.Node {
	IDs := []int{0, 1, 2, 3, 4}
	var Term int = 0

	ElectionTimeouts := []time.Duration{20 * time.Second, 2 * time.Second, 6 * time.Second, 10 * time.Second, 15 * time.Second}
	electionTimer := time.NewTimer(ElectionTimeouts[i])

	// logger_main.Info("\nCreated Node: %+v with electionTimeout: %+v\n", IDs[i], ElectionTimeouts[i])

	return node.NewNode(IDs[i], Term, node.Follower, ElectionTimeouts[i], electionTimer, 1*time.Second, time.NewTimer(1*time.Second), syslog)

}

func GetLogger(nodeName string) (*slog.Logger, *os.File) {
	file, err := os.OpenFile(
		nodeName,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0666,
	)

	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}

	return slog.New(slog.NewTextHandler(file, opts)), file
}
