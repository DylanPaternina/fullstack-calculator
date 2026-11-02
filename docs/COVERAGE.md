# Measured test coverage

Measured on **2026-11-02** with Go 1.27.3 and Vitest/coverage-v8 5.0.5.

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
| **All backend packages** | **96.5%** |

The full raw Go profile is retained as [evidence/backend-coverage.out](evidence/backend-coverage.out). The actual function report is [evidence/backend-coverage.txt](evidence/backend-coverage.txt). Regenerate the local HTML report with the command above. To render the retained profile, from `backend/` use `go tool cover '-html=../docs/evidence/backend-coverage.out' -o coverage.html`.

Table-driven tests cover arithmetic, precedence, associativity, unary signs, decimal/scientific literals, whitespace, zero, invalid syntax, unsupported input, division by zero, overflow and expression limits. They also cover powers, roots, grouping, real-domain errors, contextual percentages and combined recursion limits. HTTP tests exercise calculations and request, method, body and error validation.

The backend suite contains **191 table subcases and 11 fuzz seeds**. `TestAPIZeroResult/zero_division` requires `0/5` to return HTTP 200, valid JSON and a present/non-null `*float64` equal to zero. The error field must be absent; `error: null` is rejected. The regenerated production profile matches the retained evidence byte for byte.

The test-only comparison helper has its own regression table. It rejects NaN and either Infinity, requires exact expected zero, uses relative tolerance for finite nonzero results, and distinguishes zero, tiny values and adjacent subnormal values. Production coverage percentages do not include `_test.go` source; the helper's behavior is checked by those tests rather than inferred from application coverage.

The server entry point has no unit instrumentation coverage; startup and binding are exercised separately through HTTP and Chromium. The API's unexpected-evaluator-error fallback is not covered because the current evaluator returns only its declared error type.

Separately recorded fuzz run:

```sh
go test ./internal/evaluator '-fuzz=FuzzEvaluate' '-fuzztime=5s'
```

**PASS:** 55,684 executions, 19 new interesting inputs and no crash corpus. Seeds include powers, contextual percentages, nested grouping and root-domain failures. Fuzzing checks for panics and non-finite successful values; table-driven expectations check mathematical correctness. Regular test runs execute the 11 seeds without repeating the timed fuzz campaign.

## Frontend

Command executed from `frontend/`:

```sh
npm run test:coverage
```

**34 tests passed** across two files (calculator behavior and Fetch boundary).

| Scope | Statements | Branches | Functions | Lines |
| --- | ---: | ---: | ---: | ---: |
| `src/App.tsx` | 100% | 98.64% | 95% | 100% |
| `src/api.ts` | 100% | 100% | 100% | 100% |
| **All included source** | **100% (113/113)** | **98.88% (89/90)** | **95.45% (21/22)** | **100% (89/89)** |

The actual machine-readable summary is retained as [evidence/frontend-coverage.json](evidence/frontend-coverage.json), with source keys converted to repository-relative paths. Full local reports: `frontend/coverage/index.html` and `frontend/coverage/coverage-summary.json`.

Coverage explicitly includes `src/**/*.{ts,tsx}`, including unimported application files. It excludes test files/setup and `src/main.tsx` (the React mounting entry point). CSS and build/test configuration are outside JavaScript coverage. Mounting, CSS, and layout are exercised by real-browser checks, without claiming unit coverage for them.

Behavior tested includes number/operator/function/grouping buttons, advanced keyboard syntax, direct editing, focused-keypad Enter, clear/delete, empty/oversized input, submission, backend-returned results, errors/recovery, pending state, duplicate prevention, request cancellation, stale successes/failures, full-precision result continuation (including power and percent), explicit grouping of negative power bases, timeout, and unmount cleanup. Fetch tests cover JSON requests, structured/unstructured HTTP errors, malformed/non-JSON/non-finite responses, network failure, and cancellation.

Remaining uncovered items are the initial no-op keyboard-ref callback and the stale timer's non-current-request branch. Stale response protection is tested directly.

## Scope limits

Coverage measures execution and should be interpreted alongside the assertions. Browser integration is documented in [VALIDATION.md](VALIDATION.md). Regenerate coverage after changing application source. Generated HTML, caches, binaries and dependency directories are ignored by Git.
