package communication_impl

import (
	"consensus-kvstore/communication_iface"
	"consensus-kvstore/message"
	"consensus-kvstore/node"
	"fmt"
)

type RealNetwork struct {
	nodes map[int]*node.Node
}

func (r *RealNetwork) Init(nodeMap map[int]*node.Node) {
	r.nodes = nodeMap
}

func (r *RealNetwork) SendVoteRequest(fromID int, toID int, request message.VoteRequest) {
	fmt.Printf("\nSendVoteRequest - FromID: %d, toID: %d", request.CandidateID, toID)
}

func (r *RealNetwork) SendVoteResponse(fromID, toID int, response message.VoteResponse) {
	fmt.Printf("\nSendVoteResponse - FromID: %d, toID: %d, voteGranted: %t", fromID, toID, response.VoteGranted)
}

func (r *RealNetwork) SendAppendEntries(fromID int, toID int) {
	fmt.Printf("\nSendVoteResponse - FromID: %d, toID: %d,", fromID, toID)
}

func (r *RealNetwork) Heartbeat(LeaderID int) {
	for nodeID := range r.nodes {
		r.SendAppendEntries(LeaderID, nodeID)
	}
}

var _ communication_iface.Network = (*RealNetwork)(nil)
