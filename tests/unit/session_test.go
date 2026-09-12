package unit

import (
	"os"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/session"
)

func TestSessionStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agypp-session-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	store, err := session.NewStore(tempDir)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}

	// 1. Create Session
	sess, err := store.CreateSession("/test/workspace", "Test Session")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	if sess.ID == "" {
		t.Fatal("expected non-empty session ID")
	}

	// 2. Add Checkpoint
	cp, err := store.AddCheckpoint(sess, "Initial commit snapshot", "abc1234", "+ added file")
	if err != nil {
		t.Fatalf("failed to add checkpoint: %v", err)
	}
	if cp.ID == "" {
		t.Fatal("expected non-empty checkpoint ID")
	}

	// 3. Retrieve Session
	retrieved, err := store.Get(sess.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}
	if len(retrieved.Checkpoints) != 1 {
		t.Fatalf("expected 1 checkpoint, got %d", len(retrieved.Checkpoints))
	}
	if retrieved.Checkpoints[0].Note != "Initial commit snapshot" {
		t.Errorf("unexpected checkpoint note: %s", retrieved.Checkpoints[0].Note)
	}

	// 4. List Sessions
	list, err := store.List()
	if err != nil {
		t.Fatalf("failed to list sessions: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 session in list, got %d", len(list))
	}
}

func TestGitHelper(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	git := session.NewGitHelper(currentDir)
	if !git.IsGitRepo() {
		// Parent or current should be git repo
		t.Log("Note: not inside git repo in test")
		return
	}

	head, err := git.GetHeadCommit()
	if err != nil {
		t.Fatalf("failed to get git head: %v", err)
	}
	if head == "" {
		t.Fatal("expected non-empty head commit")
	}
}
