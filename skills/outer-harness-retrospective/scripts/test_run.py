"""Synthetic checks for the Skill's pinned acquisition entry point."""

from __future__ import annotations

import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import run


class AcquisitionTests(unittest.TestCase):
    def test_pin_is_a_full_commit(self) -> None:
        self.assertEqual(len(run.pinned_commit()), 40)

    def test_build_rejects_wrong_checkout_before_go_build(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            with patch.object(run.subprocess, "run") as command:
                with patch.object(run.subprocess, "check_output", return_value="0" * 40):
                    with self.assertRaisesRegex(ValueError, "pinned commit"):
                        run.build_agent_sessions(Path(temporary), run.pinned_commit(), {})
            self.assertFalse(any(call.args[0][0] == "go" for call in command.call_args_list))

    def test_offline_path_is_checked_by_go_helper(self) -> None:
        with patch.object(run.subprocess, "run") as command:
            command.return_value.returncode = 0
            self.assertEqual(run.run_pilot(Path("/tmp/pinned-agent-sessions"), None, {}), 0)
            self.assertEqual(
                command.call_args.args[0],
                ["go", "run", "./cmd/retro-pilot", "--agent-sessions-bin", "/tmp/pinned-agent-sessions"],
            )


if __name__ == "__main__":
    unittest.main()
