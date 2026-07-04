# skywalker_engine

The demo / skeleton application for the [**Skywalker**](https://github.com/stefanlester/Skywalker)
Go web framework. This is also the starter that `skywalker new <app>` clones — so it doubles as the
reference for how a Skywalker app is wired.

It's a small [chi](https://github.com/go-chi/chi) + [Jet](https://github.com/CloudyKit/jet) app that
demonstrates server-rendered pages and the framework's remote-filesystem layer (upload / list /
delete across MinIO, S3, SFTP, and WebDAV).

## Requirements

- **Go 1.23+**.
- The framework checked out as a **sibling directory**: `../skywalker`
  (`go.mod` has `replace github.com/stefanlester/skywalker => ../skywalker`).
- A running filesystem backend (e.g. SeaweedFS) only for the file routes; the home page and forms need nothing.

## Run it

The repo ships a run skill that builds, launches, and smoke-tests the app:

```bash
bash .claude/skills/run-skywalker-engine/smoke.sh
```

Or run it directly:

```bash
go run .            # listens on PORT from .env (default 4000)
```

Then open <http://localhost:4000>. The app boots standalone — `.env` defaults to
`SESSION_TYPE=cookie` with no database and no cache.

## Routes

| Route | Description |
|---|---|
| `GET /` | Home page |
| `GET /list-fs?fs-type=<T>` | List files in a backend (`T` = `MINIO`/`S3`/`SFTP`/`WEBDAV`) |
| `GET /files/upload?type=<T>` | Upload form |
| `POST /files/upload` | Upload a file to the chosen backend |
| `GET /delete-from-fs?fs_type=<T>&file=<name>` | Delete a file |
| `/public/*` | Static assets |

The backend is chosen at runtime via `App.FileSystems[fsType].(filesystems.FS)`, so all four are
selectable from the same UI — each is constructed only when its `*_HOST`/`*_SECRET` env var is set
(see the framework README).

## Configuration

Copy the relevant values into `.env`. To exercise the filesystem routes locally, set `MINIO_*` and
start the bundled [SeaweedFS](https://github.com/seaweedfs/seaweedfs) service (`docker compose up -d
seaweedfs`, then create the bucket — see the comment in `docker-compose.yml`). The `MINIO` backend
speaks the S3 protocol, so any S3-compatible store works (SeaweedFS, Garage, MinIO, Cloudflare R2…);
SeaweedFS is the bundled default because the MinIO community edition is no longer maintained.
`KEY` must be exactly 32 characters.

> **Security:** never commit real secrets. `.env` and `db-data/` are gitignored; if you clone an
> older revision, rotate any keys that were previously tracked.

## Notes for contributors

This tree is generated output when produced via `skywalker new`. To change app behavior, edit the
handlers/routes/views here; to change framework capabilities, work in `../skywalker`. Build/verify
with `go build ./...`, `go vet ./...`, and the run skill above.
