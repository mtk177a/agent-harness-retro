---
name: agent-harness-retro-setup
description: Prepares local access to Codex and Claude Code session records for harness retrospectives. Use when the user explicitly requests initial setup or wants to check or configure session access before a retrospective.
---

# Set up Agent Harness Retro

Work in the same local execution environment as the retrospective.

## 1. Check the CLI

- **Find the executable** — Run `agent-sessions --version` or check the absolute executable path the user supplied. Use the executable built from the commit pinned in [the installer](scripts/install-agent-sessions.sh). Its version ends in `-8bc9542a15ed`.
- **Install if needed** — Confirm the installation directory with the user. On macOS, Git and Go 1.26 or later are required. Run `sh <skill-directory>/scripts/install-agent-sessions.sh <absolute-bin-directory>`. If the target already exists, ask the user to choose another directory. Check the installed binary with `--version` and retain its absolute path.

## 2. Select the stores

Use the stores the user selected, or ask which Codex or Claude Code stores to use. Confirm each provider and absolute root before running discovery.

1. **Existing configuration** — Use an explicit `--config` path supplied by the user, or check `~/Library/Application Support/agent-sessions/config.json` on macOS. Keep each confirmed entry's provider, ID, and root together.
2. **Environment and defaults** — For a provider without a selected entry, check the first available setting in the table and confirm the resulting root with the user.

   | Provider | Environment variable | Default root |
   | --- | --- | --- |
   | Codex | `CODEX_HOME` | `~/.codex` |
   | Claude Code | `CLAUDE_CONFIG_DIR` | `~/.claude` |

3. **Custom roots** — Check the selected path with `list --limit 1 --provider <provider> --root <absolute-root>`. Include the selected config path if one is in use. If the user wants to save a working root, follow [Save a custom root](#save-a-custom-root).

## 3. Check access

Use the selected executable to run `list --limit 1 --provider <provider>` for each confirmed store. For a configured store, pass the config path and the verified `--source-instance`. For an environment or default root, pass the confirmed `--root`. Use the same source settings in the retrospective.

Check the exit status and the JSON `status`, `data.sources`, and `omissions`. Report an empty successful listing separately from an error or unsupported result. Report any coverage gaps separately from pagination. For a failed check, state which setup step is still needed.

Return the executable path, config path if used, and provider and instance ID or root for each store that can be used. Setup ends after listing. A later retrospective starts with the user's choice of evidence and review scope.

### Save a custom root

1. **Check the file** — If an existing file cannot be parsed or does not match the [configuration contract](https://github.com/mtk177a/agent-sessions/blob/8bc9542a15ed85f0b9b341886e920c76b1572d20/docs/cli-json-contract.md#configuration-structure), report the problem and leave it untouched.
2. **Confirm the entry** — Reuse an identical entry or select an unused ID for the provider. Confirm the configuration path and proposed `sources` entry with the user before writing.
3. **Save the entry** — Preserve existing entries and their IDs. Create a new config with owner-only permissions. Use the saved instance in step 3.
