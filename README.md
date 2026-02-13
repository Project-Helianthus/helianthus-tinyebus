# helianthus-tinyebus

Bootstrap-level TinyGo firmware skeleton for Helianthus eBUS experiments.

## Scope and constraints
- This repo stays firmware-focused and avoids host-side transport dependencies.
- Functional scope in this milestone is intentionally narrow: target emulation framework + identify-only presets (VR90, VR_71).
- Broader protocol stack work (full bus state machine, additional devices) remains out of scope.

## Status
This repository currently provides:
- roadmap notes and firmware package responsibilities
- placeholder interfaces for future bus/HAL wiring
- generic target-emulation framework (request matcher + response builder + timing constraints)
- deterministic virtual-time harness for unit/integration tests
- generic identify-only target profile API (`07 04`) with VR90 and VR_71 presets

## Short roadmap
- **M1 (issue #1):** Document roadmap and define package contracts only.
- **M2 (issue #8):** Add firmware-side target emulation framework and minimal VR90 profile.
- **M3:** Add minimal firmware entry wiring between `hal` and `bus` modules.
- **M4:** Introduce eBUS framing + transaction state placeholders for testing.
- **M5:** Add board-specific HAL adapters and TinyGo target validation.

## Module responsibilities
- `firmware/main.go`: TinyGo entry package and future bootstrapping point.
- `firmware/bus`: eBUS-facing interfaces/contracts (no protocol logic yet).
- `firmware/hal`: hardware abstraction interfaces (UART/timing/pins) with no concrete drivers.
- `firmware/emulation`: target-emulation framework, deterministic harness, and identify-only target profiles/presets.
- `ARCHITECTURE.md` and `CONVENTIONS.md`: high-level project constraints and coding standards.

## Target emulation timing constraints
- Precise eBUS emulation timing should run firmware-side (adapter MCU) where response jitter is predictable.
- Host-side relays such as `ebusd-tcp` are suitable for functional checks but not for cycle-accurate target emulation because scheduler/network jitter can violate tight response windows.

## Build/test
- `make test`
- `make tinygo-build` (TinyGo optional; skips if unavailable)
- `./scripts/smoke-vr90-minimal.sh vr90`
- `./scripts/smoke-vr90-minimal.sh vr71`
- `./scripts/smoke-vr90-minimal.sh all`
