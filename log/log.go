package log

import "fmt"

/*
	[x] Append a new command to the log of a node
	[x] Fetch the most recent log command entry - term, index, ID
	[x] Remove conflicting entries ()
	[x] Commit the commands, mark as leader approved
	[ ] Apply the commands (don't implement this yet)
*/

type LogEntry struct {
	Index   int
	Term    int
	Command string
}

type Log struct {
	LastLogTerm int
	CommitIndex int
	LastApplied int
	List        []LogEntry
}

func AppendEntry(Log_Node *Log, Index int, Term int, LogCommand string) {
	NewLogEntry := &LogEntry{
		Index:   Index,
		Term:    Term,
		Command: LogCommand,
	}

	Log_Node.List = append(Log_Node.List, *NewLogEntry)
}

func GetRecent(Log_Node *Log) (int, int) {
	/*
		Given a nodes log list, return the most recent (committed/uncommitted) log entry term,index
	*/
	var n int = len(Log_Node.List) - 1
	return Log_Node.List[n].Term, Log_Node.List[n].Index
}

func TruncateConflicts(Follower_Log *Log, LeaderTerm int, LeaderIndex int) bool {
	/*
		Once a leader has been chosen, all commands that are
			uncommitted and
			of older term
		need to be rolled back.

		Given a nodes log list and the latest term, deleted all commands that are uncommitted and having
		older term
	*/

	var n int = len(Follower_Log.List) - 1

	if LeaderIndex > n {
		return false
	}

	if Follower_Log.List[LeaderIndex].Term != LeaderTerm {
		Follower_Log.List = Follower_Log.List[:LeaderIndex]
	}

	return true
}

func CommitCommands(Log_Node *Log, UpToIndex int) {
	/*
		Once a leader has commmitted a command followers should also commit it.
		Commit is based on ? what is the trigger ? doesn't matter, relevant to calling function
	*/

	Log_Node.CommitIndex = UpToIndex

}

func ApplyCommands() {
	fmt.Println(("Executed command"))
}
