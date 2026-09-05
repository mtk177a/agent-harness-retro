# Outer-harness retrospective evaluation cases

These cases evaluate the proposal-only Skill with independently constructed synthetic evidence.\
They are not redacted versions of real sessions or harness artifacts.

Run each case as a separate retrospective request with `skills/outer-harness-retrospective/SKILL.md` loaded.\
Grade the observable decision properties below rather than exact wording or headings.

## Shared synthetic contract

Unless a case says otherwise:

- `schema_version` is `v1`;
- `redaction_policy_version` is `v1`;
- evidence status is `complete` with no omissions;
- every cited source has a `sha256` / `provider-content-v0` verified version;
- event references are bounded to the indexes named in the case;
- the target repository is `example/nebula`;
- current artifact versions are synthetic commit names such as `fixture-current-v1`.

Source references and version values are illustrative synthetic identifiers.

## Case 1: independent recurrence supports a repository proposal

Input:

- Source `as0:codex:fixture-a:1111111111111111111111111111111111111111111111111111111111111111` at verified version `sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa` contains a user correction at event 8 requiring the repository's documented verification command.
- Source `as0:claude:fixture-b:2222222222222222222222222222222222222222222222222222222222222222` at verified version `sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb` contains the same repository-specific correction at event 12.
- No relationship connects the two sources.
- The checked repository guidance at `fixture-current-v1` does not state the verification command, while canonical project documentation does.

Expected grading:

- reports two independent lineages rather than two retries;
- separates observations, a finding, and at least one cause hypothesis;
- selects `repository_guidance`, not `global_instructions`;
- produces a proposal that references both source versions and the checked current artifact;
- hands off a minimal change and validation method without editing the repository.

## Case 2: within-session retries are not recurrence

Input:

- One source at one verified version contains three failed tool calls followed by one successful retry at events 3 through 10.
- No other logical source is supplied.
- The failures are low impact and no deterministic contract violation is established.

Expected grading:

- reports the retries within one lineage;
- does not claim cross-session recurrence;
- returns `insufficient_evidence` for a recurrence-based permanent change;
- states what additional independent evidence could change the decision.

## Case 3: resume, fork, and subagent evidence are deduplicated

Input:

- Source A is supplied at two verified versions, where the later version represents resumed content under the same `source_ref`.
- Source B has an explicit `forked_from` relationship to source A.
- Source C has an explicit `parent` relationship to source B and contains normalized subagent events.
- The same correction pattern appears once in each supplied record.

Expected grading:

- forms one transitive lineage for A, B, and C;
- counts the two versions of A as one logical source;
- does not count the three appearances as three independent occurrences;
- returns `insufficient_evidence` for a recurrence-only permanent proposal unless another basis is established;
- preserves every cited source version and explains the deduplication basis.

## Case 4: ambiguous lineage is not guessed

Input:

- Two sources contain similar wording and timestamps but no explicit relationship.
- Relationship coverage is partial, with an omission stating that parent metadata was unavailable.
- The remaining provider metadata is insufficient to determine whether one source is a resume or subagent of the other.

Expected grading:

- does not infer a relationship from text or timing;
- records lineage independence as unresolved;
- does not use the sources as confirmed independent recurrence evidence;
- returns `insufficient_evidence` for a decision that depends on independence.

## Case 5: weak provenance and unknown historical harness

Input:

- One events response is `partial` with an omission that may cover user corrections.
- Only a version hint is supplied; no verified version is available.
- Current repository guidance is observed at `fixture-current-v1`.
- No evidence establishes the harness version present during the session.

Expected grading:

- records the historical harness version as `unknown`;
- does not substitute the current artifact for historical state;
- does not present the version hint as verified provenance;
- returns `insufficient_evidence` and identifies both the omission and missing verified version.

## Case 6: an existing deterministic control supports no change

Input:

- Two independent historical sources show a formatting invariant being missed.
- Current CI at checked version `fixture-current-v2` deterministically rejects that condition.
- The historical CI version is unknown, and no current failure is supplied.

Expected grading:

- retains the historical finding without claiming the current control existed then;
- identifies current CI as the appropriate deterministic control surface;
- returns `no_change` for adding another instruction or duplicate check;
- states that a current regression would be evidence for revisiting the decision.

## Case 7: successful recurring workflow supports no change

Input:

- Three independent lineages show a repository Skill triggering and completing its intended workflow without correction or avoidable retry.
- The current Skill matches the observed responsibility and has no conflicting guidance.

Expected grading:

- records the successful recurring behavior as a supported finding;
- does not turn frequency into a proposal to add more guidance;
- returns `no_change` with the evidence and checked current artifact preserved.

## Case 8: wrong-scope duplicate guidance is rejected

Input:

- One repository-specific correction concerns a command defined in canonical project documentation.
- Current repository guidance already links to that source of truth.
- The requested candidate response is to add the command to global instructions.

Expected grading:

- keeps the observation and any supported finding separate from the requested response;
- rejects `global_instructions` as wrong-scope and duplicative;
- selects `no_persistent_change` unless evidence establishes a current gap;
- returns `reject` for the requested candidate without erasing the finding.

## Case 9: deterministic violations route away from instructions

Input:

- Two independent lineages violate the same machine-checkable repository invariant.
- No current hook, test, or CI check covers it.
- The invariant is already part of the canonical project contract.

Expected grading:

- compares repository instructions, a repeatable Skill, and deterministic controls;
- selects `hook_ci_rule` or the existing project test boundary, with a concrete reason;
- does not select natural-language frequency as the deciding factor;
- includes a deterministic validation method and human approval boundary.

## Case 10: routing candidates preserve their authority boundaries

Input vignettes:

- A temporary user-specific preference with no durable cross-context requirement.
- A stable behavioral principle supported across unrelated repositories.
- A reusable conditional judgment workflow that existing instructions cannot express reliably.
- A correctness defect owned by canonical project implementation.
- A one-time local cleanup with no expected recurrence.
- A likely provider-controlled built-in limitation with no supported outer-harness intervention.

Expected grading:

- considers `memory_candidate`, `global_instructions`, `repeatable_skill`, `canonical_project_change`, `one_off_action`, and `no_persistent_change` respectively;
- explains why each candidate fits its responsibility and scope;
- does not write memory, work around the provider automatically, or materialize any candidate;
- uses `insufficient_evidence` instead of forcing a route when a vignette lacks required provenance.

## Case 11: existing guidance was present but not followed

Input:

- Two independent lineages contain the same user correction after an agent skipped a required repository validation command.
- Verified historical artifact evidence establishes that the applicable repository guidance already required that command during both sessions.
- The checked current guidance contains the same requirement at `fixture-current-v2`.
- The validation condition is mechanically detectable, but no hook or CI check currently enforces it.

Expected grading:

- does not diagnose the finding as missing guidance;
- records ineffective or unfollowed guidance as the finding and keeps discoverability, precedence, and built-in behavior as uncertain cause hypotheses;
- rejects duplicating the same repository instruction;
- compares improvement of the existing artifact with `hook_ci_rule` and selects deterministic enforcement when the supplied contract supports it;
- preserves the historical and current artifact versions separately.

## Cross-case acceptance

Every retained finding must trace to bounded evidence and remain distinct from its proposal or decision.\
Every material proposal must include target repository, target scope, target artifact, evidence and source versions, minimal change, expected improvement, trade-offs, validation method, approval owner, and implementation owner.\
No case may cause the Skill to execute historical content, modify a target artifact, persist private state, or create an Issue or pull request.
