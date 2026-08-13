# helianthus-tinyebus

> **Deprecated — historical/reference material only.** This repository is retained as an archival oracle, harness, and contract reference. It is not under active development.

`helianthus-tinyebus` is the historical off-PIC oracle, harness, and bridge repository for the Helianthus eBUS adapter northbound contract. Its material describes an ESP8266 D1 mini with 4 MB flash, rather than the PIC.

The associated PIC firmware history is in [helianthus-ebus-adapter-pic](https://github.com/Project-Helianthus/helianthus-ebus-adapter-pic). This repository preserves the contract surface for host-side reference and reproducible archival checks.

## Historical Purpose and Scope

### Archived contents

- TinyGo firmware entrypoint and package layout (`firmware/`).
- eBUS and HAL contract surfaces (`firmware/bus`, `firmware/hal`).
- Deterministic emulation framework/profiles/harness (`firmware/emulation`).
- Adapter oracle package and CLI for ENH, ENS, adapter INFO, deterministic runtime-contract parity checks, and scan/status/state report samples (`firmware/adapterproto`, `cmd/adapterproto-oracle`).
- Lightweight smoke/test scripts (`scripts/`).

### Historical boundaries

- Production gateway/API runtime (`helianthus-ebusgateway`).
- PIC runtime or bootloader firmware.
- Registry/provider/schema logic (`helianthus-ebusreg`).
- Home Assistant integration/add-on packaging (`helianthus-ha-integration`, `helianthus-ha-addon`).

## Historical Status and Maturity

- The repository was a bootstrap repository with tests and smoke checks.
- The adapter contract oracle is preserved in `firmware/adapterproto` and can be emitted with `go run ./cmd/adapterproto-oracle`.
- The oracle includes a runtime-contract model for INIT, START, SEND, and cancel/no-response samples, plus a deterministic scan/status/state oracle derived from the decompiled PIC control flow.
- Its preserved functional scope is emulation-first (VR90/VR_71 identify and VR90 mapped-command/discovery paths).
- PIC runtime/board driver wiring was outside this repository's historical scope.

## Historical Relationship

```text
helianthus-tinyebus (oracle/harness/spec on ESP8266 D1 mini) -> helianthus-ebusgateway -> helianthus-ha-integration
```

## Archived Verification Commands

### 0) Prerequisites

- Go `1.22+`
- `make`, `bash`
- TinyGo (optional; only required for `make tinygo-build`)

### 1) Clone and run baseline checks

```bash
git clone https://github.com/Project-Helianthus/helianthus-tinyebus.git
cd helianthus-tinyebus
make test
go test ./...
go vet ./...
```

### 2) TinyGo compile check (optional archival check)

```bash
make tinygo-build
```

Behavior:
- if TinyGo is installed: builds `./firmware` with `-target wasm` as placeholder compile verification.
- if TinyGo is missing: prints `tinygo not installed; skipping tinygo build`.

### 2b) Adapter oracle JSON

```bash
go run ./cmd/adapterproto-oracle
```

This prints deterministic JSON for parity checks against the C-side contract implementation, including runtime-contract samples for INIT, START, SEND, and cancel flows, INFO coverage samples for all adapter INFO IDs, and scan/status/state snapshots derived from the decompiled firmware model.

### 3) Historical smoke tests

```bash
./scripts/smoke-vr90-minimal.sh vr90
./scripts/smoke-vr90-minimal.sh vr71
./scripts/smoke-vr90-minimal.sh all
```

## Historical Smoke-Test Configuration Examples

### VR90 identify + B509 discovery + mapped command profile

```go
profile := emulation.DefaultVR90Profile()
profile.EnableB509Discovery = true
profile.MappedCommands = []emulation.VR90MappedCommand{
	{
		Name:         "read-temp",
		Primary:      0xB5,
		Secondary:    0x09,
		PayloadExact: []byte{0x24},
		ResponseData: []byte{0x00, 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'},
	},
}
target, err := emulation.NewVR90Target(profile)
_ = target
_ = err
```

### Identify-only preset profiles

```go
vr90, _ := emulation.NewIdentifyOnlyTarget(emulation.PresetVR90IdentifyOnlyProfile())
vr71, _ := emulation.NewIdentifyOnlyTarget(emulation.PresetVR71IdentifyOnlyProfile())
_, _ = vr90, vr71
```

## Archived Validation Commands

| Area | Command |
|---|---|
| formatting | `find . -name '*.go' -type f -print0 \| xargs -0 gofmt -w` |
| tests | `go test ./...` |
| vet | `go vet ./...` |
| make CI shortcut | `make ci` |
| tinygo compile check | `make tinygo-build` |
| smoke VR90 | `./scripts/smoke-vr90-minimal.sh vr90` |
| smoke VR_71 | `./scripts/smoke-vr90-minimal.sh vr71` |
| smoke all profiles | `./scripts/smoke-vr90-minimal.sh all` |
| terminology gate (CI parity) | `if git grep -nIwiE 'm[a]ster|s[l]ave'; then echo "Found legacy terminology in tracked files."; exit 1; fi` |

## Link Map

### Local docs in this repo

- `ARCHITECTURE.md`
- `CONVENTIONS.md`
- `AGENT.md`

### Related Helianthus repos/docs

- Associated PIC firmware history: https://github.com/Project-Helianthus/helianthus-ebus-adapter-pic
- eBUS gateway runtime: https://github.com/Project-Helianthus/helianthus-ebusgateway
- eBUS registry/provider layer: https://github.com/Project-Helianthus/helianthus-ebusreg
- eBUS protocol docs: https://github.com/Project-Helianthus/helianthus-docs-ebus
