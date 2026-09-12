package policy

import (
	"errors"
	"fmt"
	"sync"
)

// ActivationSource represents where an approval mode change originated.
type ActivationSource string

const (
	SourceCLIArg          ActivationSource = "cli_arg"
	SourceInteractiveCmd  ActivationSource = "interactive_cmd"
	SourceKeyboardShortcut ActivationSource = "keyboard_shortcut"
	SourceUntrustedModel  ActivationSource = "untrusted_model"
	SourceUntrustedMCP    ActivationSource = "untrusted_mcp"
	SourceUntrustedFile   ActivationSource = "untrusted_file"
	SourceUntrustedAgent  ActivationSource = "untrusted_agent"
)

var (
	ErrUnauthorizedActivation = errors.New("security violation: YOLO mode can only be activated by direct human action")
	ErrSandboxConflict        = errors.New("cannot silently disable sandbox when activating YOLO mode")
)

// Manager manages approval modes and enforces prompt injection defense.
type Manager struct {
	mu          sync.RWMutex
	yoloEnabled bool
	planEnabled bool
	sandbox     bool
}

// NewManager creates a new policy manager.
func NewManager(sandbox bool) *Manager {
	return &Manager{
		sandbox: sandbox,
	}
}

// EnableYolo activates YOLO mode only if requested by an authorized human source.
func (m *Manager) EnableYolo(source ActivationSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch source {
	case SourceCLIArg, SourceInteractiveCmd, SourceKeyboardShortcut:
		m.yoloEnabled = true
		m.planEnabled = false
		return nil
	default:
		return fmt.Errorf("%w: rejected source '%s'", ErrUnauthorizedActivation, source)
	}
}

// DisableYolo deactivates YOLO mode.
func (m *Manager) DisableYolo() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.yoloEnabled = false
}

// IsYoloEnabled returns whether YOLO mode is currently active.
func (m *Manager) IsYoloEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.yoloEnabled
}

// SetPlanMode sets the plan mode.
func (m *Manager) SetPlanMode(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.planEnabled = enabled
	if enabled {
		m.yoloEnabled = false
	}
}

// IsPlanEnabled returns whether plan mode is currently active.
func (m *Manager) IsPlanEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.planEnabled
}

// SandboxActive returns whether sandboxing is enabled.
func (m *Manager) SandboxActive() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sandbox
}
