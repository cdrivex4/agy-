package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
)

// Checkpoint represents a saved state snapshot (git state + notes).
type Checkpoint struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Timestamp time.Time `json:"timestamp"`
	GitCommit string    `json:"git_commit,omitempty"`
	Diff      string    `json:"diff,omitempty"`
	Note      string    `json:"note"`
}

// Session represents an AGY++ companion session record.
type Session struct {
	ID           string        `json:"id"`
	AgySessionID string        `json:"agy_session_id,omitempty"`
	Workspace    string        `json:"workspace"`
	Title        string        `json:"title"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Model        string        `json:"model,omitempty"`
	YoloEnabled  bool          `json:"yolo_enabled"`
	PlanEnabled  bool          `json:"plan_enabled"`
	Checkpoints  []*Checkpoint `json:"checkpoints,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
}

// Store manages persistence of AGY++ session metadata.
type Store struct {
	baseDir string
}

// NewStore initializes a session store in the standard user application data directory.
func NewStore(customDir ...string) (*Store, error) {
	var dir string
	if len(customDir) > 0 && customDir[0] != "" {
		dir = customDir[0]
	} else {
		dir = defaultDataDir()
	}

	sessionDir := filepath.Join(dir, "sessions")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session store directory %s: %w", sessionDir, err)
	}

	return &Store{baseDir: sessionDir}, nil
}

// CreateSession initializes and persists a new session.
func (s *Store) CreateSession(workspace, title string) (*Session, error) {
	id := generateID()
	now := time.Now().UTC()

	if title == "" {
		title = fmt.Sprintf("Session %s", now.Format("2006-01-02 15:04"))
	}

	sess := &Session{
		ID:          id,
		Workspace:   workspace,
		Title:       title,
		CreatedAt:   now,
		UpdatedAt:   now,
		Checkpoints: []*Checkpoint{},
		Tags:        []string{},
	}

	if err := s.Save(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// Save writes session data to disk.
func (s *Store) Save(sess *Session) error {
	sess.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode session %s: %w", sess.ID, err)
	}

	filePath := filepath.Join(s.baseDir, fmt.Sprintf("%s.json", sess.ID))
	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves a session by ID.
func (s *Store) Get(id string) (*Session, error) {
	filePath := filepath.Join(s.baseDir, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("session %s not found: %w", id, err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("failed to decode session %s: %w", id, err)
	}
	return &sess, nil
}

// List returns all known sessions sorted newest first.
func (s *Store) List() ([]*Session, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	var sessions []*Session
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(s.baseDir, entry.Name()))
			if err != nil {
				continue
			}
			var sess Session
			if err := json.Unmarshal(data, &sess); err == nil {
				sessions = append(sessions, &sess)
			}
		}
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	return sessions, nil
}

// AddCheckpoint records a state checkpoint for the session.
func (s *Store) AddCheckpoint(sess *Session, note, gitCommit, diff string) (*Checkpoint, error) {
	cp := &Checkpoint{
		ID:        generateID(),
		SessionID: sess.ID,
		Timestamp: time.Now().UTC(),
		GitCommit: gitCommit,
		Diff:      diff,
		Note:      note,
	}

	sess.Checkpoints = append(sess.Checkpoints, cp)
	if err := s.Save(sess); err != nil {
		return nil, err
	}
	return cp, nil
}

func defaultDataDir() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "AgyPlusPlus")
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".agypp-data"
	}
	return filepath.Join(home, ".config", "agy-plus-plus")
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
