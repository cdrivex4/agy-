package unit

import (
	"testing"

	"github.com/cdrivex4/agy-plus-plus/internal/upstream/agy"
)

func TestParseStreamLine_Init(t *testing.T) {
	line := []byte(`{"event":"init","conversation_id":"abc123","init":{"conversation_id":"abc123","permission_mode":"request-review","tools":["run_command"]}}`)
	evt, err := agy.ParseStreamLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Event != "init" {
		t.Errorf("expected event 'init', got: %s", evt.Event)
	}
	if evt.Init == nil {
		t.Fatal("expected Init to be populated")
	}
	if evt.Init.ConversationID != "abc123" {
		t.Errorf("expected conversation_id 'abc123', got: %s", evt.Init.ConversationID)
	}
}

func TestParseStreamLine_StepUpdate(t *testing.T) {
	line := []byte(`{"event":"step_update","step_update":{"conversation_id":"abc123","step_index":1,"state":"ACTIVE","step_type":"agent_response","text_delta":"hello"}}`)
	evt, err := agy.ParseStreamLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Event != "step_update" {
		t.Errorf("expected event 'step_update', got: %s", evt.Event)
	}
	if evt.StepUpdate == nil {
		t.Fatal("expected StepUpdate to be populated")
	}
	if evt.StepUpdate.TextDelta != "hello" {
		t.Errorf("expected text_delta 'hello', got: %s", evt.StepUpdate.TextDelta)
	}
	if evt.StepUpdate.StepType != "agent_response" {
		t.Errorf("expected step_type 'agent_response', got: %s", evt.StepUpdate.StepType)
	}
}

func TestParseStreamLine_Result_Success(t *testing.T) {
	line := []byte(`{"event":"result","result":{"conversation_id":"abc123","status":"SUCCESS","response":"pong\n","duration_seconds":2.1,"num_turns":1}}`)
	evt, err := agy.ParseStreamLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Event != "result" {
		t.Errorf("expected event 'result', got: %s", evt.Event)
	}
	if evt.Result == nil {
		t.Fatal("expected Result to be populated")
	}
	if evt.Result.Status != "SUCCESS" {
		t.Errorf("expected status 'SUCCESS', got: %s", evt.Result.Status)
	}
	if evt.Result.ConversationID != "abc123" {
		t.Errorf("expected conversation_id 'abc123', got: %s", evt.Result.ConversationID)
	}
}

func TestParseStreamLine_Result_Error(t *testing.T) {
	line := []byte(`{"event":"result","result":{"conversation_id":"x","status":"ERROR","response":"","error":"stream input message is missing the event field","duration_seconds":0,"num_turns":0}}`)
	evt, err := agy.ParseStreamLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Result.Status != "ERROR" {
		t.Errorf("expected ERROR status, got: %s", evt.Result.Status)
	}
	if evt.Result.Error == "" {
		t.Error("expected non-empty error string")
	}
}

func TestParseStreamLine_InvalidJSON(t *testing.T) {
	line := []byte(`not json at all`)
	_, err := agy.ParseStreamLine(line)
	if err == nil {
		t.Error("expected error parsing invalid JSON, got nil")
	}
}

func TestParseStreamLine_UnknownEvent(t *testing.T) {
	line := []byte(`{"event":"future_event","some_field":"value"}`)
	evt, err := agy.ParseStreamLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Event != "future_event" {
		t.Errorf("expected 'future_event', got: %s", evt.Event)
	}
	// All typed fields should be nil (graceful unknown event handling)
	if evt.Init != nil || evt.StepUpdate != nil || evt.Result != nil {
		t.Error("unknown event should leave all typed fields nil")
	}
}
