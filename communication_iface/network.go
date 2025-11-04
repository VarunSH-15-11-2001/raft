package communication_iface

import (
	"consensus-kvstore/message"
)

type Network interface {
	SendVoteRequest(fromID, toID int, req message.VoteRequest)
	SendAppendEntries(fromID, toID int)
	Heartbeat(LeaderID int)
}
