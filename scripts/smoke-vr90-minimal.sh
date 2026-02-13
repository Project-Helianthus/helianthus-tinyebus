#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

profile="${1:-vr90}"
shift || true

case "$profile" in
  vr90)
    test_pattern='^TestSmokeVR90MinimalQuerySet$'
    ;;
  vr71|vr_71)
    test_pattern='^TestSmokeVR71IdentifyOnlyProfile$'
    ;;
  all)
    test_pattern='^TestSmoke(VR90MinimalQuerySet|VR71IdentifyOnlyProfile)$'
    ;;
  *)
    echo "unknown profile '$profile' (use: vr90 | vr71 | all)" >&2
    exit 1
    ;;
esac

GOWORK=off go test ./firmware/emulation -run "$test_pattern" -count=1 -v "$@"
