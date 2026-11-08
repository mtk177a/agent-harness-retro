# Reporting requirements

Adapt the report to the request, keeping direct observations, findings, cause hypotheses, and proposed solutions distinct.

## Decision record

- **Scope:** State the target, evidence bounds, and material coverage limits.
- **Findings:** Explain the supported pattern and impact; separate cause hypotheses and their uncertainty from observations, and keep solutions out of findings.
- **Evidence:** Use stable references linking observations, findings, and decisions to `source_ref`, verified versions, bounded events, and the relevant checked current harness artifacts and versions.
- **Decision:** Give the outcome, selected responsibility scope and response mechanism, reasons for selection or rejection, and unresolved information that could change the decision.

For each proposal, include:

- **Change:** The target repository or global artifact, scope, and minimal intended edit.
- **Rationale:** Supporting findings and evidence, expected improvement, and material tradeoffs.
- **Handoff:** Validation method and the approval and implementation owners; use owner roles when no person is supplied, identifying the global artifact owner for a global proposal.

## Response mechanisms

| Code | Use when |
| --- | --- |
| `memory_candidate` | Temporary or user-specific knowledge needs a handoff; classification does not authorize a memory write or a custom store. |
| `global_instructions` | A stable behavioral principle is supported across materially different repository contexts. |
| `repository_guidance` | Durable agent-facing knowledge is specific to one repository. |
| `repeatable_skill` | A conditional reusable judgment or multi-step workflow is missing or misrouted. |
| `hook_ci_rule` | A condition can and should be enforced deterministically. |
| `canonical_project_change` | Correctness or knowledge belongs in project implementation, tests, or canonical documentation. |
| `one_off_action` | A local, transient response does not justify a durable mechanism. |
| `no_persistent_change` | Existing controls suffice, the behavior is provider-controlled, or durable change is not justified. |
