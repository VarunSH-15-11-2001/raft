package communication_impl

import (
	"consensus-kvstore/communication_iface"
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

func (r *InMemoryNetwork) SendVoteResponse(fromID, toID int, response message.VoteResponse) {
	r.Nodes[fromID].SystemLog.Info(fmt.Sprintf("SendVoteResponse - FromID: %d, toID: %d, voteGranted: %t", fromID, toID, response.VoteGranted))
}

func (r *InMemoryNetwork) SendAppendEntries(fromID int, toID int) {
	r.Nodes[fromID].SystemLog.Info(fmt.Sprintf("SendVoteResponse - FromID: %d, toID: %d,", fromID, toID))
}

func (r *InMemoryNetwork) Heartbeat(LeaderID int) {
	for nodeID := range r.Nodes {
		r.SendAppendEntries(LeaderID, nodeID)
	}
}

func (r *InMemoryNetwork) GetLastLogTerm(id int) int {
	return r.Nodes[id].Log_Node.LastLogTerm
}

func (r *InMemoryNetwork) GetTerm(id int) int {
	return r.Nodes[id].CurrentTerm
}

var _ communication_iface.Network = (*InMemoryNetwork)(nil)
