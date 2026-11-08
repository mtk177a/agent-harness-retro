# Language Policy

## Repository artifacts

| Artifact | Language | Role |
| --- | --- | --- |
| Canonical public repository artifacts | English | Source of truth. |
| Maintained reference translations | Japanese | The English source takes precedence if the content differs. |

Japanese references do not define independent policies, product requirements, interfaces, or agent instructions.

Japanese reference files and Issue and pull request templates are exceptions to the default English language.

## Issues and pull requests

For maintainer-created public Issues and pull requests, use these languages:

| Part | Language |
| --- | --- |
| Title | English |
| `Summary` | English |
| Other body sections | Japanese by default; English is also allowed. |
| Comments | Japanese by default; English is also allowed. |

Include a short `Summary` that identifies the public change.

External contributors are not required to write Japanese.

## Commits

- **Subject**: Use English to identify the public change.
- **Details**: Use Japanese or English as appropriate.

## Maintaining Japanese references

### Maintained files

| English source of truth | Japanese reference translation |
| --- | --- |
| [`README.md`](../README.md) | [`README.ja.md`](../README.ja.md) |
| [`AGENTS.md`](../AGENTS.md) | [`AGENTS-ja.md`](../AGENTS-ja.md) |
| [`docs/architecture.md`](architecture.md) | [`docs/ja/architecture.md`](ja/architecture.md) |
| [`docs/evaluation.md`](evaluation.md) | [`docs/ja/evaluation.md`](ja/evaluation.md) |
| [`docs/glossary.md`](glossary.md) | [`docs/ja/glossary.md`](ja/glossary.md) |
| [`docs/localization.md`](localization.md) | [`docs/ja/localization.md`](ja/localization.md) |
| [`skills/outer-harness-retrospective/SKILL.md`](../skills/outer-harness-retrospective/SKILL.md) | [`skills/outer-harness-retrospective/SKILL-ja.md`](../skills/outer-harness-retrospective/SKILL-ja.md) |
| [`skills/outer-harness-retrospective/references/lineage.md`](../skills/outer-harness-retrospective/references/lineage.md) | [`skills/outer-harness-retrospective/references/ja/lineage.md`](../skills/outer-harness-retrospective/references/ja/lineage.md) |
| [`skills/outer-harness-retrospective/references/reporting.md`](../skills/outer-harness-retrospective/references/reporting.md) | [`skills/outer-harness-retrospective/references/ja/reporting.md`](../skills/outer-harness-retrospective/references/ja/reporting.md) |
| [`.agents/skills/maintain-japanese-references/SKILL.md`](../.agents/skills/maintain-japanese-references/SKILL.md) | [`.agents/skills/maintain-japanese-references/SKILL-ja.md`](../.agents/skills/maintain-japanese-references/SKILL-ja.md) |

ADRs are outside this maintained set.\
The English ADR is authoritative; a Japanese counterpart is not required.

### Translation requirements

- **Source notice**: Link to the English source and state that it takes precedence if the content differs.
- **Meaning**: Do not add requirements, exceptions, examples, decisions, permissions, or interpretations that are absent from the English source.
- **Public interface text**: Preserve code, identifiers, commands, product names, links, and other public interface strings when translation would change their meaning or searchability.

### When to update

When adding or changing a maintained English file, use the repository-local [`maintain-japanese-references` Skill](../.agents/skills/maintain-japanese-references/SKILL.md) to create or review its Japanese reference.

| English source change | Required action |
| --- | --- |
| Meaning, requirements, behavior, scope, or other user-relevant content changes | Update the Japanese reference in the same change. |
| Meaning does not change | Leave the Japanese reference unchanged and record the reason in the pull request. |
| The source meaning or its correspondence with the translation is unclear | Report the ambiguity and leave the reference unchanged until it is resolved; do not guess. |
