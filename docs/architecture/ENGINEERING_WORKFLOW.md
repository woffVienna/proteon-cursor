# Proteon Engineering Workflow

This document describes the intended workflow for development and AI-assisted work.

## 1. Canonical Development Workflow

Prefer Make targets over ad-hoc commands.

Core root targets:

- `make setup`
- `make generate`
- `make test`
- `make check`
- `make stack-up`
- `make stack-down`

Core service targets:

- `make generate`
- `make dev`
- `make run`
- `make test`
- `make build`
- `make containerise`

## 2. Local Runtime Standard

Local development runtime standard is:

- `k3d`
- `kubectl`
- `helm`

Do not introduce docker-compose orchestration unless explicitly agreed.

## 3. AI-Assisted Workflow

### Architecture

Use ChatGPT for:

- architecture reasoning
- design intent
- option analysis
- tradeoff discussion
- implementation strategy
- architecture briefs

### Planning

Use Cursor Plan for:

- repo-specific change planning
- affected file identification
- bounded implementation steps
- architecture-to-repository translation

### Execution

Use Cursor Agent for:

- bounded code changes
- incremental implementation
- minimal diffs
- changes aligned with architecture documents and `.cursorrules`

## 4. Shared Intent Layer

The shared intent layer for ChatGPT and Cursor is:

- `ENGINEERING.md`
- `DEV.md`
- `.cursorrules`
- `docs/architecture/`

`docs/architecture/` contains the durable architecture context and design intent.

## 5. Working Style for Larger Changes

For non-trivial changes:

1. clarify the problem
2. discuss architecture in ChatGPT
3. produce/update an architecture brief
4. save the result in the repo
5. let Cursor Plan map it to repo changes
6. let Cursor Agent execute in bounded steps
7. update documentation if needed

## 6. Documentation Discipline

Do not rely on chat history as the source of truth.

Persist important outcomes in repository documents.