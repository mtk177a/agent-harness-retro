# Glossary

Terms used in the README, architecture, and retrospective reports.

- [Harness](#harness)
- [Retrospective analysis](#retrospective-analysis)
- [Responses](#responses)

## Harness

### Model

The language model used by a coding agent, such as a model provided by OpenAI or Anthropic.\
The model itself is not the harness.

### Agent harness

The mechanisms that let a language model work as an agent, such as tools, context management, orchestration, execution, and memory.

The term can refer to either of these layers:

| Layer | Controlled by |
| --- | --- |
| [Built-in harness](#built-in-harness) | The coding-agent provider. |
| [Outer harness](#outer-harness) | The user or repository owner. |

### Built-in harness

The provider-controlled harness that is part of a coding-agent product.\
Codex and Claude Code are examples of products with a built-in harness around their models.\
`agent-harness-retro` does not modify that harness.

### Outer harness

The user-controlled layer around a coding-agent product.\
This is the main review target of `agent-harness-retro`.

Examples include instructions, Skills, Hooks, settings, agent definitions, integrations, documentation, and validation.\
The [README](../README.md#review-targets) describes these mechanisms in more detail.

### Global outer harness

The part of an outer harness intended for use across multiple repositories or projects for one user or environment.

### Repository harness

The part of an outer harness owned by one repository or project.\
Examples include repository instructions, local Skills, Hooks, project documentation, validation commands, tests, and CI.

### Harness artifact

One item managed as part of an outer harness.\
It refers to the maintained source of truth, such as an instruction file, Skill, Hook configuration, settings file, evaluation asset, or repository validation script.

### Harness engineering

The work of designing, implementing, evaluating, and improving an agent harness.\
It can apply to both built-in and outer harnesses.\
`agent-harness-retro` covers one part of this work: reviewing past behavior to improve the user-controlled outer harness.

## Retrospective analysis

For this example, assume a fictional project whose current testing guide documents the correct test command.

| Term | Question it answers | Example |
| --- | --- | --- |
| [Observation](#observation) | What do the records directly show? | The agent ran the command listed in `AGENTS.md`, and it failed.<br>The current testing guide lists a different command. |
| [Finding](#finding) | What problem or useful pattern does the evidence support, and why does it matter? | `AGENTS.md` and the current testing guide give different commands, leaving conflicting instructions for the agent. |
| [Cause hypothesis](#cause-hypothesis) | What might explain the finding? | The command in `AGENTS.md` may not have been updated when the testing procedure changed. |
| [Proposal](#proposal) | What response should the owner consider? | Replace the command in `AGENTS.md` with a link to the testing guide. |

### Retrospective

A bounded analysis of past coding-agent behavior to find evidence-backed improvements to an outer harness.\
A retrospective does not imply that the harness needs to change.

### Session evidence

Normalized coding-agent session records used to review agent behavior.\
`agent-harness-retro` obtains them through `agent-sessions`, rather than parsing provider-native transcript formats.

### Logical lineage

A group of related session sources counted as one occurrence when assessing recurrence.

- **Same source:** Records with the same `source_ref`, including new versions after resume, remain in one lineage.
- **Linked sources:** Explicit parent or fork relationships connect sources into one lineage.
- **Retries:** Repeated attempts within a lineage do not count as independent recurrence.
- **Uncertain relationships:** Do not count affected sources as independent when the necessary relationship evidence is missing or ambiguous.

### Observation

A statement about what normalized session evidence or a checked current harness artifact directly shows.\
Keep it separate from an interpretation or proposed action.

### Finding

A problem, useful pattern, or consequence supported by one or more observations, with an explanation of why it matters.\
A finding does not specify a solution.

### Cause hypothesis

A possible explanation for a finding.\
It remains uncertain rather than being treated as an observation.

- **Support:** Evidence that supports the explanation.
- **Counterevidence:** Evidence that conflicts with it, when material.
- **Alternatives:** Other explanations, when material.
- **Uncertainty:** What remains unresolved.

### Evidence

The observations used to support a finding.

Evidence references can identify:

- **Sources:** Session sources and their verified versions.
- **Events:** Bounded messages or tool events.
- **Harness artifacts:** Current files or mechanisms checked during the review.
- **Validation results:** Results from deterministic checks.

### Consumer state

Private persistent state owned by `agent-harness-retro` to carry information between reviews.\
Examples include reviewed source versions, findings, evidence relationships, proposal decisions, and verification status.\
Keep it separate from provider-owned session data and `agent-sessions`.

The current Skill does not persist this state; it is a possible future capability.\
The [architecture](architecture.md#private-retrospective-state) defines its storage constraints.

## Responses

### Scope

The layer where a finding matters or a proposal applies.

- **Provider:** The built-in harness.
- **Global:** The outer harness shared across projects.
- **Repository:** The harness owned by one repository or project.
- **Individual artifact:** One harness artifact, such as a Skill or Hook configuration.
- **Project implementation:** The underlying project's code or implementation.

A scope assessment can also conclude that no persistent harness change is needed.\
Scope identifies where a response belongs; the [control surface](#control-surface) identifies the mechanism used for that response.

### Control surface

A mechanism that guides or enforces agent behavior in an outer harness.\
Examples include instructions, Skills, Hooks, settings, tests, CI, and documentation.

### Routing candidate

An option for how to respond to a finding.\
Compare candidates only after the finding and likely responsibility scope have been established.

- **Guidance:** Global instructions or repository guidance.
- **Workflow:** A Skill.
- **Enforcement:** A deterministic check.
- **Project change:** The canonical project implementation or documentation.
- **Memory handoff:** A candidate for a human to consider, without automatically writing memory.
- **One-off action:** A response limited to the current situation.
- **No persistent change:** A response that needs no lasting harness change.

The [Skill's routing table](../skills/outer-harness-retrospective/SKILL.md#scope-and-routing) defines the detailed candidates.

### Proposal

A candidate response to one or more findings.\
It names the target, expected improvement, tradeoffs, and validation needed after applying the change.\
A proposal is not an approved change.

### Decision outcome

The result of assessing a finding or candidate proposal.

| Outcome | Meaning |
| --- | --- |
| `proposal` | The evidence supports a response for the owner to consider. |
| `no_change` | The finding is supported, but no persistent response is needed or existing controls are sufficient. |
| `reject` | A candidate response is contradicted, duplicates an existing mechanism, has the wrong scope, or exceeds the permitted change scope. |
| `insufficient_evidence` | The evidence is not sufficient to establish the finding or make the required decision. |

Rejecting a proposal does not erase a supported finding.\
Missing evidence does not by itself establish that no change is needed.

### Materialization

Applying an approved proposal to an actual outer-harness artifact.\
This is separate from retrospective analysis and requires explicit user approval.\
The current Skill does not apply proposals.

### Validation

Checking whether an applied harness change has the intended effect without unacceptable regressions.\
Validation uses deterministic checks or behavioral evidence.
