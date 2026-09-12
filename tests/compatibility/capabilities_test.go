package compatibility

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

func TestParseCapabilitiesAgainstFixture(t *testing.T) {
	fixturePath := filepath.Join("..", "fixtures", "agy", "help.txt")
	content, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to read fixture file: %v", err)
	}

	caps := agy.ParseCapabilities(string(content))

	if !caps.SupportsYolo {
		t.Errorf("expected SupportsYolo to be true")
	}
	if !caps.SupportsPlan {
		t.Errorf("expected SupportsPlan to be true")
	}
	if !caps.SupportsStreamJSON {
		t.Errorf("expected SupportsStreamJSON to be true")
	}
	if !caps.SupportsInteractivePrompt {
		t.Errorf("expected SupportsInteractivePrompt to be true")
	}
	if !caps.SupportsSandbox {
		t.Errorf("expected SupportsSandbox to be true")
	}
	if !caps.SupportsMCP {
		t.Errorf("expected SupportsMCP to be true")
	}
	if !caps.SupportsModels {
		t.Errorf("expected SupportsModels to be true")
	}
}
