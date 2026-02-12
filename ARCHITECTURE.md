# Architecture

This repo is a minimal TinyGo firmware skeleton.

## Intended layout
- `firmware/` TinyGo entry point and firmware packages
  - `main.go` TinyGo entry package
  - `bus/` eBUS protocol stubs
  - `hal/` hardware abstraction stubs
  - `emulation/` target-emulation framework and deterministic harness (VR90 minimal profile)
- `.github/workflows/` CI placeholder

## Status
Repository is still intentionally minimal. `firmware/emulation` contains first functional behavior for deterministic target-emulation tests, while transport/HAL remain stubs.
