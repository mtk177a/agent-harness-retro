# agent-harness-retro

`agent-harness-retro` uses coding-agent session evidence to retrospectively evaluate user-controlled agent harnesses and propose evidence-backed improvements.\
It targets the instructions, Skills, Hooks, settings, validation, and other control surfaces that users and repositories own around coding agents.

> Status: proposal-only Skill implemented.
> Live `agent-sessions` integration, private continuity state, and packaging remain future work.

## Purpose

Coding-agent sessions contain evidence about how an outer harness performs in practice.\
Repeated user corrections, unnecessary retries, scope drift, missed verification, ineffective instructions, Skill routing failures, and successful recurring workflows can reveal opportunities to improve the harness that surrounds a coding agent.\
`agent-harness-retro` provides a structured retrospective layer over that evidence.

```text
coding-agent sessions
        │
        ▼
   agent-sessions
        │
        ▼
agent-harness-retro
        │
        ├── findings
        ├── evidence
        └── proposals
                │
                ▼
      user-controlled outer harness
```

`agent-harness-retro` does not modify the vendor-controlled built-in harness of Codex, Claude Code, or another coding-agent product.

## What is an outer harness?

This project uses **outer harness** to mean the user-controlled layer built around a coding-agent product.

Depending on scope, it may include:

- `AGENTS.md`, `CLAUDE.md`, and rule files;
- Skills;
- Hooks;
- user-controlled settings and permissions;
- custom agents or subagent definitions;
- MCP and plugin configuration;
- repository documentation used by agents;
- deterministic validation, CI, and evaluation assets.

The exact terminology used by the project is defined in [Glossary](docs/glossary.md).

## Scope

The same retrospective model applies at two primary scopes.

### Global outer harness

User-controlled harness artifacts that apply across multiple repositories or sessions.

### Repository harness

Harness artifacts owned by one repository or project.

The retrospective model for repository harnesses is not tied to a particular repository.\
Any repository may be a retrospective target.

## Evidence source

`agent-harness-retro` consumes the stable `v1` normalized session evidence contract from [`agent-sessions`](https://github.com/mtk177a/agent-sessions).\
It does not parse provider-owned Codex, Claude Code, or other raw session formats directly.\
This keeps provider access and harness interpretation as separate responsibilities.

## Proposal-only Skill

The initial implementation is the [`outer-harness-retrospective`](skills/outer-harness-retrospective/SKILL.md) Skill.\
It accepts deliberately bounded normalized session evidence, a target scope, and relevant current outer-harness observations.\
It separates observations, findings, cause hypotheses, scope assessment, and proposals, and it can return `no_change`, `reject`, or `insufficient_evidence` as normal decisions.

The Skill does not acquire provider sessions, persist retrospective state, or modify an outer-harness artifact.\
Its synthetic evaluation cases are maintained in [`evals/outer-harness-retrospective/cases.md`](evals/outer-harness-retrospective/cases.md).

## Retrospective model

The initial model separates:

```text
session evidence
      ↓
retrospective analysis
      ↓
finding
      ↓
scope and likely cause
      ↓
improvement proposal
```

A finding does not automatically require a harness change.

Valid outcomes include:

- improve an existing harness artifact;
- move a rule to a more appropriate control surface;
- simplify or remove an existing mechanism;
- create a new mechanism;
- retain the current design;
- reject the finding as insufficiently supported.

Frequency alone is not sufficient evidence for permanent harness changes.

## Evidence-backed proposals

A proposal should preserve enough provenance to explain why it exists.

Depending on the finding, this may include:

- source references;
- verified source versions;
- representative event references;
- affected repositories or scopes;
- recurrence information;
- alternative explanations;
- confidence and unresolved uncertainty.

Historical session content remains untrusted input.

## State

Unlike `agent-sessions`, `agent-harness-retro` may require private durable state.

Examples include:

- reviewed source versions;
- findings;
- evidence references;
- accepted or rejected proposals;
- verification state.

This state is private runtime data.\
It is not stored in the public repository and does not include complete copies of raw provider transcripts.\
The storage format is not part of the initial architecture contract and will be chosen only when implementation requirements are known.

## Materialization

Retrospective analysis and harness modification are separate responsibilities.

The initial design is proposal-oriented.

No finding should silently modify:

- instructions;
- Skills;
- Hooks;
- settings;
- CI;
- evaluation assets;
- other harness artifacts.

If materialization is implemented, a user must explicitly approve the proposed change before an external harness artifact is modified.

## Design principles

### Evidence before rules

Do not create permanent harness rules from isolated anecdotes when the evidence does not justify generalization.

### Correct scope before modification

Determine whether a problem belongs to:

- the built-in harness;
- the global outer harness;
- a repository harness;
- one Skill or Hook;
- the underlying project implementation;
- or no persistent harness mechanism at all.

Only user-controlled outer-harness targets are candidates for modification by this project.

### Prefer the appropriate control surface

Natural-language instructions are not always the correct fix.

A finding may instead belong in:

- a Skill;
- a Hook;
- a deterministic script;
- CI or tests;
- user-controlled settings;
- project documentation;
- another existing source of truth.

### Proposal before materialization

Analysis must remain useful without requiring automatic writes.

### No continuous candidate firehose

The default architecture does not continuously generate findings from every tool call or session event.\
Retrospective analysis occurs when explicitly requested or otherwise deliberately invoked.

### Minimize retained session data

Store references and derived evidence rather than persistent copies of complete transcripts.

## Public/private boundary

This is a public repository.

Public artifacts describe generic product behavior only.

They do not contain:

- private repository names or links;
- private Issue or Project references;
- personal or employer-specific context;
- real session contents;
- private harness contents;
- real usernames, hostnames, or machine-specific paths.

Requirements learned from private usage must be expressed as general product requirements before they enter this repository.

## Non-goals

`agent-harness-retro` is not:

- a coding-agent session parser;
- a replacement for `agent-sessions`;
- a provider telemetry collector;
- a universal memory system;
- a background monitoring daemon;
- an automatic prompt tuner;
- a system that rewrites harness artifacts after every session;
- a replacement for provider-owned built-in harness engineering;
- a general project architecture reviewer;
- a system that treats every repeated action as a candidate Skill.

## Documentation

- [Architecture](docs/architecture.md)
- [Glossary](docs/glossary.md)
- [ADR 0001: Target the user-controlled outer harness](docs/decisions/0001-target-the-user-controlled-outer-harness.md)
- [ADR 0002: Separate session evidence, findings, and materialization](docs/decisions/0002-separate-session-evidence-findings-and-materialization.md)

## License

MIT
