package unit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/tui"
	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

func TestPrintBanner_Standard(t *testing.T) {
	inst := &agy.Installation{Version: "1.2.2", BinaryPath: "agy.exe", Capabilities: &agy.Capabilities{}}
	var buf bytes.Buffer
	tui.PrintBanner(&buf, inst, false, false)
	out := buf.String()
	if !strings.Contains(out, "1.2.2") {
		t.Errorf("banner should contain AGY version, got: %s", out)
	}
	if !strings.Contains(out, "STANDARD") {
		t.Errorf("banner should show STANDARD mode, got: %s", out)
	}
}

func TestPrintBanner_Yolo(t *testing.T) {
	inst := &agy.Installation{Version: "1.2.2", BinaryPath: "agy.exe", Capabilities: &agy.Capabilities{}}
	var buf bytes.Buffer
	tui.PrintBanner(&buf, inst, true, false)
	out := buf.String()
	if !strings.Contains(out, "YOLO") {
		t.Errorf("banner should show YOLO mode, got: %s", out)
	}
}

func TestPrintBanner_Plan(t *testing.T) {
	inst := &agy.Installation{Version: "1.2.2", BinaryPath: "agy.exe", Capabilities: &agy.Capabilities{}}
	var buf bytes.Buffer
	tui.PrintBanner(&buf, inst, false, true)
	out := buf.String()
	if !strings.Contains(out, "PLAN") {
		t.Errorf("banner should show PLAN mode, got: %s", out)
	}
}
