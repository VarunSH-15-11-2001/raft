package node

import (
	"consensus-kvstore/log"
	"fmt"
	"time"
)

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

func (s NodeState) String() string {
	return [...]string{"Follower", "Candidate", "Leader"}[s]
}

func Log(message string) {
	fmt.Println(message)
}

type Node struct {
	ID               int
	currentTerm      int
	State            NodeState
	electionTimeout  time.Duration
	heartbeatTimeout time.Duration
	electionTimer    *time.Timer
	heartbeatTimer   *time.Timer
	Log_Node         *log.Log
}

func NewNode(id int, currentTerm int, state NodeState, electionTimeout time.Duration, electionTimer *time.Timer,
	heartbeatTimeout time.Duration, heartbeatTimer *time.Timer) *Node {
	return &Node{
		ID:               id,
		currentTerm:      currentTerm,
		State:            state,
		electionTimeout:  electionTimeout,
		electionTimer:    electionTimer,
		heartbeatTimeout: heartbeatTimeout,
		heartbeatTimer:   heartbeatTimer,
	}
}

func StartElectionTimeout(n *Node, done chan<- bool) {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	n.electionTimer = time.NewTimer(n.electionTimeout)
	go func() {
		<-n.electionTimer.C
		done <- true
		StartElection(n)
		fmt.Printf("election timeout for node %d, starting leader election", n.ID)
	}()
}

func StartElection(n *Node) {
	/*
		Become candidate
		request votes, nodes will reply
			Requesting votes:
				send message to others
					message:

		if it gets majority, become leader


	*/
	n.State = Candidate
	n.currentTerm++
	fmt.Printf("Node %d became candidate for term %d", n.ID, n.currentTerm)
	// RequestVotes()
	// majority:=TallyVotes()
	// if majority, MakeLeader()

}

func RequestVotes() {

}

func Run() {
	fmt.Println("node")
}
