package unit

import (
	"errors"
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/policy"
)

func TestYoloAuthorizedActivation(t *testing.T) {
	mgr := policy.NewManager(false)

	if mgr.IsYoloEnabled() {
		t.Fatal("expected YOLO to be disabled by default")
	}

	if err := mgr.EnableYolo(policy.SourceCLIArg); err != nil {
		t.Fatalf("expected CLI activation to succeed, got: %v", err)
	}

	if !mgr.IsYoloEnabled() {
		t.Fatal("expected YOLO to be enabled after authorized activation")
	}

	mgr.DisableYolo()
	if mgr.IsYoloEnabled() {
		t.Fatal("expected YOLO to be disabled after DisableYolo()")
	}
}

func TestYoloPromptInjectionDefense(t *testing.T) {
	mgr := policy.NewManager(false)

	untrustedSources := []policy.ActivationSource{
		policy.SourceUntrustedModel,
		policy.SourceUntrustedMCP,
		policy.SourceUntrustedFile,
		policy.SourceUntrustedAgent,
	}

	for _, src := range untrustedSources {
		err := mgr.EnableYolo(src)
		if err == nil {
			t.Errorf("expected error for untrusted source '%s', got nil", src)
		}
		if !errors.Is(err, policy.ErrUnauthorizedActivation) {
			t.Errorf("expected ErrUnauthorizedActivation for '%s', got: %v", src, err)
		}
		if mgr.IsYoloEnabled() {
			t.Errorf("YOLO mode must remain false when attempted by untrusted source '%s'", src)
		}
	}
}

func TestPlanModeTogglesYolo(t *testing.T) {
	mgr := policy.NewManager(false)
	_ = mgr.EnableYolo(policy.SourceCLIArg)

	if !mgr.IsYoloEnabled() {
		t.Fatal("expected YOLO to be enabled")
	}

	mgr.SetPlanMode(true)
	if mgr.IsYoloEnabled() {
		t.Fatal("activating plan mode must deactivate YOLO mode")
	}
	if !mgr.IsPlanEnabled() {
		t.Fatal("expected plan mode to be enabled")
	}
}
