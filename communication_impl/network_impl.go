package communication_impl

import (
	"consensus-kvstore/communication_iface"
	"consensus-kvstore/log"
	"consensus-kvstore/message"
	"consensus-kvstore/node"
	"fmt"
)

type InMemoryNetwork struct {
	MessageMap map[int]chan message.Envelope
	Nodes      map[int]*node.Node
	History    []message.Envelope
}

func (r *InMemoryNetwork) Init(nodeMap map[int]*node.Node) {
	r.Nodes = nodeMap
	r.MessageMap = make(map[int]chan message.Envelope)

	for id, n := range nodeMap {
		inbox := make(chan message.Envelope, 100)

		go func(n *node.Node, inbox chan message.Envelope) {
			for msg := range inbox {
				n.Inbox <- &msg
			}
		}(n, inbox)
		r.MessageMap[id] = inbox
	}
}

func (r *InMemoryNetwork) RequestVotes(candidateID int) {
	VoteRequest := message.VoteRequest{
		CandidateTerm: r.Nodes[candidateID].CurrentTerm,
		CandidateID:   candidateID,
		LastLogIndex:  len(r.Nodes[candidateID].Log_Node.List),
	}
	for to_id := range r.Nodes {
		if to_id != candidateID {
			r.SendVoteRequest(candidateID, to_id, VoteRequest)
		}
	}
}

func (r *InMemoryNetwork) SendVoteRequest(fromID int, toID int, request message.VoteRequest) {
	VoteRequest := message.Envelope{
		FromID:  fromID,
		ToID:    toID,
		Type:    "VoteRequest",
		Payload: request,
	}
	r.MessageMap[toID] <- VoteRequest
	r.History = append(r.History, VoteRequest)
	r.Nodes[fromID].SystemLog.Info(fmt.Sprintf("SendVoteRequest - FromID: %d, toID: %d", request.CandidateID, toID))
}

func (r *InMemoryNetwork) SendVoteResponse(fromID, toID int, tipe int, response message.VoteResponse) {
	var tipeString string
	switch tipe {
	case 1:
		tipeString = "VoteResponse"
	case 2:
		tipeString = "AppendResponse"
	}
	Response := message.Envelope{
		FromID:  fromID,
		ToID:    toID,
		Type:    tipeString,
		Payload: response,
	}
	select {
	case r.MessageMap[toID] <- Response:
		r.History = append(r.History, Response)
	default:
		r.Nodes[fromID].SystemLog.Warn(fmt.Sprintf("Dropped %s to %d: inbox full", tipeString, toID))
	}
}

func (r *InMemoryNetwork) SendAppendEntries(LeaderID int, nodeID int, Command string) {
	leaderNode := r.Nodes[LeaderID]

	Entries := []log.LogEntry{{
		Index:   len(leaderNode.Log_Node.List) + 1,
		Term:    leaderNode.CurrentTerm,
		Command: Command,
	},
	}
	n := len(leaderNode.Log_Node.List)
	PrevLogEntry := leaderNode.Log_Node.List[n-1]

	AppendEntryCommand := message.AppendEntriesRequest{
		LeaderID:     LeaderID,
		LeaderTerm:   leaderNode.CurrentTerm,
		Entries:      Entries,
		PrevLogIndex: PrevLogEntry.Index,
		PrevLogTerm:  PrevLogEntry.Term,
		LeaderCommit: leaderNode.Log_Node.CommitIndex,
	}

	AppendEntryMessage := message.Envelope{
		FromID:  LeaderID,
		ToID:    nodeID,
		Type:    "AppendEntry",
		Payload: AppendEntryCommand,
	}

	r.MessageMap[nodeID] <- AppendEntryMessage
	select {
	case r.MessageMap[nodeID] <- AppendEntryMessage:
		r.History = append(r.History, AppendEntryMessage)
	default:
		r.Nodes[nodeID].SystemLog.Warn(fmt.Sprintf("Dropped AppendEntry to %d: inbox full", nodeID))
	}
}

func (r *InMemoryNetwork) SendAppendEntriesResponse(fromID int, toID int, response message.AppendEntriesResponse) {
	Response := message.Envelope{
		FromID:  fromID,
		ToID:    toID,
		Type:    "AppendEntryResponse",
		Payload: response,
	}
	select {
	case r.MessageMap[toID] <- Response:
		r.History = append(r.History, Response)
	default:
		r.Nodes[fromID].SystemLog.Warn(fmt.Sprintf("Dropped AppendEntryResponse to %d: inbox full", toID))
	}
}

func (r *InMemoryNetwork) Heartbeat(LeaderID int) {
	for nodeID := range r.Nodes {
		r.SendAppendEntries(LeaderID, nodeID, "")
	}
}

func (r *InMemoryNetwork) GetLastLogTerm(id int) int {
	return r.Nodes[id].Log_Node.LastLogTerm
}

func (r *InMemoryNetwork) GetTerm(id int) int {
	return r.Nodes[id].CurrentTerm
}

func (r *InMemoryNetwork) GetNumberOfNodes() int {
	return len(r.Nodes)
}

var _ communication_iface.Network = (*InMemoryNetwork)(nil)
