package message

import (
	"consensus-kvstore/log"
)

type VoteRequest struct {
	CandidateTerm int
	CandidateID   int
	LastLogIndex  int
	LastLogTerm   int
}

type VoteResponse struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesRequest struct {
	LeaderTerm   int
	LeaderID     int
	Entries      []log.LogEntry
	PrevLogIndex int
	PrevLogTerm  int
	LeaderCommit int
}

type AppendEntriesResponse struct {
	FollowerTerm int
	Success      bool
	Index        int
}

type Envelope struct {
	FromID  int
	ToID    int
	Type    string
	Payload interface{}
}
