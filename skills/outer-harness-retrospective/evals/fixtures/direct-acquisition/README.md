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
| `source_refs` | Selected references: Codex and Claude for `case-18`, or the unsupported Claude source for `case-19`. |
| `max_events` | `40` across the selected sources. |
| `page_size` | `2`, to exercise pagination. |
| `harness` | `project/`. |
| `source_independence` | The two supported sources represent independent synthetic interactions. |

Keep assertions and grading artifacts outside the executor fixture.\
Do not run historical commands, use real session stores, or execute the target project's scripts.

## Coverage limits

The Codex legacy result bodies do not establish verified outcomes; the CLI reports `tool_outcome_unknown`.\
The fixture does not test repository discovery, real sessions, installation, or implicit Skill selection.\
Remove disposable stores and raw evaluation traces according to the evaluation plan.
