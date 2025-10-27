package log

import "testing"

func TestAppendEntry_Positive(t *testing.T) {
	LogList := []LogEntry{
		{
			Index:   1,
			Term:    1,
			Command: "Set X=5",
		},
		{
			Index:   2,
			Term:    1,
			Command: "Set X=8",
		},
	}

	Node_log := &Log{
		LastLogTerm: 1,
		CommitIndex: 2,
		LastApplied: 1,
		list:        LogList,
	}

	Log_Entry := &LogEntry{
		Index:   3,
		Term:    1,
		Command: "Set X=10",
	}

	AppendEntry(Node_log, 3, 1, "Set X=10")
	var n int = len(Node_log.list)

	var found bool = false
	for i := 0; i < n; i += 1 {
		if Node_log.list[i] == *Log_Entry {
			found = true
			break
		}
	}

	if !found {
		t.Error("Failed to append log list")
	}
}

func TestAppendEntry_Negative(t *testing.T) {
	LogList := []LogEntry{
		{
			Index:   1,
			Term:    1,
			Command: "Set X=5",
		},
		{
			Index:   2,
			Term:    1,
			Command: "Set X=8",
		},
	}

	Node_Log_Positive := &Log{
		LastLogTerm: 1,
		CommitIndex: 2,
		LastApplied: 1,
		list:        LogList,
	}

	Node_Log_Negative := &Log{
		LastLogTerm: 1,
		CommitIndex: 2,
		LastApplied: 1,
		list:        LogList,
	}

	Log_Entry := &LogEntry{
		Index:   3,
		Term:    1,
		Command: "Set X=10",
	}

	AppendEntry(Node_Log_Positive, 3, 1, "Set X=10")
	var n int = len(Node_Log_Positive.list)

	var found bool = false
	for i := 0; i < n-1; i += 1 {
		if Node_Log_Negative.list[i] == *Log_Entry {
			found = true
			break
		}
	}

	if found {
		t.Error("Unexpected append entry call")
	}
}

func TestGetRecent(t *testing.T) {
	LogList := []LogEntry{
		{
			Index:   1,
			Term:    1,
			Command: "Set X=5",
		},
		{
			Index:   2,
			Term:    1,
			Command: "Set X=8",
		},
	}

	Node_Log := &Log{
		LastLogTerm: 1,
		CommitIndex: 2,
		LastApplied: 1,
		list:        LogList,
	}

	var n int = len(Node_Log.list) - 1
	if Node_Log.list[n].Term != 1 || Node_Log.list[n].Index != 2 {
		t.Error("Wrong term or index")
	}
}

func TestTruncateConflicts(testingt *testing.T) {

}
