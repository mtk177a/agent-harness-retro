# Repository instructions

These instructions apply to agents working on `agent-harness-retro`.\
The project reviews coding-agent session evidence to evaluate the user-controlled outer harness.

## Start work

Before a material design or implementation change, read:

1. [README](README.md).
2. [Architecture](docs/architecture.md).
3. [Glossary](docs/glossary.md).
4. Any ADR directly relevant to the change in [docs/decisions](docs/decisions/).

Do not scan unrelated files merely to increase coverage.

## Design boundaries

Preserve these constraints unless an explicit architecture decision supersedes them:

- **Review target**: Evaluate the user-controlled outer harness; do not modify provider-controlled built-in harness internals.
- **Scope model**: Use the same core retrospective model for global and repository-local harnesses.
- **Separate layers**: Keep session evidence, retrospective findings, and applied harness artifacts separate.
- **Analysis**: Keep analysis separable from applying proposed changes.
- **Default execution**: Continuous finding generation for every event is not the default architecture.

### Session access

This repository owns retrospective concepts, findings, evidence relationships, proposals, and the boundary between analysis and optional harness changes.\
Provider-specific raw session access belongs to `agent-sessions`; consume its public machine-readable contract.

Before adding a direct Codex, Claude Code, or other provider transcript parser to work around missing upstream functionality, determine whether that capability belongs in `agent-sessions`.

## Information handling

### Public information boundary

This repository is public.\
Do not add:

- Private repository names or links.
- Private Issue or Project references.
- Personal or employer-specific context.
- Real coding-agent session contents.
- Real outer-harness contents from private environments.
- Credentials or authentication material.
- Real usernames, hostnames, home directories, or machine-specific paths.

A requirement learned from private usage may be included only after it is expressed as a general product requirement that does not depend on the private source.

Use synthetic repositories, sessions, harness artifacts, and identities for tests and examples.\
Never commit real private sessions as evaluation data, even after manual redaction.

### Private runtime state

Retrospective state is private runtime data.\
Do not store it in the public repository or copy raw provider transcripts into this repository or retrospective state.

Do not commit:

- Reviewed-session ledgers.
- Findings derived from private sessions.
- Proposal decisions.
- User-specific harness inventories.
- Private evidence.

Keep the persistent representation and storage format as small as the product requirements allow.\
Do not introduce high-frequency writes, background synchronization, or a database without an observed need.

## Retrospective decisions

### Assess evidence

A material finding must be traceable to evidence.\
Prefer stable references:

- **Source**: An `agent-sessions` source reference.
- **Version**: The verified source version.
- **Events**: Bounded event references.
- **Scope**: The affected scope.
- **Harness**: Current harness artifacts relevant to the finding.

Do not retain complete raw transcripts merely to make a finding auditable.\
Do not treat a generated summary as stronger evidence than the underlying observed source.

Distinguish directly observed behavior, interpretation, hypothesis, proposal, and unverified assumption.

### Select the response scope

Before proposing a harness change, determine where the response belongs.

**Possible outcomes**

- Global outer harness.
- Repository harness.
- One existing Skill.
- One Hook or deterministic check.
- Underlying project code or documentation.
- Provider-controlled built-in behavior.
- No persistent change.

**Limits on permanent rules**

- **Global rules**: Do not turn repository-local behavior into a global rule without evidence from other contexts.
- **Permanent harness rules**: Do not turn a one-off user correction into a permanent harness rule without a reason that holds beyond the original session.
- **Recurrence**: Frequency alone does not justify permanent harness guidance.

### Select the response mechanism

Compare the relevant control surfaces rather than defaulting to more instructions:

| Need | Mechanism to consider |
| --- | --- |
| Durable behavioral guidance | Instruction or rule. |
| Conditional, multi-step workflow | Skill. |
| Mechanically enforceable intervention | Hook or deterministic check. |
| Repository correctness invariant | Test or CI. |
| User-controlled execution behavior | Setting or permission. |
| Project knowledge that serves as a source of truth | Project documentation. |

Prefer changing an existing appropriate mechanism over creating a parallel one.

### Describe the proposal

Keep findings separate from proposed solutions.\
A finding should state the observed pattern and evidence before suggesting a harness artifact.\
A finding does not imply that a change is required; rejection and no-change are valid outcomes.

A proposal should identify:

- **Scope**: Where the change applies.
- **Target**: The artifact or control-surface type to change.
- **Improvement**: The expected benefit.
- **Tradeoffs**: Relevant costs or disadvantages.
- **Evidence**: What supports the proposal.
- **Validation**: The checks or evaluation needed after the change.

### Apply an approved change

Finding an improvement candidate does not authorize an external harness edit.\
Require explicit approval before writing to any external harness artifact, including through a feature that applies proposals.

Preserve the target repository's authority, instructions, review requirements, and validation commands.\
Do not bypass its source-of-truth model.

## Repository changes

### Keep the change small

Prefer the smallest mechanism that provides the required retrospective capability.

- **Premature infrastructure**: Do not build a taxonomy, scoring model, persistence system, orchestration framework, or automatic change-application layer before a concrete requirement needs it.
- **Speculative automation**: Do not add background monitoring, automatic transcript ingestion, automatic harness rewrites, or a general-purpose session database merely because they may become useful later.

### Update architecture decisions

When changing an architectural invariant:

1. **Identify the decision**: Find the affected ADR.
2. **Assess it**: Determine whether the accepted decision still holds.
3. **Revise it**: Update or supersede the ADR when needed.
4. **Align the project**: Keep the README, glossary, architecture, implementation, and tests consistent.

### Update documentation

- **Current state**: In current-state product documentation, include implementation history only when it is needed to explain an active decision.

- **Language**

  - **Public documents and rules**: English is canonical.
  - **English defaults**: Use English by default for code, identifiers, code comments, CLI and Skill public interfaces, commit summaries, and release notes.
  - **Issues and pull requests**: Follow the [language policy](docs/localization.md#issues-and-pull-requests) for the title, Summary, body, and comments.
  - **External contributors**: Do not require them to use Japanese.
  - **Duplication**: Do not require full English and Japanese copies of Issue and pull request bodies.

- **Japanese references**

  - **Maintained pairs**: Follow the file pairs and synchronization rules in the [language policy](docs/localization.md#maintaining-japanese-references).
  - **Updates**: When adding or changing a maintained English file, use the repository-local [translation maintenance Skill](.agents/skills/maintain-japanese-references/SKILL.md) to create or review its Japanese reference in the same change.
  - **Unchanged meaning**: If the English edit does not affect Japanese meaning, leave the reference unchanged and record the reason in the pull request's Validation or Risks / Follow-up section.
  - **ADRs**: Do not translate ADRs as part of the maintained set.
  - **Authority**: Japanese references, including [AGENTS-ja.md](AGENTS-ja.md), define no independent requirements or authority; the English source takes precedence.

### Validate the change

- **Evaluation policy**: Use [docs/evaluation.md](docs/evaluation.md) for Skill audit and evaluation selection, execution, and reporting.
- **Retrospective quality**: When a change affects retrospective quality, test the decision boundary that can fail rather than expanding an unrelated evaluation suite.
- **Fixtures**: Use synthetic evidence and harness fixtures.
- **Final diff**: Before reporting completion, inspect the final diff for accidental private information and machine-specific data.
