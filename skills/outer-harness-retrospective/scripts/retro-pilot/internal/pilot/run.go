package pilot

import (
	_ "embed"
	"encoding/hex"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed agent-sessions.commit
var agentSessionsCommit string

var AgentSessionsCommit = strings.TrimSpace(agentSessionsCommit)

type Caller func(args ...string) ([]byte, error)

func Run(call Caller, started time.Time) (Bundle, error) {
	started = started.UTC()
	bundle := Bundle{
		AgentSessionsCommit: AgentSessionsCommit, SchemaVersion: SchemaVersion,
		StartedAt: started.Format(time.RFC3339Nano), Cutoff: started.AddDate(0, 0, -14).Format(time.RFC3339Nano),
		Counts: map[string]ProviderCount{"codex": {}, "claude": {}}, Lineages: []Lineage{},
	}
	candidates := map[string][]candidate{}
	for _, provider := range []string{"codex", "claude"} {
		sources, err := listAll(call, provider)
		if err != nil {
			bundle.Counts = nil
			return insufficient(bundle, "list_incomplete"), nil
		}
		groups, err := selectCandidates(sources, started, provider, MaxCandidates)
		if err != nil {
			bundle.Counts = nil
			return insufficient(bundle, "cohort_ambiguous"), nil
		}
		instances := map[string]bool{}
		for _, group := range groups {
			for _, source := range group.Sources {
				instances[source.Identity.SourceInstance] = true
			}
		}
		if len(instances) > 1 {
			bundle.Counts = nil
			return insufficient(bundle, "multiple_source_instances"), nil
		}
		candidates[provider] = groups
		count := bundle.Counts[provider]
		count.Candidates = len(groups)
		bundle.Counts[provider] = count
	}
	remaining := map[string][]candidate{}
	for _, provider := range []string{"codex", "claude"} {
		groups := candidates[provider]
		i := 0
		for i < len(groups) && bundle.Counts[provider].Adopted < MinPerProvider {
			if groups[i].ambiguous {
				i++
				count := bundle.Counts[provider]
				count.Excluded++
				bundle.Counts[provider] = count
				continue
			}
			lineage, size, overflow, err := fetchCandidate(call, groups[i], MaxEvents-bundle.EventCount)
			i++
			if overflow {
				return insufficient(bundle, "minimum_exceeds_event_budget"), nil
			}
			if err != nil {
				count := bundle.Counts[provider]
				count.Excluded++
				bundle.Counts[provider] = count
				continue
			}
			bundle.Lineages = append(bundle.Lineages, lineage)
			bundle.EventCount += size
			count := bundle.Counts[provider]
			count.Adopted++
			bundle.Counts[provider] = count
		}
		if bundle.Counts[provider].Adopted < MinPerProvider {
			return insufficient(bundle, "minimum_not_met"), nil
		}
		remaining[provider] = groups[i:]
	}
	type next struct {
		provider string
		group    candidate
	}
	var rest []next
	for _, provider := range []string{"codex", "claude"} {
		for _, group := range remaining[provider] {
			rest = append(rest, next{provider, group})
		}
	}
	sort.Slice(rest, func(i, j int) bool {
		if rest[i].group.latest.Equal(rest[j].group.latest) {
			return rest[i].group.key < rest[j].group.key
		}
		return rest[i].group.latest.After(rest[j].group.latest)
	})
	for _, item := range rest {
		if bundle.EventCount == MaxEvents {
			break
		}
		if item.group.ambiguous {
			count := bundle.Counts[item.provider]
			count.Excluded++
			bundle.Counts[item.provider] = count
			continue
		}
		lineage, size, overflow, err := fetchCandidate(call, item.group, MaxEvents-bundle.EventCount)
		if overflow {
			break
		}
		count := bundle.Counts[item.provider]
		if err != nil {
			count.Excluded++
		} else {
			bundle.Lineages = append(bundle.Lineages, lineage)
			bundle.EventCount += size
			count.Adopted++
		}
		bundle.Counts[item.provider] = count
	}
	bundle.Status = "ready"
	return bundle, nil
}

func insufficient(bundle Bundle, reason string) Bundle {
	bundle.Status = "insufficient_evidence"
	bundle.Reason = reason
	return bundle
}

func callEnvelope(call Caller, operation string, args ...string) (Envelope, error) {
	raw, err := call(append([]string{operation}, args...)...)
	if err != nil {
		return Envelope{}, errors.New("agent-sessions call failed")
	}
	return validateEnvelope(operation, raw)
}

func listAll(call Caller, provider string) ([]Source, error) {
	var sources []Source
	cursor := ""
	seenCursor := map[string]bool{}
	seenSource := map[string]bool{}
	for {
		args := []string{"--provider", provider, "--limit", strconv.Itoa(PageLimit)}
		if cursor != "" {
			args = append(args, "--cursor", cursor)
		}
		envelope, err := callEnvelope(call, "list", args...)
		if err != nil {
			return nil, err
		}
		if err := completeOrPagination(envelope); err != nil {
			return nil, err
		}
		if len(envelope.Data.Sources) == 0 && envelope.Page.HasMore {
			return nil, errors.New("empty nonterminal list page")
		}
		for _, source := range envelope.Data.Sources {
			if source.Identity.Provider != provider || source.Identity.SourceInstance == "" || source.Identity.SourceRef == "" || source.LastInteractionAt == "" || seenSource[source.Identity.SourceRef] {
				return nil, errors.New("invalid list source")
			}
			seenSource[source.Identity.SourceRef] = true
			sources = append(sources, source)
		}
		if !envelope.Page.HasMore {
			return sources, nil
		}
		cursor = envelope.Page.NextCursor
		if seenCursor[cursor] {
			return nil, errors.New("repeated list cursor")
		}
		seenCursor[cursor] = true
	}
}

func fetchCandidate(call Caller, group candidate, remaining int) (Lineage, int, bool, error) {
	lineage := Lineage{Sources: []SelectedSource{}}
	size := 0
	for _, listed := range group.Sources {
		ref := listed.Identity.SourceRef
		shown, err := callEnvelope(call, "show", ref)
		if err != nil || shown.Status != "complete" || shown.Data.Source == nil || !sameSource(listed, *shown.Data.Source) {
			return Lineage{}, 0, false, errors.New("source changed or incomplete")
		}
		before, err := verify(call, ref)
		if err != nil {
			return Lineage{}, 0, false, err
		}
		events := []EventSummary{}
		cursor := ""
		seenCursor := map[string]bool{}
		for {
			limit := remaining + 1 - size
			if limit > PageLimit {
				limit = PageLimit
			}
			if limit < 1 {
				return Lineage{}, 0, true, nil
			}
			args := []string{"--limit", strconv.Itoa(limit)}
			if cursor != "" {
				args = append(args, "--cursor", cursor)
			}
			args = append(args, ref)
			page, err := callEnvelope(call, "events", args...)
			if err != nil || page.Data.SourceRef != ref {
				return Lineage{}, 0, false, errors.New("events unavailable")
			}
			if err := completeOrPagination(page); err != nil {
				return Lineage{}, 0, false, err
			}
			if len(page.Data.Events) > limit || len(page.Data.Events) == 0 && page.Page.HasMore {
				return Lineage{}, 0, false, errors.New("invalid events page")
			}
			for _, event := range page.Data.Events {
				summary, err := summarizeEvent(event, uint64(len(events)))
				if err != nil {
					return Lineage{}, 0, false, err
				}
				events = append(events, summary)
				size++
			}
			if size > remaining {
				return Lineage{}, 0, true, nil
			}
			if !page.Page.HasMore {
				break
			}
			cursor = page.Page.NextCursor
			if seenCursor[cursor] {
				return Lineage{}, 0, false, errors.New("repeated events cursor")
			}
			seenCursor[cursor] = true
		}
		after, err := verify(call, ref)
		if err != nil || before != after {
			return Lineage{}, 0, false, errors.New("verified version changed")
		}
		publicSource := SourceSummary{
			Identity: shown.Data.Source.Identity, Kind: shown.Data.Source.Kind,
			LastInteractionAt: shown.Data.Source.LastInteractionAt,
			Relationships:     shown.Data.Source.Relationships,
		}
		lineage.Sources = append(lineage.Sources, SelectedSource{Source: publicSource, VerifiedVersion: before, Events: events})
	}
	return lineage, size, false, nil
}

func verify(call Caller, ref string) (Version, error) {
	envelope, err := callEnvelope(call, "verify", ref)
	if err != nil || envelope.Status != "complete" || envelope.Data.SourceRef != ref || envelope.Data.VerifiedVersion == nil {
		return Version{}, errors.New("verification unavailable")
	}
	version := *envelope.Data.VerifiedVersion
	if version.Algorithm != "sha256" || version.Basis != "provider-content-v0" || !strings.HasPrefix(version.Value, "sha256:") {
		return Version{}, errors.New("unsupported verified version")
	}
	digest, err := hex.DecodeString(strings.TrimPrefix(version.Value, "sha256:"))
	if err != nil || len(digest) != 32 {
		return Version{}, errors.New("invalid verified version")
	}
	return version, nil
}

func sameSource(a, b Source) bool {
	return a.Identity == b.Identity && a.LastInteractionAt == b.LastInteractionAt && reflect.DeepEqual(a.Relationships, b.Relationships)
}

func validEventKind(kind string) bool {
	switch kind {
	case "message", "tool_call", "tool_result", "error":
		return true
	default:
		return false
	}
}
