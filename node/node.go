package node

import (
	"consensus-kvstore/communication_iface"
	"consensus-kvstore/log"
	"consensus-kvstore/message"
	"fmt"
	"log/slog"
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
	CurrentTerm      int
	State            NodeState
	electionTimeout  time.Duration
	heartbeatTimeout time.Duration
	electionTimer    *time.Timer
	heartbeatTimer   *time.Timer
	Log_Node         *log.Log
	Inbox            chan *message.Envelope
	Voted            bool
	SystemLog        *slog.Logger
}

func NewNode(id int, CurrentTerm int, state NodeState, electionTimeout time.Duration, electionTimer *time.Timer,
	heartbeatTimeout time.Duration, heartbeatTimer *time.Timer, syslog *slog.Logger) *Node {
	Log_Object := &log.Log{
		LastLogTerm: 0,
		CommitIndex: 0,
		LastApplied: 0,
		List:        []log.LogEntry{},
	}
	syslog.Info(fmt.Sprintf("Created node %d", id))
	return &Node{
		ID:               id,
		CurrentTerm:      CurrentTerm,
		State:            state,
		electionTimeout:  electionTimeout,
		electionTimer:    electionTimer,
		heartbeatTimeout: heartbeatTimeout,
		heartbeatTimer:   heartbeatTimer,
		Log_Node:         Log_Object,
		Inbox:            make(chan *message.Envelope, 100),
		Voted:            false,
		SystemLog:        syslog,
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
		n.SystemLog.Info(fmt.Sprintf("election timeout for node %d, starting leader election", n.ID))
		// fmt.Printf("\nelection timeout for node %d, starting leader election", n.ID)
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
	n.CurrentTerm++
	n.SystemLog.Info(fmt.Sprintf("Node %d became candidate for term %d", n.ID, n.CurrentTerm))
	// fmt.Printf("\nNode %d became candidate for term %d", n.ID, n.CurrentTerm)
	net.RequestVotes(n.ID)
	// net.RequestVotes(candiadte, net)
	// majority:=TallyVotes()
	// if majority, MakeLeader()

}

/*
func RequestVotes() {
	for nodes in nodemaP:
		SendVoteRequest()
}
*/

// func RequestVotes(candidate *Node, net communication_iface.Network) {
// 	VoteRequest := &message.VoteRequest{
// 		CandidateTerm: candidate.currentTerm,
// 		CandidateID:   candidate.ID,
// 		LastLogIndex:  len(candidate.Log_Node.List),
// 	}
// 	// send vote requests to ALL nodes in netmap
// 	net.SendVoteRequest(candidate.ID, 1, *VoteRequest)
// }

func (n *Node) Listener(net communication_iface.Network) {
	go func() {
		for {
			msg, ok := <-n.Inbox
			if ok {
				n.HandleMessage(msg, net)
				n.SystemLog.Info(fmt.Sprintf("node %d received a message!\n", n.ID))
			}
		}

	}()

}

func (n *Node) HandleMessage(msg *message.Envelope, net communication_iface.Network) {
	if msg.Type == "VoteRequest" {
		VoteResponse := message.VoteResponse{
			Term:        n.CurrentTerm,
			VoteGranted: false,
		}

		req := msg.Payload.(message.VoteRequest)
		if req.CandidateTerm < n.Log_Node.LastLogTerm {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: candidate id:%d term:%d < my term %d", n.ID, req.CandidateID, req.CandidateTerm, n.CurrentTerm))
			net.SendVoteResponse(n.ID, msg.FromID, VoteResponse)
		}

		if n.Voted {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: already voted in term %d", n.ID, n.CurrentTerm))
			net.SendVoteResponse(n.ID, msg.FromID, VoteResponse)
			return
		}

		if req.LastLogTerm < n.Log_Node.LastLogTerm ||
			(req.LastLogTerm == n.Log_Node.LastLogTerm && req.LastLogIndex < len(n.Log_Node.List)) {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: candidate pair log term, log index:%d, %d < my pair %d,%d", n.ID, req.LastLogTerm, req.LastLogIndex, n.Log_Node.LastLogTerm, len(n.Log_Node.List)))
			net.SendVoteResponse(n.ID, msg.FromID, VoteResponse)
			return
		}

		n.Voted = true
		VoteResponse.VoteGranted = true
		n.SystemLog.Info(fmt.Sprintf("[Node %d] Accept vote", n.ID))
		net.SendVoteResponse(n.ID, msg.FromID, VoteResponse)
		return
	}
}

func Run() {
	fmt.Println("node")
}
