package pilot

import (
	"encoding/json"
	"errors"
	"sort"
	"time"
)

const (
	SchemaVersion  = "v2"
	MaxCandidates  = 10
	MinPerProvider = 3
	MaxEvents      = 2000
	PageLimit      = 100
)

type Identity struct {
	Provider                  string `json:"provider"`
	SourceInstance            string `json:"source_instance"`
	ProviderSourceFingerprint string `json:"provider_source_fingerprint"`
	SourceRef                 string `json:"source_ref"`
}

type Relationship struct {
	Kind      string `json:"kind"`
	SourceRef string `json:"source_ref"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Source struct {
	Identity          Identity       `json:"identity"`
	Kind              string         `json:"kind"`
	LastInteractionAt string         `json:"last_interaction_at"`
	Relationships     []Relationship `json:"relationships"`
	Metadata          []Metadata     `json:"metadata"`
}

type Version struct {
	Algorithm string `json:"algorithm"`
	Basis     string `json:"basis"`
	Value     string `json:"value"`
}

type Page struct {
	Limit      int    `json:"limit"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type Omission struct {
	Code  string `json:"code"`
	Scope string `json:"scope"`
}

type Envelope struct {
	SchemaVersion string `json:"schema_version"`
	CLIVersion    string `json:"cli_version"`
	Operation     string `json:"operation"`
	Status        string `json:"status"`
	Data          struct {
		Sources         []Source          `json:"sources"`
		Source          *Source           `json:"source"`
		SourceRef       string            `json:"source_ref"`
		Events          []json.RawMessage `json:"events"`
		VerifiedVersion *Version          `json:"verified_version"`
	} `json:"data"`
	Page      *Page      `json:"page"`
	Omissions []Omission `json:"omissions"`
}

type SelectedSource struct {
	Source          SourceSummary  `json:"source"`
	VerifiedVersion Version        `json:"verified_version"`
	Events          []EventSummary `json:"events"`
}

// SourceSummary excludes provider-native IDs and free-form metadata by construction.
type SourceSummary struct {
	Identity          Identity       `json:"identity"`
	Kind              string         `json:"kind"`
	LastInteractionAt string         `json:"last_interaction_at"`
	Relationships     []Relationship `json:"relationships"`
}

// EventSummary is deliberately content-free; v2 recorded bodies have no redaction guarantee.
type EventSummary struct {
	Index            uint64 `json:"index"`
	Kind             string `json:"kind"`
	TimeState        string `json:"time_state"`
	RecordedAt       string `json:"recorded_at,omitempty"`
	Role             string `json:"role,omitempty"`
	InputState       string `json:"input_state,omitempty"`
	ContentState     string `json:"content_state,omitempty"`
	Outcome          string `json:"outcome,omitempty"`
	CorrelationState string `json:"correlation_state,omitempty"`
}

type Lineage struct {
	Sources []SelectedSource `json:"sources"`
}

type ProviderCount struct {
	Candidates int `json:"candidates"`
	Adopted    int `json:"adopted"`
	Excluded   int `json:"excluded"`
}

type Bundle struct {
	AgentSessionsCommit string                   `json:"agent_sessions_commit"`
	SchemaVersion       string                   `json:"schema_version"`
	StartedAt           string                   `json:"started_at"`
	Cutoff              string                   `json:"cutoff"`
	Status              string                   `json:"status"`
	Reason              string                   `json:"reason,omitempty"`
	Counts              map[string]ProviderCount `json:"counts,omitempty"`
	EventCount          int                      `json:"event_count"`
	Lineages            []Lineage                `json:"lineages"`
}

type candidate struct {
	Sources   []Source
	latest    time.Time
	key       string
	ambiguous bool
}

func selectCandidates(sources []Source, started time.Time, provider string, maximum int) ([]candidate, error) {
	cutoff := started.AddDate(0, 0, -14)
	parent := map[string]string{}
	var eligible []Source
	for _, source := range sources {
		if source.Identity.Provider != provider || source.Identity.SourceRef == "" {
			return nil, errors.New("invalid source identity")
		}
		at, err := parseCanonicalTime(source.LastInteractionAt)
		if err != nil {
			return nil, errors.New("invalid source time")
		}
		if !at.Before(cutoff) && at.Before(started) {
			eligible = append(eligible, source)
			if _, exists := parent[source.Identity.SourceRef]; exists {
				return nil, errors.New("duplicate source")
			}
			parent[source.Identity.SourceRef] = source.Identity.SourceRef
		}
	}
	var find func(string) string
	find = func(ref string) string {
		root := parent[ref]
		if root != ref {
			parent[ref] = find(root)
		}
		return parent[ref]
	}
	union := func(a, b string) {
		if _, ok := parent[b]; !ok {
			parent[b] = b
		}
		ra, rb := find(a), find(b)
		if ra != rb {
			if ra < rb {
				parent[rb] = ra
			} else {
				parent[ra] = rb
			}
		}
	}
	ambiguous := map[string]bool{}
	for _, source := range eligible {
		for _, relation := range source.Relationships {
			if relation.Kind != "parent" && relation.Kind != "forked_from" || relation.SourceRef == "" {
				ambiguous[source.Identity.SourceRef] = true
				continue
			}
			union(source.Identity.SourceRef, relation.SourceRef)
		}
	}
	groups := map[string]*candidate{}
	for _, source := range eligible {
		root := find(source.Identity.SourceRef)
		group := groups[root]
		if group == nil {
			group = &candidate{key: source.Identity.SourceRef}
			groups[root] = group
		}
		group.Sources = append(group.Sources, source)
		group.ambiguous = group.ambiguous || ambiguous[source.Identity.SourceRef]
		at, _ := time.Parse(time.RFC3339Nano, source.LastInteractionAt)
		if at.After(group.latest) {
			group.latest = at
		}
		if source.Identity.SourceRef < group.key {
			group.key = source.Identity.SourceRef
		}
	}
	ordered := make([]candidate, 0, len(groups))
	for _, group := range groups {
		sort.Slice(group.Sources, func(i, j int) bool { return group.Sources[i].Identity.SourceRef < group.Sources[j].Identity.SourceRef })
		ordered = append(ordered, *group)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].latest.Equal(ordered[j].latest) {
			return ordered[i].key < ordered[j].key
		}
		return ordered[i].latest.After(ordered[j].latest)
	})
	if len(ordered) > maximum {
		ordered = ordered[:maximum]
	}
	return ordered, nil
}

func parseCanonicalTime(value string) (time.Time, error) {
	at, err := time.Parse(time.RFC3339Nano, value)
	if err != nil || at.UTC().Format(time.RFC3339Nano) != value || value == "" {
		return time.Time{}, errors.New("invalid canonical time")
	}
	return at, nil
}
