# AGENTS.md

This file defines repository-specific instructions for agents working on `agent-harness-retro`.

## Purpose

`agent-harness-retro` retrospectively evaluates user-controlled outer harnesses using coding-agent session evidence.

It owns retrospective concepts, findings, evidence relationships, proposals, and the boundary between analysis and optional harness modification.

It does not own provider-specific raw session access.\
Use `agent-sessions` for that responsibility.

## Read first

Before making a material design or implementation change, read:

1. `README.md`
2. `docs/architecture.md`
3. `docs/glossary.md`
4. any ADR directly relevant to the change

Do not scan unrelated files merely to increase coverage.

## Public information boundary

This is a public repository.

Do not add:

- private repository names or links;
- private Issue or Project references;
- personal or employer-specific context;
- real coding-agent session contents;
- real outer-harness contents from private environments;
- credentials or authentication material;
- real usernames, hostnames, home directories, or machine-specific paths.

A requirement learned from private usage may be incorporated only after expressing it as a general product requirement that does not depend on the private source.

Tests and examples must use synthetic repositories, sessions, harness artifacts, and identities.

## Core architecture constraints

Preserve the following unless an explicit architecture decision supersedes them:

- Provider-specific session access belongs to `agent-sessions`.
- Raw provider transcripts are not copied into this repository or retrospective state.
- The retrospective target is the user-controlled outer harness, not vendor-controlled built-in harness internals.
- Global and repository-local harness scopes use the same core retrospective model.
- Session evidence, retrospective findings, and materialized harness artifacts remain separate layers.
- A finding does not imply that a change is required.
- Frequency alone does not justify permanent harness guidance.
- Analysis and materialization remain separable.
- External harness artifacts are not modified without explicit approval.
- Continuous per-event finding generation is not the default architecture.
- Private runtime state is not stored in the public repository.

Do not add background monitoring, automatic transcript ingestion, automatic harness rewrites, or a general-purpose session database merely because they may become useful later.

## Evidence

A material finding must remain traceable to evidence.

Prefer stable references such as:

- `agent-sessions` source references;
- verified source versions;
- bounded event references;
- affected scope;
- current harness artifacts relevant to the finding.

Do not preserve complete raw transcripts merely to make a finding auditable.

Do not treat a generated summary as stronger evidence than the underlying observed source.

Distinguish:

- directly observed behavior;
- interpretation;
- hypothesis;
- proposal;
- unverified assumption.

## Harness scope

Before proposing a harness change, determine the appropriate scope.

Possible outcomes include:

- global outer harness;
- repository harness;
- one existing Skill;
- one Hook or deterministic check;
- underlying project code or documentation;
- provider-controlled built-in behavior;
- no persistent change.

Do not promote repository-local behavior into a global rule without cross-context evidence.

Do not promote one-off user correction into a permanent harness rule without a reason that survives the original session.

## Control-surface selection

Do not default to adding instructions.

When proposing a fix, compare the relevant control surfaces.

For example:

- durable behavioral guidance → instruction or rule;
- conditional multi-step workflow → Skill;
- mechanically enforceable intervention → Hook or deterministic check;
- repository correctness invariant → test or CI;
- user-controlled execution behavior → setting or permission;
- source-of-truth project knowledge → project documentation.

Prefer modifying an existing appropriate mechanism over creating a parallel one.

## Findings and proposals

Keep findings separate from proposed solutions.

A finding should state the observed pattern and evidence before prescribing a harness artifact.

A proposal should identify:

- target scope;
- target artifact or control-surface type;
- expected improvement;
- relevant tradeoffs;
- evidence supporting the proposal;
- validation or evaluation needed after the change.

Rejection and no-change are valid outcomes.

## Materialization

- Do not modify an external harness artifact merely because the retrospective found a candidate improvement.
- If materialization support exists, require explicit approval before write operations.
- Preserve the target repository's own authority, instructions, review requirements, and validation commands.
- Do not bypass a target repository's source-of-truth model.

## Runtime state

Retrospective state is private runtime data.

Do not commit:

- reviewed-session ledgers;
- findings derived from private sessions;
- proposal decisions;
- user-specific harness inventories;
- private evidence.

The persistent representation and storage format must remain as small as the product requirements allow.

Do not introduce high-frequency writes, background synchronization, or a database without an observed need.

## Dependency boundary

Consume coding-agent session data through the public machine-readable contract of `agent-sessions`.

Do not add direct Codex, Claude Code, or other provider transcript parsers to this repository to work around missing upstream functionality without first determining whether the required capability belongs in `agent-sessions`.

## Language

English is the canonical language for public repository documents and rules.
Maintained Japanese files are reference translations only; they do not define independent requirements or authority.
If a Japanese reference differs from its English canonical source, the English source takes precedence.

The maintained file pairs and synchronization rules are defined in `docs/localization.md`.
When adding or changing a maintained English canonical file, use the repository-local `.agents/skills/maintain-japanese-references` Skill to create or review its Japanese reference in the same change.
If the English edit does not affect Japanese meaning, leave the reference unchanged and record the reason in the pull request's Validation or Risks / Follow-up section.
Do not translate ADRs as part of this maintained set.

English remains the default for:

- code and identifiers;
- code comments;
- CLI or Skill public interfaces;
- commit message summaries;
- release notes;
- canonical README, documentation, glossary, ADRs, Skills, and repository policy files.

Maintainer-created public Issues and pull requests use an English title and a short English `Summary`.
Their other body sections and comments are Japanese by default and may be written in English when useful.
Do not require external contributors to use Japanese, and do not require a full English and Japanese duplication of the body.

`AGENTS-ja.md` is a reference translation of this file, not an independent instruction source.

## Change discipline

Prefer the smallest mechanism that establishes the required retrospective capability.

Do not build a taxonomy, scoring model, persistence system, orchestration framework, or automatic materialization layer before a concrete requirement needs it.

When changing an architectural invariant:

1. identify the affected ADR;
2. determine whether the accepted decision still holds;
3. update or supersede the ADR when needed;
4. keep README, glossary, architecture, implementation, and tests aligned.

Do not add implementation history to current-state product documentation unless it is necessary to explain an active decision.

## Validation

Use synthetic evidence and harness fixtures.

Never use real private sessions as committed evaluation data, even after manual redaction.

When a change affects retrospective quality, test the decision boundary that can fail rather than expanding an unrelated evaluation suite.

Always inspect the final diff for accidental private information and machine-specific data before reporting completion.
