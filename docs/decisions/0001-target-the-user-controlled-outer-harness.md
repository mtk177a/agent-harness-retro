# ADR 0001: Target the user-controlled outer harness

- Status: Accepted
- Date: 2026-09-01

## Context

An *agent harness* is the set of mechanisms that lets a language model work as a coding agent.\
Those mechanisms can belong to different owners:

| Layer | Controlled by | Role |
| --- | --- | --- |
| Built-in harness | The coding-agent provider. | Product internals, such as context management, tool execution, and orchestration. |
| Outer harness | The user or repository owner. | Agent-facing mechanisms added around the product. |

Retrospective analysis needs clear ownership and scope boundaries to avoid two mistakes:

- **Provider limitations**: Adding inappropriate user instructions to compensate for provider-controlled behavior.
- **Repository-local patterns**: Turning behavior from one repository into a global rule without sufficient evidence.

## Decision

### Review target

`agent-harness-retro` targets the **user-controlled outer harness**: the agent-facing mechanisms that a user or repository owner controls around a coding-agent product.

The target model has two primary scopes:

| Scope | Applies to |
| --- | --- |
| Global outer harness | Mechanisms shared across repositories or projects. |
| Repository harness | Mechanisms owned by one repository or project. |

The project does not modify provider-controlled built-in harness internals.

### Target artifacts

Possible target artifacts include:

- Instructions and rule files.

- Skills.

- Hooks.

- **Execution controls**

  - User-controlled settings.
  - Permissions.

- **Agent definitions**

  - Custom agents.
  - Subagents.

- **Integration configuration**

  - MCP configuration.
  - Plugin configuration.

- Agent-facing repository documentation.

- **Validation**

  - Deterministic validation.
  - Tests.
  - CI.
  - Evaluation assets.

### Proposal judgment

- **Provider limitations**: A likely built-in harness limitation may be reported as a finding without automatically proposing an outer-harness workaround.
- **Response selection**: Before proposing a persistent harness change, determine the most appropriate scope and control surface.
- **Recurrence**: A repeated pattern alone does not justify a global rule, repository rule, or new Skill.

## Consequences

### Positive

- **Authority**: The boundary of what the project can modify is clear.
- **Shared model**: Global and repository-local reviews use the same conceptual model.
- **Provider limitations**: Provider-controlled behavior is less likely to lead to inappropriate user-level workarounds.
- **Local scope**: Repository-local behavior is less likely to become global guidance.
- **No change**: A finding can result in no harness change when another layer is responsible.
- **Portability**: The project applies across coding-agent products without claiming control over their internals.

### Negative

- **Fix coverage**: Some observed failures cannot be fixed directly, and the project cannot guarantee correction of built-in harness limitations.
- **Scope assessment**: Determining the response scope is part of retrospective judgment.
- **Uncertainty**: Interactions between user-controlled configuration and provider-controlled behavior can be difficult to separate conclusively.

These limits are preferable to treating every observed behavior as user-editable configuration.

## Alternatives considered

| Alternative | Decision | Reason |
| --- | --- | --- |
| Treat the entire agent harness as one writable system. | Rejected. | Users do not control provider internals; merging the layers would obscure authority and sources of truth. |
| Target repository harnesses only. | Rejected. | Many instructions, Skills, Hooks, and settings are intentionally global and need the same retrospective logic. |
| Target only global user configuration. | Rejected. | Many useful improvements belong to one repository and should not become global defaults. |
| Avoid the term *outer harness*. | Rejected for the internal conceptual model. | A specific term is needed to distinguish user-controlled configuration from the built-in harness; the repository name remains `agent-harness-retro` for clarity and discoverability. |

## Revisit conditions

Revisit this decision if:

- **Terminology**: Coding-agent products expose a formally supported user-editable harness layer with a clearer industry-standard name.
- **Ownership**: Provider and user ownership boundaries materially change.
- **Scope model**: The global/repository scope model cannot represent a recurring real-world target.
