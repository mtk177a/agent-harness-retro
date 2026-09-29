# Evaluation and Skill Audit

This policy governs evaluation of the repository's maintained Skills and changes that affect retrospective decisions.
It adapts the change-scoped evaluation rules used by [`mtk177a/skills` Issue #49](https://github.com/mtk177a/skills/issues/49) to this repository's public and private evidence boundary.

## Audit decision

For each material Skill audit, record the target revision and working-tree state, the intended behavior, affected users and client, the guidance surfaces and their authority, observed evidence, findings, alternatives, and remaining uncertainty.
Check responsibility and scope, selection triggers and coexistence, source precedence, evidence quality, safety, context cost, portability, and maintenance only where they matter to the Skill.
Check the English canonical Skill and the naturalness and semantic fidelity of its Japanese reference independently, even if the English text did not change.
Classify each issue as `change_required`, `no_change`, `separate_issue`, or `unverified`; support every material finding with evidence, impact, confidence, and a way to falsify it.
Do not turn a static audit into a claim about model behavior.

## Select evaluation work

Identify the changed responsibility, a plausible failure, and the observation that would detect it before selecting cases.
Use the least evidence needed for that decision:

| Change | Default path |
| --- | --- |
| Documentation, format, or meaning-preserving text | `static-only` |
| Skill instructions, output contract, or safety boundary | `targeted-candidate` |
| Name, description, invocation, or adjacent responsibility boundary | `targeted-routing` |
| Known regression, major redesign, or ambiguous candidate result | `baseline-comparison` |
| Explicit client or environment support claim | `target-environment` |

Use a baseline, a no-Skill comparison, repetition, or another client only if its result could change the decision.
Do not run the complete suite after every edit.
When existing `evals.json` or `triggers.json` definitions are migrated, migrate the whole affected set together; format migration alone needs static validation, not a mandatory model run.

## Executable definitions

Keep behavior cases in `skills/<name>/evals/evals.json` and routing cases in `triggers.json`.
Use the Agent Skills executable envelope with `skill_name` and `evals` and the case fields `id`, `prompt`, optional `files`, `expected_output`, and `assertions`; routing cases use `expected_handlers`.
The case prompt and fixture contain only executor input.
Assertions and expected output are hidden from the executor and used for grading.
Use synthetic identities, repositories, harness artifacts, and sessions; never commit a real session, even after manual redaction.
Keep raw prompts, outputs, traces, and disposable fixtures outside the public repository.
Run `python3 scripts/check_repository.py` to validate the maintained pairs, Skill metadata, executable definitions, and any current compact report without a model call.
For changes to the bundled acquisition scripts, run `python3 -m unittest discover -s skills/outer-harness-retrospective/scripts -p 'test_*.py'` and `go test ./...` from `skills/outer-harness-retrospective/scripts/retro-pilot`.\
Verify a pinned source build without reading real sessions.

## Run and record

Before any model call, write an evaluation plan naming the candidate revision and file hashes, affected responsibility, path, selected cases and conditions, client, model, sandbox, and maximum model calls.
The repository does not bundle a model runner; a compatible runner must implement this plan, isolated execution, direct routing observations, and compact reporting contract.\
The `mtk177a/skills` `scripts/run_skill_evaluation.py` runner with its adjacent `evaluation_contract.py` is one compatible implementation and accepts this repository through `--root`.\
Use a temporary directory for its plan, run, and grading artifacts; do not add that checkout as a runtime dependency of the pilot helper.
Review the call count and make no unplanned calls.
Run cases in disposable fixtures with the intended Skill isolated from personal copies and unrelated plugins; if isolation cannot be established, record the run as unverified.
For routing, grade observed Skill loading, not a final response's claim that a Skill was used.
Grade only the selected case-condition pairs against their assertions, and distinguish `pass`, `fail`, `inconclusive`, and `error`.
Keep the compact result tied to the evaluated revision and input hashes; record the stopping reason, unrun cases, and unverified boundaries.
Only add cases or rerun when a specific remaining uncertainty would change the decision.

The public evaluation record may include synthetic case IDs, counts, verdicts, hashes, and explanations that contain no private session material.
Private pilot evidence and model outputs from a real retrospective stay outside the repository.
No evaluation result grants authority to approve a proposal, edit an external harness, or post to an Issue.
