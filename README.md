# Full-stack calculator

A classic calculator built for a Software Engineer Intern assessment. React builds and edits an expression; Go alone parses and evaluates it. Supports `+`, `-`, `*`, `/`, standard precedence, left associativity, decimals, unary signs, and whitespace.

## Run locally

Install **Go 1.26+**, **Node.js 24 LTS** (22.13+ in the 22.x line or 26+ also supported by the dependency engine constraints), and npm. Tested with Go 1.27.3, Node 24.22.1, and npm 12.2.1. No database, environment file, or backend dependency installation is required.

From the repository root, start the API in one terminal:

```sh
cd backend
go run ./cmd/server
```

In a second terminal:

```sh
cd frontend
npm ci
npm run dev
```

Open **http://127.0.0.1:5173**. The API binds to **127.0.0.1:8080** by default. Vite proxies `/api` to that address, so browser requests stay same-origin without permissive CORS.

`ADDR` configures the Go listen address, e.g. `ADDR=127.0.0.1:8090 go run ./cmd/server` on POSIX or `$env:ADDR='127.0.0.1:8090'; go run ./cmd/server` in PowerShell. Update the proxy target in `frontend/vite.config.ts` if changing the API port.

## Use

- Click the numeric/operator buttons or type into the editable expression field. Use `*` and `/` for multiplication and division.
- `Enter` or `=` calculates, including when a keypad button has focus; `Escape` clears; `Backspace` deletes. Space activates the focused button. Other input fields, links, editable regions, buttons outside the calculator, and modified shortcuts retain native behavior.
- After a result, an operator button/global shortcut continues from the full backend value; a number starts a new expression. Editing the expression field or deleting edits the original expression. After an error, edit and retry.
- Pending calculations disable equals. Editing or clearing aborts and invalidates the pending request. A request identity check rejects stale responses even when cancellation is ignored by the transport. Requests time out after 10 seconds.

## API

`POST /api/calculate` with JSON:

```json
{"expression":"5 + 3 * 2"}
```

Response (`200`):

```json
{"result":11}
```

```sh
curl -i http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"5 + 3 * 2"}'
curl -i http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"5 / 0"}'
```

PowerShell equivalent (avoids shell-specific curl quoting):

```powershell
Invoke-RestMethod -Uri http://127.0.0.1:8080/api/calculate -Method Post -ContentType application/json -Body '{"expression":"5 + 3 * 2"}'
```

Division by zero returns `400`:

```json
{"error":{"code":"division_by_zero","message":"Cannot divide by zero"}}
```

| Status | Error codes | Meaning |
| --- | --- | --- |
| 400 | `invalid_request` | Empty/invalid JSON, missing/null/non-string expression, unknown fields, extra JSON values |
| 400 | `invalid_expression` | Empty expression, unsupported character, malformed number or syntax |
| 400 | `division_by_zero` | Division by positive or negative zero |
| 400 | `numeric_overflow` | Literal or intermediate result cannot be represented as finite float64 |
| 400 | `expression_too_large` | More than 1,024 UTF-8 bytes |
| 413 | `request_too_large` | Body exceeds 8,192 bytes when decoded/read |
| 405 | `method_not_allowed` | Endpoint requires POST; includes `Allow: POST` |
| 404 | `not_found` | Unknown endpoint |
| 500 | `internal_error` | Defensive fallback for unexpected evaluator failures |

All application responses are JSON. The server applies 5-second header, 10-second read/write, and 60-second idle timeouts, and a 16 KiB header limit. Standard `net/http` connection-level errors (such as malformed HTTP or timeouts) can terminate a request before the JSON handler runs. Malformed oversized JSON may fail validation with 400 before the size boundary is read; valid oversized JSON/whitespace is rejected with 413. JSON decoding follows Go's standard field matching and last-value behavior for duplicate keys.

## Structure and decisions

```text
backend/
  cmd/server/main.go                 Single HTTP service; address/timeouts
  internal/api/api.go                JSON transport and request validation
  internal/evaluator/evaluator.go    Restricted arithmetic parser; no HTTP
frontend/
  src/App.tsx                       Calculator state, keyboard and presentation
  src/api.ts                        Fetch and runtime response validation
  src/styles.css                    Responsive layout and focus states
  e2e/calculator.spec.ts             Real-browser integration checks
scripts/smoke.mjs                    Real HTTP checks, API and Vite proxy
docs/                               Coverage, validation evidence, original prompt
PROMPTS.md                          Actual AI instructions and usage disclosure
```

- **React + strict TypeScript:** small declarative UI with built-in state; runtime JSON validation protects the actual HTTP boundary. No frontend arithmetic evaluator or `eval`.
- **Go + standard-library HTTP:** one independently runnable service, no third-party Go dependencies. `httptest` exercises transport independently of a running server.
- **Restricted recursive-descent parser:** expression → term → number. Repetition and unary signs use loops, so input cannot create deep recursion. Work and memory are bounded by the 1,024-byte input limit. No expression library, database, service wrappers, or general-purpose execution.
- **Separation:** presentation, HTTP communication, transport validation, and mathematics each have a clear responsibility and independent tests.
- **Dependencies:** React 19.3.1, Vite 8.3.5, TypeScript 5.9.3, Vitest/coverage-v8 5.0.5, React Testing Library 16.3.3, ESLint 10.14.0, Playwright 1.64.1. `frontend/package-lock.json` pins the full dependency tree; use `npm ci`.

## Numeric policy and scope

- Arithmetic uses IEEE-754 `float64`. Decimal arithmetic is approximate: the API returns `0.30000000000000004` for `0.1+0.2`. This is a general calculator, not a financial/arbitrary-precision tool.
- Results display up to **12 significant digits** without unnecessary trailing zeros. The unrounded backend value remains in state and is used when continuing a calculation; its full string is also available as the result tooltip. Presentation rounding never feeds arithmetic.
- Decimal literals accept `12`, `.5`, `1.`, and scientific notation such as `1e-7`, needed to round-trip very small/large backend results. `e` is entered through the expression field; it is not exponentiation.
- Unary plus/minus may repeat: `2--3 = 5`, `--3 = 3`. Adjacent numbers, implicit multiplication, hexadecimal, commas, and non-ASCII operator glyphs are rejected. Button labels `×`, `÷`, and `−` insert ASCII operators.
- Division by zero and overflow at **any intermediate step** produce controlled errors. NaN and Infinity are never successful results. Ordinary IEEE underflow can round to zero; finite subnormal values are accepted. Negative zero may display as zero.
- The frontend limits input to 1,024 JavaScript characters; Go enforces the authoritative 1,024 UTF-8 byte limit, so multibyte input may reach it earlier.
- Parentheses, powers, square root, and percentage are intentionally omitted to keep the mandatory grammar small and well-tested. Docker is omitted: two local processes and a lockfile are sufficient for this scope.

## Tests and build

```sh
cd backend
gofmt -w .
go vet ./...
go test ./...
go test ./... '-coverprofile=coverage.out'
go tool cover '-func=coverage.out'
go tool cover '-html=coverage.out' -o coverage.html
go test ./internal/evaluator '-fuzz=FuzzEvaluate' '-fuzztime=5s'
go build ./cmd/server
```

Quote the complete Go coverage flags as shown, including in PowerShell.

```sh
cd frontend
npm run lint
npm run typecheck
npm test
npm run test:coverage
npm run build
```

Frontend coverage: `frontend/coverage/index.html`; backend coverage: `backend/coverage.html`. Actual results, included files, and limits are in [docs/COVERAGE.md](docs/COVERAGE.md). Build output is `frontend/dist/`. `npm run preview` previews that static build only; the Go API remains a separate service. For a deployed build, configure a same-origin reverse proxy for `/api`; Vite's dev proxy is not a production server.

With **both servers running**, from the repository root:

```sh
node scripts/smoke.mjs
cd frontend
npx playwright install chromium
npm run test:integration
```

Three normal browser tests verify real calculations, server errors and recovery, button/keyboard interaction, 1440/768/375 px layouts, focus, and an injected transport failure. A fourth test deliberately requires the API to be stopped. Keep Vite running, stop Go, then run:

```sh
CALCULATOR_BACKEND_DOWN=1 npm run test:integration -- --grep 'stopped backend'
```

PowerShell: `$env:CALCULATOR_BACKEND_DOWN='1'; npm run test:integration -- --grep 'stopped backend'`, then `Remove-Item Env:CALCULATOR_BACKEND_DOWN`. Restart Go afterward. The normal run explicitly skips this opt-in test. Test failure traces and viewport screenshots are generated in `frontend/test-results/` and ignored by Git.

## Validation and limitations

See [docs/VALIDATION.md](docs/VALIDATION.md) for executed checks and requirement evidence. Browser checks used Chromium on Windows; other engines, real mobile devices, and screen-reader behavior were not exhaustively tested. Semantics, labels, live announcements, native keyboard activation, visible focus, and responsive sizing are implemented. No public deployment, remote repository, or hosting configuration is included.

AI-assisted implementation and the exact original master prompt are disclosed in [PROMPTS.md](PROMPTS.md). Publication requires explicit owner authorization.
