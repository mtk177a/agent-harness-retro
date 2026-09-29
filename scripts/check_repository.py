#!/usr/bin/env python3
"""Check this repository's maintained Skill, evaluation assets, and references."""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from pathlib import Path


PAIRS = (
    ("README.md", "README.ja.md"),
    ("AGENTS.md", "AGENTS-ja.md"),
    ("docs/architecture.md", "docs/ja/architecture.md"),
    ("docs/evaluation.md", "docs/ja/evaluation.md"),
    ("docs/glossary.md", "docs/ja/glossary.md"),
    ("docs/localization.md", "docs/ja/localization.md"),
    ("skills/outer-harness-retrospective/SKILL.md", "skills/outer-harness-retrospective/SKILL-ja.md"),
    (".agents/skills/maintain-japanese-references/SKILL.md", ".agents/skills/maintain-japanese-references/SKILL-ja.md"),
)
SKILL = "outer-harness-retrospective"
SKILL_DIR = Path("skills") / SKILL


def require(condition: bool, message: str) -> None:
    if not condition:
        raise ValueError(message)


def load_json(path: Path) -> dict:
    value = json.loads(path.read_text(encoding="utf-8"))
    require(isinstance(value, dict), f"{path}: top level must be an object")
    return value


def check_pairs(root: Path) -> None:
    for english, japanese in PAIRS:
        source, reference = root / english, root / japanese
        require(source.is_file() and reference.is_file(), f"missing translation pair: {english}")
        text = reference.read_text(encoding="utf-8")
        require(Path(english).name in text, f"missing English source link: {japanese}")
        require("英語" in text and ("優先" in text or "正式" in text), f"missing reference notice: {japanese}")


def check_skill(root: Path) -> None:
    text = (root / SKILL_DIR / "SKILL.md").read_text(encoding="utf-8")
    require(text.startswith("---\n"), "Skill frontmatter is missing")
    header = text.split("---", 2)[1]
    require(f"name: {SKILL}" in header, "Skill name is incorrect")
    require("description: " in header, "Skill description is missing")


def check_evaluations(root: Path) -> None:
    for filename, routing in (("evals.json", False), ("triggers.json", True)):
        path = root / SKILL_DIR / "evals" / filename
        document = load_json(path)
        require(set(document) <= {"skill_name", "evals", "execution"}, f"{path}: unknown top-level fields")
        require(document.get("skill_name") == SKILL, f"{path}: wrong Skill name")
        cases = document.get("evals")
        require(isinstance(cases, list) and cases, f"{path}: cases must be nonempty")
        seen: set[str] = set()
        for case in cases:
            require(isinstance(case, dict), f"{path}: case must be an object")
            identifier = str(case.get("id", ""))
            require(identifier and identifier not in seen, f"{path}: missing or duplicate case ID")
            seen.add(identifier)
            require(isinstance(case.get("prompt"), str) and bool(case["prompt"].strip()), f"{identifier}: prompt missing")
            if routing:
                handlers = case.get("expected_handlers")
                require(isinstance(handlers, list) and all(isinstance(x, str) for x in handlers), f"{identifier}: expected_handlers missing")
                require(len(handlers) == len(set(handlers)), f"{identifier}: duplicate handler")
            else:
                require("expected_handlers" not in case, f"{identifier}: routing field in behavior case")
                assertions = case.get("assertions", [])
                expected = case.get("expected_output", "")
                require(bool(assertions) or isinstance(expected, str) and bool(expected.strip()), f"{identifier}: grading requirement missing")
                require(isinstance(assertions, list), f"{identifier}: assertions must be a list")
                require(all(isinstance(a, str) and bool(a.strip()) or isinstance(a, dict) and isinstance(a.get("text"), str) and bool(a["text"].strip()) for a in assertions), f"{identifier}: invalid assertion")


def check_report(root: Path) -> None:
    path = root / SKILL_DIR / "evals" / "report.json"
    if not path.exists():
        return
    report = load_json(path)
    require(report.get("schema_version") == 2 and report.get("skill") == SKILL, "invalid evaluation report identity")
    candidate = report.get("candidate", {}).get("files", {})
    require(isinstance(candidate, dict), "invalid candidate manifest")
    for name, metadata in candidate.items():
        relative = Path(name)
        require(not relative.is_absolute() and ".." not in relative.parts, "unsafe report path")
        file = root / SKILL_DIR / relative
        require(file.is_file() and not file.is_symlink(), f"stale report file: {name}")
        digest = "sha256:" + hashlib.sha256(file.read_bytes()).hexdigest()
        require(metadata.get("sha256") == digest, f"stale report hash: {name}")
    for name, digest in report.get("evaluation_inputs", {}).get("evaluation_files", {}).items():
        relative = Path(name)
        require(not relative.is_absolute() and ".." not in relative.parts, "unsafe evaluation input path")
        file = root / relative
        require(file.is_file() and "sha256:" + hashlib.sha256(file.read_bytes()).hexdigest() == digest, f"stale evaluation input: {name}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parents[1])
    parser.add_argument("--ignore-report-for-skill")
    args = parser.parse_args()
    root = args.root.resolve()
    try:
        check_pairs(root)
        check_skill(root)
        check_evaluations(root)
        if args.ignore_report_for_skill != SKILL:
            check_report(root)
    except (OSError, ValueError, json.JSONDecodeError) as error:
        print(f"repository check failed: {error}", file=sys.stderr)
        return 1
    print("repository checks passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
