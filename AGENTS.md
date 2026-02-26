# AGENTS

This repository is part of the **Helianthus Multi-Protocol HVAC Gateway Platform**.

## Dual-AI Operating Model

All development follows the dual-AI orchestrator protocol defined in the workspace-root [`AGENTS.md`](../AGENTS.md):

- **Orchestrator:** Claude Code — orchestration, hard dev (complexity 7–10), angry tester, deep consultant
- **Co-Pilot:** Codex — adversarial planning, easy dev (complexity 1–6), code review, second opinions
- Phases: Adversarial Planning → Smart Routing → Dual Code Review
- Hard rules: one issue/PR per repo, squash+merge only, doc-gate, transport-gate, MCP-first

See the root AGENTS.md for the full protocol, routing tables, system prompts, and invariants.

---

## Repo-Specific Rules

## AGENT

Bootstrap-only repository. Keep changes minimal and avoid adding firmware features
unless explicitly requested.

### Quick commands
- `go test ./...`
- `tinygo build ./firmware` (placeholder)

### Notes
- Docs live in `README.md`, `ARCHITECTURE.md`, and `CONVENTIONS.md`.
- Local notes belong in `AGENT-local.md` (gitignored).
- React (emoji) to every review comment and reply with status when actioned.
