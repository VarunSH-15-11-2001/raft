package node

import (
	"testing"
	"time"
)

func TestNodeConstructor(t *testing.T) {
	electionTimeout := 3 * time.Second
	heartbeatTimeout := 5 * time.Second
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	n := NewNode(1, 0, 0, electionTimeout, electionTimer, heartbeatTimeout, heartbeatTimer)

	if n.ID != 1 {
		t.Errorf("Expected ID=1 but got %d", n.ID)
	}
	if n.State != Follower {
		t.Errorf("Expected State=follower but got %s", n.State)
	}
	if n.electionTimeout != electionTimeout {
		t.Errorf("expected electionTimeout=%v, got %v", electionTimeout, n.electionTimeout)
	}
}

func TestNodeTimer(t *testing.T) {
	electionTimeout := 100 * time.Millisecond
	heartbeatTimeout := 200 * time.Millisecond
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	n := NewNode(1, 0, Follower, electionTimeout, electionTimer, heartbeatTimeout, heartbeatTimer)

	done := make(chan bool, 1)
	StartElectionTimeout(n, done)

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected election timeout did not fire")
	}
}

func TestStartElection(t *testing.T) {
	electionTimeout := 100 * time.Millisecond
	heartbeatTimeout := 200 * time.Millisecond
	electionTimer := time.NewTimer(electionTimeout)
	heartbeatTimer := time.NewTimer(heartbeatTimeout)
	n := NewNode(1, 0, Follower, electionTimeout, electionTimer, heartbeatTimeout, heartbeatTimer)
	prevTerm := n.currentTerm
	StartElection(n)

	if n.State != Candidate {
		t.Errorf("Node=%d failed to enter candidate state", n.ID)
	}
	if n.currentTerm != prevTerm+1 {
		t.Errorf("Node=%d term did not update. Previous=%d, after=%d", n.ID, prevTerm, n.currentTerm)
	}
}
