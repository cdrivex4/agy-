package unit

import (
	"strings"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/version"
)

func TestVersionString(t *testing.T) {
	str := version.String()
	if !strings.Contains(str, "AGY++") {
		t.Errorf("expected 'AGY++' in version string, got: %s", str)
	}
	if !strings.Contains(str, version.BaselineAgy) {
		t.Errorf("expected baseline AGY '%s' in version string, got: %s", version.BaselineAgy, str)
	}
}
