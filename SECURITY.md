# Security model

This repo ships three local stdio MCP servers and a Claude Code skill. The threat model is small and worth being explicit about.

## Trust boundary

Each MCP server runs as a **subprocess of the calling agent** (Claude Code, or whichever client launched it). The server inherits the agent's user, environment, and filesystem permissions. There is **no privilege boundary between the agent and the MCP server** — anything the agent could do directly with `os.ReadFile` / `os.WriteFile`, the server can do via its tool calls.

This is the standard local-stdio MCP shape. It is **not** a sandbox. Don't run servers from this repo against an untrusted agent or one with elevated privileges you wouldn't grant the user.

## What each server does

- **cal_log** writes append-only events to `$CAL_LOG_PATH` (default `~/.cal_log/calibration.jsonl`). The contents are caller-supplied prediction records. **No network calls.**
- **metacog** reads a caller-supplied JSONL substrate path, runs in-process audits, returns results. **No network calls. No writes.**
- **schemaforge** reads a caller-supplied corpus path, makes calls to the **Anthropic API** using `ANTHROPIC_API_KEY` from env, and writes per-round outputs under a caller-supplied `output_dir`. **Sends corpus content to api.anthropic.com.**

## What this means for you

- **Don't run schemaforge on private/sensitive data** without first checking your Anthropic data-retention policy. The corpus content goes through the API the same way any other prompt does.
- **`cal_log` JSONL contains the predictions and verdicts you record**, including any text you pass via the `prediction` or `notes` fields. Treat the file as you'd treat any local log: don't put anything in there you wouldn't put in a debug log.
- **The MCP servers will read or write any path the calling agent supplies.** They don't validate that the path is "appropriate" — that's the calling agent's responsibility (and ultimately yours, since you launched the agent).

## Reporting an issue

If you find a real security issue (privilege escalation, secret leakage, command injection in tool arguments, etc.), please open a GitHub Issue with `[security]` in the title or email the repo author directly via the GitHub profile. There's no formal SLA on fixes — this is research-stage software shipped by one person.

## What is explicitly out of scope

- Sandboxing the MCP servers from each other or from the host (use OS-level sandboxing if you need that).
- Auditing the prompts in `mcp_servers/schemaforge/prompts.go` for safety against prompt-injection attacks via corpus content (the loop is designed for trusted corpora; running it on attacker-controlled content has not been threat-modeled).
- Network-side security of the Anthropic API itself.
