package session

import (
	"bytes"
	"os/exec"
	"strings"
)

// GitHelper provides lightweight inspection of workspace git status for checkpoints.
type GitHelper struct {
	Workspace string
}

// NewGitHelper creates a GitHelper for the given workspace.
func NewGitHelper(workspace string) *GitHelper {
	return &GitHelper{Workspace: workspace}
}

// IsGitRepo checks if the workspace is inside a git work tree.
func (g *GitHelper) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = g.Workspace
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// GetHeadCommit returns the current HEAD commit hash (short).
func (g *GitHelper) GetHeadCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = g.Workspace
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// GetDiff returns the uncommitted git diff in the workspace.
func (g *GitHelper) GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	cmd.Dir = g.Workspace
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	return out.String(), nil
}

// GetStatus returns the short git status.
func (g *GitHelper) GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "-s")
	cmd.Dir = g.Workspace
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
