package unit

import (
	"strings"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/commands"
	"github.com/cdrivex4/agy-plus-plus/internal/session"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

func TestCommandsRegistry(t *testing.T) {
	reg := commands.NewRegistry()

	// 1. Lookup
	cmd, ok := reg.Lookup("help")
	if !ok {
		t.Fatal("expected 'help' command to exist")
	}
	if cmd.Name != "help" {
		t.Errorf("expected name 'help', got: %s", cmd.Name)
	}

	// 2. Lookup alias
	cmdAlias, ok := reg.Lookup("h")
	if !ok || cmdAlias.Name != "help" {
		t.Fatal("expected 'h' alias to map to 'help'")
	}

	// 3. Status command handler
	statusCmd, ok := reg.Lookup("status")
	if !ok {
		t.Fatal("expected 'status' command")
	}

	ctx := &commands.CommandContext{
		Installation: &agy.Installation{
			BinaryPath: "agy.exe",
			Version:    "1.2.2",
			Capabilities: &agy.Capabilities{
				SupportsYolo: true,
			},
		},
		Session: &session.Session{
			ID:        "test-session-123",
			Workspace: "d:\\Dev\\test",
		},
		YoloEnabled: true,
	}

	out, err := statusCmd.Handler(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-session-123") {
		t.Errorf("expected session ID in status output, got: %s", out)
	}
	if !strings.Contains(out, "YOLO Mode:     ENABLED") {
		t.Errorf("expected YOLO enabled in status output, got: %s", out)
	}

	// 4. YOLO toggle handler
	yoloCmd, _ := reg.Lookup("yolo")
	var toggledVal bool
	ctx.SetYolo = func(enable bool) error {
		toggledVal = enable
		ctx.YoloEnabled = enable
		return nil
	}

	out, err = yoloCmd.Handler(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error toggling yolo: %v", err)
	}
	if !strings.Contains(out, "YOLO MODE DISABLED") {
		t.Errorf("expected YOLO disabled response, got: %s", out)
	}
	if toggledVal != false {
		t.Errorf("expected toggledVal to be false")
	}
}
