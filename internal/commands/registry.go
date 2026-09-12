package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

// Command represents an interactive slash command.
type Command struct {
	Name               string
	Aliases            []string
	Description        string
	RequiredCapability string
	Handler            func(args []string) error
}

// Registry stores registered commands.
type Registry struct {
	commands map[string]*Command
}

// NewRegistry initializes the command registry with core commands.
func NewRegistry() *Registry {
	r := &Registry{
		commands: make(map[string]*Command),
	}
	r.registerDefaults()
	return r
}

// Register adds a command to the registry.
func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

// Lookup finds a command by name or alias.
func (r *Registry) Lookup(name string) (*Command, bool) {
	name = strings.TrimPrefix(name, "/")
	cmd, ok := r.commands[name]
	return cmd, ok
}

// List returns a sorted list of unique commands.
func (r *Registry) List() []*Command {
	seen := make(map[string]bool)
	var list []*Command
	for _, cmd := range r.commands {
		if !seen[cmd.Name] {
			seen[cmd.Name] = true
			list = append(list, cmd)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

// CheckCapability verifies that the given command is supported by the installed backend.
func (r *Registry) CheckCapability(cmd *Command, caps *agy.Capabilities) error {
	if cmd.RequiredCapability == "" {
		return nil
	}
	if caps == nil || !caps.Supports(cmd.RequiredCapability) {
		return fmt.Errorf("command '/%s' is unavailable with the current AGY backend (missing capability: %s)",
			cmd.Name, cmd.RequiredCapability)
	}
	return nil
}

func (r *Registry) registerDefaults() {
	r.Register(&Command{
		Name:        "help",
		Aliases:     []string{"h", "?"},
		Description: "Show available commands and usage help",
	})
	r.Register(&Command{
		Name:        "status",
		Description: "Show active session state, model, and permission mode",
	})
	r.Register(&Command{
		Name:               "yolo",
		Aliases:            []string{"y"},
		Description:        "Toggle YOLO mode (automatic permission approval)",
		RequiredCapability: "yolo",
	})
	r.Register(&Command{
		Name:               "plan",
		Description:        "Toggle plan mode (read-only planning guardrails)",
		RequiredCapability: "plan",
	})
	r.Register(&Command{
		Name:        "diagnostics",
		Description: "Inspect system, AGY installation, and compatibility status",
	})
	r.Register(&Command{
		Name:        "quit",
		Aliases:     []string{"exit", "q"},
		Description: "Exit the AGY++ session",
	})
}
