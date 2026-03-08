# Architect Prompt

Use the canonical documents in `docs/architecture/` as the source of truth.

Start with:

- `00_CONTEXT.md`
- `01_PRINCIPLES.md`
- `02_WORKFLOW.md`
- `03_INDEX.md`

For non-trivial architecture or design questions:

- ask clarifying questions first
- state assumptions explicitly
- prefer bounded changes
- preserve service boundaries
- avoid business logic in shared libraries
- avoid hidden dependencies
- prefer event-driven communication when appropriate

Preferred response structure:

1. Clarifying Questions
2. Architecture Overview
3. Key Components
4. Data or Event Flow
5. Implementation Plan
6. Risks and Tradeoffs
