# Full-stack calculator

A classic calculator for a Software Engineer Intern assessment. React builds and edits an expression; Go alone parses and evaluates it. The project prioritizes correctness, clear responsibilities and reproducible tests over additional infrastructure.

[Local setup](#local-setup) · [API](#api-usage) · [Design decisions](#design-decisions) · [Tests](#testing-and-coverage) · [Docker](#docker-deployment) · [AI disclosure](#10-ai-assisted-development)

## Features and technology

- Addition, subtraction, multiplication and division, with normal mathematical precedence.
- Powers, square roots, contextual percentages, parentheses, decimals, scientific notation and unary signs.
- Traditional keypad plus editable expression input, physical keyboard support, clear/delete, loading and recoverable errors.
- Responsive layout, associated labels, visible keyboard focus and accessible result/error announcements.

| Layer | Technology |
| --- | --- |
| Frontend | React 19, strict TypeScript 5.9, Fetch; Vite 8 |
| Backend | Go and standard-library HTTP/JSON; no third-party Go dependencies |
| Tests | Go testing/httptest, Vitest 5, React Testing Library, Playwright/Chromium |
| Optional deployment | Two Docker Compose services: Nginx frontend and Go backend |

[frontend/package-lock.json](frontend/package-lock.json) pins the installed dependency tree; [frontend/package.json](frontend/package.json) and [backend/go.mod](backend/go.mod) define tool requirements.

## Project structure

```text
backend/
  cmd/server/main.go                 HTTP startup, address and timeouts
  internal/api/                     JSON handler and HTTP tests
  internal/evaluator/               Arithmetic parser and mathematical tests
frontend/
  src/App.tsx                       Input, keyboard, request state and presentation
  src/api.ts                        Fetch and external response validation
  src/styles.css                    Responsive layout and focus states
  e2e/calculator.spec.ts             Browser integration checks
scripts/smoke.mjs                    Real HTTP checks through Go and Vite
compose.yaml                        Two-service production configuration
docs/                               API, mathematics, coverage and validation evidence
PROMPTS.md                          Actual development prompts and AI disclosure
```

## Requirements

Install **Go 1.26+**, **Node.js 24 LTS** and npm. The dependency constraints also permit Node 22.13+ in the 22.x line or Node 26+. Recorded validation used Go 1.27.3, Node 24.22.1 and npm 12.2.1. Docker with the Compose plugin is optional for local development.

No database, environment file or backend dependency installation is required. Browser tests need Playwright's Chromium download, shown below.

## Local setup

```sh
git clone https://github.com/DylanPaternina/fullstack-calculator.git
cd fullstack-calculator
```

**Backend:** from the repository root in terminal 1:

```sh
cd backend
go run ./cmd/server
```

**Frontend:** from the repository root in terminal 2:

```sh
cd frontend
npm ci
npm run dev
```

Open **http://127.0.0.1:5173**. Go listens at **127.0.0.1:8080**; Vite proxies `/api` there, keeping browser requests same-origin.

`ADDR` changes the Go listen address: `ADDR=127.0.0.1:8090 go run ./cmd/server` on POSIX, or `$env:ADDR='127.0.0.1:8090'; go run ./cmd/server` in PowerShell. If changing that port, also update the development proxy target in `frontend/vite.config.ts`.

### Using the calculator

Click keys or edit the expression field. Use `*`, `/`, `^` and lowercase `sqrt(...)`; the square-root key inserts `sqrt(`, which needs a closing `)`.

`Enter` or `=` calculates, `Escape` clears and `Backspace` removes one character. Space activates a focused button. Other editable regions, links, unrelated buttons and modified shortcuts retain native behavior.

After a result, a binary operator or `%` continues from the full backend value; numbers, `(` and the root key start a new expression. A negative power base is grouped, so continuing -2 with power 2 submits `(-2)^2`. Direct editing/deleting edits the original expression, including characters within `sqrt(`.

While a request is pending, equals is disabled. Editing/clearing cancels and invalidates it; request identity guards prevent stale responses from replacing current state. Requests time out after 10 seconds. After an error, edit and retry.

## API usage

`POST /api/calculate` accepts `{"expression":"5 + 3 * 2"}` and returns HTTP 200 with `{"result":11}`. Errors use `{"error":{"code":"division_by_zero","message":"Cannot divide by zero"}}`, with a corresponding non-success status.

Examples below use POSIX shell quoting:

```sh
curl -i http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"5 + 3 * 2"}'
curl -i http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"5 / 0"}'
curl http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"2^3+sqrt(16)"}'
curl http://127.0.0.1:8080/api/calculate -H 'Content-Type: application/json' -d '{"expression":"200+10%"}'
```

Expected: 200/result 11, 400/division_by_zero, 200/result 12 and 200/result 220. PowerShell equivalent:

```powershell
Invoke-RestMethod -Uri http://127.0.0.1:8080/api/calculate -Method Post -ContentType application/json -Body '{"expression":"5 + 3 * 2"}'
```

| Expression | HTTP | Result or error |
| --- | ---: | --- |
| `5+3` / `5+3*2` | 200 | 8 / 11 |
| `2^3` / `sqrt(81)` | 200 | 8 / 9 |
| `200+10%` / `0/5` | 200 | 220 / 0, with the result property present |
| `5/0` | 400 | `division_by_zero` |
| `5+*2` | 400 | `invalid_expression` |

See the [API reference](docs/API.md) for every error code, HTTP/body limits, JSON-decoding assumptions and container-port examples.

## Mathematical behavior

The engine supports conventional arithmetic precedence, right-associative powers and explicitly defined contextual percentages. The decisions below explain those rules and the float64/display policy. The [mathematical reference](docs/MATHEMATICS.md) preserves the full examples, supported literals, domain conventions, percentage composition and numeric limitations.

## Design Decisions

### Project Approach

My goal was to build a full-stack calculator that demonstrates engineering judgment through correctness, maintainability, testing and clear responsibilities. I prioritized reliable functionality, a usable interface, predictable errors and meaningful tests over unnecessary complexity.

After completing the core requirements, I extended the application with optional mathematical operations and Docker support. These were incremental additions to the existing architecture, with regression checks. The decisions here describe my technical direction; implementation and automated validation were AI-assisted.

### 1. Frontend — React and TypeScript

I selected React and TypeScript for an interactive interface with explicit data types and maintainable state. React built on my previous experience; TypeScript adds compile-time checks and clarifies the frontend/backend contract. Runtime response validation still protects the external JSON boundary.

I used Vite for a straightforward development workflow and kept React's built-in state management. My priorities were responsive controls, mouse/keyboard interaction, predictable loading/errors/interrupted requests and a small component structure without unnecessary dependencies.

### 2. Backend — Go

I selected Go because the assessment preferred it and its standard library suits a small REST service. Using `net/http` keeps dependencies and framework overhead low, with transport logic that can be understood and tested independently.

I had more previous experience with React than Go, so I deliberately kept the backend architecture simple and focused on idiomatic Go practices.

### 3. Classic Calculator Interface

I chose a traditional keypad and editable expression rather than a two-input operation form. Familiar numeric/operator controls, direct editing and keyboard shortcuts support complete expressions with clear result feedback. I prioritized responsive behavior and accessible focus over visual features unrelated to the calculation workflow.

### 4. Backend-Based Expression Evaluation

I chose to keep all mathematical evaluation in Go. React constructs/displays expressions and server results, giving arithmetic a single source of truth without duplicating rules in two languages.

Complete expressions such as `5 + 3 * 2 = 11` require correct precedence. A restricted recursive-descent parser fits this small grammar without a general-purpose expression engine or `eval`; input-length and recursion limits bound its work. Presentation, Fetch, HTTP validation and mathematics remain separate responsibilities.

### 5. Advanced Mathematical Operations

Once the basic calculator was complete, I added powers, roots, parentheses and percentages through the same parser. I defined the supported behavior explicitly:

| Expression | Result | Decision |
| --- | ---: | --- |
| `2^3^2` | 512 | Powers associate right to left |
| `-2^2` / `(-2)^2` | -4 / 4 | Power precedes the leading sign; grouping changes the base |
| `2^-3` / `sqrt(81)` | 0.125 / 9 | Signed exponents and square roots are supported |
| `200+10%` / `200*10%` | 220 / 20 | Direct additive percentages use the left base; multiplication uses a fraction |

I selected contextual percentages for familiar percentage calculations. Each operand carries a value and a direct-percentage marker; signs/groups preserve that context, while completed binary operations and roots clear it. Thus `200+10%+5%=231`, but `200+10%*2=200.2`. The grammar remains restricted, with controlled syntax, overflow, nesting and domain errors.

### 6. REST API Design

I chose one endpoint, `POST /api/calculate`, accepting an expression and returning a numeric result or structured error. It supports basic and advanced operations without separate endpoints per function. I prioritized consistent JSON, explicit validation, meaningful HTTP statuses and understandable messages. The evaluator remains independent of HTTP and can be tested directly.

### 7. Numerical Precision

I selected Go's IEEE-754 `float64` for a general-purpose calculator, keeping arbitrary-precision complexity outside this assessment's scope. It is not a financial-precision tool: `0.1+0.2` returns `0.30000000000000004` from Go.

I distinguished numerical precision from presentation. Results display up to 12 significant digits (`0.3` in that example), while the full backend value is retained for continuation and the result tooltip. Display rounding never feeds later arithmetic. Division by zero, overflow and non-finite results are rejected; ordinary underflow may round to zero.

### 8. Testing and Quality Assurance

I treated automated testing as a core requirement. The suites cover mathematics/precedence, invalid expressions and numeric edges, HTTP contracts, React/keyboard interaction, network failures, recovery, full-stack integration and regressions after advanced features.

Go testing/`httptest`, Vitest/React Testing Library, Playwright and HTTP smoke checks provide complementary evidence. Coverage measures execution but does not replace meaningful assertions. Independent technical reviews informed targeted corrections without a redesign, including precision comparisons, focus spacing, dependency maintenance and explicit zero-result serialization. Recorded results come from automated tool execution.

### 9. Docker and Deployment Architecture

I chose two Compose services: a production React build served by Nginx and a separate Go backend. Nginx forwards API requests through Docker's internal network, preserving the frontend/backend boundary and providing one way to start the whole application.

Only the frontend port is exposed to the host; browser requests use a same-origin proxy. The Go runtime is a static executable running without root privileges. This keeps deployment reproducible without unrelated infrastructure.

### 10. AI-Assisted Development

I used OpenAI Codex as an engineering assistant for implementation, testing, documentation and refinement. My involvement focused on requirements, behavior, technical direction, quality expectations, reviewing results and directing improvements from technical feedback.

I made specific decisions about the interface, backend expression evaluation, precedence, percentage semantics, optional operations and deployment. AI assistance accelerated implementation; automated tests and reviews validated the resulting work. This is not represented as exclusively manual coding or personally performed manual testing. The genuine development instructions are recorded in [PROMPTS.md](PROMPTS.md); private review reports/transcripts remain excluded.

### Final Engineering Principle

I evaluated major decisions against the goal of a correct, understandable, maintainable and testable application, keeping the architecture proportional to the assessment. Dependency maintenance remains targeted: the lockfile and three consumer-scoped minimatch overrides preserve current primary tooling, with compatibility checks and limits in [VALIDATION.md](docs/VALIDATION.md).

## Testing and coverage

From the repository root, run backend checks:

```sh
cd backend
go vet ./...
go test ./... '-coverprofile=coverage.out'
go tool cover '-func=coverage.out'
go tool cover '-html=coverage.out' -o coverage.html
go build ./cmd/server
```

Quote the complete coverage flags as shown, including in PowerShell. HTML output: `backend/coverage.html`.

From the repository root in another terminal, run frontend checks:

```sh
cd frontend
npm run lint
npm run typecheck
npm run test:coverage
npm run build
```

`npm test` runs the same unit suite without coverage. HTML output: `frontend/coverage/index.html`; production assets: `frontend/dist/`.

| Measured scope | Result |
| --- | --- |
| Go tests | 191 table subcases and 11 fuzz seeds pass |
| Go statement coverage | 96.5% total; evaluator 100%, HTTP 97.6% |
| React/Fetch unit tests | 34 pass |
| Included frontend coverage | Statements/lines 100%; branches 98.88%; functions 95.45% |

[Coverage details and raw evidence](docs/COVERAGE.md) explain exclusions: the Go entry point is 0% in the unit profile; React mounting, CSS and configuration are outside JS unit coverage. These areas have real-process/browser checks, without claiming unit instrumentation for them. [Validation results](docs/VALIDATION.md) distinguish executed checks from remaining limits.

With **Go and Vite both running**, from the repository root:

```sh
node scripts/smoke.mjs
cd frontend
npx playwright install chromium
npm run test:integration
```

Five normal browser tests cover arithmetic, advanced controls, recovery, keyboard, focused-input spacing, negative-result continuation and 375/768/1440 px layouts. A sixth requires Go to be stopped while Vite keeps running:

```sh
CALCULATOR_BACKEND_DOWN=1 npm run test:integration -- --grep 'stopped backend'
```

In PowerShell use `$env:CALCULATOR_BACKEND_DOWN='1'; npm run test:integration -- --grep 'stopped backend'`, then `Remove-Item Env:CALCULATOR_BACKEND_DOWN`. Restart Go for normal use. The normal suite intentionally skips this opt-in case. Traces/screenshots are ignored under `frontend/test-results/`.

## Docker deployment

With Docker Desktop/Engine running, execute from the repository root:

```sh
docker compose config
docker compose up --build -d
docker compose ps
docker compose exec -T frontend nginx -t
docker compose logs --tail=100
```

Open **http://127.0.0.1:8088**. The Node 24 build produces React assets served by Nginx; `/api/` forwards to `backend:8080` preserving the path. The Go 1.27 builder produces a static executable in an unprivileged scratch runtime. Only the frontend port is published, on loopback. Build contexts exclude local tools, credentials and generated files.

```sh
curl http://127.0.0.1:8088/
curl http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"2^3+sqrt(16)"}'
docker compose down
```

The calculation returns 12. Additional Docker-port examples are in the [API reference](docs/API.md). Initial startup may briefly return 502 before Go is ready. `down` removes this project's containers/network and retains images; it does not uninstall Docker.

**Runtime verified on 2026-11-02:** both services built/started, Nginx validation passed, and HTTP plus Chromium checked production assets, calculations, zero serialization and controlled errors through port 8088. No unexpected application errors or service restarts were observed. Test services were stopped afterward. See [container evidence](docs/evidence/docker-validation.json).

Base images follow the [Go](https://hub.docker.com/_/golang), [Node](https://hub.docker.com/_/node) release lines and [Nginx](https://hub.docker.com/_/nginx) stable Alpine tag; digests are not pinned. For a static preview, `npm run preview` does not replace the separately running Go API. A deployed build needs a same-origin `/api` reverse proxy; Vite's development proxy is not a production server.

## Repository and delivery

Repository: [DylanPaternina/fullstack-calculator](https://github.com/DylanPaternina/fullstack-calculator). Both layers, unit tests, coverage evidence and development prompts are included. Use `git rev-parse HEAD` to identify the checked-out version.

| Preserved milestone | Commit | Contents |
| --- | --- | --- |
| [assessment-basic](https://github.com/DylanPaternina/fullstack-calculator/tree/assessment-basic) | `7054b3b` | Original four-operation calculator, tests and documentation |
| [assessment-advanced](https://github.com/DylanPaternina/fullstack-calculator/tree/assessment-advanced) | `c55f2ef` | Targeted fixes, powers, roots, percentages and Compose |

The original history and tags are retained; main includes follow-up maintenance. Publishing source is separate from deploying the running application, which has not been publicly hosted. Local release updates require explicit publication authorization before push.

Browser validation used Chromium on Windows. Other engines, real touch devices, screen readers, load performance and exhaustive image security have not been comprehensively verified. The documented float64 and percentage policies remain intentional limitations, not claims of universal calculator semantics.
