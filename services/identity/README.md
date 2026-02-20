# Identity service

## Runtime config model

Identity resolves configuration from environment variables at startup in
`internal/platform/config`.

- General env loading is provided by `libs/platform/config`.
- Shared loader returns a typed config envelope; service-specific values live
  under `Config.Service` (for identity: `Config.Service.JWT`).
- For host-mode development, `.env.local` is loaded.
- Key names stay identical across environments (`PORT`, `DB_DSN`, `JWT_ISSUER`, ...).

Resolution order:

1. Existing process env
2. Value from `.env.local` (host-mode development)
3. Built-in service default (if defined)

## Port convention

- Local host run (`make run` / `make dev`): service listens on `8081`
- Kubernetes container: service listens on `8081`

This keeps service-internal ports stable and avoids clashes with local host
processes.

## Local run

`make run` and `make dev` load `.env.local` via the shared loader.

## Kubernetes run

Kubernetes deployment uses the Helm chart in `infra/k8s/charts/identity`.
The chart provides runtime keys via ConfigMap.
