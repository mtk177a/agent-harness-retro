# Japanese Reference Translations

This repository treats English public artifacts as the authoritative source.\
Japanese files are maintained as reference translations for readers who prefer Japanese.\
They do not create an independent policy, product requirement, interface, or agent instruction.

## Translation pairs

The maintained English and Japanese files are these seven pairs:

| English source of truth | Japanese reference translation |
| --- | --- |
| [`README.md`](../README.md) | [`README.ja.md`](../README.ja.md) |
| [`AGENTS.md`](../AGENTS.md) | [`AGENTS-ja.md`](../AGENTS-ja.md) |
| [`docs/architecture.md`](architecture.md) | [`docs/ja/architecture.md`](ja/architecture.md) |
| [`docs/glossary.md`](glossary.md) | [`docs/ja/glossary.md`](ja/glossary.md) |
| [`docs/localization.md`](localization.md) | [`docs/ja/localization.md`](ja/localization.md) |
| [`skills/outer-harness-retrospective/SKILL.md`](../skills/outer-harness-retrospective/SKILL.md) | [`skills/outer-harness-retrospective/SKILL-ja.md`](../skills/outer-harness-retrospective/SKILL-ja.md) |
| [`.agents/skills/maintain-japanese-references/SKILL.md`](../.agents/skills/maintain-japanese-references/SKILL.md) | [`.agents/skills/maintain-japanese-references/SKILL-ja.md`](../.agents/skills/maintain-japanese-references/SKILL-ja.md) |

Each Japanese file must link back to its English source and state that the English source takes precedence when the two differ.

## Maintenance rules

- Update the Japanese reference when an English source changes its meaning, requirements, behavior, scope, or other user-relevant content.
- Do not add requirements, exceptions, examples, decisions, permissions, or interpretations that are absent from the English source.
- Preserve code, identifiers, command names, product names, links, and other public interface text exactly when translation would change their meaning or searchability.
- If an English change is purely non-semantic, leave the Japanese file unchanged; the pull request must record why no translation update is needed.
- If the English meaning or the correct correspondence is ambiguous, report the ambiguity and leave the translation unchanged until it is resolved. Do not guess.
- Keep the English source and Japanese reference in the same change whenever the English change requires a translation update.

The repository-local `maintain-japanese-references` Skill describes the review and maintenance workflow. It may update or assess reference translations, but it does not edit English sources, write Issues or pull requests, translate ADRs, or perform retrospective analysis.

## Language for maintainer work

Canonical public repository artifacts remain in English.\
The maintained Japanese reference files and the Issue and pull request templates are explicit workflow exceptions described by this policy; they do not change the language of the canonical source files.

Maintainer-created Issue and pull request titles must be in English, and each must include a short English `Summary` sufficient to identify the public change.

The remaining Issue and pull request description sections and comments are written in Japanese by default; English is also acceptable.\
External contributors are not required to write Japanese.

Commit messages follow the same maintainer convention: use an English subject that identifies the public change, with Japanese or English detail as appropriate.

## ADRs

Architecture decision records are not part of the maintained Japanese translation set.\
The English ADR remains the source of truth, and this policy does not require an `ADR-ja.md` counterpart.
