# Measured test coverage

Measured on **2026-11-02** with Go 1.27.3 and Vitest/coverage-v8 5.0.5. These are executed results, not coverage targets.

## Backend

Commands executed from `backend/`:

```sh
go test ./... '-coverprofile=coverage.out'
go tool cover '-func=coverage.out'
go tool cover '-html=coverage.out' -o coverage.html
```

| Scope | Statement coverage |
| --- | ---: |
| `internal/evaluator` | 100.0% |
| `internal/api` | 97.6% |
| `cmd/server` | 0.0% |
| **All backend packages** | **94.8%** |

The full raw Go profile is retained as [evidence/backend-coverage.out](evidence/backend-coverage.out). The actual function report is [evidence/backend-coverage.txt](evidence/backend-coverage.txt). Regenerate the local HTML report with the command above. To render the retained profile, from `backend/` use `go tool cover '-html=../docs/evidence/backend-coverage.out' -o coverage.html`.

Table-driven tests verify arithmetic, precedence, associativity, unary signs, decimals/scientific literals, whitespace, zero, invalid/truncated expressions, unsupported input, division by zero, intermediate overflow, literal overflow, and exact expression boundaries. HTTP tests cover status codes, JSON contracts, null/missing/wrong-type fields, unknown fields, malformed/empty/multiple JSON values, body limits (including trailing whitespace), methods, and unknown routes.

The server entry point has no unit instrumentation coverage. Its actual startup and binding were exercised by real HTTP and Chromium tests, which are not included in this profile. The API's defensive unexpected-evaluator-error branch is not forced: the evaluator currently returns only its declared error type. The report does not hide these gaps or substitute integration execution for unit coverage.

Additional fuzz run:

```sh
go test ./internal/evaluator '-fuzz=FuzzEvaluate' '-fuzztime=5s'
```

**PASS**, 71,764 executions, 64 new interesting inputs in that run. The fuzz property checks for panics and non-finite successful values; it is not a mathematical oracle. Table-driven expected results verify mathematical correctness separately. No crash corpus was generated.

## Frontend

Command executed from `frontend/`:

```sh
npm run test:coverage
```

**30 tests passed** across two files (calculator behavior and Fetch boundary).

| Scope | Statements | Branches | Functions | Lines |
| --- | ---: | ---: | ---: | ---: |
| `src/App.tsx` | 100% | 98.55% | 94.44% | 100% |
| `src/api.ts` | 100% | 100% | 100% | 100% |
| **All included source** | **100% (106/106)** | **98.82% (84/85)** | **95% (19/20)** | **100% (84/84)** |

The actual machine-readable summary is retained as [evidence/frontend-coverage.json](evidence/frontend-coverage.json), with source keys converted to repository-relative paths. Full local reports: `frontend/coverage/index.html` and `frontend/coverage/coverage-summary.json`.

Coverage explicitly includes `src/**/*.{ts,tsx}`, including unimported application files. It excludes test files/setup and `src/main.tsx` (the React mounting entry point). CSS and build/test configuration are outside JavaScript coverage. Mounting, CSS, and layout are exercised by real-browser checks, without claiming unit coverage for them.

Behavior tested includes number/operator buttons, direct editing, physical keyboard and focused-keypad Enter, clear/delete, empty/oversized input, submission, backend-returned results, errors/recovery, pending state, duplicate prevention, request cancellation, stale successes/failures, full-precision result continuation, timeout, and unmount cleanup. Fetch tests cover JSON requests, structured/unstructured HTTP errors, malformed/non-JSON/non-finite responses, network failure, and cancellation.

Remaining uncovered items are the initial no-op keyboard-ref callback and the stale timer's non-current-request branch. Stale response protection is tested directly. No superficial tests were added solely to execute unreachable initialization paths.

## Scope limits

Coverage proves execution, not correctness by itself. Browser integration is documented separately in [VALIDATION.md](VALIDATION.md). No universal 100% threshold is imposed. Snapshot profiles reflect this implementation; rerun coverage after changing application source. Generated HTML, caches, binaries, and dependency directories are ignored by Git.
