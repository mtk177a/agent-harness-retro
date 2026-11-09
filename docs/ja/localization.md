# 言語方針

> これは [英語版 `localization.md`](../localization.md) の日本語参考訳である。
> 英語版が正式版であり、内容が異なる場合は英語版を優先する。

## リポジトリの成果物

| 成果物 | 言語 | 役割 |
| --- | --- | --- |
| 正式な公開成果物 | 英語 | 正式な情報源。 |
| 保守する参考訳 | 日本語 | 内容が異なる場合は英語版を優先する。 |

日本語参考訳は、英語版とは別の方針、製品要件、インターフェース、エージェント向け指示を定めない。

英語を既定とし、日本語参考訳と Issue・Pull Request のテンプレートを例外として扱う。

## Issue と Pull Request

保守者が作成する公開 Issue と Pull Request では、次の言語を使う。

| 箇所 | 言語 |
| --- | --- |
| タイトル | 英語 |
| `Summary` | 英語 |
| その他の本文セクション | 日本語が既定。英語も使用できる。 |
| コメント | 日本語が既定。英語も使用できる。 |

`Summary` を設け、公開する変更を識別できる短い内容にする。

外部のコントリビューターに日本語で書くことは求めない。

## コミット

- **件名**: 公開する変更を識別できる英語の件名にする。
- **詳細**: 必要に応じて、日本語または英語で書く。

## 日本語参考訳の保守

### 保守対象

| 正式な英語版 | 日本語参考訳 |
| --- | --- |
| [`README.md`](../../README.md) | [`README.ja.md`](../../README.ja.md) |
| [`AGENTS.md`](../../AGENTS.md) | [`AGENTS-ja.md`](../../AGENTS-ja.md) |
| [`docs/architecture.md`](../architecture.md) | [`docs/ja/architecture.md`](architecture.md) |
| [`docs/evaluation.md`](../evaluation.md) | [`docs/ja/evaluation.md`](evaluation.md) |
| [`docs/glossary.md`](../glossary.md) | [`docs/ja/glossary.md`](glossary.md) |
| [`docs/localization.md`](../localization.md) | [`docs/ja/localization.md`](localization.md) |
| [`skills/outer-harness-retrospective/SKILL.md`](../../skills/outer-harness-retrospective/SKILL.md) | [`skills/outer-harness-retrospective/SKILL-ja.md`](../../skills/outer-harness-retrospective/SKILL-ja.md) |
| [`skills/outer-harness-retrospective/references/lineage.md`](../../skills/outer-harness-retrospective/references/lineage.md) | [`skills/outer-harness-retrospective/references/ja/lineage.md`](../../skills/outer-harness-retrospective/references/ja/lineage.md) |
| [`skills/outer-harness-retrospective/references/reporting.md`](../../skills/outer-harness-retrospective/references/reporting.md) | [`skills/outer-harness-retrospective/references/ja/reporting.md`](../../skills/outer-harness-retrospective/references/ja/reporting.md) |
| [`.agents/skills/maintain-japanese-references/SKILL.md`](../../.agents/skills/maintain-japanese-references/SKILL.md) | [`.agents/skills/maintain-japanese-references/SKILL-ja.md`](../../.agents/skills/maintain-japanese-references/SKILL-ja.md) |

ADR は、この保守対象に含めない。\
英語の ADR を正式版とし、対応する日本語版の作成は求めない。

### 翻訳の要件

- **正式版の案内**: 英語版へのリンクを付け、内容が異なる場合は英語版を優先すると明記する。
- **意味**: 英語版にない要件、例外、例、判断、権限、解釈を追加しない。
- **公開インターフェースの文字列**: 翻訳によって意味や検索性が変わる場合は、コード、識別子、コマンド、製品名、リンク、その他の公開インターフェースの文字列を保持する。

### 更新時の判断

保守対象の英語版を追加・変更するときは、リポジトリ専用の [`maintain-japanese-references` Skill](../../.agents/skills/maintain-japanese-references/SKILL-ja.md) を使って、日本語参考訳を作成または確認する。

| 英語版の変更 | 必要な対応 |
| --- | --- |
| 意味、要件、動作、範囲、その他の利用者に関係する内容が変わる | 同じ変更内で日本語参考訳を更新する。 |
| 意味が変わらない | 日本語参考訳は変更せず、その理由を Pull Request に記録する。 |
| 英語版の意味や、訳文との対応が曖昧である | 曖昧な点を報告し、解決するまで参考訳を変更しない。推測で訳さない。 |
