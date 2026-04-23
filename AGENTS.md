# AGENTS

This repository is part of the **Helianthus Multi-Protocol HVAC Gateway Platform**.

## Dual-AI Operating Model

All development follows the dual-AI orchestrator protocol defined in the workspace-root [`AGENTS.md`](../AGENTS.md):

- **Role binding:** `ORCHESTRATOR` and `CO_PILOT` are portable roles. The current workspace default is Claude as orchestrator and Codex as co-pilot, but this repo must remain swap-ready.
- **Co-Pilot use:** Use the co-pilot only for adversarial/cooperative reasoning roles such as planner, bounded developer, reviewer, and second-opinion consultant. Do not spend Claude MCP or any equivalent co-pilot runtime on file reads, globs, grep, polling, or routine repo inspection.
- **Fallback:** If the preferred co-pilot is unavailable, throttled, or not integrated, the active orchestrator spawns fresh-context agents on the available runtime and keeps the same supervision contract.
- Phases: Adversarial Planning → Smart Routing → Dual Code Review
- Hard rules: one issue/PR per repo, squash+merge only, doc-gate, transport-gate, MCP-first

See the root AGENTS.md for the full protocol, routing tables, portable role prompts, and invariants. When running under Codex local orchestration, use the workspace-root skills `helianthus-orchestrator-supervision` and `helianthus-review-watch` as the portable supervision contract.

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
