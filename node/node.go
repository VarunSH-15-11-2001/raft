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
	Votes            int
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
		Votes:            0,
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
		/*
			Receiver implementation:
			1. Reply false if term < currentTerm (§5.1)
			2. If votedFor is null or candidateId, and candidate’s log is at
			least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
		*/

		VoteResponse := message.VoteResponse{
			Term:        n.CurrentTerm,
			VoteGranted: false,
		}

		req := msg.Payload.(message.VoteRequest)
		if req.CandidateTerm < n.Log_Node.LastLogTerm {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: candidate id:%d term:%d < my term %d", n.ID, req.CandidateID, req.CandidateTerm, n.CurrentTerm))
			net.SendVoteResponse(n.ID, msg.FromID, 1, VoteResponse)
		}

		if n.Voted {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: already voted in term %d", n.ID, n.CurrentTerm))
			net.SendVoteResponse(n.ID, msg.FromID, 1, VoteResponse)
			return
		}

		if req.LastLogTerm < n.Log_Node.LastLogTerm ||
			(req.LastLogTerm == n.Log_Node.LastLogTerm && req.LastLogIndex < len(n.Log_Node.List)) {
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Reject vote: candidate pair log term, log index:%d, %d < my pair %d,%d", n.ID, req.LastLogTerm, req.LastLogIndex, n.Log_Node.LastLogTerm, len(n.Log_Node.List)))
			net.SendVoteResponse(n.ID, msg.FromID, 1, VoteResponse)
			return
		}

		n.Voted = true
		VoteResponse.VoteGranted = true
		n.SystemLog.Info(fmt.Sprintf("[Node %d] Accept vote", n.ID))
		net.SendVoteResponse(n.ID, msg.FromID, 1, VoteResponse)
		return
	} else if msg.Type == "VoteResponse" {
		resp := msg.Payload.(message.VoteResponse)

		if resp.Term > n.CurrentTerm {
			n.CurrentTerm = resp.Term
			n.State = Follower
			n.Voted = false
			return
		}

		if resp.VoteGranted && n.State == Candidate {
			n.Votes++

			if n.Votes > net.GetNumberOfNodes()/2 {
				n.State = Leader
				n.SystemLog.Info(fmt.Sprintf("[Node %d] Became leader for term %d", n.ID, n.CurrentTerm))
				// net.HeartBeat(n.ID)
			}
		}

	} else if msg.Type == "AppendEntry" {
		/*
			Receiver implementation:
				[x] 1. Reply false if term < currentTerm (§5.1)
				[x] 2. Reply false if log doesn’t contain an entry at prevLogIndex whose term matches prevLogTerm (§5.3)
					(consistency check to confirm that the follower is atleast upto date with leader)
				[x] 3. If an existing entry conflicts with a new one (same index but different terms),
					delete the existing entry and all that follow it (§5.3)
				[x] 4. Append any new entries not already in the log
				[x] 5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index of last new entry)
		*/

		req := msg.Payload.(message.AppendEntriesRequest)
		AppendEntryPayload := message.AppendEntriesResponse{
			FollowerTerm: n.CurrentTerm,
			Success:      false,
			Index:        0,
		}

		if n.CurrentTerm > req.LeaderTerm {
			net.SendAppendEntriesResponse(n.ID, req.LeaderID, AppendEntryPayload)
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Didn't append because my term %d is newer than leader %d", n.ID, n.CurrentTerm, req.LeaderTerm))
			return
		}

		if req.PrevLogIndex > len(n.Log_Node.List)-1 || n.Log_Node.List[req.PrevLogIndex].Term != req.PrevLogTerm {
			net.SendAppendEntriesResponse(n.ID, req.LeaderID, AppendEntryPayload)
			n.SystemLog.Info(fmt.Sprintf("[Node %d] Didn't append because I don't have an entry at prevLogIndex %d with term %d", n.ID, req.PrevLogIndex, req.PrevLogTerm))
			return
		}

		conflictIndex := -1
		for i, entry := range req.Entries {
			targetIndex := req.PrevLogIndex + i + 1

			if targetIndex > len(n.Log_Node.List) {
				conflictIndex = i
				break
			}

			if targetIndex < len(n.Log_Node.List) {
				if n.Log_Node.List[targetIndex].Term != entry.Term {
					n.Log_Node.List = n.Log_Node.List[:targetIndex]
					break
				}
			}
		}

		if conflictIndex != -1 {
			n.Log_Node.List = append(n.Log_Node.List, req.Entries[conflictIndex:]...)
		} else {
			startIndex := req.PrevLogIndex + 1
			if startIndex < len(req.Entries) {
				missing := req.Entries[startIndex-req.PrevLogIndex:]
				n.Log_Node.List = append(n.Log_Node.List, missing...)
			}
		}

		if req.LeaderCommit > n.Log_Node.CommitIndex {
			lastEntryCommit := n.Log_Node.List[len(n.Log_Node.List)-1].Index
			n.Log_Node.CommitIndex = min(req.LeaderCommit, lastEntryCommit)
			n.SystemLog.Info(
				fmt.Sprintf(
					"[Node %d] Updated commitIndex to %d (leaderCommit=%d, lastEntry=%d)",
					n.ID, n.Log_Node.CommitIndex, req.LeaderCommit, lastEntryCommit,
				),
			)
		}

		AppendEntryPayload.Success = true
		AppendEntryPayload.Index = len(n.Log_Node.List) - 1
		net.SendAppendEntriesResponse(n.ID, req.LeaderID, AppendEntryPayload)
		return
	}
}

func Run() {
	fmt.Println("node")
}
