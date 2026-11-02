# AI assistance and actual prompts

This project was implemented with **OpenAI Codex (GPT-6)**. AI assistance produced application code, tests, documentation, and review fixes. It must not be represented as exclusively manual work. The exact runtime model variant is not available in the project record.

## Actual user instructions

The following goal was supplied verbatim (the `/goal` command introduced it):

> Build a complete, clean, maintainable, and thoroughly tested React + TypeScript and Go calculator that satisfies the technical assessment requirements. Verify all mandatory functionality, tests, coverage, integration, and documentation. Do not publish without explicit authorization.

The complete, original implementation instruction is preserved verbatim in [docs/MASTER_PROMPT.md](docs/MASTER_PROMPT.md), copied from the supplied text attachment. It includes the functional contract, architecture, testing, correction workflow, senior-review instructions, publication boundaries, and definition of done. It is the actual prompt, not a reconstructed summary.

No additional user correction prompts or separate AI reviewer prompts were used during this implementation. The agent reviewed its own code under the master prompt's Phase 6 directive; no independent reviewer is claimed. User preferences required responses in the user's language; application text and repository documentation follow the master prompt's English requirement.

## Actions performed (summary, not prompts)

- Inspected the empty workspace and local tooling, selected compatible releases from official registries, and downloaded project-local Go/npm tooling. Go's archive checksum was checked against the official manifest.
- Implemented a Go arithmetic parser, JSON API, React calculator, and Fetch boundary.
- Used Go tests/httptest, fuzzing, Vitest, React Testing Library, ESLint, TypeScript, Vite, Node HTTP smoke checks, and Playwright/Chromium against real local processes.
- Corrected findings from test/static-check execution: editable-region keyboard isolation, ambiguous test selectors, Fetch error cause preservation, and TypeScript test-query options. Review also tightened Enter behavior after button clicks, unmount cancellation, and error-description associations. Corrected environment issues with Go temporary directories and PowerShell flag parsing; reran affected checks.
- Reviewed generated viewport screenshots and code responsibilities, parser bounds/precedence, numeric handling, stale-request safety, dependency scope, and repository contents.

Tools included terminal/patch tools, official registry queries/web lookup, local file/image inspection, the Sites local-preview guide, and the Codex preview panel. No Sites registration, publication, or external messaging occurred. Actual executed validation and remaining limits are recorded in [docs/VALIDATION.md](docs/VALIDATION.md) and [docs/COVERAGE.md](docs/COVERAGE.md).

Internal model reasoning and system/developer instructions are not part of this disclosure. No hidden reasoning or invented development prompts are included.
