# ADR 0002: Separate session evidence, findings, and materialization

- Status: Accepted
- Date: 2026-09-01

## Context

Retrospective improvement involves several kinds of data with different owners and lifecycles.

Coding-agent session records are owned by their providers and accessed through `agent-sessions`.

Retrospective analysis may need durable private state such as:

- source versions already reviewed;
- findings;
- evidence relationships;
- proposal decisions;
- verification state.

Actual outer-harness artifacts have separate sources of truth, often in user configuration or repositories.

Combining these layers would create several problems.

- A session-access database would gain retrospective semantics that unrelated consumers do not need.
- A finding database could become a shadow source of truth for actual harness configuration.
- Automatic materialization could turn uncertain retrospective interpretation directly into persistent instructions or enforcement.
- Continuous event-level persistence could also create a large, noisy queue of low-value candidate findings.

## Decision

`agent-harness-retro` will keep three layers separate.

### Session evidence

Session evidence is obtained through the public machine-readable contract of `agent-sessions`.\
`agent-harness-retro` does not own or duplicate complete provider transcripts.

### Retrospective state

`agent-harness-retro` owns private state required to support retrospective continuity.

This may include:

- reviewed source versions;
- findings;
- evidence references;
- proposals;
- accepted or rejected decisions;
- verification state.

The state format is not fixed by this ADR.\
It must remain private and separate from the public repository.

Persistent writes occur only as part of deliberate retrospective or decision operations.\
The default architecture does not use continuous per-event candidate generation, background synchronization, or periodic analysis.

### Materialized outer harness

Actual instructions, Skills, Hooks, settings, CI, and other harness artifacts remain owned by their existing source of truth.

A proposal does not modify them automatically.

Materialization, if implemented, requires explicit approval and must respect the target's own authority, source-of-truth model, permissions, and validation process.

The initial retrospective workflow may stop at proposal generation.

## Evidence retention

Findings should reference source evidence rather than persist complete transcripts.

Where practical, retained references should include a logical source reference and verified source version.

A bounded derived summary may be stored for usability, but generated summaries do not replace source provenance.

When a private finding becomes a permanent public harness improvement, prefer creating a synthetic regression case that expresses the general failure instead of embedding private session content.

## Consequences

### Positive

- `agent-sessions` remains a reusable session-access substrate.
- Retrospective state can evolve without changing provider access.
- Actual harness artifacts retain one clear source of truth.
- Analysis can remain useful even if automatic writes are never implemented.
- Uncertain findings do not become permanent instructions automatically.
- Private evidence can remain private while generalized regression tests become public.
- Continuous low-value candidate accumulation is avoided by default.

### Negative

- Retrospective state needs its own storage and lifecycle.
- A proposal and the materialized target may diverge until explicitly applied.
- Some evidence may become unavailable if the provider deletes the original session.
- Materialization requires an additional approval and validation step.

These costs preserve clearer ownership and safer change boundaries.

## Alternatives considered

### Store retrospective state in `agent-sessions`

Rejected.

Processing and findings are consumer-specific semantics and would reduce the reuse value of the session-access layer.

### Copy complete session transcripts into retrospective storage

Rejected.

It increases privacy, retention, storage, and synchronization responsibilities beyond what retrospective provenance requires.

### Automatically modify the harness after every retrospective

Rejected.

Retrospective conclusions may be uncertain, incorrectly scoped, or better solved through a different control surface.

### Continuously persist one candidate finding per event or tool failure

Rejected as the default.

High-frequency candidate generation creates noise and review debt and encourages frequency to substitute for judgment.

## Revisit conditions

Revisit this decision if:

- a real workflow demonstrates that proposal-only retrospective cannot remain usable;
- evidence retention cannot remain adequate without bounded source snapshots;
- the state model requires a more formal persistence contract;
- a reliable, explicitly authorized materialization workflow has been proven useful.
