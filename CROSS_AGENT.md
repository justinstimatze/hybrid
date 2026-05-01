# Cross-agent portability

This repository targets **Anthropic Claude Code** as the primary platform. The pattern itself — alternating LLM and deterministic layers in a mutually-constraining cycle — is agent-agnostic; the skill content (`skills/hybrid-loops/`) and the three MCP servers under `mcp_servers/` (cal_log, metacog, schemaforge) are model-agnostic. What varies between agent ecosystems is the manifest format and discovery path.

Stub manifests for other major coding agents are included as a friendly gesture; they are not actively tested and may need adjustment to match each agent's current plugin spec.

## Claude Code (primary, tested)

```
.claude-plugin/
├── plugin.json
└── marketplace.json
```

Install via:

```bash
# Add this repo as a marketplace
/plugin marketplace add justinstimatze/hybrid

# Install
/plugin install hybrid-loops@hybrid-loops
```

Or symlink `skills/hybrid-loops/` into `~/.claude/skills/hybrid-loops/` for direct skill access without plugin install.

## OpenAI Codex (stub)

```
.codex-plugin/
└── plugin.json
```

The skill content at `skills/hybrid-loops/SKILL.md` is portable. The Codex plugin spec evolves; the stub at `.codex-plugin/plugin.json` is a starting point. **PRs from Codex users to update the manifest are welcome.**

## Cursor (stub)

```
.cursor-plugin/
└── plugin.json
```

Same approach: the skill content is portable; the manifest is a starting point. **PRs from Cursor users welcome.**

## Gemini (stub)

```
gemini-extension.json
```

Single-file extension manifest at the repo root. Same caveat: stub. **PRs welcome.**

## What's portable

- `skills/hybrid-loops/SKILL.md` — the skill body, agent-agnostic
- `skills/hybrid-loops/references/*.md` — reference files, agent-agnostic
- `mcp_servers/cal_log/` — standard MCP server (stdio JSON-RPC over the Anthropic-defined Model Context Protocol). Any agent supporting MCP can wire it up; the install path differs per agent.

## What's not (yet) portable

- The plugin manifests themselves are agent-specific
- Hook integration (when added) will be agent-specific
- Slash-command paths will be agent-specific

## Filing portability bugs

If you're using hybrid-loops on a non-Claude agent and the manifest needs updating, please open an issue or PR with the diff. Cross-agent maintenance is a community contribution; the maintainer's primary platform is Claude Code.
