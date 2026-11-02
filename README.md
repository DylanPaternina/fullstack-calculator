# Full-stack calculator

A classic calculator built for a Software Engineer Intern assessment. React builds and edits an expression; Go alone parses and evaluates it. Supports `+`, `-`, `*`, `/`, `^`, `sqrt(...)`, contextual `%`, parentheses, decimals, unary signs, and whitespace.

## Repository and version history

Source: [DylanPaternina/fullstack-calculator](https://github.com/DylanPaternina/fullstack-calculator).

| Milestone | Commit | Contents |
| --- | --- | --- |
| [assessment-basic](https://github.com/DylanPaternina/fullstack-calculator/tree/assessment-basic) | `7054b3b` | Original four-operation calculator, tests and documentation |
| [assessment-advanced](https://github.com/DylanPaternina/fullstack-calculator/tree/assessment-advanced) | `c55f2ef` | Targeted test/focus fixes, powers, roots, contextual percentages and Docker Compose |

The basic version was published first, followed by the advanced version using the existing commit history. `main` contains the advanced implementation and subsequent publication documentation. The tags retain the two implementation milestones for comparison.

```sh
git clone https://github.com/DylanPaternina/fullstack-calculator.git
cd fullstack-calculator
```

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

- Click the numeric/operator buttons or type into the editable expression field. Use `*` and `/` for multiplication/division, `^` for powers, and lowercase `sqrt(...)` for square roots. The `√` button inserts `sqrt(`; close it with `)`.
- `Enter` or `=` calculates, including when a keypad button has focus; `Escape` clears; `Backspace` deletes. Space activates the focused button. Other input fields, links, editable regions, buttons outside the calculator, and modified shortcuts retain native behavior.
- After a result, a binary operator or `%` continues from the full backend value; numbers, `(` and `√` start a new expression. A negative result used as a power base is parenthesized: continuing from -2 with `^2` submits `(-2)^2`. Editing the expression field or deleting edits the original expression. Backspace removes one character, including inside `sqrt(`. After an error, edit and retry.
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
| 400 | `domain_error` | Negative square root or a power outside the supported real-number domain |
| 400 | `expression_too_large` | More than 1,024 UTF-8 bytes |
| 400 | `expression_too_complex` | More than 64 simultaneously nested groups, square roots, and power right-hand sides |
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
- **Restricted recursive-descent parser:** expression → term → unary → power → postfix → primary/number. Additive/multiplicative operations and unary signs use loops; recursive grouping and powers share a depth limit of 64. Work and memory are bounded by that limit and the 1,024-byte input limit. No expression library, database, service wrappers, or general-purpose execution.
- **Percentage context:** each parsed operand carries its numeric value and a `directPercent` marker. Unary signs and grouping preserve that marker; binary operations and `sqrt` clear it. This distinguishes a direct percentage from a term merely containing one, without an AST, text replacement, or a second evaluation pass.
- **Separation:** presentation, HTTP communication, transport validation, and mathematics each have a clear responsibility and independent tests.
- **Dependencies:** React 19.3.1, Vite 8.3.5, TypeScript 5.9.3, Vitest/coverage-v8 5.0.5, React Testing Library 16.3.3, ESLint 10.14.0, Playwright 1.64.1. `frontend/package-lock.json` pins the full dependency tree; use `npm ci`.

## Numeric policy and scope

- Arithmetic uses IEEE-754 `float64`. Decimal arithmetic is approximate: the API returns `0.30000000000000004` for `0.1+0.2`. This is a general calculator, not a financial/arbitrary-precision tool.
- Results display up to **12 significant digits** without unnecessary trailing zeros. The unrounded backend value remains in state and is used when continuing a calculation; its full string is also available as the result tooltip. Presentation rounding never feeds arithmetic.
- Decimal literals accept `12`, `.5`, `1.`, and scientific notation such as `1e-7`, needed to round-trip very small/large backend results. `e` is entered through the expression field; it is not exponentiation.
- Unary plus/minus may repeat: `2--3 = 5`, `--3 = 3`. Adjacent numbers, implicit multiplication, hexadecimal, commas, and non-ASCII operator glyphs are rejected. Button labels `×`, `÷`, and `−` insert ASCII operators.
- Division by zero and overflow at **any intermediate step** produce controlled errors. NaN and Infinity are never successful results. Ordinary IEEE underflow can round to zero; finite subnormal values are accepted. Negative zero may display as zero.
- The frontend limits input to 1,024 JavaScript characters; Go enforces the authoritative 1,024 UTF-8 byte limit, so multibyte input may reach it earlier.

## Advanced operations and precedence

| Expression | Result | Rule |
| --- | ---: | --- |
| `2^3` | 8 | Power |
| `2^3^2` | 512 | Powers associate right to left: `2^(3^2)` |
| `-2^2` / `(-2)^2` | -4 / 4 | Power binds more strongly than a unary sign on its left |
| `2^-3` | 0.125 | Signed exponents are supported |
| `0^0` | 1 | Convention follows Go `math.Pow` |
| `sqrt(81)` | 9 | Square root of an expression |
| `2^3+sqrt(16)` | 12 | Functions and powers combine with arithmetic |
| `(2+3)*4` | 20 | Parentheses group expressions |
| `2^50%` | ≈1.41421356237 | Postfix `%` binds before power: `2^(50/100)` |

`sqrt(-1)`, `(-2)^0.5`, malformed syntax, unclosed groups, and overflow return controlled HTTP 400 errors. Negative bases require integer exponents under the supported real-valued `math.Pow` policy. Zero raised to a negative exponent uses `division_by_zero`. No NaN or Infinity is returned as a successful JSON result. Implicit multiplication (`2(3)`), other function names, and uppercase `SQRT` are unsupported.

### Contextual percentages

For a **direct** percentage on the right of `+` or `-`, use the already accumulated left value as its base. For multiplication, division and standalone values, percentage is ordinary division by 100.

| Expression | Result |
| --- | ---: |
| `50%` / `25%` | 0.5 / 0.25 |
| `200+10%` / `200-10%` | 220 / 180 |
| `200*10%` / `200/10%` | 20 / 2000 |
| `200+10%+5%` | 231 (220 plus 5% of 220) |
| `25%+1` | 1.25 |
| `sqrt(25)*20%` | 1 |
| `200+10%*2` | 200.2 (the right operand is multiplication, not a direct percentage) |
| `200+(-10%)` / `200+0.5%` | 180 / 201 |
| `200+(10%)` / `200+(5+5)%` | 220 / 220 |
| `50%%` / `(50%)%` | 0.005 / 0.005 (each `%` divides by 100 again) |
| `200+10%%` | 200.2 |
| `200+(10%+5%)` | 200.105 (the inner sum is evaluated first, using its own left base) |
| `200+2^50%` / `200+(2^3)%` | ≈201.414213562 / 216 |
| `200+sqrt(100%)` / `200+sqrt(25)%` | 201 / 210 |

Parentheses and unary signs preserve a direct percentage; a completed binary operation or square root does not. A `%` after a whole grouped result marks that result as a new direct percentage. Invalid prefix `%2` is rejected; `200/0%` is division by zero. These rules are intentionally contextual calculator semantics, not modulo or a global replacement of `%` with `/100`.

Example requests use the unchanged endpoint:

```sh
curl http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"2^3+sqrt(16)"}'
curl http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"200+10%"}'
```

Results are `{"result":12}` and `{"result":220}`.

## Focused correctness and accessibility fixes

The mathematical test comparator explicitly rejects NaN and either Infinity, requires exact zero when zero is expected (including documented underflow), and uses relative tolerance for finite nonzero values. Regression cases verify rejection of incorrect near-zero results and protect tiny/subnormal values without making ordinary decimal comparisons exact.

The expression input reserves space between its label and focus outline. The visible, contrasting keyboard indicator is retained; browser tests check that its outer top edge stays below the label at 375, 768 and 1440 px. The original layout and colors are preserved.

## Docker Compose (two services)

Install Docker with the Compose plugin separately; it is not required for normal local development. From the repository root:

```sh
docker compose config
docker compose up --build -d
docker compose ps
docker compose logs
```

Open **http://127.0.0.1:8088**. The frontend uses a Node 24 multi-stage production build and Nginx to serve assets. Nginx forwards `/api/` to `http://backend:8080` without removing the `/api` prefix. Fetch stays relative to `/api/calculate`; Docker service names never appear in browser requests. The backend uses a Go 1.27 builder and an unprivileged static binary in a `scratch` runtime, listening on `0.0.0.0:8080`. Only the frontend port is published on the host; Go stays internal to the Compose network. Build contexts allow only relevant source/configuration files, excluding tools, credentials and generated artifacts.

Base tags follow their Go/Node release line and Nginx stable Alpine image; dependency versions are fixed by `npm ci` and the lockfile. Image digests are not pinned. Official image references: [Go](https://hub.docker.com/_/golang), [Node](https://hub.docker.com/_/node), [Nginx](https://hub.docker.com/_/nginx).

Verify the real containers:

```sh
docker compose exec frontend nginx -t
curl http://127.0.0.1:8088/
curl http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"2^3+sqrt(16)"}'
curl http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"200+10%"}'
curl -i http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"sqrt(-1)"}'
docker compose down
```

Expect result 12, result 220, and an HTTP 400 `domain_error`, respectively. Initial startup can briefly return 502 before Go is ready; retry after both processes start. `docker compose down` stops/removes the two services.

**Current status: runtime-verified on 2026-11-02.** Docker Compose built and started both services; `nginx -t` passed. Actual HTTP and Chromium checks through `127.0.0.1:8088` verified production assets, original/advanced calculations, a present zero result, division-by-zero errors and recovery. Both services remained running without restarts or unexpected application errors. Test containers were removed with `docker compose down`, retaining the images. Measured evidence and limits are in [docs/VALIDATION.md](docs/VALIDATION.md).

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

Five normal browser tests verify real original/advanced calculations, negative-result power continuation, server errors and recovery, button/keyboard interaction, 1440/768/375 px layouts, focus separation, and an injected transport failure. A sixth test deliberately requires the API to be stopped. Keep Vite running, stop Go, then run:

```sh
CALCULATOR_BACKEND_DOWN=1 npm run test:integration -- --grep 'stopped backend'
```

PowerShell: `$env:CALCULATOR_BACKEND_DOWN='1'; npm run test:integration -- --grep 'stopped backend'`, then `Remove-Item Env:CALCULATOR_BACKEND_DOWN`. Restart Go afterward. The normal run explicitly skips this opt-in test. Test failure traces and viewport screenshots are generated in `frontend/test-results/` and ignored by Git.

## Validation and limitations

See [docs/VALIDATION.md](docs/VALIDATION.md) for executed checks and requirement evidence. Browser checks used Chromium on Windows; other engines, real mobile devices, and screen-reader behavior were not exhaustively tested. Semantics, labels, live announcements, keyboard activation, visible focus, and responsive sizing are implemented. The final maintenance update scopes the patched `minimatch@11.0.1` to its three inspected development-tool consumers and adds an HTTP regression requiring a present, non-null zero result with no error field. Full regression and Docker runtime checks passed. The source is published on GitHub; the running application has not been publicly deployed.

AI-assisted implementation and the exact original master prompt are disclosed in [PROMPTS.md](PROMPTS.md). GitHub repository publication was explicitly authorized by the owner; application hosting is a separate step.
