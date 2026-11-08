# Outer harness retrospective evaluations

The behavior cases in `evals.json` preserve the 14 original synthetic decision boundaries and add three `agent-sessions` v2 boundaries.
The routing cases in `triggers.json` distinguish retrospective analysis from provider parsing, implementation, and general code review.
All cases use invented sources and harness artifacts.
The executor receives each case's `prompt`; assertions and expected handlers remain grading data.

| Responsibility | Failure exposed by | Cases | Grade |
| --- | --- | --- | --- |
| Independent recurrence and lineage | Retry or fork counted as recurrence | `case-01` to `case-04`, `case-13` | Inspect lineage basis and decision |
| Provenance and historical state | Missing version or old harness inferred from current one | `case-05`, `case-11`, `case-12` | Inspect evidence gate and historical version |
| Scope and control surface | Duplicate or wrong-scope guidance | `case-06` to `case-10` | Inspect considered routes and outcome |
| Untrusted content and proposal boundary | Historical instruction treated as authority | `case-14`, `case-17` | Inspect actions and report content |
| v2 uncertainty | Unknown outcome or unmatched result treated as fact | `case-15` | Inspect supported conclusion |
| Content-free pilot handoff | Metadata treated as behavioral evidence | `case-16` | Require a bounded reviewed observation |
| Skill discovery | Acquisition request misses the Skill or another task invokes it | `triggers.json` | Observe actual Skill loading |

Use [`docs/evaluation.md`](../../../docs/evaluation.md) to select and plan a change-scoped run.
Run `python3 scripts/check_repository.py` from the repository root for the local static checks.
The cases do not authorize access to a real session, writes to a target harness, or public posting.

The audit and the earlier revision-bound routing result are in [`audit.md`](audit.md) and [`history/2026-09-29-routing-report.json`](history/2026-09-29-routing-report.json).\
The focused acquisition routing result for the Skill revision evaluated on 2026-09-29 is archived in [`history/2026-09-29-acquisition-routing-report.json`](history/2026-09-29-acquisition-routing-report.json); naming the Skill without loading it did not pass the observed-loading requirement.
Model observations belong in a compact report tied to the candidate revision; a passing case cannot establish behavior across other prompts, models, or environments.
