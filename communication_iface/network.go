package communication_iface

import (
	"consensus-kvstore/message"
)

type Network interface {
	RequestVotes(candidateID int)
	SendVoteRequest(fromID, toID int, req message.VoteRequest)
	SendAppendEntries(LeaderID int, toID int, Command string)
	SendVoteResponse(fromID, toID int, tipe int, response message.VoteResponse)
	SendAppendEntriesResponse(fromID int, toID int, response message.AppendEntriesResponse)
	Heartbeat(LeaderID int)
	GetLastLogTerm(id int) int
	GetTerm(id int) int
	GetNumberOfNodes() int
}
