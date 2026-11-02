# Validation record

Executed on **2026-11-02**, Windows amd64. Tools: Go 1.27.3, Node 24.22.1, npm 12.2.1, React 19.3.1, TypeScript 5.9.3, Vite 8.3.5, Vitest 5.0.5, ESLint 10.14.0, Playwright 1.64.1 / Chromium 156.

## Final quality gates

| Command/check | Actual outcome |
| --- | --- |
| `gofmt -w .` / formatting inspection | Go files formatted |
| `go vet ./...` | Pass |
| `go test ./...` (also executed with coverage) | All packages pass |
| `go test ./... '-coverprofile=coverage.out'` | Pass; 94.8% total statements |
| `go tool cover '-func=coverage.out'` | Function report generated |
| `go tool cover '-html=coverage.out' -o coverage.html` | HTML report generated |
| `go test ./internal/evaluator '-fuzz=FuzzEvaluate' '-fuzztime=5s'` | Pass; 71,764 executions |
| `go build -o ../.tools/calculator-server.exe ./cmd/server` | Pass; executable started for integration |
| `npm run lint` | Pass; no lint errors |
| `npm run typecheck` and production build's `tsc --noEmit` | Pass |
| `npm run test:coverage` | 30 tests pass; 100% lines, 98.82% branches |
| `npm run build` | Pass; generated `dist/` |
| `node scripts/smoke.mjs` | Pass; 24 real HTTP requests across API and Vite proxy |
| `npm run test:integration` | 3 Chromium tests pass; intentional stopped-backend test skipped in normal run |
| `CALCULATOR_BACKEND_DOWN=1 npm run test:integration -- --grep 'stopped backend'` | 1 test passes with Go actually stopped |
| README curl examples | Observed 200/result 11 and 400/division-by-zero JSON |

The production bundle was approximately 225.3 kB JavaScript (70.8 kB gzip) and 3.1 kB CSS (1.3 kB gzip). No unsupported performance claim or premature optimization was used.

## Requirement evidence

| Mandatory area | Evidence |
| --- | --- |
| Addition, subtraction, multiplication, division | Evaluator cases, HTTP smoke cases, Chromium calculations |
| Precedence, left associativity, negatives, unary signs, decimals | Independent Go expected-result tests; real-browser negative/decimal/precedence checks |
| All arithmetic in Go | `App.tsx` constructs expressions; `api.ts` posts them; `evaluator.go` parses and calculates. Component test deliberately displays a mocked backend result that differs from local arithmetic |
| Invalid expressions, division by zero, overflow, no non-finite JSON | Evaluator/handler cases; real HTTP controlled errors; frontend response validation |
| Classic calculator, buttons, editable expression, result | RTL interaction tests and Chromium button/keyboard flow |
| Keyboard, clear/delete, loading/error states | RTL behavior tests; real-browser keyboard and error/recovery flows |
| Duplicate submissions, outdated responses | Deferred-promise RTL tests, request abort signals, identity checks and timeout test |
| REST/JSON, validation, size/method limits | `httptest` table and independent requests against running Go |
| Responsive/accessibility basics | Semantic controls, labels, described errors, live regions, focus CSS; browser focus/navigation and viewport checks |
| Tests/coverage/static checks/build | Executed quality gates above and retained coverage evidence |
| Setup, API, decisions, assumptions, AI prompts | `README.md`, `PROMPTS.md`, original `docs/MASTER_PROMPT.md`, `docs/COVERAGE.md` |
| Local version control / publication boundary | Local Git repository; generated tooling/dependencies excluded; no remote created or push performed |

## Actual full-stack integration

Started the compiled Go binary at `127.0.0.1:8080` and Vite at `127.0.0.1:5173`. The Chromium arithmetic test exercised **browser → React → Fetch → Vite proxy → real Go evaluator → JSON → rendered result** without mocked calculation responses.

Verified `5+3`, `5-3`, `5*3`, `10/4`, `5+3*2`, `10-6/2`, `-3*2`, `2*-3`, and `0.1+0.2`; `5/0` and `5+*2` produced HTTP 400 and visible errors. The same browser session recovered and calculated again through numeric/operator buttons and physical keyboard input. No page JavaScript errors occurred in that calculation flow.

The HTTP smoke script independently verified nine successful calculations plus three controlled errors at each of the API and proxy addresses (24 requests). It includes numeric overflow. This is separate from mocked component tests.

Failure checks were distinguished explicitly:

1. A browser route abort injected a transport failure, verified its visible error, removed the injection, and verified a successful real calculation.
2. Go was actually stopped while Vite remained running. A new browser calculation received the proxy's **502** response, displayed the invalid-server-response error, and re-enabled equals. This test passed separately. Go was then restarted and the normal browser and HTTP suites passed again.

## Responsive and accessibility review

Chromium checked **1440, 768, and 375 px** widths at 900 px viewport height. There was no horizontal overflow with either the empty UI or a large scientific-notation result. Equals remained in view. Full-page screenshots were generated in `frontend/test-results/` and visually inspected at desktop/mobile sizes. Longer pages may scroll vertically.

Tab navigation reaches the labeled expression input and semantic buttons; a real-browser assertion checks visible focus styling. Space activates focused buttons; Enter calculates within the calculator. Input errors have `aria-invalid`, linked descriptive text, and a live alert; status and result use live semantics. Keyboard shortcuts leave unrelated editable/interactive elements alone.

These checks are not a comprehensive accessibility audit. VoiceOver/NVDA, other browser engines, real mobile devices, touch hardware, and every zoom/font preference were not tested.

## Review and resolved findings

The implementation agent reviewed parser grammar, bounds, intermediate overflow, API decoding, native keyboard interactions, request lifecycle, rounding/continuation, runtime JSON validation, dependency scope, and documentation. This was a self-review, not a separate independent review.

Resolved findings included isolation of editable regions, Enter after a keypad click, cancellation identity on unmount, described errors, preserving the underlying cause of Fetch failures, and meaningful malformed multiplicative-expression tests. Test-only findings (ambiguous status queries and an unsupported query option) were corrected without weakening validation.

Initial environment failures were diagnosed and corrected before recording passes: missing Go/npm were installed only under ignored `.tools/`; a Windows temporary-directory permission error was resolved with project-local `GOCACHE`/`GOTMPDIR`; coverage flags were quoted to avoid PowerShell splitting the output filename; local networking required execution outside network isolation. The proxy's actual stopped-backend status was 502, and its test was corrected accordingly.

## Deliberate scope and remaining limits

Parentheses, exponentiation, square root, percentage, Docker, database, history, authentication, remote hosting, and extra services were omitted. The README specifies float64 approximation, underflow, result formatting, repeated unary signs, decoder behavior, input limits, and the development proxy boundary. The service is intended for local assessment; external deployment would require its own hosting/reverse-proxy configuration.

No mandatory functionality remains unverified within the stated scope. Unit coverage does not cover server bootstrap or the impossible-at-present internal-error fallback; integration coverage is reported separately. Publication remains subject to explicit owner authorization.
