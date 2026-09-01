<!-- docs/decisions/0001-target-the-user-controlled-outer-harness.md -->

# ADR 0001: Target the user-controlled outer harness

- Status: Accepted
- Date: 2026-09-01

## Context

The term *agent harness* can refer broadly to the mechanisms that turn a language model into an operational coding agent.\
A coding-agent product already contains vendor-controlled harness behavior, including context management, tool execution, orchestration, and other product internals.\
Users and repositories can also build another layer around that product through instructions, Skills, Hooks, settings, integrations, documentation, tests, and CI.

Retrospective analysis needs a clear authority boundary.\
Without one, an observed failure could incorrectly become a user instruction intended to compensate for behavior that is actually controlled by the provider, or a repository-local pattern could be promoted into a global rule without sufficient evidence.

## Decision

`agent-harness-retro` will target the **user-controlled outer harness**.

The project defines the outer harness as the layer of agent-facing mechanisms controlled by the user or repository owner around a coding-agent product.

The target model includes two primary scopes:

- global outer harness;
- repository harness.

Possible target artifacts include:

- instructions and rule files;
- Skills;
- Hooks;
- user-controlled settings and permissions;
- custom agent or subagent definitions;
- MCP and plugin configuration;
- agent-facing repository documentation;
- deterministic validation;
- tests, CI, and evaluation assets.

The project does not modify vendor-controlled built-in harness internals.

When retrospective evidence indicates a likely built-in harness limitation, that may be reported as a finding rather than converted automatically into an outer-harness workaround.

Before proposing a persistent harness change, the retrospective must determine the most appropriate scope and control surface.

A repeated pattern does not by itself justify a global rule, repository rule, or new Skill.

## Consequences

### Positive

- The project has a clear modification-authority boundary.
- Global and repository-local retrospective analyses use the same conceptual model.
- Provider-controlled behavior is less likely to create inappropriate user-level workarounds.
- Repository-local behavior is less likely to leak into global guidance.
- Findings can result in no harness change when another layer is responsible.
- The project remains applicable across coding-agent products without claiming control over their internals.

### Negative

- Some observed failures cannot be fixed directly by this project.
- Scope resolution becomes part of retrospective judgment.
- User-controlled configuration and provider-controlled behavior may sometimes interact in ways that are difficult to separate conclusively.
- The project cannot guarantee correction of built-in harness limitations.

These costs are preferable to treating every observed behavior as writable user configuration.

## Alternatives considered

### Treat the entire agent harness as one writable system

Rejected.

The user does not control vendor product internals, and merging the two layers would obscure authority and source-of-truth boundaries.

### Target repository harnesses only

Rejected.

Many important instructions, Skills, Hooks, and settings are intentionally global and require the same retrospective logic.

### Target only global user configuration

Rejected.

Many useful improvements are repository-specific and should not become global defaults.

### Avoid the term outer harness

Rejected for the internal conceptual model.

A more specific term is needed to distinguish user-controlled configuration from the vendor-provided built-in harness.

The repository name remains `agent-harness-retro` for clarity and discoverability.

## Revisit conditions

Revisit this decision if:

- coding-agent products expose a formally supported user-editable harness layer with a clearer industry-standard name;
- provider and user ownership boundaries materially change;
- the global/repository scope model cannot represent a recurring real-world target.
