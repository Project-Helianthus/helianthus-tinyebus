# AGENTS

## Repository status

**DEPRECATED — historical oracle, harness, and bridge-contract reference.** This
repository is not active production firmware and must not be used to start new
firmware features, board support, or deployments. The associated PIC firmware
history is publicly retained in
[helianthus-ebus-adapter-pic](https://github.com/Project-Helianthus/helianthus-ebus-adapter-pic).

Preserve the deterministic oracle and harness behavior. Permitted changes are
narrow archival maintenance, reproducibility, historical-compatibility, or
security fixes with a concrete issue and acceptance criteria.

## Working rules

1. Work from an issue-specific branch named `issue/<number>-<slug>`; keep one active issue and PR per repository.
2. Keep changes minimal and targeted. Do not introduce new production firmware behavior or expand the historical contract without explicit operator direction.
3. Do not connect to, flash, provision, or mutate hardware; do not handle credentials or irreversible operations without explicit operator approval at the time of action.
4. Preserve deterministic oracle output and existing archival fixtures unless the issue explicitly changes a verified historical contract.
5. Keep terminology inclusive and avoid unrelated rewrites.
6. Externally visible eBUS protocol or contract changes require a companion update in the public [eBUS documentation repository](https://github.com/Project-Helianthus/helianthus-docs-ebus).
7. Address blocking review findings against the exact PR head before merge consideration. Use squash merge only when all required checks and reviews pass.

## Validation

Run the checks applicable to the changed files:

- Go behavior: `go test ./...` and `go vet ./...`.
- Archived firmware compile evidence when firmware code changes: install TinyGo,
  run `tinygo build -o /tmp/tinyebus-firmware.hex -target=pico firmware/adapterproto`,
  and require a successful native build. `make tinygo-build` may be used as a
  local convenience check, but its successful "tinygo not installed; skipping"
  result is not compile evidence and does not satisfy this gate.
- Oracle or profile changes: `go run ./cmd/adapterproto-oracle` and the affected `./scripts/smoke-vr90-minimal.sh` profile.
- Documentation-only changes: Markdown/link validation and `git diff --check`.

`make ci` remains the repository CI shortcut. These instructions are
self-contained: they require no workspace-root file, parent checkout, local
path convention, or private URL.
