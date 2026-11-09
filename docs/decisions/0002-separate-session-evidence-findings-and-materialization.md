# ADR 0002: Separate session evidence, findings, and materialization

- Status: Accepted
- Date: 2026-09-01

## Context

Session records, retrospective state, and actual harness artifacts have different owners and lifecycles.\
Combining them would create four problems:

- **Session access**: A session-access database would have to represent retrospective concepts that other consumers do not need.
- **Configuration authority**: A finding database could become a second source of truth for actual harness configuration.
- **Automatic changes**: Uncertain interpretations could become permanent instructions or enforcement without a separate decision.
- **Review workload**: Continuous event-level persistence could create a large, noisy queue of low-value candidate findings.

## Decision

Keep session evidence, retrospective state, and actual outer-harness artifacts separate.

### Session evidence

- **Ownership**: Coding-agent providers own the session records.
- **Access**: Obtain evidence through the public machine-readable contract of `agent-sessions`.
- **Copies**: `agent-harness-retro` does not own or duplicate complete provider transcripts.

### Retrospective state

`agent-harness-retro` owns private state needed to carry information between reviews.

- **Possible records**

  - Reviewed source versions.
  - Findings.
  - Evidence relationships.
  - Evidence references.
  - Proposals.
  - Accepted or rejected decisions.
  - Verification state.

- **Format**: This ADR does not fix the persistent state format.

- **Privacy**: Keep state private and separate from the public repository.

- **Writes**: Write persistent state only as part of deliberate retrospective or decision operations.

- **Default execution**

  The default architecture does not use:

  - Continuous per-event candidate generation.
  - Background synchronization.
  - Periodic analysis.

### Materialized outer harness

Instructions, Skills, Hooks, settings, CI, and other harness artifacts retain their existing sources of truth.

- **Proposals**: A proposal does not modify those artifacts automatically.
- **Approval**: If proposal application is implemented, require explicit approval before applying a proposal.
- **Target requirements**: Any implemented proposal application must respect the target's authority, source-of-truth model, permissions, and validation process.
- **Initial workflow**: The initial retrospective workflow may stop at proposal generation.

## Evidence retention

- **Source evidence**: Findings should reference source evidence rather than store complete transcripts.
- **Provenance**: Where practical, retained references should include a logical source reference and verified source version.
- **Summaries**: A bounded derived summary may be stored for usability, but it does not replace source provenance.
- **Public improvements**: When a private finding leads to a permanent public harness improvement, prefer a synthetic regression case that expresses the general failure instead of embedding private session content.

## Consequences

### Positive

- **Reuse**: `agent-sessions` remains reusable for session access.
- **Independent state**: Retrospective state can evolve without changing provider access.
- **Configuration authority**: Actual harness artifacts retain one clear source of truth.
- **Proposal-only use**: Analysis can remain useful even if automatic writes are never implemented.
- **Uncertain findings**: Uncertain findings do not automatically become permanent instructions.
- **Privacy**: Private evidence can stay private while generalized regression tests become public.
- **Review workload**: Low-value candidates do not accumulate continuously by default.

### Negative

- **State maintenance**: Retrospective state needs its own storage and lifecycle.
- **Divergence**: A proposal and the actual target may differ until the proposal is explicitly applied.
- **Evidence availability**: Evidence may become unavailable if the provider deletes the original session.
- **Application process**: Applying a proposal requires an additional approval and validation step.

These costs preserve clear ownership and safer change boundaries.

## Alternatives considered

| Alternative | Decision | Reason |
| --- | --- | --- |
| Store retrospective state in `agent-sessions`. | Rejected. | Consumer-specific processing and findings would reduce the reuse value of session access. |
| Copy complete session transcripts into retrospective storage. | Rejected. | Privacy, retention, storage, and synchronization responsibilities would exceed what evidence provenance requires. |
| Automatically modify the harness after every retrospective. | Rejected. | Conclusions may be uncertain, have the wrong scope, or be better addressed through another control surface. |
| Continuously store one candidate finding per event or tool failure. | Rejected as the default. | High-frequency generation creates noise and review debt and encourages frequency to replace judgment. |

## Revisit conditions

Revisit this decision if:

- **Workflow limits**: A real workflow shows that proposal-only retrospectives cannot remain usable.
- **Evidence retention**: Adequate evidence cannot be retained without bounded source snapshots.
- **State contract**: The state model needs a more formal persistence contract.
- **Proposal application**: A reliable workflow with explicit approval has been proven useful.
