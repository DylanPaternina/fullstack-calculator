# AI-assisted development and prompts

The project used **OpenAI Codex (GPT-6)** for application code, tests, documentation and targeted corrections, under the owner's technical direction. The owner defined requirements, mathematical behavior, architectural choices and review priorities.

## Development prompts

This record focuses on the substantive implementation prompts. The linked documents preserve their original wording.

| Stage | Original prompt | Scope |
| --- | --- | --- |
| Initial implementation | [Master prompt](docs/MASTER_PROMPT.md) | React/TypeScript and Go calculator, architecture, API, testing, documentation and acceptance criteria |
| Advanced functionality | [Final improvements](docs/FINAL_IMPROVEMENTS_PROMPT.md) | Targeted correctness/accessibility fixes, powers, square roots, contextual percentages and two-service Docker Compose |
| Maintenance | [Final maintenance](docs/MAINTENANCE_PROMPT.md) | Safe minimatch remediation, explicit HTTP zero-result regression and full release validation |

## Implementation and validation

AI assistance supported the Go evaluator and JSON API, React interface, automated tests, documentation and focused review corrections. Automated checks used Go testing/httptest, Vitest, React Testing Library, ESLint, TypeScript, Vite, Playwright and HTTP smoke tests.

The owner's engineering rationale is documented in the [README's Design Decisions](README.md#design-decisions). Executed checks, measured results and verification limits are recorded in [VALIDATION.md](docs/VALIDATION.md) and [COVERAGE.md](docs/COVERAGE.md).
