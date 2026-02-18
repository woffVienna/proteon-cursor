# Identity service

## Runtime config model

Identity resolves configuration from environment variables at startup in
`internal/platform/config`.

- Runtime mode is controlled by `RUNTIME_MODE` (`local`, `docker`, `cloud`).
- General env loading is provided by `libs/platform/config`.
- Shared loader returns a typed config envelope; service-specific values live
  under `Config.Service` (for identity: `Config.Service.JWT`).
- Runtime mode selects one env file:
  - `local` -> `.env.local`
  - `docker` -> `.env.docker`
  - `cloud` -> `.env.cloud` (optional; cloud usually injects env directly)
- Key names stay identical across environments (`PORT`, `DB_DSN`, `JWT_ISSUER`, ...).

Resolution order:

1. Existing process env (already set by shell/docker/cloud)
2. Value from selected `.env.<mode>` file
3. Built-in service default (if defined)

## Port convention

- Local host run (`make run` / `make dev`): service listens on `8081`
- Docker container: service still listens on `8081`
- Docker host published port: `9091:8081`

This keeps service-internal ports stable and avoids clashes with local host
processes.

## Local run

`make run` and `make dev` start with `RUNTIME_MODE=local` by default and load
`.env.local` via the shared loader.

## Docker run

`tools/docker/compose.services.yml` uses `env_file` with
`services/identity/.env.docker` and sets `RUNTIME_MODE=docker`.
