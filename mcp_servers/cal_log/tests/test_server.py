#!/usr/bin/env python3
# /// script
# requires-python = ">=3.10"
# dependencies = ["mcp>=1.0", "pytest"]
# ///
"""Direct unit tests for cal_log core logic. Bypasses MCP transport and tests the
underlying functions on a temp file. Doesn't test stdio JSON-RPC routing — that's
covered by the mcp library's own tests.
"""

import json
import os
import tempfile
import time
from pathlib import Path

# Use a temp file before importing the module
_tmpdir = tempfile.mkdtemp(prefix="cal_log_test_")
os.environ["CAL_LOG_PATH"] = str(Path(_tmpdir) / "calibration.jsonl")

import sys
sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
import server  # noqa: E402


def _reset() -> None:
    if server.DB_PATH.exists():
        server.DB_PATH.unlink()


def test_predict_returns_id_and_due_by() -> None:
    _reset()
    out = server.predict(
        loop="t1",
        input_hash="h1",
        prediction={"flag": "X"},
        model_id="claude-sonnet-4-6",
        verdict_due_in_days=7,
    )
    assert "prediction_id" in out
    assert "verdict_due_by" in out
    assert out["verdict_due_by"] > time.time()


def test_resolve_then_hit_rate() -> None:
    _reset()
    p1 = server.predict(loop="t2", input_hash="h", prediction={}, model_id="m")
    p2 = server.predict(loop="t2", input_hash="h", prediction={}, model_id="m")
    p3 = server.predict(loop="t2", input_hash="h", prediction={}, model_id="m")
    server.resolve(p1["prediction_id"], "confirmed")
    server.resolve(p2["prediction_id"], "confirmed")
    server.resolve(p3["prediction_id"], "refuted")
    hr = server.hit_rate("t2", window_days=30)
    assert hr["total_resolved"] == 3
    assert abs(hr["hit_rate"] - 2 / 3) < 1e-9
    assert hr["verdict_breakdown"] == {"confirmed": 2, "refuted": 1}


def test_resolve_unknown_id_returns_error() -> None:
    _reset()
    out = server.resolve("not-a-real-id", "confirmed")
    assert "error" in out


def test_resolve_twice_returns_error() -> None:
    _reset()
    p = server.predict(loop="t3", input_hash="h", prediction={}, model_id="m")
    server.resolve(p["prediction_id"], "confirmed")
    out = server.resolve(p["prediction_id"], "refuted")
    assert "error" in out


def test_list_pending_orders_by_due_date() -> None:
    _reset()
    p1 = server.predict(loop="t4", input_hash="h", prediction={}, model_id="m", verdict_due_in_days=10)
    p2 = server.predict(loop="t4", input_hash="h", prediction={}, model_id="m", verdict_due_in_days=1)
    p3 = server.predict(loop="t4", input_hash="h", prediction={}, model_id="m", verdict_due_in_days=5)
    pending = server.list_pending(loop="t4")
    assert len(pending) == 3
    assert pending[0]["prediction_id"] == p2["prediction_id"]
    assert pending[1]["prediction_id"] == p3["prediction_id"]
    assert pending[2]["prediction_id"] == p1["prediction_id"]


def test_resolved_excluded_from_pending() -> None:
    _reset()
    p = server.predict(loop="t5", input_hash="h", prediction={}, model_id="m")
    server.resolve(p["prediction_id"], "confirmed")
    pending = server.list_pending(loop="t5")
    assert len(pending) == 0


def test_hit_rate_zero_resolved() -> None:
    _reset()
    server.predict(loop="t6", input_hash="h", prediction={}, model_id="m")
    hr = server.hit_rate("t6")
    assert hr["total_resolved"] == 0
    assert hr["hit_rate"] is None


def test_stats_aggregates_across_loops() -> None:
    _reset()
    p_a1 = server.predict(loop="loop_a", input_hash="h", prediction={}, model_id="m")
    p_a2 = server.predict(loop="loop_a", input_hash="h", prediction={}, model_id="m")
    p_b1 = server.predict(loop="loop_b", input_hash="h", prediction={}, model_id="m")
    server.resolve(p_a1["prediction_id"], "confirmed")
    server.resolve(p_b1["prediction_id"], "refuted")
    s = server.stats()
    assert s["total_predictions"] == 3
    assert s["by_loop"]["loop_a"]["total"] == 2
    assert s["by_loop"]["loop_a"]["resolved"] == 1
    assert s["by_loop"]["loop_a"]["hit_rate"] == 1.0
    assert s["by_loop"]["loop_b"]["resolved"] == 1
    assert s["by_loop"]["loop_b"]["hit_rate"] == 0.0


if __name__ == "__main__":
    import pytest
    pytest.main([__file__, "-v"])
