---
name: run-skywalker-engine
description: Build, launch, run, serve, and smoke-test the skywalker_engine Go web app (chi router + Jet templates, the Skywalker framework demo). Use when asked to run, start, serve, boot, or screenshot skywalker_engine, or to verify the demo app compiles and its routes respond.
---

# Run skywalker_engine

`skywalker_engine` (Go module `myapp`) is the demo web app for the **Skywalker** framework:
a chi router serving server-rendered Jet templates, with a remote-filesystem (MinIO) upload/list/delete
demo. It is **server-rendered HTML**, so the driver is **curl**, not a browser. The committed driver
[`.claude/skills/run-skywalker-engine/smoke.sh`](smoke.sh) builds it, launches it, probes its routes,
and tears it down — that is the programmatic handle on the running app.

Paths below are relative to the repo root (`skywalker_engine/`). The default host here is Windows/PowerShell;
the driver runs under the Bash tool (Git Bash). Go 1.25 is installed (module declares `go 1.18`).

## Prerequisites

- **Go** (1.18+; 1.25 verified) and **curl** (preinstalled in Git Bash).
- The **framework checked out as a sibling**: `../skywalker` (go.mod has `replace github.com/stefanlester/skywalker => ../skywalker`). On this Windows host the actual folder `..\Skywalker` resolves fine.
- No database, Redis, or MinIO required to boot — `.env` ships with `SESSION_TYPE=cookie`, empty `DATABASE_TYPE`, empty `CACHE`. (MinIO is only touched by the FS routes; see Gotchas.)

## Run (agent path) — use this

One command does everything (build, launch on :4000, smoke-test, stop):

```bash
bash .claude/skills/run-skywalker-engine/smoke.sh
```

Verified output:

```
== build ==
== launch (logs -> /tmp/skywalker_engine.log) ==
== smoke ==
  [OK]   GET /                          200
  [OK]   GET /files/upload              200
  [OK]   GET /list-fs (no fs)           200
  [OK]   GET /public/ (static)          200
  [OK]   GET /nope (404)                404
  [OK]   POST /files/upload             400 (nosurf CSRF rejects token-less POST)

SMOKE PASSED
```

Flags: `--no-build` (reuse the existing binary), `--external` with `BASE=http://host:port` (probe a server you started yourself). Exit code is 0 only if every probe passes.

### Launch it yourself and poke individual routes

```bash
go build -o skywalker_engine_app.exe .   # from repo root
./skywalker_engine_app.exe &             # logs: "Listening on port 4000"
curl -s http://localhost:4000/ | grep Celeritas
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:4000/files/upload
taskkill //F //IM skywalker_engine_app.exe   # stop it (Windows)
```

The home page renders "Celeritas / Go build something awesome". `/files/upload` serves the upload form
(`action="/files/upload"`, a `csrf_token` hidden field, `type="file"`). `/list-fs` with no query renders
an empty listing.

## Run (human path)

`go run .` from the repo root, then open `http://localhost:4000/` in a browser. Ctrl-C to stop.
Same thing the driver does, minus the automated probes.

## Gotchas (battle scars from this container)

- **POST routes are CSRF-protected (`nosurf`).** A token-less `POST /files/upload` returns **400** — that's expected, not a bug. To drive a real upload you must first `GET /files/upload`, scrape the `csrf_token` hidden input **and** the session cookie, then POST both as `multipart/form-data` with a `formFile` field and `upload-type=MINIO`.
- **You must launch from the repo root.** The app calls `os.Getwd()` and loads `./.env` plus `./views` from there. Run it anywhere else and config/templates won't load. The driver `cd`s to root for you.
- **MinIO-backed routes need a MinIO server** at `127.0.0.1:9000` (`docker-compose up -d` brings up the configured `testbucket`). Affected: `/list-fs?fs-type=MINIO`, `/test-minio`, the actual upload POST, and `/delete-from-fs`. The smoke test deliberately avoids these so it needs zero services. (`GET /files/upload?type=MINIO` still returns 200 — it only renders the form; the MinIO call happens on submit.)
- **Case-sensitive module path.** The `replace => ../skywalker` resolves on Windows (case-insensitive FS) where the folder is `Skywalker`. On Linux it won't — `ln -s Skywalker skywalker` in the parent dir, or fix the path in go.mod.
- **`.env` carries `KEY` (must be exactly 32 chars)** and MinIO creds. It is in the working tree; the home/upload-form/list paths boot with cookie sessions and no DB.
- **`/skywalker` and `/api` return 404 at their root** — they're mount points with no index route. Not an error.

## Troubleshooting

- **App exits with code 1 immediately, log shows a port bind error** → port 4000 is already taken by a stale run. `taskkill //F //IM skywalker_engine_app.exe`, then retry.
- **`SMOKE FAILED`** → read `/tmp/skywalker_engine.log` for the startup stack trace (missing `.env`, template parse error, etc.).
- **Build error `cannot find module ../skywalker`** (Linux) → case sensitivity; symlink as above.
