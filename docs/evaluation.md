# Evaluation policy

This policy applies to the repository's maintained Skills and changes that affect retrospective decisions.\
It applies the [change-scoped evaluation rules from `mtk177a/skills` Issue #49](https://github.com/mtk177a/skills/issues/49) within this repository's public and private evidence boundary.

- [Audit a Skill](#audit-a-skill)
- [Select evaluation work](#select-evaluation-work)
- [Prepare cases](#prepare-cases)
- [Run static checks](#run-static-checks)
- [Run model evaluations](#run-model-evaluations)
- [Record results](#record-results)

## Audit a Skill

### Audit record

For each material Skill audit, record:

| Field | What to record |
| --- | --- |
| Revision | The version being audited. |
| Working tree | The working-tree state. |
| Intended behavior | What the Skill should do. |
| Users | Who the Skill affects. |
| Client | The client that uses the Skill. |
| Guidance surfaces | Where the relevant instructions are defined. |
| Authority | Which instructions have authority. |
| Evidence | What was observed. |
| Findings | The supported problems or useful patterns. |
| Alternatives | Other responses considered. |
| Uncertainty | What remains unknown. |

### Audit checks

Check only the aspects that matter to the Skill:

- **Responsibility:** What the Skill owns.
- **Scope:** Where it applies.
- **Selection:** When it should be selected.
- **Coexistence:** How it works alongside other Skills.
- **Source precedence:** Which source takes priority.
- **Evidence quality:** How well the evidence supports the findings.
- **Safety:** The relevant safety concerns.
- **Context cost:** The context the Skill requires.
- **Portability:** Its use across clients or environments.
- **Maintenance:** The work needed to maintain it.

Review the canonical English Skill separately from its Japanese reference.\
Check the reference for faithful meaning and natural Japanese, even if the English text has not changed.

### Audit results

Classify each issue:

| Classification | Meaning |
| --- | --- |
| `change_required` | A change is needed. |
| `no_change` | No change is needed for this issue. |
| `separate_issue` | Handle the issue separately from the current work. |
| `unverified` | The issue has not been verified. |

For each material finding, include:

- **Evidence:** What supports the finding.
- **Impact:** Why it matters.
- **Confidence:** How strongly the evidence supports it.
- **Falsification:** A check that could disprove it.

Do not treat a static audit as confirmation of actual model behavior.

## Select evaluation work

Before selecting cases, identify:

- **Changed responsibility:** What the change affects.
- **Possible failure:** What could go wrong.
- **Detection:** What observation would reveal the failure.

Use the least evidence needed for the decision.

| Change | Default evaluation path |
| --- | --- |
| Documentation, format, or meaning-preserving text | `static-only` |
| Skill instructions, output contract, or safety boundary | `targeted-candidate` |
| Name, description, invocation, or adjacent responsibility boundary | `targeted-routing` |
| Known regression, major redesign, or ambiguous candidate result | `baseline-comparison` |
| Explicit client or environment support claim | `target-environment` |

Add any of the following only if its result could change the decision:

- A baseline comparison.
- A comparison without the Skill.
- Repeated runs.
- Another client.

Do not run the complete suite after every edit.

## Prepare cases

### Definition files

Behavior cases check what the Skill does; routing cases check whether it is selected and loaded.

| Case type | File |
| --- | --- |
| Behavior | `skills/<name>/evals/evals.json` |
| Routing | `skills/<name>/evals/triggers.json` |

Use the Agent Skills executable format:

| Purpose | Fields |
| --- | --- |
| Envelope (top level) | `skill_name`, `evals` |
| Case identity | `id` |
| Executor input | `prompt`, optional `files` |
| Behavior grading | `expected_output`, `assertions` |
| Routing grading | `expected_handlers` |

Put only executor input in the case prompt and fixture.\
Hide assertions and expected output from the executor and use them for grading.

### Test data

Use synthetic identities, repositories, harness artifacts, and sessions.\
Never commit a real session, even after manual redaction.

### Definition migration

When migrating existing `evals.json` or `triggers.json` definitions, migrate the whole affected set together.\
A format-only migration needs static validation; it does not require a model run.

## Run static checks

### Repository checks

Run from the repository root:

```sh
python3 scripts/check_repository.py
```

This checks maintained translation pairs, Skill metadata, executable definitions, and any current compact report without calling a model.

### Acquisition script checks

When changing the bundled acquisition scripts, also perform these checks:

- **Python tests**

  Run from the repository root:

  ```sh
  python3 -m unittest discover -s skills/outer-harness-retrospective/scripts -p 'test_*.py'
  ```

- **Go tests**

  Run from `skills/outer-harness-retrospective/scripts/retro-pilot`:

  ```sh
  go test ./...
  ```

- **Pinned build:** Verify that the pinned source builds without reading real sessions.

## Run model evaluations

### Evaluation plan

Before any model call, write a plan that identifies:

| Field | What to specify |
| --- | --- |
| Candidate revision | The version to evaluate. |
| File hashes | The hashes of the candidate files. |
| Responsibility | The responsibility affected by the change. |
| Evaluation path | The selected evaluation path. |
| Cases and conditions | The case-condition pairs to run. |
| Client | The client used for the run. |
| Model | The model used for the run. |
| Sandbox | The execution sandbox. |
| Call limit | The maximum number of model calls. |

### Runner

This repository does not bundle a model runner.\
Use a compatible runner that supports:

- the evaluation plan;
- isolated execution;
- direct observation of Skill loading for routing;
- the compact reporting contract.

The `mtk177a/skills` runner `scripts/run_skill_evaluation.py`, with its adjacent `evaluation_contract.py`, is one compatible implementation.\
It accepts this repository through `--root`.

Keep its plan, run, and grading artifacts in a temporary directory.\
Do not add that checkout as a runtime dependency of the pilot helper.

### Execution

- **Call count:** Review the number of model calls and make no unplanned calls.
- **Fixtures:** Run cases in disposable fixtures.
- **Isolation:** Isolate the intended Skill from personal copies and unrelated plugins.
- **Unverified isolation:** If isolation cannot be established, record the run as unverified.

### Grading

Grade only the selected case-condition pairs against their assertions.\
Distinguish `pass`, `fail`, `inconclusive`, and `error`.

For routing, grade the observed loading of the Skill.\
Do not use the final response's claim that the Skill was used as proof of loading.

### Additional evaluation

Add cases or rerun only when a specific remaining uncertainty could change the decision.

## Record results

### Compact record

Tie the result to the evaluated revision and input hashes.\
Include:

- **Stopping reason:** Why the run ended.
- **Unrun cases:** Which cases were not run.
- **Unverified boundaries:** What the run did not verify.

### Information handling

| Record | Where it belongs |
| --- | --- |
| Synthetic case IDs, counts, verdicts, hashes, and explanations with no private session material | May be included in the public evaluation record. |
| Raw prompts, outputs, traces, and disposable fixtures | Outside the public repository. |
| Private pilot evidence and model outputs from a real retrospective | Outside the public repository. |

### Authority

An evaluation result does not authorize any of the following:

- approving a proposal;
- editing an external harness;
- posting to an Issue.
