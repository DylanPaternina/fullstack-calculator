# FULL-STACK CALCULATOR — FINAL MAINTENANCE AND RELEASE VALIDATION

## ROLE

Act as a Senior Full-Stack Engineer, Security-Conscious Dependency Maintainer, and Technical Code Reviewer.

The Full-Stack Calculator is already implemented, tested, documented, and published.

The latest independent review awarded an indicative score of 98/100 and found no confirmed functional bugs.

Your task is to resolve the two remaining review observations without introducing regressions or unnecessary architectural changes.

**This is a targeted maintenance task, not a new development phase.**

## OBJECTIVES

### PRV-01 — Resolve the minimatch security advisory

The independent review identified:

- Package: `minimatch`
- Installed version during review: `10.2.6`
- Advisory: `GHSA-g3hv-75wg-m249`
- First patched version: `11.0.1`
- Dependency type: Transitive development dependency
- Affected tooling: ESLint-related dependency chain
- Risk assessed for this project: Low
- Production dependency audit during review: Zero findings

The review found 11 propagated moderate npm audit entries associated with one underlying advisory.

Do not interpret these as 11 independent root vulnerabilities.

### Required investigation

1. Inspect the current `frontend/package.json` and `package-lock.json`.
2. Determine the currently installed minimatch version.
3. Execute `npm ls minimatch`.
4. Execute `npm explain minimatch`.
5. Run the current npm security audit.
6. Check whether the vulnerability is already resolved.
7. Identify which parent dependencies require the affected version.
8. Investigate compatible patched dependency versions.

### Remediation strategy

Prefer, in order:

1. A compatible update to the parent development dependencies.
2. A targeted transitive dependency update that preserves compatibility.
3. A narrowly scoped npm override, only if compatibility is verified.

Do not automatically force `minimatch@11` into a dependency tree that requires incompatible APIs.

Do not run `npm audit fix --force`.

Do not downgrade ESLint or other primary tooling merely to make the audit report disappear.

Do not remove lint or security tooling.

If no safe fix is available, preserve the working dependency tree and clearly document the remaining risk instead of introducing a breaking change.

### Required verification

After applying a compatible change:

```bash
npm ci
npm ls minimatch
npm run lint
npm run typecheck
npm run test:coverage
npm run build
npm audit
npm audit --omit=dev
```

Run these commands in the frontend directory.

Inspect the complete resolved dependency tree.

If the audit still reports vulnerabilities, distinguish the original advisory from any unrelated or newly discovered findings.

Do not declare PRV-01 resolved unless the relevant affected dependency is genuinely replaced with a safe, compatible version and validation succeeds.

---

## PRV-02 — Add explicit HTTP regression coverage for zero results

The independent review identified an informational improvement, not a confirmed bug.

Affected files:

- `backend/internal/api/api_test.go`
- Potentially `backend/internal/api/advanced_test.go`

The existing API correctly returns:

```json
{
  "result": 0
}
```

However, the committed HTTP tests do not contain one dedicated case verifying all aspects of zero-result serialization.

### Required change

Add a focused regression test using:

```json
{
  "expression": "0/5"
}
```

Verify:

1. HTTP status is exactly `200`.
2. Response is valid JSON.
3. The `result` field exists.
4. The `result` value is exactly `0`.
5. The response does not include an error field.
6. The result is not confused with a missing or null value.

Use an appropriate typed response structure with `*float64`, or an equivalent explicit JSON field-presence check.

Follow the project's existing table-driven testing conventions.

Do not introduce a new testing framework.

Do not modify the production handler unless the new regression test reveals a genuine defect.

Do not refactor unrelated tests.

### Required validation

From the backend directory, execute:

```bash
go test ./...
go vet ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Verify that existing tests remain successful.

The new test should prevent future accidental omission of the `result` field when its value is zero.

---

## 3. FULL REGRESSION VERIFICATION

After completing PRV-01 and PRV-02, verify the unchanged calculator behavior.

### Backend

- Basic arithmetic.
- Operator precedence.
- Exponentiation.
- Square roots.
- Contextual percentages.
- Invalid expressions.
- Division by zero.
- Numeric domain errors.
- Numeric overflow.
- HTTP response consistency.

### Frontend

- Calculator buttons.
- Physical keyboard.
- Expression editing.
- Loading state.
- Error recovery.
- Backend response handling.
- Responsive behavior.
- Existing accessibility behavior.

### Integration

Run the existing smoke and integration tests.

Use the actual project scripts where available.

Do not modify production code to satisfy an unrelated test unless an actual defect is demonstrated.

---

## 4. DOCKER VERIFICATION

Docker Desktop is now installed and working on the local Windows computer.

The repository already includes a two-service Docker Compose configuration.

Do not recreate or redesign the Docker infrastructure.

Check:

```bash
docker compose config
docker compose up --build -d
docker compose ps
```

Verify the actual frontend service name, then execute the appropriate Nginx configuration test.

Test the application through the existing mapped frontend port.

At minimum, verify:

| Expression | Expected result |
|---|---|
| `5+3*2` | 11 |
| `2^3+sqrt(16)` | 12 |
| `200+10%` | 220 |
| `5/0` | Controlled HTTP error |
| `0/5` | 0 |

Check that the frontend and backend communicate correctly through the reverse proxy.

Inspect service logs for unexpected errors.

After verification, stop only the calculator's services using:

```bash
docker compose down
```

Do not remove unrelated containers, volumes, or images.

If Docker fails, distinguish environmental problems from application defects.

Do not alter Docker configuration unless a reproducible defect demonstrates that a change is required.

---

## 5. CHANGE RESTRICTIONS

This is a minimal maintenance update.

Allowed changes:

- Frontend dependency manifest and lockfile, when justified.
- Focused backend API regression tests.
- Necessary updates to existing validation documentation.
- Accurate recording of this AI-assisted maintenance prompt.

Do not:

- Rewrite the Go expression evaluator.
- Redesign the React interface.
- Add calculator features.
- Introduce new frameworks.
- Modify mathematical semantics.
- Alter the API contract.
- Change Docker architecture unnecessarily.
- Remove meaningful tests.
- Inflate coverage artificially.
- Commit private audit reports.
- Modify unrelated files.

Review the existing Git status before editing.

Preserve unrelated uncommitted changes.

---

## 6. DOCUMENTATION AND AI TRANSPARENCY

Update documentation only where necessary to reflect actual maintenance results.

If the dependency tree changes, record the relevant maintenance decision concisely.

If coverage values change, update the corresponding report with real results.

Because the employer requires disclosure of prompts used during development, record this maintenance instruction accurately in the existing prompt history.

Do not fabricate earlier prompts.

Do not include private review transcripts or the independent audit report in the public repository.

---

## 7. GIT AND RELEASE SAFETY

The project has already been published.

You may prepare a focused local commit after successful validation.

Do not amend or rewrite earlier commits.

Do not change the repository remote.

Do not push to GitHub without explicit authorization.

Before finishing:

1. Review the complete diff.
2. Confirm only intended files changed.
3. Verify Git formatting and status.
4. Confirm no credentials or private audit materials were added.
5. Verify that documentation matches the implementation.

---

## 8. FINAL TECHNICAL REPORT

Provide a concise report in Spanish.

Include:

### PRV-01

- Original dependency path.
- Previous vulnerable version.
- Updated version, if changed.
- Selected remediation approach.
- Reasons for the chosen approach.
- Final npm audit findings.
- Compatibility verification.

### PRV-02

- New regression test location.
- Assertions implemented.
- Test execution result.
- Whether production code required modification.

### General validation

- Backend tests.
- Frontend tests.
- TypeScript.
- Lint.
- Build.
- Coverage.
- Full-stack integration.
- Docker Compose verification.

### Final status

Classify each observation as:

- Resolved and verified.
- Mitigated with documented residual risk.
- Unresolved.
- Not verified.

State whether the application remains ready for technical evaluation.

Do not assign an artificial new review score.

Do not claim that the project is fully verified if critical checks were not executed.

## FINAL DIRECTIVE

Resolve PRV-01 safely and add the focused PRV-02 regression test.

Preserve the existing calculator's behavior, architecture, and code quality.

Verify the entire application after the changes, including the functioning Docker environment.

Avoid unrelated improvements.

When the maintenance work is complete, stop before pushing and request explicit publication authorization.