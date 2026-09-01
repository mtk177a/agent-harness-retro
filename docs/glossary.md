# Glossary

This glossary defines project-specific terminology used by `agent-harness-retro`.\
The terms are intended to keep vendor-controlled coding-agent systems separate from user-controlled harness configuration.

## Model

The underlying language model used by a coding agent.\
Examples include model families provided by OpenAI, Anthropic, or other vendors.\
A model is not itself the harness.

## Agent harness

A broad term for the system that turns a model into an operational agent by supplying tools, context, orchestration, execution, memory, and other control mechanisms.\
Because the term can refer both to vendor-controlled products and user-controlled layers, this project uses more specific terms where the distinction matters.

## Built-in harness

The vendor-controlled agent harness provided as part of a coding-agent product.\
For this project, Codex and Claude Code are examples of products containing built-in harness behavior around their underlying models.\
`agent-harness-retro` does not attempt to modify vendor-controlled built-in harness internals.

## Outer harness

The user-controlled layer built around a coding-agent product.\
An outer harness may include instructions, Skills, Hooks, settings, agent definitions, integrations, documentation, validation, and other mechanisms controlled by the user or repository owner.\
The outer harness is the primary retrospective target of this project.

## Global outer harness

Outer-harness artifacts intended to apply across multiple repositories or projects for one user or environment.

## Repository harness

The repository-scoped portion of an outer harness.\
It consists of agent-facing mechanisms owned by one repository or project, such as repository instructions, local Skills, Hooks, documentation, validation commands, tests, and CI.

## Harness artifact

One concrete source-of-truth item that participates in an outer harness.\
Examples include an instruction file, Skill, Hook configuration, settings file, evaluation asset, or repository validation script.

## Harness engineering

The activity of designing, implementing, evaluating, and improving an agent harness.\
The term is broader than `agent-harness-retro` and may apply to both built-in and outer harnesses.\
This project focuses on retrospective improvement of user-controlled outer harnesses.

## Session evidence

Normalized observations from coding-agent sessions used as retrospective evidence.\
`agent-harness-retro` obtains session evidence through `agent-sessions` rather than provider-native transcript formats.

## Retrospective

A bounded analysis of past coding-agent behavior intended to identify evidence-backed opportunities to improve an outer harness.\
A retrospective does not imply that a harness change will be made.

## Finding

A supported observation produced by retrospective analysis.\
A finding describes what appears to have happened and why it matters.\
It is distinct from a proposed solution.

## Evidence

The observations supporting a finding.\
Evidence may reference session sources, verified source versions, bounded events, current harness artifacts, or deterministic validation results.

## Proposal

A candidate response to one or more findings.\
A proposal identifies an intended target, expected effect, tradeoffs, and validation requirements.\
A proposal is not an approved change.

## Materialization

The act of applying an approved proposal to an actual outer-harness artifact.\
Materialization is separate from retrospective analysis and requires explicit authorization.

## Validation

Deterministic or behavioral evidence used to determine whether a materialized harness change produces the intended effect without unacceptable regressions.

## Scope

The layer at which a finding or proposal applies.

Typical scopes include:

- built-in harness;
- global outer harness;
- repository harness;
- one harness artifact;
- underlying project implementation;
- no persistent harness change.

## Control surface

The mechanism through which an outer-harness behavior is expressed or enforced.\
Examples include instructions, Skills, Hooks, settings, tests, CI, and documentation.

## Consumer state

Private persistent state owned by `agent-harness-retro`, such as reviewed source versions, findings, evidence relationships, proposal decisions, and verification state.\
Consumer state is separate from provider-owned session data and `agent-sessions`.
