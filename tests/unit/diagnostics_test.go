package unit

import (
	"strings"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/diagnostics"
)

func TestDiagnosticsRedaction(t *testing.T) {
	raw := "User session with ya29.a0AfH6SMDI839210-fake-oauth-token and API key AIzaSyD9x82910abcdefghij1234567890123"
	redacted := diagnostics.Redact(raw)

	if strings.Contains(redacted, "ya29.a0AfH6SMDI839210") {
		t.Errorf("expected OAuth token to be redacted, got: %s", redacted)
	}
	if strings.Contains(redacted, "AIzaSyD9x82910") {
		t.Errorf("expected API key to be redacted, got: %s", redacted)
	}
	if !strings.Contains(redacted, "[REDACTED]") {
		t.Errorf("expected '[REDACTED]' placeholder in output, got: %s", redacted)
	}
}
