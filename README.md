# helianthus-tinyebus

Bootstrap-level TinyGo firmware skeleton for Helianthus eBUS experiments.

## Status
No firmware behavior is implemented yet. This repository currently provides:
- roadmap notes
- package responsibilities
- placeholder interfaces for future bus/HAL wiring

## Short roadmap
- **M1 (issue #1):** Document roadmap and define package contracts only.
- **M2:** Add minimal firmware entry wiring between `hal` and `bus` modules.
- **M3:** Introduce eBUS framing + transaction state placeholders for testing.
- **M4:** Add board-specific HAL adapters and TinyGo target validation.

## Module responsibilities
- `firmware/main.go`: TinyGo entry package and future bootstrapping point.
- `firmware/bus`: eBUS-facing interfaces/contracts (no protocol logic yet).
- `firmware/hal`: hardware abstraction interfaces (UART/timing/pins) with no concrete drivers.
- `ARCHITECTURE.md` and `CONVENTIONS.md`: high-level project constraints and coding standards.

## Build/test (placeholder)
- `go test ./...`
- `tinygo build ./firmware` (if TinyGo is available locally)
