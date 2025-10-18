package message

import (
	"consensus-kvstore/node"
)

type RequestMessage struct {
	FromID       int
	Term         int
	ToID         int
	LastLogIndex int
	LastLogTerm  int
}

type ResponseMessage struct {
	FromID       int
	Term         int
	VotedGranted bool
}

func CreateVoteRequest(from_node *node.Node, to_node *node.Node) *RequestMessage {
	requestMessage := &RequestMessage{
		FromID:       from_node.ID,
		Term:         from_node.currentTerm,
		ToID:         to_node.ID,
		LastLogIndex: from_node.LastLogIndex,
		LastLogTerm:  from_node.LastLogIndex,
	}
	return requestMessage
}
