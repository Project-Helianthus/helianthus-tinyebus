# helianthus-tinyebus

TinyGo firmware skeleton for Helianthus eBUS experiments.

## Current status
- Firmware scope stays intentionally narrow: deterministic target emulation + identify-only presets (`VR90`, `VR_71`).
- `firmware/bus` and `firmware/hal` are contracts only (no runtime protocol engine, no board drivers).
- `firmware/emulation` is the only functional subsystem and is covered by Go tests + smoke scripts.

## Prerequisites
- Go `1.22+` (module baseline in `go.mod`).
- `make` and `bash`.
- TinyGo (optional for now, required for `make tinygo-build` and future board flashing).

## Toolchain setup
1. Install Go `1.22+`.
2. Install TinyGo from https://tinygo.org/getting-started/ (optional but recommended).
3. Verify toolchain:
   - `go version`
   - `tinygo version` (if installed)
   - `make --version`

## Build + flash workflow
1. Run unit tests:
   - `make test`
2. Run TinyGo compile check (placeholder artifact):
   - `make tinygo-build`
   - Current behavior compiles `./firmware` to a temporary `wasm` output and deletes it.
3. Flash workflow status:
   - Board-specific flashing is **not wired yet** (no concrete HAL implementation and `firmware/main.go` is bootstrap-only).
   - For future board targets, flash shape will be:
     - `tinygo flash -target <tinygo-target> ./firmware`

### Supported target status
| Target | Build status | Flash status | Notes |
| --- | --- | --- | --- |
| `wasm` | ✅ via `make tinygo-build` | N/A | Placeholder compile check used for TinyGo availability. |
| Hardware TinyGo targets | 🚧 pending | 🚧 pending | No board adapter/entry wiring yet. |

## Firmware module map
| Module | Responsibility | Status |
| --- | --- | --- |
| `firmware/main.go` | TinyGo entrypoint | Bootstrap placeholder (`main()` intentionally empty). |
| `firmware/bus` | eBUS transport/engine contracts | Interface-only; no protocol logic yet. |
| `firmware/hal` | UART/clock/board abstraction contracts | Interface-only; no board drivers yet. |
| `firmware/emulation/framework.go` | Rule matcher + response builder + timing constraints | Functional and unit-tested. |
| `firmware/emulation/identify_only.go` | Generic identify-only target profile (`07 04`) | Functional with presets and tests. |
| `firmware/emulation/vr90.go` | VR90 convenience profile wrapper | Functional (`NewVR90Target`, defaults, tests). |
| `firmware/emulation/harness.go` | Deterministic virtual-time query harness | Functional and unit-tested. |
| `scripts/smoke-vr90-minimal.sh` | Smoke test runner for profile-specific identify checks | Functional (`vr90`, `vr71`, `all`). |

## Emulation status
| Profile | Address | Identify (`PB=0x07`, `SB=0x04`) | Status | Smoke coverage |
| --- | --- | --- | --- | --- |
| `VR90` | `0x15` | ✅ implemented | Identify-only | `TestSmokeVR90MinimalQuerySet` |
| `VR_71` | `0x26` | ✅ implemented | Identify-only | `TestSmokeVR71IdentifyOnlyProfile` |
| Other commands | N/A | ❌ not implemented | Returns `ErrNoMatchingRule` | Covered by negative tests |

## Smoke validation commands
- `./scripts/smoke-vr90-minimal.sh vr90`
- `./scripts/smoke-vr90-minimal.sh vr71`
- `./scripts/smoke-vr90-minimal.sh all`
- `make smoke-vr90` (defaults to `vr90`)

Expected result for successful runs:
- the selected `TestSmoke...` test(s) run
- output ends with `PASS` and `ok github.com/d3vi1/helianthus-tinyebus/firmware/emulation ...`

## Timing note
- Precise eBUS emulation timing belongs on firmware/MCU-side execution where jitter is predictable.
- Host-side relays (for example `ebusd-tcp`) are useful for functional checks but not cycle-accurate timing validation.

## Concise roadmap
- ✅ **M1 / issue #1:** package contracts + initial docs.
- ✅ **M2 / issue #8:** target-emulation framework + identify-only profile path.
- 🔜 **M3:** wire bootstrap path between `hal` and `bus`.
- 🔜 **M4:** add eBUS framing + transaction state placeholders.
- 🔜 **M5:** add board-specific HAL adapters and real TinyGo target flashing validation.
