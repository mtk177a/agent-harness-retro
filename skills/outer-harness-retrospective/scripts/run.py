#!/usr/bin/env python3
"""Acquire bounded pilot evidence with a pinned agent-sessions build."""

from __future__ import annotations

import argparse
import os
import re
import subprocess
import sys
import tempfile
from pathlib import Path


SCRIPTS = Path(__file__).resolve().parent
PILOT = SCRIPTS / "retro-pilot"
PIN_FILE = PILOT / "internal" / "pilot" / "agent-sessions.commit"
UPSTREAM = "https://github.com/mtk177a/agent-sessions.git"


def pinned_commit() -> str:
    commit = PIN_FILE.read_text(encoding="ascii").strip()
    if not re.fullmatch(r"[0-9a-f]{40}", commit):
        raise ValueError("invalid agent-sessions commit pin")
    return commit


def build_environment(directory: Path) -> dict[str, str]:
    environment = os.environ.copy()
    environment.update(
        GOTOOLCHAIN="local",
        GOPROXY="off",
        GOWORK="off",
        GOFLAGS="",
        GOMODCACHE=str(directory / "go-modules"),
        GOCACHE=str(directory / "go-build"),
    )
    return environment


def build_agent_sessions(directory: Path, commit: str, environment: dict[str, str]) -> Path:
    source = directory / "agent-sessions"
    binary = directory / "agent-sessions-bin"
    subprocess.run(["git", "init", "--quiet", str(source)], check=True, env=environment)
    subprocess.run(
        ["git", "-C", str(source), "fetch", "--quiet", "--depth=1", UPSTREAM, commit],
        check=True,
        env=environment,
    )
    subprocess.run(
        ["git", "-C", str(source), "checkout", "--quiet", "--detach", "FETCH_HEAD"],
        check=True,
        env=environment,
    )
    revision = subprocess.check_output(
        ["git", "-C", str(source), "rev-parse", "HEAD"], env=environment, text=True
    ).strip()
    if revision != commit:
        raise ValueError("agent-sessions checkout does not match the pinned commit")
    status = subprocess.check_output(
        ["git", "-C", str(source), "status", "--porcelain"], env=environment, text=True
    )
    if status:
        raise ValueError("agent-sessions checkout is not clean")
    subprocess.run(
        ["go", "build", "-buildvcs=true", "-o", str(binary), "./cmd/agent-sessions"],
        check=True,
        cwd=source,
        env=environment,
    )
    return binary


def run_pilot(binary: Path, config: Path | None, environment: dict[str, str]) -> int:
    command = ["go", "run", "./cmd/retro-pilot", "--agent-sessions-bin", str(binary)]
    if config is not None:
        command.extend(["--config", str(config)])
    return subprocess.run(command, cwd=PILOT, env=environment, check=False).returncode


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--config", type=Path, help="absolute path to an existing agent-sessions config")
    parser.add_argument(
        "--agent-sessions-bin",
        type=Path,
        help="optional offline executable override; its clean pinned build is still checked",
    )
    args = parser.parse_args()
    if args.config is not None and not args.config.is_absolute():
        parser.error("--config must be an absolute path")
    if args.agent_sessions_bin is not None and not args.agent_sessions_bin.is_absolute():
        parser.error("--agent-sessions-bin must be an absolute path")
    try:
        commit = pinned_commit()
        with tempfile.TemporaryDirectory(prefix="outer-harness-retrospective-") as temporary:
            directory = Path(temporary)
            environment = build_environment(directory)
            binary = args.agent_sessions_bin
            if binary is None:
                binary = build_agent_sessions(directory, commit, environment)
            return run_pilot(binary, args.config, environment)
    except (OSError, subprocess.CalledProcessError, ValueError) as error:
        print(f"evidence acquisition failed: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
