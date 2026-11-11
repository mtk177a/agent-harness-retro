# Direct acquisition fixture

All identities, paths, conversations, tool activity, and harness files in this fixture are synthetic.

| Directory | Use |
| --- | --- |
| `codex-home` | Codex `0.149.1` legacy records. |
| `claude-home` | Claude Code `2.1.228` records. |
| `unsupported-home` | The synthetic Claude conversation labeled with an unverified `9.9.9` format. |
| `project` | Current repository harness artifacts to inspect. |

## Prepare a run

1. **Copy the inputs**: Copy the provider stores and `project` into a disposable fixture outside this repository.
2. **Provide the tool**: Use an available `agent-sessions` executable and record its source revision and binary hash.
3. **Configure the sources**: Write an explicit `agent-sessions` configuration whose roots point only to those copies.
4. **Select references**: Use the public `list` operation to obtain source references for the selected case.
5. **Supply the request**: Write `request.json` using the fields below and run the selected evaluation prompt under the repository's [evaluation policy](../../../../../docs/evaluation.md).

| Request field | Value |
| --- | --- |
| `agent_sessions_bin` | Absolute path to the available executable. |
| `config` | Absolute path to the explicit configuration. |
| `source_refs` | Codex and Claude for `case-18`; unsupported Claude for `case-19`; supported Claude only for `case-20` and `case-21`. |
| `max_events` | `40` for `case-18` and `case-19`; `8` for `case-20`; `3` for `case-21`, across all reads including retries. |
| `page_size` | `2`, to exercise pagination. |
| `harness` | `project/`. |
| `source_independence` | The two supported sources represent independent synthetic interactions. |

Keep assertions and grading artifacts outside the executor fixture.\
Do not run historical commands, use real session stores, or execute the target project's scripts.

## Acquisition boundary conditions

Use separate disposable copies for `case-20` and `case-21`.\
The executor receives the ordinary review prompt and `request.json`; keep the conditions and assertions below outside its fixture.

| Case | Controller action | Observation |
| --- | --- | --- |
| `case-20` | After the first completed `events` operation, replace the final synthetic Claude assistant text with `The synthetic record was revised during review.` before the next operation. | Real CLI verification before and after the read returns different versions; the event count remains eight. |
| `case-21` | Leave the source unchanged and set a total budget of three events with page size two. | The second page can return only one event; five events remain unread. |

The evaluation controller changes only the disposable synthetic store.\
It must not alter the executable, substitute responses, or tell the executor that a revision mismatch is expected.\
Observe each public CLI operation separately so that the source mutation happens between operations, not during a single CLI read.\
Record the mutation checkpoint, real CLI responses, requested limits, returned event indexes and final judgments.\
Current harness contradictions may support independent findings even when the selected conversation cannot support the requested conclusion.

## Repository search conditions

For `case-22` and `case-23`, use a build from the merged [repository search change](https://github.com/mtk177a/agent-sessions/pull/27).\
Initialize the disposable `project` as a Git repository and commit its synthetic files.\
Supply its absolute path as `repository` in `request.json` instead of supplying `source_refs`.\
Set `search_page_size` to `1`, `page_size` to `2`, and `max_events` to `40`.

| Case | Prepare the disposable stores |
| --- | --- |
| `case-22` | Put the Codex session in a linked worktree by setting its header `cwd` to that worktree. Set the supported Claude session's row `cwd` to `project`. Add another Codex session with its own ID and a `cwd` in an independent repository. Configure both stores. |
| `case-23` | Configure only the supported Claude store. Set every recorded `cwd` in it to an absent directory with no matching worktree registration. Set `provider` to `claude` in the request and pass `--provider claude` on every search page. |

Check the real `list --repository` responses before evaluating, and keep their expected counts and source references outside the executor input.\
Record the search cursors, membership omissions, and any sources opened for body reads.

## Coverage limits

The Codex legacy result bodies do not establish verified outcomes; the CLI reports `tool_outcome_unknown`.\
The fixture uses synthetic sessions and a supplied executable with explicit Skill invocation.\
Remove disposable stores and raw evaluation traces according to the evaluation plan.
