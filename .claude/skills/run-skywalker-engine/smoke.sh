#!/usr/bin/env bash
# Driver for the skywalker_engine web app: build -> launch -> smoke-test -> teardown.
# This is the harness the run-skywalker-engine skill points at. Every probe here was
# run against the actual app; it is the programmatic handle on the running server.
#
# Usage (run from anywhere; it cd's to the repo root itself):
#   bash .claude/skills/run-skywalker-engine/smoke.sh             # build + launch + probe + stop
#   bash .claude/skills/run-skywalker-engine/smoke.sh --no-build  # reuse existing binary
#   BASE=http://localhost:4000 bash .../smoke.sh --external       # probe an already-running server
#
# Exit code is 0 only if every probe passed.
set -u

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../../.." && pwd)"          # <unit> = skywalker_engine repo root
BASE="${BASE:-http://localhost:4000}"

# Host is Windows; Go emits a .exe. Stay portable if someone runs this on Linux.
BIN="skywalker_engine_app"
case "$(uname -s)" in MINGW*|MSYS*|CYGWIN*) BIN="$BIN.exe";; esac

PID=""; EXTERNAL=0; BUILD=1
for a in "$@"; do case "$a" in
  --no-build) BUILD=0 ;;
  --external) EXTERNAL=1 ;;
esac; done

cd "$ROOT"  # the app reads ./.env via os.Getwd() — MUST launch from repo root
cleanup() { [ -n "$PID" ] && kill "$PID" 2>/dev/null; }
trap cleanup EXIT

if [ "$BUILD" = 1 ]; then
  echo "== build =="
  go build -o "$BIN" . || { echo "BUILD FAILED"; exit 1; }
fi

if [ "$EXTERNAL" = 0 ]; then
  echo "== launch (logs -> /tmp/skywalker_engine.log) =="
  ./"$BIN" >/tmp/skywalker_engine.log 2>&1 &
  PID=$!
  for _ in $(seq 1 40); do
    curl -s -o /dev/null "$BASE/" 2>/dev/null && break
    sleep 0.5
  done
fi

fail=0
probe() { # path  want-code  label
  local code; code=$(curl -s -o /tmp/sw_body.html -w "%{http_code}" "$BASE$1")
  if [ "$code" = "$2" ]; then printf "  [OK]   %-30s %s\n" "$3" "$code"
  else printf "  [FAIL] %-30s got %s want %s\n" "$3" "$code" "$2"; fail=1; fi
}

echo "== smoke =="
probe "/"                          200 "GET /"
grep -q "Celeritas" /tmp/sw_body.html || { echo "  [FAIL] home missing 'Celeritas'"; fail=1; }
probe "/files/upload"              200 "GET /files/upload"
probe "/list-fs"                   200 "GET /list-fs (no fs)"
probe "/public/images/celeritas.jpg" 200 "GET /public/ (static)"
probe "/nope"                      404 "GET /nope (404)"

# nosurf CSRF: a token-less POST must be rejected with 400 (this is expected, not a bug)
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/files/upload")
if [ "$code" = "400" ]; then printf "  [OK]   %-30s %s (nosurf CSRF rejects token-less POST)\n" "POST /files/upload" "$code"
else printf "  [FAIL] %-30s got %s want 400\n" "POST /files/upload" "$code"; fail=1; fi

echo
[ "$fail" = 0 ] && echo "SMOKE PASSED" || { echo "SMOKE FAILED — see /tmp/skywalker_engine.log"; }
exit $fail
