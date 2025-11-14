package communication_iface

import (
	"consensus-kvstore/message"
)

type Network interface {
	RequestVotes(candidateID int)
	SendVoteRequest(fromID, toID int, req message.VoteRequest)
	SendAppendEntries(fromID, toID int)
	SendVoteResponse(fromID, toID int, response message.VoteResponse)
	Heartbeat(LeaderID int)
	GetLastLogTerm(id int) int
	GetTerm(id int) int
}
