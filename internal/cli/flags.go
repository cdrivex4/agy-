package cli

import (
	"flag"
)

// ParsedFlags holds all CLI flags parsed before execution.
type ParsedFlags struct {
	ShowVersion       bool
	YoloMode          bool
	PlanMode          bool
	Sandbox           bool
	PrintPrompt       string
	InteractivePrompt string
	OutputFormat      string
	Continue          bool
	ConversationID    string
	Model             string
	Project           string
	Effort            string
	ExtraArgs         []string
}

// ParseFlags parses CLI arguments into a structured ParsedFlags object.
// It uses a lenient parse mode so unknown flags bubble up as ExtraArgs.
func ParseFlags(args []string) (*ParsedFlags, error) {
	fs := flag.NewFlagSet("agy++", flag.ContinueOnError)

	showVersion := fs.Bool("version", false, "")
	fs.Bool("v", false, "")

	yoloLong := fs.Bool("yolo", false, "")
	yoloShort := fs.Bool("y", false, "")

	planMode := fs.Bool("plan", false, "")
	sandbox := fs.Bool("sandbox", false, "")

	printShort := fs.String("p", "", "")
	printLong := fs.String("print", "", "")
	promptAlias := fs.String("prompt", "", "")

	interactiveShort := fs.String("i", "", "")
	interactiveLong := fs.String("prompt-interactive", "", "")

	continueLong := fs.Bool("continue", false, "")
	continueShort := fs.Bool("c", false, "")

	conversationID := fs.String("conversation", "", "")
	model := fs.String("model", "", "")
	project := fs.String("project", "", "")
	effort := fs.String("effort", "", "")
	outputFormat := fs.String("output-format", "text", "")

	// Permissive: ignore unknown flags rather than aborting
	if err := fs.Parse(args); err != nil {
		// non-fatal — extra passthrough args collected below
	}

	// Merge aliases
	print := *printShort
	if *printLong != "" {
		print = *printLong
	}
	if *promptAlias != "" {
		print = *promptAlias
	}

	interactive := *interactiveShort
	if *interactiveLong != "" {
		interactive = *interactiveLong
	}

	return &ParsedFlags{
		ShowVersion:       *showVersion,
		YoloMode:          *yoloLong || *yoloShort,
		PlanMode:          *planMode,
		Sandbox:           *sandbox,
		PrintPrompt:       print,
		InteractivePrompt: interactive,
		OutputFormat:      *outputFormat,
		Continue:          *continueLong || *continueShort,
		ConversationID:    *conversationID,
		Model:             *model,
		Project:           *project,
		Effort:            *effort,
		ExtraArgs:         fs.Args(),
	}, nil
}
