# Architecture

This repo is an off-PIC oracle/harness/spec layer for the Helianthus eBUS adapter contract. It runs on an ESP8266 D1 mini with 4 MB flash and is not the PIC runtime.

## Intended layout
- `firmware/` TinyGo entry point and firmware packages
  - `main.go` TinyGo entry package
  - `bus/` eBUS protocol stubs
  - `hal/` hardware abstraction stubs
  - `emulation/` target-emulation framework and deterministic harness (VR90 minimal profile)
- `firmware/adapterproto/` ENH/ENS/info contract oracle, deterministic runtime-contract model, scan/status/state report samples, and parity helpers for the PIC northbound protocol
- `cmd/adapterproto-oracle/` CLI that emits deterministic JSON for host-side parity checks
- `.github/workflows/` CI placeholder

## Status
Repository is still intentionally minimal for on-device behavior. `firmware/emulation` contains deterministic target-emulation tests, while `firmware/adapterproto` defines the host-side oracle for ENH, ENS, adapter INFO payload parsing, a deterministic runtime-contract model, and scan/status/state report samples derived from the decompiled firmware control flow. The PIC runtime/bootloader are intentionally kept out of this repository.
