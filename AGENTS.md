# Repository Guidelines

## Project Context & Goals

This repository is a mentor-driven exercise: build a simplified “cloud drive” (Yandex.Disk-like) service that can handle large files.

- Operations: upload, download, list, delete user files.
- Limits: files up to **1 GB**; stream I/O and avoid buffering whole files in memory.
- Multi-user: several users, but **no explicit auth yet** (user identification mechanism is TBD).
- Process: tasks are intentionally abstract; requirement clarification and documenting assumptions is part of the work.
- Mentor brief: `docs/mentor-brief.md`.

## Project Structure & Module Organization

- `cmd/disk-service/`: service entrypoint (`main.go`).
- `internal/`: application code not meant to be imported externally.
  - `internal/lib/`: shared helpers (example: `internal/lib/logger/sl`).
  - `internal/config/`: YAML/env configuration loading.
  - `internal/repository/`: storage abstractions and errors.
  - `internal/repository/postgres/`: Postgres-backed repository implementation.
- `config/`: local YAML configs (example: `config/local.yaml`).
- `docs/`: design notes and brief specs.

## Build, Test, and Development Commands

This repo uses plain Go tooling (no Makefile).

- Run locally (requires `CONFIG_PATH`): `CONFIG_PATH=config/local.yaml go run ./cmd/disk-service`
- Prereqs: a running Postgres reachable via `database_url` in your config.
- Build: `go build ./...`
- Tests: `go test ./...`
- Quick checks: `go vet ./...`
- Format: `gofmt -w .`
- Dependency hygiene: `go mod tidy`

## Coding Style & Naming Conventions

- Formatting is Go-standard: run `gofmt` before pushing.
- Indentation: tabs/spaces are handled by `gofmt` (don’t hand-format).
- Packages: short, lowercase names (e.g., `config`, `repository`).
- Exported identifiers: `CamelCase` with doc comments when public outside the package.

## Testing Guidelines

- Test files: `*_test.go`; test funcs: `TestXxx`.
- Prefer table-driven unit tests; avoid hitting a real database by default.
- If you add DB integration tests, gate them behind a build tag (example: `//go:build integration`) and document how to run them.

## Commit & Pull Request Guidelines

- Commit messages follow the existing history: short, imperative, sentence case (e.g., `Add slog logger`).
- PRs should include:
  - What changed + why (link an issue if applicable).
  - How to run/verify (commands and any required config changes).
  - Notes on schema/config updates (update `config/local.yaml` when config fields change).

## Configuration & Security Tips

- `CONFIG_PATH` is required at runtime; `database_url` must be set in the YAML.
- Don’t commit real credentials. Use local values in `config/local.yaml` and keep secrets in environment variables (see `.gitignore` for ignored `.env*` patterns).
