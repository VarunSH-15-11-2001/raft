package node

import (
	"consensus-kvstore/communication_iface"
	"consensus-kvstore/log"
	"consensus-kvstore/message"
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
	Log_Object := &log.Log{
		LastLogTerm: 0,
		CommitIndex: 0,
		LastApplied: 0,
		List:        []log.LogEntry{},
	}
	return &Node{
		ID:               id,
		currentTerm:      currentTerm,
		State:            state,
		electionTimeout:  electionTimeout,
		electionTimer:    electionTimer,
		heartbeatTimeout: heartbeatTimeout,
		heartbeatTimer:   heartbeatTimer,
		Log_Node:         Log_Object,
	}
}

func StartElectionTimeout(n *Node, done chan<- bool, net communication_iface.Network) {
	if n.electionTimer != nil {
		n.electionTimer.Stop()
	}
	n.electionTimer = time.NewTimer(n.electionTimeout)
	go func() {
		<-n.electionTimer.C
		done <- true
		StartElection(n, net)
		fmt.Printf("\nelection timeout for node %d, starting leader election", n.ID)
	}()
}

func StartElection(n *Node, net communication_iface.Network) {
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
	fmt.Printf("\nNode %d became candidate for term %d", n.ID, n.currentTerm)
	RequestVotes(n, net)

	// majority:=TallyVotes()
	// if majority, MakeLeader()

}

func RequestVotes(candidate *Node, net communication_iface.Network) {
	VoteRequest := &message.VoteRequest{
		CandidateTerm: candidate.currentTerm,
		CandidateID:   candidate.ID,
		LastLogIndex:  len(candidate.Log_Node.List),
	}
	net.SendVoteRequest(candidate.ID, 1, *VoteRequest)
}

func Run() {
	fmt.Println("node")
}
