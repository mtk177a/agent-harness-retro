package pilot

import (
	"strings"
	"testing"
	"time"
)

func TestRunExcludesTruncatedMessage(t *testing.T) {
	event := syntheticEvent(0)
	event["message"].(map[string]any)["truncated"] = true
	assertIncompleteEventExcluded(t, event)
}

func TestRunExcludesTruncatedToolContent(t *testing.T) {
	for _, kind := range []string{"tool_call", "tool_result"} {
		t.Run(kind, func(t *testing.T) {
			content := map[string]any{"format": "text", "text": "synthetic", "truncated": true}
			assertIncompleteEventExcluded(t, syntheticToolEvent(t, kind, content))
		})
	}
}

func TestRunExcludesOmittedToolContent(t *testing.T) {
	for _, kind := range []string{"tool_call", "tool_result"} {
		t.Run(kind, func(t *testing.T) {
			content := map[string]any{"format": "text", "text": "synthetic", "truncated": false, "omitted": true}
			assertIncompleteEventExcluded(t, syntheticToolEvent(t, kind, content))
		})
	}
}

func TestRunExcludesInvalidJSONToolContent(t *testing.T) {
	for _, kind := range []string{"tool_call", "tool_result"} {
		t.Run(kind, func(t *testing.T) {
			content := map[string]any{"format": "json", "text": "{", "truncated": false}
			assertIncompleteEventExcluded(t, syntheticToolEvent(t, kind, content))
		})
	}
}

func TestRunAcceptsEmptyMessage(t *testing.T) {
	event := syntheticEvent(0)
	event["message"].(map[string]any)["text"] = ""
	assertEventAccepted(t, event)
}

func TestRunAcceptsCompleteToolContent(t *testing.T) {
	for _, kind := range []string{"tool_call", "tool_result"} {
		for _, value := range []struct {
			name, format, text string
		}{
			{"empty_text", "text", ""},
			{"empty_object", "json", "{}"},
			{"empty_array", "json", "[]"},
			{"empty_json_string", "json", `""`},
		} {
			t.Run(kind+"/"+value.name, func(t *testing.T) {
				content := map[string]any{"format": value.format, "text": value.text, "truncated": false, "omitted": false}
				assertEventAccepted(t, syntheticToolEvent(t, kind, content))
			})
		}
	}
}

func TestRunAcceptsAbsentToolContent(t *testing.T) {
	for _, kind := range []string{"tool_call", "tool_result"} {
		t.Run(kind, func(t *testing.T) {
			event := syntheticToolEvent(t, kind, nil)
			payload := event[kind].(map[string]any)
			if kind == "tool_call" {
				payload["input_state"] = "absent"
				delete(payload, "input")
			} else {
				payload["content_state"] = "absent"
				delete(payload, "content")
			}
			assertEventAccepted(t, event)
		})
	}
}

func assertEventAccepted(t *testing.T, event map[string]any) {
	t.Helper()
	bundle := runWithEvent(t, event)
	count := bundle.Counts["codex"]
	if bundle.Status != "ready" || count.Adopted != 3 || count.Excluded != 0 {
		t.Fatalf("status=%s reason=%s codex=%+v", bundle.Status, bundle.Reason, count)
	}
}

func syntheticToolEvent(t *testing.T, kind string, content map[string]any) map[string]any {
	t.Helper()
	event := syntheticEvent(0)
	delete(event, "message")
	event["kind"] = kind
	switch kind {
	case "tool_call":
		event[kind] = map[string]any{"category": "shell", "input_state": "available", "input": content}
	case "tool_result":
		event[kind] = map[string]any{"outcome": "success", "correlation_state": "matched", "content_state": "available", "content": content}
	default:
		t.Fatalf("unsupported synthetic tool event: %s", kind)
	}
	return event
}

func assertIncompleteEventExcluded(t *testing.T, event map[string]any) {
	t.Helper()
	bundle := runWithEvent(t, event)
	count := bundle.Counts["codex"]
	if bundle.Status != "insufficient_evidence" || bundle.Reason != "minimum_not_met" || count.Adopted != 2 || count.Excluded != 1 {
		t.Fatalf("status=%s reason=%s codex=%+v", bundle.Status, bundle.Reason, count)
	}
}

func runWithEvent(t *testing.T, event map[string]any) Bundle {
	t.Helper()
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		ref := args[len(args)-1]
		if args[0] == "events" && strings.HasPrefix(ref, "as0:codex:") && strings.HasSuffix(ref, "1") {
			// A complete envelope must not hide an incomplete event body.
			return response("events", "complete", map[string]any{"source_ref": ref, "events": []any{event}}, map[string]any{"limit": 100, "has_more": false}, []any{}), nil
		}
		return base(args...)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	return bundle
}
