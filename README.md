# agent-harness-retro

`agent-harness-retro` is a Skill that reviews past work with coding agents such as Codex and Claude Code.\
It compares selected conversations and tool activity with your current instructions and configuration.\
It proposes changes with supporting evidence and a way to check their effect.

**Contents**

- [Capabilities](#capabilities)
- [Review targets](#review-targets)
- [Review workflow](#review-workflow)
- [Review criteria](#review-criteria)
- [Information handling](#information-handling)
- [Acquisition script maintenance](#acquisition-script-maintenance)
- [Documentation](#documentation)
- [License](#license)

## Capabilities

### When to use it

- **Investigate repeated problems:** Find out why you keep giving an agent the same correction and whether the harness needs to change.
- **Reuse successful workflows:** Identify a useful way of working to document or reuse.

### Example: an outdated test command

Suppose an agent follows a test command in `AGENTS.md`, the command fails, and you have to explain the correct command.\
In this fictional project, the current testing guide already documents the correct command.

- **Input**

  - **Session evidence:** The relevant conversation and command results obtained through `agent-sessions`.
  - **Current instructions:** `AGENTS.md`.
  - **Current procedure:** The testing guide.

- **Possible proposal:** Replace the outdated command in `AGENTS.md` with a link to the testing guide, citing the failed command and the conflicting instructions as evidence.

- **Suggested validation:** On a representative coding task, check whether the agent follows the link and runs the documented test command.

The Skill does not apply the proposed edit; you decide whether to make it.

### Availability

| Item | Status | Details |
| --- | --- | --- |
| Retrospective Skill | Implemented | [`outer-harness-retrospective`](skills/outer-harness-retrospective/SKILL.md) analyzes evidence and proposes improvements. |
| Evidence acquisition script | Implemented | Bundled under the Skill's `scripts/` directory to acquire metadata. |
| Executable evaluation cases | Available | The Skill's [`evals`](skills/outer-harness-retrospective/evals/README.md) directory contains cases that use fictional data. |
| Review continuity | Future work | The current Skill does not save findings between reviews. |

## Review targets

### What is an outer harness?

An **outer harness** is the set of instructions, tools, and settings that you or a repository owner control around a coding agent.

- **Instructions:** `AGENTS.md`, `CLAUDE.md`, and other rule files.

- **Skills:** Reusable workflows.

- **Hooks:** Actions that run at specified workflow points.

- **Execution settings:** Settings that you can change to control execution.

- **Permissions:** The operations you can allow an agent to perform.

- **Agent definitions**

  - Custom agents.
  - Subagents.

- **Integration settings**

  - MCP configuration.
  - Plugin configuration.

- **Project documentation:** Repository documents that agents read.

- **Behavior checks**

  - Validation scripts.
  - CI.
  - Evaluation cases.

The coding-agent provider controls the **built-in harness**, such as its internal context management and tool runtime.\
This project can report a problem in that layer, but it cannot change the provider's internals.\
The [Glossary](docs/glossary.md) defines these terms in more detail.

### Application scope

| Scope | Where it applies | Example |
| --- | --- | --- |
| Global outer harness | Multiple repositories or sessions. | A shared Skill for preparing pull requests in different projects. |
| Repository harness | One repository or project. | An `AGENTS.md` file that points to the project's test instructions. |

## Review workflow

A review runs when someone explicitly requests or deliberately starts it.\
The Skill does not monitor sessions in the background.

### Prepare the inputs

- **Session evidence:** Records selected for a defined review.
- **Review target:** The repository or global artifact to examine, such as a shared Skill.
- **Current harness:** Relevant files or observations that identify what was checked.

```mermaid
flowchart TD
    source["agent-sessions v2"] --> helper["Acquisition script"]
    helper --> metadata["Metadata"]
    source --> excerpts["Selected, privacy-checked excerpts"]
    metadata --> review["Retrospective Skill"]
    excerpts --> review
    scope["Review target and scope"] --> review
    harness["Current instructions and configuration"] --> review
    review --> result["Findings and decisions"]
```

### Acquire session evidence

The pilot script prepares a limited set of session metadata for the Skill's analysis.

#### Requirements

- Python 3
- Git
- Go 1.26 or later
- Network access to the public `agent-sessions` repository for a normal run

#### Run the script

##### Normal run

Run from the repository root:

```sh
python3 skills/outer-harness-retrospective/scripts/run.py
```

##### Offline run

For an offline run, add `--agent-sessions-bin <absolute-path>` to use an existing executable instead of fetching and building one.

```sh
python3 skills/outer-harness-retrospective/scripts/run.py \
  --agent-sessions-bin /absolute/path/to/agent-sessions
```

The script checks the executable against the commit in `agent-sessions.commit` and requires a clean source build.

##### Common options

- **Configuration file:** For a normal or offline run, add `--config <absolute-path>` to use an existing `agent-sessions` configuration.

#### Evidence selection

**Review period**

The pilot fixes the review period when it starts and uses the preceding 14 days for Codex and Claude Code.

**Session grouping**

To avoid counting a resumed or forked conversation as a separate occurrence, the script groups explicitly related sessions into a **logical lineage**.

**Count limits**

| Counted item | Limit |
| --- | --- |
| Candidate lineages examined | At most ten for each provider. |
| Complete, independent lineages required | At least three from each provider. |
| Included events, such as messages and tool calls | At most 2,000 in total. |

These numbers limit this pilot run; they do not define how many occurrences justify a permanent rule.

**Insufficient evidence**

The pilot returns `insufficient_evidence` when it cannot establish a complete selection, for example because:

- **Unknown review period:** The review period cannot be established.
- **Unsupported response schema:** An `agent-sessions` response is missing `schema_version` or uses a schema other than `v2`.
- **Missing records:** Required records are missing.
- **Missing event details:** Required event details are missing.
- **Unverified source versions:** Source versions cannot be verified.
- **Changed source versions:** A source version changes during collection.
- **Mixed recording sources:** Selected records for the same provider come from multiple recording sources (`source_instance`).
- **Too few lineages:** Too few complete, independent lineages fit within the limits above.

#### Use the acquisition result

**Output**

The `v2` interface does not guarantee automatic removal of private information from event bodies.\
The script therefore writes only metadata to a temporary file, without event bodies.

| Output | Contents |
| --- | --- |
| Private temporary file outside this repository | Evidence metadata, such as source references, versions, event indexes, and event types. |
| Standard output | Acquisition status.<br>Aggregate counts. |
| Standard error | The location of the private metadata file. |

**Handoff steps**

1. **Supply the metadata:** Pass the script's metadata file to the Skill.

2. **Supply current harness information:** Provide relevant files or observations that identify what was checked.

3. **Select excerpts:** Choose the session excerpts to use in the analysis.

4. **Check privacy:** Check the selected excerpts for private information.

5. **Supply the excerpts**

   Pass the checked excerpts to the Skill separately from the metadata.\
   Metadata alone cannot show what an agent did or why.

6. **Remove the temporary file:** Delete the metadata file after the review.

### Read the analysis result

#### Analysis records

- **Observation:** What the records or checked files directly show.
- **Finding:** What problem or useful pattern those observations support, and why it matters.
- **Cause hypothesis:** What might explain the finding, including uncertainty.
- **Scope assessment:** Where the cause or response belongs, such as a repository instruction or a provider-controlled feature.
- **Proposal:** A possible change, its evidence, and a way to check its effect.

#### Decision outcomes

| Outcome | Meaning |
| --- | --- |
| `proposal` | The evidence supports a proposed response. |
| `no_change` | No change is needed. |
| `reject` | A suggested response should be rejected. |
| `insufficient_evidence` | The evidence does not support a decision. |

#### Proposed changes

- **Improve:** Revise an existing harness artifact.
- **Move:** Transfer a rule to a mechanism better suited to the problem.
- **Simplify:** Reduce the complexity of an existing mechanism.
- **Remove:** Delete an unnecessary mechanism.
- **Create:** Add a new mechanism.

### Approval before changes

A request to review a harness does not authorize changes to its files or settings.\
Any future feature that applies proposals must obtain explicit user approval before modifying an external harness.

## Review criteria

### Assess the evidence

Compare the session evidence with current instructions and tools before proposing a change.

An isolated example does not justify a permanent rule unless the evidence supports applying it beyond that example.\
Frequency alone does not justify a permanent rule.

### Choose the response scope

Match the response scope to the cause of the finding:

- **Provider:** The provider-controlled built-in harness.
- **Global harness:** Artifacts shared across repositories.
- **Repository harness:** Artifacts owned by one repository.
- **Individual artifact:** One existing harness artifact, such as a Skill or Hook.
- **Project implementation:** The project's code or other implementation.
- **One-off response:** An action limited to the current situation.
- **No change:** No permanent harness mechanism is needed.

A finding about a provider-controlled problem can be reported without proposing an outer-harness workaround.

### Choose the response mechanism

Adding natural-language instructions is not always the best response.

| Need | Mechanism to consider |
| --- | --- |
| A stable behavioral principle needs to apply across tasks. | Instructions. |
| A reusable task needs conditional judgment or multiple steps. | A Skill. |
| An action must run at a specific workflow point. | A Hook. |
| A condition can be checked without model judgment. | A deterministic validation script. |
| Code changes must preserve project requirements. | Tests. |
| Project checks must run automatically when code changes. | CI. |
| Execution behavior needs to change. | Execution settings. |
| The operations allowed for an agent need to change. | Permissions. |
| Project information is missing or outdated. | The canonical project documentation. |

Prefer updating an existing appropriate mechanism over creating a parallel one.

### Explain the proposal evidence

Select evidence details to suit the finding:

- **Source references:** Which sessions support the finding.
- **Verified source versions:** Which versions of those records were checked.
- **Event references:** Which messages or tool events show the behavior.
- **Affected scope:** Which repositories, instructions, or tools are involved.
- **Recurrence:** Whether the behavior occurred in independent interactions or repeated attempts within one interaction.
- **Alternative explanations:** What else might account for the behavior.
- **Confidence:** How strongly the evidence supports the conclusion.
- **Uncertainty:** What remains unknown.

For a proposed rule, explain why it is needed, where it applies, and what it is expected to improve.

## Information handling

### Session record source

[`agent-sessions`](https://github.com/mtk177a/agent-sessions) reads coding-agent records and presents them in a common format.\
`agent-harness-retro` uses its stable `v2` interface to read messages, tool activity, and information about their sources.\
It does not parse provider-specific raw formats or collect provider telemetry.

### Instructions in past conversations

Past session content is evidence to examine, not instructions for the reviewing agent to follow.

### Keep review data private

Do not add real session contents, private harness files, or private findings to this public repository.\
Use synthetic data for examples and tests.\
The [public information boundary](AGENTS.md#public-information-boundary) lists the full publication rules.

The [Architecture](docs/architecture.md#private-retrospective-state) describes possible future records and their storage constraints.

## Acquisition script maintenance

### Launcher operation

During a normal run, the [`run.py`](skills/outer-harness-retrospective/scripts/run.py) launcher:

1. **Fetch source:** Downloads the commit recorded in [`agent-sessions.commit`](skills/outer-harness-retrospective/scripts/retro-pilot/internal/pilot/agent-sessions.commit).

2. **Build:** Creates an executable in a private temporary directory.

3. **Start the helper:** Runs the bundled Go helper with the executable.

4. **Verify the build**

   The helper checks the executable's Go build information:

   - **Commit:** It matches the recorded commit.
   - **Source tree:** It was unchanged when built.

5. **Remove temporary files:** Deletes the downloaded source and built executable after the run.

The launcher does not install a global executable.\
The upstream source is [`mtk177a/agent-sessions`](https://github.com/mtk177a/agent-sessions), under the MIT license.

### Update the pinned version

1. **Review the upstream changes**

   - **Interface:** Check changes to the public contract.
   - **License:** Check changes to the license.

2. **Change the pinned commit:** Update [`agent-sessions.commit`](skills/outer-harness-retrospective/scripts/retro-pilot/internal/pilot/agent-sessions.commit).

3. **Verify the change**

   Before using real sessions, run the checks in [Evaluation and Skill audit](docs/evaluation.md):

   - Go checks.
   - Python checks.

## Documentation

- [Architecture](docs/architecture.md)
- [Evaluation and Skill audit](docs/evaluation.md)
- [Glossary](docs/glossary.md)
- [Localization policy](docs/localization.md)
- [ADR 0001: Target the user-controlled outer harness](docs/decisions/0001-target-the-user-controlled-outer-harness.md)
- [ADR 0002: Separate session evidence, findings, and materialization](docs/decisions/0002-separate-session-evidence-findings-and-materialization.md)

## License

MIT
