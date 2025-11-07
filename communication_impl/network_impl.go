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

func (r *InMemoryNetwork) SendVoteRequest(fromID int, toID int, request message.VoteRequest) {
	VoteRequest := message.Envelope{
		FromID:  fromID,
		ToID:    toID,
		Type:    "VoteRequest",
		Payload: request,
	}
	r.MessageMap[toID] <- VoteRequest
	r.History = append(r.History, VoteRequest)
	fmt.Printf("\nSendVoteRequest - FromID: %d, toID: %d", request.CandidateID, toID)
}

func (r *InMemoryNetwork) SendVoteResponse(fromID, toID int, response message.VoteResponse) {
	fmt.Printf("\nSendVoteResponse - FromID: %d, toID: %d, voteGranted: %t", fromID, toID, response.VoteGranted)
}

func (r *InMemoryNetwork) SendAppendEntries(fromID int, toID int) {
	fmt.Printf("\nSendVoteResponse - FromID: %d, toID: %d,", fromID, toID)
}

func (r *InMemoryNetwork) Heartbeat(LeaderID int) {
	for nodeID := range r.Nodes {
		r.SendAppendEntries(LeaderID, nodeID)
	}
}

var _ communication_iface.Network = (*InMemoryNetwork)(nil)
