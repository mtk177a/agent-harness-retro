# Architecture

## Purpose

`agent-harness-retro` retrospectively evaluates user-controlled coding-agent harnesses using normalized session evidence.

Its responsibilities are limited to:

- selecting relevant session evidence;
- observing recurring or material behavior;
- relating behavior to the current outer harness;
- recording findings;
- determining likely scope and control surface;
- producing evidence-backed improvement proposals;
- tracking retrospective decisions when durable state is required.

Provider-specific session acquisition is outside this boundary.

## System boundary

```text
Provider-owned sessions
        │
        ▼
   agent-sessions
        │
        │ normalized evidence
        ▼
agent-harness-retro
──────────────────────────────────

evidence selection
current harness observation
retrospective analysis
findings
scope resolution
proposal generation
decision state

        │
        ▼

User-controlled outer harness
──────────────────────────────────

global outer harness
repository harness
```

| Actor                 | Owns                                                          |
| --------------------- | ------------------------------------------------------------- |
| Coding-agent provider | Raw session records and built-in harness                      |
| `agent-sessions`      | Provider-neutral read-only session access                     |
| `agent-harness-retro` | Retrospective findings, evidence relationships, and proposals |
| Harness owner         | Actual outer-harness artifacts and final change decisions     |

## Target boundary

### Built-in harness

The vendor-controlled system that turns a model into a coding-agent product.

Examples may include provider-owned:

- system prompts;
- context management;
- internal orchestration;
- tool runtime;
- product-internal memory;
- internal routing.

`agent-harness-retro` may observe behavior that appears related to the built-in harness, but it does not modify that layer.\
A built-in limitation may be reported as such rather than converted into an inappropriate outer-harness workaround.

### Outer harness

The user-controlled layer surrounding the built-in coding-agent harness.\
This is the primary target of `agent-harness-retro`.

### Global outer harness

User-controlled artifacts intended to affect multiple repositories or projects.

Possible artifacts include:

- global instructions;
- global Skills;
- user Hooks;
- user-controlled settings and permissions;
- reusable agent definitions;
- shared MCP or plugin configuration.

### Repository harness

Agent-facing control surfaces owned by one repository or project.

Possible artifacts include:

- repository instructions;
- repository-local Skills;
- repository Hooks and settings;
- agent-facing project documentation;
- deterministic repository commands;
- tests and CI;
- evaluation assets.

The same retrospective model applies to global and repository-local scopes.

## Input contract

`agent-harness-retro` consumes session evidence through `agent-sessions`.

It must not depend directly on provider-native Codex, Claude Code, or other raw transcript formats.

The required evidence contract is expected to include concepts equivalent to:

- stable source reference;
- provider;
- source instance;
- source version;
- session relationship;
- ordered messages and tool events;
- completeness and omission information.

The exact upstream schema is owned by `agent-sessions`.

## Current harness observation

A retrospective is not performed against session evidence alone.

The analyzer also needs enough current harness context to determine whether:

- relevant guidance already exists;
- a Skill should have triggered;
- a Hook or deterministic check already covers the failure;
- a finding has already been addressed;
- a proposed change would duplicate another source of truth.

Harness observation should remain target-scoped.\
Do not load every global and repository artifact merely because it exists.\
Read only the artifacts relevant to a candidate finding or the requested retrospective scope.

## Retrospective pipeline

The conceptual pipeline is:

```text
select evidence
      │
      ▼
observe behavior
      │
      ▼
identify finding
      │
      ▼
check current harness
      │
      ▼
resolve likely cause and scope
      │
      ▼
compare control surfaces
      │
      ▼
produce proposal or no-change decision
```

The pipeline is conceptual rather than a requirement for one implementation process.

## Findings

A finding records a material observation supported by session evidence.

Possible examples include:

- recurring user correction;
- repeated unnecessary retry;
- scope drift;
- verification gap;
- ineffective instruction;
- Skill routing or trigger failure;
- inappropriate delegation;
- recurring successful workflow;
- duplicated or conflicting harness guidance.

The initial implementation should not freeze a broad universal taxonomy before real retrospective use demonstrates that one is necessary.

### Finding requirements

A material finding should identify:

- what was observed;
- evidence supporting the observation;
- affected scope;
- relevant current harness context;
- uncertainty or alternative explanations when material.

A finding does not prescribe a solution.

## Root cause and scope

Repeated behavior can originate from different layers.

For example:

```text
same observed failure
        │
        ├── built-in harness limitation
        ├── global outer-harness guidance
        ├── repository harness
        ├── one Skill or Hook
        ├── project implementation
        └── one-off session context
```

The retrospective must not convert every failure into an instruction change.\
Scope resolution is a core responsibility.

## Control surfaces

A proposal selects a control surface only after the finding and scope have been established.

Potential control surfaces include:

| Need                                    | Typical control surface  |
| --------------------------------------- | ------------------------ |
| Durable behavioral guidance             | Instruction or rule      |
| Conditional judgment workflow           | Skill                    |
| Mechanical intervention                 | Hook                     |
| Repository invariant                    | Test, checker, or CI     |
| User-controlled execution behavior      | Setting or permission    |
| Canonical project knowledge             | Repository documentation |
| Existing behavior is already sufficient | No change                |

This table is guidance rather than a rigid mapping.

## Proposals

A proposal is a candidate response to one or more findings.

A proposal should state:

- target scope;
- target artifact or control-surface type;
- expected effect;
- evidence;
- known tradeoffs;
- validation or evaluation required after materialization.

The system must support a no-change outcome.

Repeated behavior alone is not enough to justify a proposal.

## Evidence model

A finding should reference source evidence rather than duplicate complete transcripts.

Where practical, evidence should use:

- source reference;
- verified source version;
- bounded event reference.

A derived summary may be retained for usability but does not replace the underlying evidence relationship.

If provider-owned evidence is later unavailable, retained evidence should remain the minimum needed to understand the accepted finding or regression contract.

## Private retrospective state

Unlike `agent-sessions`, `agent-harness-retro` may require persistent consumer-owned state.

Potential state includes:

- source versions already reviewed for retrospective purposes;
- findings;
- evidence relationships;
- proposal decisions;
- materialization state;
- verification state.

This state is:

- private;
- user-local by default;
- separate from the public repository;
- separate from provider-owned transcripts;
- separate from `agent-sessions`.

The initial architecture does not prescribe SQLite or another specific storage format.\
A storage implementation should be chosen only when the required query and update patterns are known.

### Write behavior

Runtime state is updated only through deliberate retrospective or decision operations.

The architecture does not require:

- background synchronization;
- continuous event ingestion;
- a file watcher;
- per-tool-call candidate persistence;
- periodic session analysis.

Idle write activity should remain zero.

## Analysis cadence

The default model is explicit retrospective analysis.

A retrospective may target:

- selected sessions;
- a repository and time range;
- new source versions since the previous retrospective;
- evidence related to one harness artifact;
- another deliberately bounded scope.

The system does not continuously create a candidate finding for every observed event.\
This avoids converting normal agent activity into an ever-growing review queue.

## Materialization boundary

Retrospective analysis and outer-harness modification are separate layers.

```text
finding
   │
   ▼
proposal
   │
   ▼
explicit approval
   │
   ▼
materialization
   │
   ▼
validation / evaluation
```

The initial system may stop at proposal generation.

If materialization support is later implemented:

- explicit user approval is required;
- the target repository remains authoritative;
- target repository instructions and permissions apply;
- writes use the target's real source of truth rather than generated deployment output;
- the resulting change is validated through the target's normal mechanisms.

Materialization must not silently rewrite a global or repository harness.

## Evaluation after change

A retrospective-derived change should be verified at the layer where the failure can recur.

Examples include:

- Skill routing failure → Skill trigger or behavior evaluation;
- missing deterministic validation → checker or CI test;
- instruction failure → representative behavioral scenario;
- Hook change → deterministic Hook behavior test.

When practical, a private session finding should be converted into a synthetic regression case rather than embedding private evidence into a public repository.

## Public/private boundary

This public repository contains:

- generic architecture;
- generic retrospective behavior;
- synthetic test fixtures;
- public implementation.

It does not contain:

- real session evidence;
- private harness snapshots;
- private finding ledgers;
- personal context;
- private repository identifiers.

Private observations may influence the product only after being generalized into standalone requirements or synthetic tests.

## Non-goals

The architecture does not currently attempt to provide:

- provider-specific session parsing;
- continuous session monitoring;
- provider telemetry;
- automatic prompt optimization;
- autonomous harness rewriting;
- a universal scoring model for agent quality;
- a universal taxonomy of agent failures;
- a general-purpose code review system;
- a general repository architecture audit;
- cloud synchronization;
- a web dashboard;
- a multi-user collaboration service.

These capabilities require independent evidence before being added.
