package agy

import "encoding/json"

// StreamEvent is the top-level envelope for all stream-json events.
type StreamEvent struct {
	Event      string          `json:"event"`
	Raw        json.RawMessage `json:"-"`

	// Populated by event type
	Init       *InitEvent       `json:"init,omitempty"`
	StepUpdate *StepUpdateEvent `json:"step_update,omitempty"`
	Result     *ResultEvent     `json:"result,omitempty"`
}

// InitEvent is emitted once per session on startup.
type InitEvent struct {
	ConversationID string   `json:"conversation_id"`
	PermissionMode string   `json:"permission_mode"`
	Tools          []string `json:"tools"`
}

// StepUpdateEvent is emitted as the agent produces output.
type StepUpdateEvent struct {
	ConversationID string  `json:"conversation_id"`
	StepIndex      int     `json:"step_index"`
	State          string  `json:"state"` // ACTIVE, DONE
	StepType       string  `json:"step_type"` // user_input, agent_response, tool_call
	TextDelta      string  `json:"text_delta,omitempty"`
	DurationSecs   float64 `json:"duration_seconds,omitempty"`
}

// ResultEvent is the final event for each turn.
type ResultEvent struct {
	ConversationID  string  `json:"conversation_id"`
	Status          string  `json:"status"` // SUCCESS, ERROR
	Response        string  `json:"response"`
	Error           string  `json:"error,omitempty"`
	DurationSecs    float64 `json:"duration_seconds"`
	NumTurns        int     `json:"num_turns"`
}

// ParseStreamLine parses one NDJSON line from agy --output-format stream-json.
func ParseStreamLine(line []byte) (*StreamEvent, error) {
	var env struct {
		Event      string          `json:"event"`
		Init       *InitEvent      `json:"init"`
		StepUpdate *StepUpdateEvent `json:"step_update"`
		Result     *ResultEvent    `json:"result"`
	}
	if err := json.Unmarshal(line, &env); err != nil {
		return nil, err
	}
	return &StreamEvent{
		Event:      env.Event,
		Init:       env.Init,
		StepUpdate: env.StepUpdate,
		Result:     env.Result,
	}, nil
}
