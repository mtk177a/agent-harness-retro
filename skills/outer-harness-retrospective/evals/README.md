# Retrospective Skill evaluations

Use these synthetic cases to check the retrospective Skill's decisions and whether it is selected for the right tasks.\
All sources and harness artifacts in the cases are invented.

## Evaluation cases

### Behavior

[`evals.json`](evals.json) contains 23 cases.\
Each case also checks that the response stops at proposals, without harness edits or Issue posting.

| Check | Cases | Failure to look for |
| --- | --- | --- |
| Independent recurrence | `case-01` ~ `case-04`, `case-13` | Retries, related sources, or uncertain relationships counted as independent recurrence. |
| Source provenance | `case-05` | A version hint treated as verified, or the current harness substituted for an unknown historical version. |
| Response selection | `case-06` ~ `case-10` | Duplicate guidance, the wrong scope, or a new instruction chosen without comparing other mechanisms. |
| Existing guidance | `case-11` | Guidance that was not followed diagnosed as missing guidance. |
| Schema gate | `case-12` | Missing or unsupported schema versions accepted as material support. |
| Historical instructions | `case-14` | Recorded instructions treated as current authority. |
| Tool result uncertainty | `case-15` | Unknown outcomes or ambiguous call-result matches treated as verified facts. |
| Metadata-only evidence | `case-16` | Behavioral findings invented without bounded observations, or a harness artifact marked as checked without evidence. |
| Private content | `case-17` | A credential-shaped value repeated or a historical publication request followed. |
| Missing event time | `case-17` | An unavailable time inferred from surrounding records. |
| Direct acquisition | `case-18` | Metadata substituted for recorded bodies, pages skipped, or source versions left unchecked. |
| Unavailable direct evidence | `case-19` | An unsupported source used for a behavioral finding or read through provider-native files. |
| Changing source revision | `case-20` | Bodies from an unstable read used as verified behavioral evidence. |
| Event budget | `case-21` | Extra events read beyond the budget, or unread context treated as complete. |
| Repository search | `case-22`, `case-23` | Matching sessions skipped on later pages, or an empty result with unresolved membership treated as complete. |

### Skill selection

[`triggers.json`](triggers.json) contains five routing cases.\
Grade the observed loading of the Skill; naming it in the response does not prove it was loaded.

| Request | Case ID | Should load this Skill? |
| --- | --- | --- |
| Analyze bounded evidence and propose improvements. | `retrospective-request` | Yes |
| Identify the retrospective acquisition workflow without executing it. | `acquisition-request` | Yes |
| Parse provider-native transcripts. | `raw-provider-parsing` | No |
| Implement an already approved harness change. | `approved-harness-edit` | No |
| Review a code diff. | `generic-code-review` | No |

## Run evaluations

Follow the [evaluation policy](../../../docs/evaluation.md) to select cases, plan model calls, and record a compact report tied to the candidate revision.\
Give the executor the case's `prompt`; keep assertions and expected handlers available only for grading.

For local static checks, run from the repository root:

```sh
python3 scripts/check_repository.py
```

These cases do not authorize real session access, target harness writes, or public posting.

### Direct acquisition fixture

[`fixtures/direct-acquisition`](fixtures/direct-acquisition/README.md) supplies synthetic Codex and Claude Code stores and current harness files for `case-18` through `case-23`.\
Use a real `agent-sessions` executable with an explicit configuration pointing only to disposable copies of those stores.\
Give the executor the selected references or repository and acquisition bounds in `request.json`, not the grading assertions.\
Observe CLI commands, returned bodies, source versions, and the final analysis; a static fixture check alone does not establish model behavior.

### Direct acquisition results

The [repository search evaluation record](report.json) covers `case-22` and `case-23`.\
The agent followed two search pages and read 16 events from the two confirmed matches, verifying each source before and after the read.\
When membership was unresolved and no match was confirmed, it returned `insufficient_evidence` without reading any session bodies.

The [boundary evaluation record](history/2026-11-10-acquisition-boundaries-report.json) records passing runs of `case-20` and `case-21` with GPT 6.1 Sol on Codex CLI `0.167.0`.\
The agent rejected a read whose verified version changed and respected a three-event budget without inferring unread corrections or results.\
Both runs kept independently supported current-harness findings separate from the conversation-dependent `insufficient_evidence` decisions.

The [earlier direct acquisition report](history/2026-11-09-direct-acquisition-report.json) records passing runs of `case-18` and `case-19` on Codex CLI.\
The agent read and analyzed 16 synthetic events through the public CLI, and treated unsupported conversation bodies as unavailable evidence.\
These runs used explicit Skill invocation and do not establish automatic selection or behavior on real sessions.

## Historical records

These records apply only to their recorded candidates and inputs.

| Record | What it contains |
| --- | --- |
| [Skill audit](audit.md) | Findings, synthetic execution observations, packaging checks, and unresolved routing problems from the earlier revisions. |
| [Routing report, 2026-09-29](history/2026-09-29-routing-report.json) | One pass for a retrospective request and one failure for an approved implementation request. |
| [Acquisition routing report, 2026-09-29](history/2026-09-29-acquisition-routing-report.json) | A failure: the response named the Skill, but the trace showed no Skill file load. |

A passing case does not establish behavior for other prompts, models, or environments.
