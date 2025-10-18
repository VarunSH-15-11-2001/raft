package log

type LogEntry struct {
	Index   int
	Term    int
	Command string
}

type Log struct {
	LastLogIndex int
	LastLogTerm  int
	CommitIndex  int
	LastApplied  int
	list         []LogEntry
}

func AppendEntry(Log_Node *Log, Index int, Term int, LogCommand string) {
	NewLogEntry := &LogEntry{
		Index:   Index,
		Term:    Term,
		Command: LogCommand,
	}

	Log_Node.list = append(Log_Node.list, *NewLogEntry)
}
