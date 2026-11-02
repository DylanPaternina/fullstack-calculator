# Validation record

Validation recorded on **2026-11-02** for application revision `2a8bae7`.

Environment: Windows amd64; Go 1.27.3; Node 24.22.1; npm 12.2.1; TypeScript 5.9.3; Vite 8.3.5; Vitest 5.0.5; ESLint 10.14.0; Playwright 1.64.1 / Chromium 156. Container validation used Docker Desktop 4.97.0, Linux Engine 29.9.2 on WSL 2 and Compose 5.6.1.

## Quality checks

| Command/check | Recorded result |
| --- | --- |
| `go fmt ./internal/api`; `gofmt -l .` | Formatting clean |
| `go vet ./...` | Pass |
| `go test ./...` and `go test ./... '-coverprofile=coverage.out'` | Pass: 191 table subcases and 11 fuzz seeds |
| `go tool cover '-func=coverage.out'` | 96.5% total statements; evaluator 100%, API 97.6%, server entry point 0% |
| `go tool cover '-html=coverage.out' -o coverage.html` | HTML report generated |
| `go build ./cmd/server` | Pass; compiled backend exercised in integration checks |
| `npm ci` in a clean checkout | Pass |
| `npm ls minimatch`, `npm explain minimatch`, `npm ls --all --json` | Valid dependency tree; minimatch resolves to 11.0.1 |
| `npm audit` and `npm audit --omit=dev` | Both exit 0; zero findings |
| `npm run lint` | Pass |
| `npm run typecheck` | Pass |
| `npm run test:coverage` | 34 tests pass; statements/lines 100%, branches 98.88%, functions 95.45% |
| `npm run build` | Pass; production assets generated |
| `node scripts/smoke.mjs` | Pass: 56 real HTTP requests through Go and Vite |
| `npm run test:integration` | Five normal Chromium tests pass |
| Opt-in stopped-backend integration test | Pass in a separate run with Go stopped and Vite active; expected HTTP 502 and visible error |

[Coverage details](COVERAGE.md) explain the instrumented files, exclusions and retained profiles. The regenerated Go production profile matches the versioned evidence byte for byte; frontend coverage values also match the retained summary.

## Backend regression checks

The suites cover basic arithmetic, precedence, powers, roots, contextual percentages, grouping, signed/scientific literals, invalid syntax, division by zero, domain errors, overflow and exact input/depth boundaries.

The test comparison helper rejects NaN and Infinity, requires exact zero when zero is expected and uses relative tolerance for finite nonzero values. Dedicated cases cover tiny and subnormal results.

`TestAPIZeroResult/zero_division` requires `0/5` to return HTTP 200 and valid JSON with a present, non-null result equal to zero. A `*float64` distinguishes zero from an absent/null field; `json.RawMessage` detects any error property, including `error: null`.

Expression length is limited to 1,024 UTF-8 bytes and request bodies to 8,192 bytes. Recursive groups, square roots and power right-hand sides share a depth limit of 64. Flat arithmetic and unary-sign sequences are iterative. The [mathematical reference](MATHEMATICS.md) defines the supported semantics and numeric limits.

## Frontend and full-stack integration

The local integration path is **React → Fetch → Vite proxy → Go → JSON → result display**, using Go at `127.0.0.1:8080` and Vite at `127.0.0.1:5173`.

Unit and browser checks cover buttons, physical keyboard input, expression editing, loading, cancellation, duplicate prevention, timeouts, stale responses and recovery after errors. Continuing from a negative result preserves the complete power base: `(-2)^2` returns 4, while `(-2)^0.5` returns a controlled domain error. Display rounding does not change the value used for continuation.

The smoke script checks 22 calculations and six controlled errors through each local endpoint. Browser transport-failure injection verifies the error state and successful recovery after removing the injected failure. A separate test with Go actually stopped verifies the Vite 502 response and an enabled retry control.

Responsive checks at 375, 768 and 1440 px found no horizontal overflow and verified that the input focus outline stays below its label. Keyboard navigation, labels and result/error announcements are exercised by the existing tests.

Thirty-three representative API requests were also checked across direct Go, Vite and Nginx: arithmetic, precedence, powers, roots, percentages, zero, division by zero, malformed expressions and domain errors. The README's PowerShell example returned 11. The [API reference](API.md) documents the request/response contract and expected errors.

## Dependency validation

The development dependency tree previously resolved minimatch 10.2.6 through ESLint 10.14.0, `@eslint/config-array` 0.23.5 and `@typescript-eslint/typescript-estree` 8.72.1. The [GHSA-g3hv-75wg-m249 advisory](https://github.com/advisories/GHSA-g3hv-75wg-m249) appeared as 11 propagated moderate package entries from one root issue; the production-only audit had no findings.

At validation time, the current parent versions still required minimatch 10 and no patched 10.x release was available. Three consumer-scoped overrides select minimatch **11.0.1**, retaining the primary tooling versions. The lockfile includes the required nonrecursive matcher dependency. Reassess the overrides when upstream packages adopt a patched version.

Compatibility checks cover the consumed constructor, `match(path, partial)`, `pattern` property, `minimatch()` function, CommonJS/ESM exports and supported Node engines. Of 14,196 comparisons, configured project patterns agreed; 18 differences occurred on two exploratory patterns containing `**/..`, outside the project configuration. Compatibility is established for the current tooling and patterns, rather than every possible glob expression.

ESLint retained the same 10 files and ignore behavior and detected a deliberate unused-variable probe. The TypeScript project-service default-project glob path also passed. Following `npm ci`, the full resolved tree was valid, and both development-inclusive and production-only audits returned zero findings. The frontend Docker builder also completed installation with zero findings. [Dependency evidence](evidence/dependency-audit.json) records the snapshot; advisory data may change over time.

## Docker runtime

The existing configuration defines two services: Nginx serving the production frontend at `127.0.0.1:8088`, and Go listening internally at `backend:8080`. Nginx preserves `/api/calculate` when forwarding requests; the backend has no published host port.

| Command/check | Recorded result |
| --- | --- |
| `docker compose config` | Pass |
| `docker compose up --build -d` | Both images built and services started; valid cached layers reused |
| `docker compose ps` | Both services running |
| `docker compose exec -T frontend nginx -t` | Syntax/configuration test passed |
| Frontend HTML and production JS/CSS | HTTP 200 |
| API through the frontend port | Eight cases passed with expected JSON and HTTP status |
| Chromium production UI | Eight real API cases, recovery, Enter and three viewport widths passed |
| Application logs/container state | No unexpected errors; zero restarts and empty error states |
| `docker compose down` | Test containers/network removed; images retained |

The eight API cases were `5+3*2 → 11`, `2^3 → 8`, `sqrt(81) → 9`, `2^3+sqrt(16) → 12`, `200+10% → 220`, `200*10% → 20`, `0/5 → 0`, and `5/0 → HTTP 400 division_by_zero`. Container checks used the published frontend port, exercising the reverse proxy and internal networking. No health checks are configured.

[Container evidence](evidence/docker-validation.json) retains the measured results. The configuration also passed YAML/schema and build-path checks. Base image tags follow their release lines; image digests are not pinned.

## Documentation and verification scope

Setup commands were exercised from a clean local checkout using the documented tools. Relative links and anchors were checked, the original API examples were preserved, and the development-prompt documents were compared with their source text.

Coverage applies to the scopes listed in [COVERAGE.md](COVERAGE.md). Real-process/browser checks cover startup and layout without adding them to unit coverage. Browser verification used Chromium on Windows; other engines, physical touch devices, screen readers, load behavior, the Go race detector and exhaustive image-security scanning are outside this validation record.
