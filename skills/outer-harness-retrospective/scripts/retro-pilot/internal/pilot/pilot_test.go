package pilot

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestEnvelopeRejectsUnsupportedVersionBeforeData(t *testing.T) {
	for _, value := range []string{"v1", "v3", ""} {
		input := map[string]any{"schema_version": value, "cli_version": "dev", "operation": "list", "status": "complete", "data": map[string]any{"sources": []any{}}, "page": map[string]any{"limit": 100, "has_more": false}, "omissions": []any{}}
		raw, _ := json.Marshal(input)
		if _, err := validateEnvelope("list", raw); err == nil {
			t.Fatalf("accepted schema_version=%q", value)
		}
	}
}

func TestRunNeverRetainsRecordedContent(t *testing.T) {
	bundle, err := Run(syntheticCaller(), time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "ready" {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
	encoded, err := json.Marshal(bundle)
	if err != nil || strings.Contains(string(encoded), "synthetic-private-token") {
		t.Fatalf("private content retained: %s", encoded)
	}
}

func TestRunDoesNotRetainSourceMetadata(t *testing.T) {
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] != "show" {
			return base(args...)
		}
		raw, _ := base(args...)
		var value map[string]any
		json.Unmarshal(raw, &value)
		data := value["data"].(map[string]any)
		source := data["source"].(map[string]any)
		source["metadata"] = []any{map[string]any{"name": "opaque", "value": "synthetic-private-metadata"}}
		return json.Marshal(value)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "ready" {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
	encoded, _ := json.Marshal(bundle)
	if strings.Contains(string(encoded), "synthetic-private-metadata") {
		t.Fatal("metadata retained")
	}
}

func TestV2IncompleteEventExcludesCandidate(t *testing.T) {
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] == "events" && strings.HasSuffix(args[len(args)-1], "1") {
			return response("events", "partial", map[string]any{"source_ref": args[len(args)-1], "events": []any{syntheticEvent(0)}}, map[string]any{"limit": 100, "has_more": false}, []any{map[string]any{"code": "tool_outcome_unknown", "scope": "event", "message": "unknown"}}), nil
		}
		return base(args...)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "insufficient_evidence" || bundle.Reason != "minimum_not_met" {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
}

func TestEventSummaryDropsToolBodies(t *testing.T) {
	for _, event := range []string{
		`{"index":0,"kind":"tool_call","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_call":{"category":"shell","input_state":"available","input":{"format":"text","text":"synthetic-private-input","truncated":false}},"metadata":[]}`,
		`{"index":0,"kind":"tool_result","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_result":{"outcome":"success","correlation_state":"matched","content_state":"available","content":{"format":"text","text":"synthetic-private-result","truncated":false}},"metadata":[]}`,
	} {
		summary, err := summarizeEvent([]byte(event), 0)
		if err != nil {
			t.Fatal(err)
		}
		encoded, _ := json.Marshal(summary)
		if strings.Contains(string(encoded), "synthetic-private") {
			t.Fatalf("content retained: %s", encoded)
		}
	}
}

func TestEventSummaryRejectsAvailableStateWithoutBody(t *testing.T) {
	for _, event := range []string{
		`{"index":0,"kind":"tool_call","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_call":{"category":"shell","input_state":"available"},"metadata":[]}`,
		`{"index":0,"kind":"tool_result","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_result":{"outcome":"success","correlation_state":"matched","content_state":"available"},"metadata":[]}`,
	} {
		if _, err := summarizeEvent([]byte(event), 0); err == nil {
			t.Fatalf("accepted missing body: %s", event)
		}
	}
}

func TestEventSummaryRejectsStatesThatRequireOmissions(t *testing.T) {
	for _, event := range []string{
		`{"index":0,"kind":"message","time_state":"unavailable","message":{"role":"user","text":"synthetic","truncated":false},"metadata":[]}`,
		`{"index":0,"kind":"tool_result","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_result":{"outcome":"unknown","correlation_state":"matched","content_state":"absent"},"metadata":[]}`,
		`{"index":0,"kind":"tool_result","time_state":"available","recorded_at":"2026-09-20T01:00:00Z","tool_result":{"outcome":"success","correlation_state":"ambiguous","content_state":"absent"},"metadata":[]}`,
	} {
		if _, err := summarizeEvent([]byte(event), 0); err == nil {
			t.Fatalf("accepted incomplete state: %s", event)
		}
	}
}

func TestRunRejectsMissingSourceTime(t *testing.T) {
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] != "list" {
			return base(args...)
		}
		raw, _ := base(args...)
		var value map[string]any
		json.Unmarshal(raw, &value)
		sources := value["data"].(map[string]any)["sources"].([]any)
		delete(sources[0].(map[string]any), "last_interaction_at")
		return json.Marshal(value)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "insufficient_evidence" || bundle.Reason != "list_incomplete" {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
}

func TestCohortGroupsExplicitLineageAndLimitsCandidates(t *testing.T) {
	start := time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)
	sources := []Source{
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:a"}, LastInteractionAt: "2026-09-20T09:00:00Z", Relationships: []Relationship{{Kind: "parent", SourceRef: "as0:codex:one:root"}}},
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:b"}, LastInteractionAt: "2026-09-20T10:00:00Z", Relationships: []Relationship{{Kind: "forked_from", SourceRef: "as0:codex:one:root"}}},
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:c"}, LastInteractionAt: "2026-09-20T10:00:00Z"},
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:old"}, LastInteractionAt: "2026-09-01T00:00:00Z"},
	}
	groups, err := selectCandidates(sources, start, "codex", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || len(groups[0].Sources) != 2 || groups[0].Sources[0].Identity.SourceRef != "as0:codex:one:a" {
		t.Fatalf("groups = %#v", groups)
	}
}

func TestUnknownRelationshipExcludesOnlyItsCandidate(t *testing.T) {
	started := time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)
	sources := []Source{
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:a"}, LastInteractionAt: "2026-09-20T10:00:00Z", Relationships: []Relationship{{Kind: "unknown", SourceRef: "as0:codex:one:z"}}},
		{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: "as0:codex:one:b"}, LastInteractionAt: "2026-09-20T09:00:00Z"},
	}
	groups, err := selectCandidates(sources, started, "codex", 10)
	if err != nil || len(groups) != 2 || !groups[0].ambiguous || groups[1].ambiguous {
		t.Fatalf("groups=%#v err=%v", groups, err)
	}
}

func TestRunAcceptsThreeCompleteIndependentSessionsPerProvider(t *testing.T) {
	started := time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)
	bundle, err := Run(syntheticCaller(), started)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Status != "ready" || bundle.EventCount != 6 || bundle.Counts["codex"].Adopted != 3 || bundle.Counts["claude"].Adopted != 3 {
		t.Fatalf("bundle = %#v", bundle)
	}
}

func TestRunRejectsIncompleteCohort(t *testing.T) {
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] == "list" && argumentAfter(args, "--provider") == "codex" {
			raw, _ := base(args...)
			var value map[string]any
			json.Unmarshal(raw, &value)
			value["status"] = "partial"
			value["omissions"] = []any{map[string]any{"code": "source_time_unavailable", "scope": "source"}}
			return json.Marshal(value)
		}
		return base(args...)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "insufficient_evidence" || bundle.Reason != "list_incomplete" {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
}

func TestRunExcludesChangedVerifiedVersion(t *testing.T) {
	base := syntheticCaller()
	checks := 0
	call := func(args ...string) ([]byte, error) {
		if args[0] == "verify" && strings.HasPrefix(args[len(args)-1], "as0:codex:") && strings.HasSuffix(args[len(args)-1], "1") {
			checks++
			if checks == 2 {
				return response("verify", "complete", map[string]any{"source_ref": args[len(args)-1], "verified_version": Version{Algorithm: "sha256", Basis: "provider-content-v0", Value: "sha256:" + strings.Repeat("b", 64)}}, nil, []any{}), nil
			}
		}
		return base(args...)
	}
	bundle, err := Run(call, time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC))
	if err != nil || bundle.Status != "insufficient_evidence" || bundle.Counts["codex"].Excluded != 1 || bundle.Counts["codex"].Adopted != 2 {
		t.Fatalf("bundle=%#v err=%v", bundle, err)
	}
}

func TestFetchCandidateStopsAtRemainingPlusOne(t *testing.T) {
	ref := "as0:codex:one:synthetic"
	source := Source{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: ref}, LastInteractionAt: "2026-09-20T10:00:00Z", Relationships: []Relationship{}, Metadata: []Metadata{}}
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] == "show" {
			return response("show", "complete", map[string]any{"source": source}, nil, []any{}), nil
		}
		if args[0] == "events" {
			if argumentAfter(args, "--limit") != "2" {
				t.Fatalf("limit=%q", argumentAfter(args, "--limit"))
			}
			return response("events", "complete", map[string]any{"source_ref": ref, "events": []any{syntheticEvent(0), syntheticEvent(1)}}, map[string]any{"limit": 2, "has_more": false}, []any{}), nil
		}
		return base(args...)
	}
	_, _, overflow, err := fetchCandidate(call, candidate{Sources: []Source{source}}, 1)
	if err != nil || !overflow {
		t.Fatalf("overflow=%v err=%v", overflow, err)
	}
}

func TestFetchCandidateCollectsAllEventPages(t *testing.T) {
	ref := "as0:codex:one:synthetic"
	source := Source{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: ref}, LastInteractionAt: "2026-09-20T10:00:00Z", Relationships: []Relationship{}, Metadata: []Metadata{}}
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] == "show" {
			return response("show", "complete", map[string]any{"source": source}, nil, []any{}), nil
		}
		if args[0] == "events" {
			if argumentAfter(args, "--cursor") == "next" {
				return response("events", "complete", map[string]any{"source_ref": ref, "events": []any{syntheticEvent(1)}}, map[string]any{"limit": 100, "has_more": false}, []any{}), nil
			}
			return response("events", "partial", map[string]any{"source_ref": ref, "events": []any{syntheticEvent(0)}}, map[string]any{"limit": 100, "has_more": true, "next_cursor": "next"}, []any{map[string]any{"code": "pagination", "scope": "page"}}), nil
		}
		return base(args...)
	}
	lineage, size, overflow, err := fetchCandidate(call, candidate{Sources: []Source{source}}, 100)
	if err != nil || overflow || size != 2 || len(lineage.Sources[0].Events) != 2 {
		t.Fatalf("lineage=%#v size=%d overflow=%v err=%v", lineage, size, overflow, err)
	}
}

func TestFinalEventsPageMustBeComplete(t *testing.T) {
	ref := "as0:codex:one:synthetic"
	source := Source{Identity: Identity{Provider: "codex", SourceInstance: "one", SourceRef: ref}, LastInteractionAt: "2026-09-20T10:00:00Z", Relationships: []Relationship{}, Metadata: []Metadata{}}
	base := syntheticCaller()
	call := func(args ...string) ([]byte, error) {
		if args[0] == "show" {
			return response("show", "complete", map[string]any{"source": source}, nil, []any{}), nil
		}
		if args[0] == "events" {
			return response("events", "partial", map[string]any{"source_ref": ref, "events": []any{map[string]any{"index": 0, "kind": "message"}}}, map[string]any{"limit": 100, "has_more": false}, []any{map[string]any{"code": "output_redacted", "scope": "event"}}), nil
		}
		return base(args...)
	}
	_, _, overflow, err := fetchCandidate(call, candidate{Sources: []Source{source}}, 100)
	if err == nil || overflow {
		t.Fatalf("overflow=%v err=%v", overflow, err)
	}
}

func syntheticCaller() Caller {
	return func(args ...string) ([]byte, error) {
		operation := args[0]
		switch operation {
		case "list":
			provider := argumentAfter(args, "--provider")
			var sources []Source
			for i := 0; i < 3; i++ {
				ref := fmt.Sprintf("as0:%s:one:%064d", provider, i+1)
				sources = append(sources, Source{Identity: Identity{Provider: provider, SourceInstance: "one", SourceRef: ref}, LastInteractionAt: fmt.Sprintf("2026-09-20T%02d:00:00Z", i+1), Relationships: []Relationship{}, Metadata: []Metadata{}})
			}
			return response("list", "complete", map[string]any{"sources": sources}, map[string]any{"limit": 100, "has_more": false}, []any{}), nil
		case "show":
			ref := args[len(args)-1]
			provider := strings.Split(ref, ":")[1]
			i := int(ref[len(ref)-1] - '0')
			source := Source{Identity: Identity{Provider: provider, SourceInstance: "one", SourceRef: ref}, LastInteractionAt: fmt.Sprintf("2026-09-20T%02d:00:00Z", i), Relationships: []Relationship{}, Metadata: []Metadata{}}
			return response("show", "complete", map[string]any{"source": source}, nil, []any{}), nil
		case "verify":
			return response("verify", "complete", map[string]any{"source_ref": args[len(args)-1], "verified_version": Version{Algorithm: "sha256", Basis: "provider-content-v0", Value: "sha256:" + strings.Repeat("a", 64)}}, nil, []any{}), nil
		case "events":
			return response("events", "complete", map[string]any{"source_ref": args[len(args)-1], "events": []any{map[string]any{"index": 0, "kind": "message", "time_state": "available", "recorded_at": "2026-09-20T01:00:00Z", "message": map[string]any{"role": "user", "text": "synthetic-private-token", "truncated": false}, "metadata": []any{}}}}, map[string]any{"limit": 100, "has_more": false}, []any{}), nil
		}
		return nil, fmt.Errorf("unexpected operation")
	}
}

func response(operation, status string, data any, page any, omissions []any) []byte {
	envelope := map[string]any{"schema_version": "v2", "cli_version": "dev", "operation": operation, "status": status, "data": data, "omissions": omissions}
	if page != nil {
		envelope["page"] = page
	}
	raw, _ := json.Marshal(envelope)
	return raw
}

func syntheticEvent(index int) map[string]any {
	return map[string]any{"index": index, "kind": "message", "time_state": "available", "recorded_at": "2026-09-20T01:00:00Z", "message": map[string]any{"role": "user", "text": "synthetic-private-token", "truncated": false}, "metadata": []any{}}
}

func argumentAfter(args []string, name string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == name {
			return args[i+1]
		}
	}
	return ""
}
