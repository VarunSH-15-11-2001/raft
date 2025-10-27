package log

import "fmt"

/*
	[x] Append a new command to the log of a node
	[ ] Fetch the most recent log command entry - term, index, ID
	[ ] Remove conflicting entries ()
	[ ] Commit the commands, mark as leader approved
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
	list        []LogEntry
}

func AppendEntry(Log_Node *Log, Index int, Term int, LogCommand string) {
	NewLogEntry := &LogEntry{
		Index:   Index,
		Term:    Term,
		Command: LogCommand,
	}

	Log_Node.list = append(Log_Node.list, *NewLogEntry)
}

func GetRecent(Log_Node *Log) (int, int) {
	/*
		Given a nodes log list, return the most recent (committed/uncommitted) log entry term,index
	*/
	var n int = len(Log_Node.list) - 1
	return Log_Node.list[n].Term, Log_Node.list[n].Index
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

	var n int = len(Follower_Log.list) - 1

	if LeaderIndex > n {
		return false
	}

	if Follower_Log.list[LeaderIndex].Term != LeaderTerm {
		Follower_Log.list = Follower_Log.list[:LeaderIndex]
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
