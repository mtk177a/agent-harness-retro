package pilot

import (
	"encoding/json"
	"errors"
)

func validateEnvelope(operation string, raw []byte) (Envelope, error) {
	var gate struct {
		SchemaVersion string `json:"schema_version"`
	}
	if err := json.Unmarshal(raw, &gate); err != nil {
		return Envelope{}, errors.New("invalid JSON response")
	}
	if gate.SchemaVersion != SchemaVersion {
		return Envelope{}, errors.New("unsupported response version")
	}
	var shape struct {
		CLIVersion string          `json:"cli_version"`
		Data       json.RawMessage `json:"data"`
		Page       json.RawMessage `json:"page"`
		Omissions  json.RawMessage `json:"omissions"`
	}
	if err := json.Unmarshal(raw, &shape); err != nil || shape.CLIVersion == "" || len(shape.Data) == 0 || string(shape.Data) == "null" || len(shape.Omissions) == 0 || string(shape.Omissions) == "null" {
		return Envelope{}, errors.New("missing required response fields")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(shape.Data, &fields); err != nil {
		return Envelope{}, errors.New("invalid operation data")
	}
	var required string
	switch operation {
	case "list":
		required = "sources"
	case "show":
		required = "source"
	case "events":
		required = "events"
	case "verify":
		required = "verified_version"
	default:
		return Envelope{}, errors.New("unknown operation")
	}
	if len(fields[required]) == 0 || string(fields[required]) == "null" {
		return Envelope{}, errors.New("missing operation data")
	}
	var envelope Envelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return Envelope{}, errors.New("invalid response shape")
	}
	if envelope.Operation != operation {
		return Envelope{}, errors.New("unexpected operation")
	}
	switch envelope.Status {
	case "complete":
		if len(envelope.Omissions) != 0 {
			return Envelope{}, errors.New("complete response has omissions")
		}
	case "partial":
		if len(envelope.Omissions) == 0 {
			return Envelope{}, errors.New("partial response lacks omissions")
		}
	case "unsupported", "error":
		return Envelope{}, errors.New("operation unavailable")
	default:
		return Envelope{}, errors.New("unknown response status")
	}
	if (operation == "list" || operation == "events") && (envelope.Page == nil || envelope.Page.Limit < 1 || envelope.Page.Limit > PageLimit) {
		return Envelope{}, errors.New("missing page")
	}
	return envelope, nil
}

func completeOrPagination(envelope Envelope) error {
	if envelope.Page == nil {
		return errors.New("missing page")
	}
	if envelope.Page.HasMore {
		if envelope.Page.NextCursor == "" || envelope.Status != "partial" || len(envelope.Omissions) != 1 || envelope.Omissions[0].Code != "pagination" {
			return errors.New("incomplete page with unexpected omissions")
		}
		return nil
	}
	if envelope.Page.NextCursor != "" || envelope.Status != "complete" {
		return errors.New("incomplete final page")
	}
	return nil
}

func summarizeEvent(raw json.RawMessage, expected uint64) (EventSummary, error) {
	var event struct {
		Index      *uint64 `json:"index"`
		Kind       string  `json:"kind"`
		TimeState  string  `json:"time_state"`
		RecordedAt string  `json:"recorded_at"`
		Message    *struct {
			Role      string  `json:"role"`
			Text      *string `json:"text"`
			Truncated *bool   `json:"truncated"`
		} `json:"message"`
		ToolCall *struct {
			Category   string          `json:"category"`
			Action     string          `json:"action"`
			InputState string          `json:"input_state"`
			Input      json.RawMessage `json:"input"`
		} `json:"tool_call"`
		ToolResult *struct {
			Outcome          string          `json:"outcome"`
			CorrelationState string          `json:"correlation_state"`
			ContentState     string          `json:"content_state"`
			Content          json.RawMessage `json:"content"`
		} `json:"tool_result"`
		Error *struct {
			Category string  `json:"category"`
			Message  *string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &event) != nil || event.Index == nil || *event.Index != expected {
		return EventSummary{}, errors.New("invalid event index")
	}
	summary := EventSummary{Index: expected, Kind: event.Kind, TimeState: event.TimeState, RecordedAt: event.RecordedAt}
	if !validEventKind(event.Kind) {
		return EventSummary{}, errors.New("invalid event kind")
	}
	if event.TimeState != "available" {
		return EventSummary{}, errors.New("event time unavailable")
	}
	if _, err := parseCanonicalTime(event.RecordedAt); err != nil {
		return EventSummary{}, errors.New("invalid event time")
	}
	payloads := 0
	if event.Message != nil {
		payloads++
	}
	if event.ToolCall != nil {
		payloads++
	}
	if event.ToolResult != nil {
		payloads++
	}
	if event.Error != nil {
		payloads++
	}
	if payloads != 1 {
		return EventSummary{}, errors.New("invalid event payload")
	}
	switch event.Kind {
	case "message":
		if event.Message == nil || (event.Message.Role != "user" && event.Message.Role != "assistant") || event.Message.Text == nil || event.Message.Truncated == nil {
			return EventSummary{}, errors.New("invalid message")
		}
		summary.Role = event.Message.Role
	case "tool_call":
		if event.ToolCall == nil || event.ToolCall.Category == "" || !validContentState(event.ToolCall.InputState) || !validRecordedContent(event.ToolCall.InputState, event.ToolCall.Input) {
			return EventSummary{}, errors.New("invalid tool call")
		}
		if event.ToolCall.InputState == "unavailable" || event.ToolCall.InputState == "unsupported" {
			return EventSummary{}, errors.New("tool input unavailable")
		}
		summary.InputState = event.ToolCall.InputState
	case "tool_result":
		if event.ToolResult == nil || !validContentState(event.ToolResult.ContentState) || !validRecordedContent(event.ToolResult.ContentState, event.ToolResult.Content) || (event.ToolResult.Outcome != "success" && event.ToolResult.Outcome != "failure" && event.ToolResult.Outcome != "unknown") || (event.ToolResult.CorrelationState != "matched" && event.ToolResult.CorrelationState != "unmatched" && event.ToolResult.CorrelationState != "ambiguous") {
			return EventSummary{}, errors.New("invalid tool result")
		}
		if event.ToolResult.Outcome == "unknown" || event.ToolResult.CorrelationState != "matched" || event.ToolResult.ContentState == "unavailable" || event.ToolResult.ContentState == "unsupported" {
			return EventSummary{}, errors.New("tool result incomplete")
		}
		summary.Outcome, summary.CorrelationState, summary.ContentState = event.ToolResult.Outcome, event.ToolResult.CorrelationState, event.ToolResult.ContentState
	case "error":
		if event.Error == nil || event.Error.Category == "" || event.Error.Message == nil {
			return EventSummary{}, errors.New("invalid error event")
		}
	}
	return summary, nil
}

func validContentState(state string) bool {
	return state == "available" || state == "absent" || state == "unavailable" || state == "unsupported"
}

func validRecordedContent(state string, raw json.RawMessage) bool {
	if state != "available" {
		return len(raw) == 0
	}
	var content struct {
		Format    string  `json:"format"`
		Text      *string `json:"text"`
		Truncated *bool   `json:"truncated"`
	}
	return json.Unmarshal(raw, &content) == nil && (content.Format == "text" || content.Format == "json") && content.Text != nil && content.Truncated != nil
}
