# Validation record — final improvements

Executed on **2026-11-02**, Windows amd64. Go 1.27.3, Node 24.22.1, npm 12.2.1, TypeScript 5.9.3, Vite 8.3.5, Vitest 5.0.5, ESLint 10.14.0, Playwright 1.64.1 / Chromium 156. No application dependency versions changed.

## Baseline and focused fixes

Before extending the evaluator, existing Go tests passed. The existing 30 frontend tests passed, and the frontend build passed (the build completed after the isolated focus-spacing adjustment). No pre-existing test failure was observed. Historical coverage was not used as evidence for the update.

The mathematical comparator was fixed and its regression tests passed **before** changing the production evaluator. Tests explicitly reject NaN, positive/negative Infinity, an incorrect `1e-13` when zero is expected, loss of tiny values, and unequal subnormals. Exact zero is required for zero expectations; finite nonzero results use relative tolerance. Existing approximate decimal cases still pass.

The input now has 6 px of extra top spacing. Chromium checks the actual bounding boxes and computed outline: the outer top of the focus ring must remain below the label at 375, 768 and 1440 px. Screenshots were generated and visually inspected. The visible focus ring, contrast and existing visual direction remain intact.

## Executed quality gates

| Command/check | Actual outcome |
| --- | --- |
| `gofmt -w` on changed Go files; formatting inspection | Pass |
| `go vet ./...` | Pass |
| `go test ./...` and `go test ./... '-coverprofile=coverage.out'` | Pass; all original and advanced cases |
| `go tool cover '-func=coverage.out'` | 96.5% total statements; evaluator 100%, HTTP 97.6%, server entry point 0% |
| `go tool cover '-html=coverage.out' -o coverage.html` | HTML report generated |
| `go test ./internal/evaluator '-fuzz=FuzzEvaluate' '-fuzztime=5s'` | Pass; 55,684 executions, 19 new interesting inputs, no crash |
| `go build -o ../.tools/calculator-server.exe ./cmd/server` | Pass; this binary used in actual integration |
| `npm run lint` | Pass |
| `npm run typecheck` / build's `tsc --noEmit` | Pass |
| `npm run test:coverage` | 34 tests pass; 100% lines, 98.88% branches |
| `npm run build` | Pass; generated production assets |
| `npm run test:integration` | 5 normal Chromium tests pass; stopped-backend test intentionally skipped in this run |
| `node scripts/smoke.mjs` | 56 real HTTP requests pass: 22 calculations and 6 controlled errors at each of Go and Vite |

The initial browser run exposed an outdated Tab-order expectation after inserting the new advanced row. The test now walks all five new controls before Clear, retaining assertions for every focused control. The focus-separation checks themselves passed. YAML schema validation initially used the wrong JSON Schema draft; the validator was changed to Ajv's draft-2020-12 implementation to match the official schema, without disabling schema validation.

A fresh independent static engineering review identified a continuation defect: a negative returned result used as the next power base lost its grouping. A new unit regression first failed with actual `-2^2` versus expected `(-2)^2`. The syntax construction was corrected, and the full unit suite and real-browser regressions passed: continuing -2 with exponent 2 returns 4, while exponent 0.5 produces an HTTP 400 `domain_error`. This change does not evaluate powers in React. Private evaluation reports and transcripts are excluded; this record describes the technical change and actual validation only.

## Local full-stack integration

Compiled Go listens at `127.0.0.1:8080`; Vite at `127.0.0.1:5173`. The real Chromium flow crosses **React → Fetch → Vite proxy → Go → JSON → result display**, with no mocked calculation responses in the arithmetic tests.

Verified original arithmetic, precedence, signed values and decimals, plus `2^3`, `2^3^2`, `-2^2`, `(-2)^2`, `2^-3`, `sqrt(81)`, `sqrt(2)`, `2^3+sqrt(16)`, contextual `200±10%`, `200*10%`, `200/10%`, chained percentages, multiplication containing a percentage, roots with percentages and grouped arithmetic. Invalid/root-domain/power-domain/percentage-zero expressions return HTTP 400 and visible errors. Buttons and physical keyboard construct the supported syntax; successful requests after errors still work.

Existing unit tests retain deferred-response checks proving stale successes/failures cannot replace current state, duplicate submissions are prevented, and editing/clear/unmount cancel requests. The real-browser fault-injection test aborts transport, checks the visible failure, removes the injection and successfully recalculates through Go. The opt-in stopped-backend test was also executed separately with Go actually stopped and Vite running: **1 test passed**, observing HTTP 502, a visible error and an enabled retry button. Go was restarted for the final power-continuation checks.

## Numeric semantics and parser safety

The complete acceptance examples are represented in Go tests, including all powers/root/percentage examples from the implementation instructions. Percentages carry an evaluated number plus a direct-percentage marker; grouping and unary signs preserve it, while completed binary operations and square roots clear it. The README defines signs, decimals, repeated percentages, grouping, powers and roots in this model.

Length remains 1,024 UTF-8 bytes; body limit remains 8,192 bytes. Recursive groups, roots and power right-hand sides share a limit of 64; unary-sign and flat arithmetic repetitions remain iterative. Tests check exact limits and combined over-limit cases. NaN, Infinity, overflow, invalid syntax and domain errors remain controlled. Fuzz checks finiteness/panics, separately from the expected-result mathematical tests.

## Docker — implemented but not runtime-verified

The machine has no Docker executable. Docker was **not installed**. No image build, container startup, Compose execution, containerized integration, or `nginx -t` is claimed.

Implemented files: `compose.yaml`, backend/frontend multi-stage Dockerfiles, both context-specific `.dockerignore` allowlists, and `frontend/nginx.conf`. Exactly two services are declared. Frontend serves static assets via Nginx at host `127.0.0.1:8088`; `/api/` is proxied to `backend:8080` with the full API path preserved. Backend binds `0.0.0.0:8080` internally and has no published host port. Final backend runtime is an unprivileged static binary; final frontend runtime contains Nginx and built assets rather than Node development tools.

Actual static checks executed using ignored project-local validation tools:

- YAML parsed with duplicate-key detection.
- Document validated by Ajv draft-2020-12 against the [official Compose schema](https://github.com/compose-spec/compose-spec/blob/main/schema/compose-spec.json).
- Both service names/build contexts, COPY input existence, frontend host mapping, internal Go bind address and absence of backend host ports asserted.
- Nginx proxy destination/path-preservation and browser-relative Fetch path checked against each other.
- Referenced production `frontend/dist/index.html` exists after the normal local build.
- Base image release lines checked against official Go, Node and Nginx image listings.
- Go cross-compiled successfully for Linux amd64 with `CGO_ENABLED=0`; the ELF output was checked. This is a local cross-build, not a Docker build or runtime test.

These checks cannot replace Docker's builder, Nginx's own parser or runtime networking. Follow the README's `docker compose config`, `up --build -d`, `ps`, `logs`, `exec frontend nginx -t`, curl and `down` commands when Docker is available. Image tags are release-line/stable tags, not immutable digest pins.

## Coverage and verification limits

[COVERAGE.md](COVERAGE.md) contains the regenerated figures, raw evidence locations and meaningful coverage gaps. No application source was newly excluded to raise coverage. Go bootstrap and its defensive internal-error fallback lack unit coverage; React mounting/CSS are outside the JS unit profile and are checked through the browser where available.

Viewports 375/768/1440 were checked for horizontal overflow and focused input geometry; the updated full-page screenshots were inspected. Other browser engines, real touch devices and screen readers were not exhaustively tested. No race-detector claim is made. Docker runtime verification remains **pending**.

The private historical audit/report and evaluation transcripts are deliberately excluded. The README documents technical fixes, and PROMPTS.md retains actual implementation instructions without including private audit materials. These implementation checks preceded publication; the owner subsequently authorized the GitHub repository, which preserves the basic and advanced implementation milestones. Source publication provides no additional Docker runtime verification.
