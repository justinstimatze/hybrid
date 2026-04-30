#!/usr/bin/env python3
# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "mcp>=1.0",
# ]
# ///
"""cal_log — calibration logger MCP server for hybrid-loop projects.

Per-evaluator prediction logging with rolling-window hit-rate aggregation.
Append-only event log at $CAL_LOG_PATH (default: ~/.cal_log/calibration.jsonl).

Tools:
- predict(loop, input_hash, prediction, model_id, ...) -> {prediction_id, verdict_due_by}
- resolve(prediction_id, verdict, verdict_source) -> resolved record
- hit_rate(loop, window_days) -> {hit_rate, total_resolved, verdict_breakdown}
- list_pending(loop?, limit) -> unresolved predictions ordered by due date
- list_recent(loop?, limit) -> recent predictions (resolved or not)

The hybrid claim: an evaluator that can't show its hit-rate is theater.
"""

from __future__ import annotations

import json
import os
import time
import uuid
from pathlib import Path
from typing import Any

from mcp.server.fastmcp import FastMCP

DB_PATH = Path(os.environ.get("CAL_LOG_PATH", str(Path.home() / ".cal_log" / "calibration.jsonl")))
DB_PATH.parent.mkdir(parents=True, exist_ok=True)

mcp = FastMCP("cal_log")


def _read_events() -> list[dict[str, Any]]:
    if not DB_PATH.exists():
        return []
    events: list[dict[str, Any]] = []
    with DB_PATH.open() as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                events.append(json.loads(line))
            except json.JSONDecodeError:
                continue
    return events


def _append_event(event: dict[str, Any]) -> None:
    with DB_PATH.open("a") as f:
        f.write(json.dumps(event, separators=(",", ":")) + "\n")


def _fold(events: list[dict[str, Any]]) -> dict[str, dict[str, Any]]:
    """Fold the event stream into current state per prediction_id."""
    state: dict[str, dict[str, Any]] = {}
    for ev in events:
        pid = ev.get("prediction_id")
        if not pid:
            continue
        if ev.get("event") == "predict":
            state[pid] = {
                "prediction_id": pid,
                "ts": ev.get("ts"),
                "loop": ev.get("loop"),
                "lens_or_reasoner": ev.get("lens_or_reasoner"),
                "input_hash": ev.get("input_hash"),
                "prediction": ev.get("prediction"),
                "model_id": ev.get("model_id"),
                "schema_version": ev.get("schema_version"),
                "verdict_due_by": ev.get("verdict_due_by"),
                "verdict": None,
                "verdict_source": None,
                "verdict_ts": None,
            }
        elif ev.get("event") == "resolve" and pid in state:
            state[pid]["verdict"] = ev.get("verdict")
            state[pid]["verdict_source"] = ev.get("verdict_source")
            state[pid]["verdict_ts"] = ev.get("verdict_ts")
    return state


@mcp.tool()
def predict(
    loop: str,
    input_hash: str,
    prediction: dict[str, Any],
    model_id: str,
    schema_version: int = 1,
    verdict_due_in_days: int = 7,
    lens_or_reasoner: str = "reasoner",
) -> dict[str, Any]:
    """Log a typed LLM evaluator's prediction. Returns prediction_id for later resolve().

    `loop` should be a stable identifier for the project / evaluator (e.g. "slimemold-claims",
    "recruiter-fit-scorer"). `input_hash` is any stable hash of the input the prediction was made on.
    `prediction` is the structured record the evaluator emitted.
    """
    pid = str(uuid.uuid4())
    now = time.time()
    event = {
        "event": "predict",
        "prediction_id": pid,
        "ts": now,
        "loop": loop,
        "lens_or_reasoner": lens_or_reasoner,
        "input_hash": input_hash,
        "prediction": prediction,
        "model_id": model_id,
        "schema_version": schema_version,
        "verdict_due_by": now + verdict_due_in_days * 86400,
    }
    _append_event(event)
    return {
        "prediction_id": pid,
        "verdict_due_by": event["verdict_due_by"],
    }


@mcp.tool()
def resolve(
    prediction_id: str,
    verdict: str,
    verdict_source: str = "manual",
) -> dict[str, Any]:
    """Mark a prediction's verdict. Common verdicts: 'confirmed', 'refuted', 'partial', 'irrelevant'.
    `verdict_source` describes how the verdict was determined (e.g. 'user_pushback', 'metric_X', 'manual').
    """
    state = _fold(_read_events())
    if prediction_id not in state:
        return {"error": f"prediction_id {prediction_id} not found"}
    if state[prediction_id]["verdict"] is not None:
        return {"error": f"prediction_id {prediction_id} already resolved", "existing": state[prediction_id]}
    event = {
        "event": "resolve",
        "prediction_id": prediction_id,
        "ts": time.time(),
        "verdict": verdict,
        "verdict_source": verdict_source,
        "verdict_ts": time.time(),
    }
    _append_event(event)
    state[prediction_id]["verdict"] = verdict
    state[prediction_id]["verdict_source"] = verdict_source
    state[prediction_id]["verdict_ts"] = event["verdict_ts"]
    return state[prediction_id]


@mcp.tool()
def hit_rate(loop: str, window_days: int = 30) -> dict[str, Any]:
    """Compute hit-rate and verdict breakdown for `loop` in the past `window_days`.
    Hit rate = (verdict == 'confirmed') / total_resolved. Returns None if zero resolved.
    """
    cutoff = time.time() - window_days * 86400
    state = _fold(_read_events())
    in_loop = [
        r for r in state.values()
        if r.get("loop") == loop
        and r.get("verdict") is not None
        and (r.get("verdict_ts") or 0) >= cutoff
    ]
    total = len(in_loop)
    if total == 0:
        return {
            "loop": loop,
            "window_days": window_days,
            "total_resolved": 0,
            "hit_rate": None,
            "verdict_breakdown": {},
        }
    breakdown: dict[str, int] = {}
    for r in in_loop:
        v = r["verdict"]
        breakdown[v] = breakdown.get(v, 0) + 1
    confirmed = breakdown.get("confirmed", 0)
    return {
        "loop": loop,
        "window_days": window_days,
        "total_resolved": total,
        "hit_rate": confirmed / total,
        "verdict_breakdown": breakdown,
    }


@mcp.tool()
def list_pending(loop: str | None = None, limit: int = 50) -> list[dict[str, Any]]:
    """List unresolved predictions, ordered by verdict_due_by (oldest first).
    Useful for finding what needs verdicts next.
    """
    state = _fold(_read_events())
    pending = [r for r in state.values() if r.get("verdict") is None]
    if loop:
        pending = [r for r in pending if r.get("loop") == loop]
    pending.sort(key=lambda r: r.get("verdict_due_by") or 0)
    return pending[:limit]


@mcp.tool()
def list_recent(loop: str | None = None, limit: int = 50) -> list[dict[str, Any]]:
    """List recent predictions (resolved or not), most recent first."""
    state = _fold(_read_events())
    records = list(state.values())
    if loop:
        records = [r for r in records if r.get("loop") == loop]
    records.sort(key=lambda r: r.get("ts") or 0, reverse=True)
    return records[:limit]


@mcp.tool()
def stats() -> dict[str, Any]:
    """Top-level summary across all loops: counts, hit-rate where computable."""
    state = _fold(_read_events())
    by_loop: dict[str, dict[str, Any]] = {}
    for r in state.values():
        loop = r.get("loop") or "unknown"
        b = by_loop.setdefault(loop, {"total": 0, "resolved": 0, "confirmed": 0})
        b["total"] += 1
        if r.get("verdict") is not None:
            b["resolved"] += 1
            if r["verdict"] == "confirmed":
                b["confirmed"] += 1
    for loop, b in by_loop.items():
        b["hit_rate"] = (b["confirmed"] / b["resolved"]) if b["resolved"] else None
    return {
        "db_path": str(DB_PATH),
        "total_predictions": len(state),
        "by_loop": by_loop,
    }


if __name__ == "__main__":
    mcp.run()
