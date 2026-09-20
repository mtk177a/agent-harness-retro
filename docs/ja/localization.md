# 日本語参考訳

> これは [英語版 `localization.md`](../localization.md) の日本語参考訳です。\
> 英語版が正式版であり、内容に差異がある場合は英語版を優先します。

このリポジトリでは、英語の公開成果物を正式な情報源として扱います。\
日本語ファイルは、日本語で読みたい人のための参考訳として維持します。\
日本語版は、独立した方針、製品要件、インターフェース、エージェント向けの指示を定めません。

## 翻訳の対応関係

保守する英語版と日本語参考訳は、次の 7 組です。

| 正式な英語版 | 日本語参考訳 |
| --- | --- |
| [`README.md`](../../README.md) | [`README.ja.md`](../../README.ja.md) |
| [`AGENTS.md`](../../AGENTS.md) | [`AGENTS-ja.md`](../../AGENTS-ja.md) |
| [`docs/architecture.md`](../architecture.md) | [`docs/ja/architecture.md`](architecture.md) |
| [`docs/glossary.md`](../glossary.md) | [`docs/ja/glossary.md`](glossary.md) |
| [`docs/localization.md`](../localization.md) | [`docs/ja/localization.md`](localization.md) |
| [`skills/outer-harness-retrospective/SKILL.md`](../../skills/outer-harness-retrospective/SKILL.md) | [`skills/outer-harness-retrospective/SKILL-ja.md`](../../skills/outer-harness-retrospective/SKILL-ja.md) |
| [`.agents/skills/maintain-japanese-references/SKILL.md`](../../.agents/skills/maintain-japanese-references/SKILL.md) | [`.agents/skills/maintain-japanese-references/SKILL-ja.md`](../../.agents/skills/maintain-japanese-references/SKILL-ja.md) |

各日本語ファイルには、対応する英語版へのリンクと、内容に差異がある場合は英語版を優先する旨を記します。

## 保守規則

- 英語版の意味、要件、動作、範囲、その他の利用者に関係する内容が変わった場合は、日本語参考訳を更新します。
- 英語版にない要件、例外、例、判断、権限、解釈を追加しません。
- コード、識別子、コマンド名、製品名、リンク、その他の公開インターフェースの文字列は、翻訳によって意味や検索性が変わる場合、元の表記を正確に保持します。
- 英語版の変更が意味に影響しない場合、日本語ファイルは変更しません。\
  その場合は、日本語版の更新が不要な理由を Pull Request に記録しなければなりません。
- 英語版の意味や正しい対応関係が曖昧な場合は、その曖昧さを報告し、解決するまで参考訳を変更しません。\
  推測で訳してはいけません。
- 英語版の変更に参考訳の更新が必要な場合は、同じ変更内で英語版と日本語版を更新します。

リポジトリ専用の `maintain-japanese-references` Skill は、参考訳の確認と保守の手順を定めます。\
参考訳を更新または確認できますが、英語版の編集、Issue や Pull Request の執筆、ADR の翻訳、振り返り分析は行いません。

## 保守者の作業で使う言語

正式な公開成果物は、引き続き英語で作成します。\
保守する日本語参考訳と Issue・Pull Request のテンプレートは、この方針で明示する作業上の例外です。\
これらの例外は、正式版の言語を変えません。

保守者が作成する Issue と Pull Request には、英語のタイトルと、公開する変更を識別できる短い英語の `Summary` を付けなければなりません。

Issue と Pull Request のそれ以外の本文セクションとコメントは、日本語を既定とします。\
英語も使用でき、外部のコントリビューターに日本語で書くことは求めません。

コミットメッセージは、公開する変更を識別できる英語の件名を使うという保守者向け規則に従います。\
詳細は、必要に応じて日本語または英語で書けます。

## ADR

アーキテクチャ決定記録は、保守する日本語参考訳の対象ではありません。\
英語の ADR が正式版であり、この方針は対応する `ADR-ja.md` を要求しません。
