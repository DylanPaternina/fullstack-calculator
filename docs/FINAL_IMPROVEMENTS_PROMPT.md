# FULL-STACK CALCULATOR — FINAL IMPROVEMENTS
## Codex Implementation Prompt | React + TypeScript + Go

### ROLE AND MISSION

Act as a Senior Full-Stack Software Engineer specializing in Go, React, TypeScript, mathematical expression parsing, automated testing, and maintainable software architecture.

You are working on an **existing, functioning technical assessment project** for a Software Engineer Intern position.

The application already implements:

- A classic React + TypeScript calculator.
- A Go backend with a mathematical expression evaluator.
- A REST API.
- Basic arithmetic operations and operator precedence.
- Automated tests and coverage reporting.
- Project documentation.

The application previously received a favorable technical review, but that evaluation is historical and must not be treated as evidence that the updated application is correct.

**Do not rebuild the application from scratch.**

Your mission is to correct two identified defects and implement all four previously optional deliverables.

### Six mandatory objectives

1. **R01:** Correct overly permissive mathematical test assertions.
2. **R02:** Fix the focus-outline overlap in the calculator interface.
3. **Exponentiation:** Implement the `^` operator.
4. **Square root:** Implement `sqrt(expression)`.
5. **Percentage:** Implement conventional contextual percentage behavior.
6. **Docker:** Provide a complete two-service deployment using Docker Compose.

All six objectives are required for this update.

### Engineering priorities

1. Mathematical correctness.
2. Preservation of existing functionality.
3. Clean, idiomatic, maintainable code.
4. Meaningful automated testing.
5. Reliable frontend/backend integration.
6. Complete implementation of the six objectives.
7. Documentation and reproducibility.
8. Visual polish and nonessential optimizations.

Do not introduce architectural complexity unless it is necessary to implement these requirements correctly.

---

## 1. AUTONOMY AND WORKSPACE PROTECTION

You may autonomously:

- Inspect the existing repository.
- Modify project files.
- Extend the parser and frontend.
- Add project-local dependencies when justified.
- Run tests, builds, formatting, and static analysis.
- Add and update documentation.
- Prepare local Git commits.

You must not:

- Recreate or replace the entire project.
- Modify unrelated projects.
- Discard existing uncommitted user changes.
- Install global tools without authorization.
- Install Docker without authorization.
- Delete unrelated files.
- Change repository remotes or visibility.
- Push or publish the project without explicit approval.
- Claim verification that has not actually occurred.

Make reasonable engineering decisions autonomously.

Only request clarification when a genuine blocker cannot be resolved from the code, the requirements, or a safe documented assumption.

---

## 2. INITIAL REPOSITORY INSPECTION

Before modifying anything:

1. Inspect the repository structure.
2. Check Git status and existing changes.
3. Read the current Go expression evaluator and its tests.
4. Inspect the HTTP API.
5. Inspect React components, keyboard handlers, and calculator state.
6. Inspect the existing CSS and accessibility behavior.
7. Review the frontend and backend test configuration.
8. Read README.md, PROMPTS.md, and coverage documentation.
9. Identify available local development and testing tools.

### Audit findings

A previous technical audit identified:

**R01 — Mathematical test assertions**

Expected location:

`backend/internal/evaluator/evaluator_test.go`

The current comparison logic may incorrectly accept NaN or small incorrect values, especially when exact zero is expected.

**R02 — Focus-outline overlap**

Expected location:

`frontend/src/styles.css`

The input's focus outline visually overlaps the `EXPRESSION` label.

The original `TECHNICAL_REVIEW.md` is stored outside the current project.

Do not assume that this report is accessible.

If it can be located through an explicitly provided path, inspect it without copying private audit materials into the repository.

Do not search unrelated personal directories.

If it is unavailable, use the defect descriptions above, independently verify the relevant code, and clearly disclose that the original audit report was not inspected.

Do not block unrelated implementation work merely because the report is unavailable.

### Preserve the existing API

The current contract is equivalent to:

`POST /api/calculate`

Request:

```json
{
  "expression": "5+3*2"
}
```

Response:

```json
{
  "result": 11
}
```

Preserve this contract.

Do not replace it with separate `a`, `b`, and `operation` fields.

Preserve the existing error-response structure unless a small, justified extension is necessary.

### Initial baseline

Before implementing new functionality:

- Run existing relevant tests.
- Record current test and build results.
- Identify any pre-existing failures.
- Inspect existing coverage data.
- Record the current state without assuming historical coverage is still accurate.

Give a short implementation plan, then proceed autonomously.

---

## 3. CORRECTION R01 — MATHEMATICAL TEST ASSERTIONS

Inspect the existing mathematical comparison helper and its usage.

Correct assertions that can falsely accept incorrect results.

### Required behavior

1. Explicitly reject unexpected NaN.
2. Explicitly reject unexpected positive or negative Infinity.
3. Require exact zero where the operation is mathematically and representationally expected to produce exact zero.
4. Use appropriate floating-point tolerances when approximate comparisons are justified.
5. Prevent absolute tolerances from automatically accepting incorrect tiny nonzero results.
6. Check subnormal and very small finite values.
7. Preserve appropriate comparisons for ordinary decimal arithmetic.

Do not apply exact comparisons indiscriminately to calculations subject to floating-point rounding.

### Verification

Create regression tests demonstrating that the comparison helper rejects the previously problematic values.

The test suite itself must pass when the helper correctly identifies those differences.

Do not merely modify test thresholds until the existing tests pass.

Do not change the production mathematical engine unless an independently verified defect requires it.

---

## 4. CORRECTION R02 — ACCESSIBLE FOCUS INDICATOR

Inspect the current focus styling and the surrounding layout.

Fix the overlap between the focus outline and the `EXPRESSION` label using the smallest justified change.

Possible solutions include:

- Adjusting the spacing between label and input.
- Modifying outline offset.
- Using an appropriate inset focus indicator.

Preserve:

- Visible keyboard focus.
- Readable contrast.
- Existing visual identity.
- Responsive layout.
- Existing interaction behavior.

Do not remove or hide the focus indicator.

Do not redesign the calculator.

### Verification

Inspect representative viewport widths:

- 375 px.
- 768 px.
- 1440 px.

Use real browser or screenshot verification if available.

If visual verification cannot be performed, disclose the limitation rather than claiming the issue is visually resolved.

---

## 5. EXTEND THE EXISTING MATHEMATICAL PARSER

The existing Go evaluator is the single source of truth for mathematical calculations.

Extend its current architecture rather than introducing a second evaluator.

The updated evaluator must support:

- Existing basic arithmetic.
- Standard operator precedence.
- Unary signs.
- Parentheses.
- Right-associative exponentiation.
- Square-root expressions.
- Contextual percentage operations.

### Parser requirements

Preserve a small and understandable grammar.

Use the existing parser structure where practical.

A lightweight syntax tree or equivalent representation is acceptable if required to distinguish contextual percentage operations from ordinary numeric values.

Do not introduce an unnecessarily generic expression framework.

Never evaluate expressions as executable source code.

Avoid third-party parser libraries unless an actual limitation makes them necessary.

### Required precedence behavior

| Expression | Expected result |
|---|---|
| `2+3*4` | 14 |
| `(2+3)*4` | 20 |
| `2^3^2` | 512 |
| `-2^2` | -4 |
| `(-2)^2` | 4 |
| `2^-3` | 0.125 |

Exponentiation must associate from right to left.

Exponentiation must bind more strongly than unary negation on its left, while negative exponents must remain supported.

Postfix percentage must have a well-defined precedence that supports expressions such as `2^50%` and `200*10%`.

Do not rely on accidental parser behavior for these cases.

### Parser safety

Enforce reasonable restrictions on:

- Expression length.
- Nesting depth.
- Invalid tokens.
- Incomplete expressions.
- Unexpected input.
- Numeric overflow.
- Non-finite results.

Avoid unbounded recursive nesting and excessive computational work.

Preserve existing input limits unless an extension requires a justified adjustment.

---

## 6. EXPONENTIATION

Implement:

`^`

Use Go's `math.Pow` or an equivalent standard-library operation.

### Required examples

| Expression | Expected result |
|---|---|
| `2^3` | 8 |
| `5^2` | 25 |
| `10^0` | 1 |
| `2^-3` | 0.125 |
| `2^3^2` | 512 |
| `-2^2` | -4 |
| `(-2)^2` | 4 |
| `2+3^2` | 11 |
| `2*3^2` | 18 |

### Edge cases

Explicitly handle:

- Missing exponents.
- Incomplete expressions.
- Negative exponents.
- Zero raised to a negative exponent.
- Negative bases raised to non-integer exponents outside the real-number domain.
- Numerical overflow.
- NaN and Infinity.
- Nested powers.

Define `0^0 = 1`, consistently with Go's `math.Pow`, and document the convention.

A result that cannot be represented as a finite supported real number must produce a controlled error.

### Frontend

Add a clearly recognizable exponentiation button, displaying `^` or `xʸ`.

The button must insert valid expression syntax and work with the existing calculator state.

Support typing `^` from the physical keyboard.

Do not calculate powers in React.

---

## 7. SQUARE ROOT

Implement:

`sqrt(expression)`

Use Go's `math.Sqrt`.

### Required examples

| Expression | Expected result |
|---|---|
| `sqrt(81)` | 9 |
| `sqrt(16)` | 4 |
| `sqrt(0)` | 0 |
| `sqrt(2)` | Approximately 1.41421356237 |
| `sqrt(9)+2` | 5 |
| `sqrt(25)*2` | 10 |

### Error handling

Reject:

- `sqrt(-1)` and other negative arguments.
- Missing arguments.
- Empty parentheses.
- Malformed expressions.
- Unclosed parentheses.
- Non-finite results.

Support square roots combined with ordinary arithmetic and exponentiation.

### Frontend

Add a `√` button.

The visible button may insert the supported `sqrt(` syntax into the expression display.

Users must be able to complete the expression naturally using the existing parentheses controls.

Preserve correct Backspace and Clear behavior.

Do not perform square-root calculations in React.

---

## 8. CONTEXTUAL PERCENTAGE — REQUIRED SEMANTICS

Implement `%` using a documented **conventional contextual calculator interpretation**.

This is a deliberate requirement.

Do not silently replace it with a purely unary percentage implementation.

### Required behavior

For addition and subtraction, a directly applied percentage of the left-hand accumulated value must be calculated relative to that value.

For multiplication, division, and standalone percentages, `%` behaves as a division by 100.

### Mandatory examples

| Expression | Expected result |
|---|---|
| `50%` | 0.5 |
| `25%` | 0.25 |
| `200+10%` | 220 |
| `200-10%` | 180 |
| `200*10%` | 20 |
| `200/10%` | 2000 |
| `200+10%+5%` | 231 |
| `25%+1` | 1.25 |
| `sqrt(25)*20%` | 1 |

### Semantic rules

For direct percentage operands:

- `A + B%` means `A + (A × B / 100)`.
- `A - B%` means `A - (A × B / 100)`.
- `A × B%` means `A × (B / 100)`.
- `A ÷ B%` means `A ÷ (B / 100)`.
- Standalone `B%` means `B / 100`.

For chained addition/subtraction, use the already evaluated left-hand result as the percentage base.

For example:

`200+10%+5%`

First calculate 220, then add 5% of 220, producing 231.

The contextual rule applies when the right-hand operand of addition/subtraction is itself a direct percentage expression, not automatically to an arbitrary expression containing a percentage somewhere inside it.

Respect ordinary operator precedence.

For example:

`200+10%*2 = 200.2`

Here the right-hand operand of addition is a multiplication expression, so `10%` contributes its numeric value of `0.1` before multiplication.

Document this distinction clearly.

The implementation must distinguish percentage syntax from its evaluated numeric value where necessary.

Do not use fragile text replacement as a substitute for correct expression semantics.

### Additional edge cases

Define and test behavior for:

- Negative percentages.
- Decimal percentages.
- Repeated percentages.
- Percentages in parentheses.
- Percentages combined with powers.
- Percentages combined with square roots.
- Invalid percentage placement.
- Division by zero resulting from percentage expressions.

If a rarely used combination is ambiguous, either define it explicitly with tests or reject it consistently.

Do not silently produce mathematically inconsistent behavior.

### Frontend

Add a `%` button.

Support the `%` physical keyboard key.

Preserve existing editing, clearing, result display, and error recovery.

All contextual percentage calculations must occur in Go.

---

## 9. API INTEGRATION AND ERROR HANDLING

Preserve:

`POST /api/calculate`

The request continues to contain one expression string.

Example:

```json
{
  "expression": "2^3+sqrt(16)"
}
```

Expected response:

```json
{
  "result": 12
}
```

Preserve existing structured error responses.

Reuse or extend error codes consistently for:

- Invalid expression.
- Invalid request.
- Division by zero.
- Mathematical domain errors.
- Numeric overflow.
- Unsupported syntax.
- Expression complexity limits.

Expected mathematical errors must not become HTTP 500 responses.

Do not introduce additional calculator endpoints unnecessarily.

Avoid breaking existing API consumers.

---

## 10. FRONTEND INTEGRATION

Preserve the existing React interface and its visual design.

Add controls for:

- Exponentiation.
- Square root.
- Percentage.
- Opening and closing parentheses, if not already present.

### Interaction requirements

All new controls must:

- Work with mouse and touch.
- Support appropriate keyboard input.
- Insert valid expression syntax.
- Integrate with the existing state-management logic.
- Respect the current display behavior.
- Preserve Backspace and Clear.
- Preserve behavior after successful calculations.
- Preserve error recovery.
- Avoid duplicate submissions.
- Avoid stale HTTP responses overwriting newer calculator state.

Maintain the existing responsive and accessible interface.

Reorganize the button layout only as much as necessary.

Do not replace the calculator with a new form or redesign the entire interface.

Do not duplicate backend mathematical logic in frontend code.

---

## 11. DOCKER — TWO SERVICES WITH COMPOSE

Implement the requested Docker configuration.

The local machine currently does **not** have Docker installed.

Do not install Docker automatically.

### Required architecture

Use exactly two application services:

**Frontend service**

- Build the React application for production.
- Serve the compiled static assets through Nginx.
- Provide a reverse proxy for `/api` requests to the backend.
- Expose the web interface on a documented host port.

**Backend service**

- Build the existing Go application.
- Run the existing API.
- Listen on an interface reachable by the frontend container.
- Be accessible to the frontend through Docker Compose networking.

### Networking

Use Docker Compose service-name resolution.

For example, the frontend's Nginx configuration may forward `/api` traffic to:

`http://backend:8080`

Ensure that the reverse proxy preserves the intended API paths.

Do not configure browser-side Fetch calls to target `http://backend:8080` directly.

The browser should make relative requests such as:

`/api/calculate`

Keep the backend port internal unless external exposure is actually necessary.

### Required files

Create or update, as appropriate:

- Backend Dockerfile.
- Frontend Dockerfile.
- `compose.yaml`.
- `.dockerignore` files.
- Minimal Nginx configuration.
- README deployment instructions.

Prefer multi-stage builds.

Use compatible base images and avoid unnecessary build tools in runtime images.

Use existing package lockfiles for reproducible frontend dependency installation.

Do not modify the normal local development workflow.

### Docker execution

The intended commands should be documented, including:

```bash
docker compose up --build -d
docker compose ps
docker compose logs
docker compose down
```

When Docker becomes available, the application must be tested in a real container environment.

### Current verification constraints

Because Docker is not installed:

1. Generate the Docker configuration.
2. Inspect the generated files for consistency.
3. Validate the YAML and configuration using available tools.
4. Verify that the referenced paths and build artifacts exist.
5. Check that service names, ports, build contexts, and proxy paths agree.
6. Run normal local frontend/backend tests.
7. Document the exact Docker commands needed for subsequent execution.

If a validation tool is unavailable, state that limitation.

**Never report a successful Docker build, container startup, or containerized end-to-end test unless it actually occurred.**

Classify Docker runtime verification as pending until it can be performed.

---

## 12. TESTING AND REGRESSION REQUIREMENTS

Preserve every existing meaningful test.

Extend tests to cover the new features and audit corrections.

### Backend tests

Test:

- Original arithmetic operations.
- Existing precedence.
- Unary operators.
- Negative and decimal values.
- Exponentiation.
- Right-associative powers.
- Negative exponents.
- Square roots.
- Mathematical domain errors.
- Contextual percentages.
- Combined advanced expressions.
- Parentheses and nesting.
- Parser depth limits.
- Invalid tokens and incomplete input.
- Numeric overflow.
- Non-finite results.

Include the complete acceptance examples from Sections 5–8.

### Combined-expression tests

At minimum:

| Expression | Expected result |
|---|---|
| `2+3*4` | 14 |
| `(2+3)*4` | 20 |
| `2^3+sqrt(16)` | 12 |
| `sqrt(25)*20%` | 1 |
| `200+10%` | 220 |
| `200-10%` | 180 |
| `200+10%+5%` | 231 |
| `2^3^2` | 512 |
| `-2^2` | -4 |
| `(-2)^2` | 4 |

### API tests

Verify:

- Valid advanced expressions.
- Correct JSON result structure.
- Existing arithmetic requests.
- Domain errors.
- Invalid syntax.
- Correct HTTP status codes.
- Consistent structured errors.

### Frontend tests

Verify:

- Exponentiation button.
- Square-root button.
- Percentage button.
- Parentheses.
- Keyboard input.
- Expression display.
- API request formatting.
- Result rendering.
- Clear and Backspace.
- Loading and error states.
- Recovery from failures.
- Existing calculator interactions.

Mock the HTTP boundary for isolated frontend tests.

Do not falsely claim that these mocked tests verify the real backend.

### Integration verification

Run the actual local frontend and Go backend.

Verify successful requests involving each advanced operation.

Use browser-level testing if available.

Preserve checks for disconnected servers, stale responses, and mathematical errors.

### Regression policy

After each important implementation stage, rerun the relevant tests.

Do not accept a change that enables new features but breaks previous behavior.

---

## 13. AUTOMATED QUALITY GATES

Run the existing applicable project commands.

### Go

Execute:

```bash
go vet ./...
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Run `gofmt` on changed Go files and verify formatting.

Execute these commands in the appropriate backend directory.

### React

Execute or configure equivalent project scripts for:

```bash
npm run lint
npm run typecheck
npm run test:coverage
npm run build
```

Use the project's real script names. Do not pretend that a script exists if it has not been configured.

### Coverage

Historical coverage figures may be used as context but must not be treated as current results or mandatory minimums.

Regenerate coverage after the changes.

Inspect meaningful untested branches, not just aggregate percentages.

Update `docs/COVERAGE.md` with the actual results.

Do not exclude relevant source files merely to improve the coverage percentage.

### Failed validation

If a test or build fails:

1. Identify the root cause.
2. Fix the implementation or incorrect test.
3. Rerun the affected checks.
4. Document any unresolved failure.

Do not suppress TypeScript, lint, or test errors to create an artificial passing result.

---

## 14. CLEAN CODE REVIEW

After implementation, perform a focused technical review.

### Go

Check:

- Parser readability.
- Correct precedence and associativity.
- Explicit mathematical errors.
- Reasonable function size and cohesion.
- Unnecessary allocations or duplication.
- Correct handling of recursion and input limits.
- Appropriate standard-library usage.
- No redundant evaluation engine.

### React

Check:

- Minimal, understandable state.
- Correct integration of new buttons.
- Predictable editing behavior.
- No duplicated mathematical calculations.
- Correct TypeScript types.
- No unnecessary hooks or dependencies.
- Preserved accessibility and responsive design.

### Docker

Check:

- Understandable configuration.
- Appropriate build stages.
- Reproducible dependency installation.
- Correct frontend-to-backend routing.
- Minimal unnecessary runtime dependencies.

Do not perform cosmetic refactoring unrelated to a concrete problem.

Do not introduce design patterns merely to demonstrate architectural sophistication.

---

## 15. DOCUMENTATION AND AI TRANSPARENCY

Update the existing English README.

Document:

- All supported operations.
- Exponentiation syntax and precedence.
- Square-root syntax.
- Contextual percentage rules.
- Parentheses.
- Mathematical limitations.
- Error handling.
- API examples.
- Local development setup.
- Docker Compose setup.
- Commands to verify the containerized application.
- Current Docker verification limitations.
- Testing and coverage.
- Relevant design decisions.

Include examples for:

- `2^3`
- `sqrt(81)`
- `200+10%`
- `200*10%`
- `2^3+sqrt(16)`

Explain significant implementation decisions, particularly the choice of contextual percentage behavior.

### PROMPTS.md

The employer requires actual prompts used during development.

Update `PROMPTS.md` with the real instructions used in this phase.

Include this implementation prompt accurately, using the exact received text where accessible.

Preserve earlier genuine prompts already recorded.

Do not invent earlier messages or reconstruct missing transcripts as though they were exact originals.

If an exact prompt is unavailable, report that limitation.

Do not automatically commit the private audit report or unrelated internal evaluation material.

Document the actual role of AI assistance truthfully.

---

## 16. GIT SAFETY AND PUBLICATION

Preserve the existing Git repository and its history.

Before finishing:

- Review the diff.
- Check for unintended changes.
- Check `.gitignore`.
- Check Docker build exclusions.
- Check documentation.
- Ensure no credentials or private files are included.
- Prepare meaningful local commits if appropriate.

Do not create or modify remote repositories.

Do not push.

Do not publish.

The hosting service will be selected later.

---

## 17. INDEPENDENT FINAL REVIEW

After implementing the changes and running the available checks, perform a fresh engineering review.

If an existing independent review agent is available and its use is appropriate, it may be used to provide a separate assessment.

Do not treat the previous audit score as proof that the updated project is correct.

Evaluate:

1. R01 and R02 corrections.
2. Mathematical correctness.
3. Exponentiation and precedence.
4. Square-root behavior.
5. Contextual percentage semantics.
6. Parser safety and edge cases.
7. Original functionality and regressions.
8. API compatibility.
9. Frontend interaction and accessibility.
10. Test quality and coverage.
11. Docker configuration.
12. Documentation.
13. Maintainability and unnecessary complexity.

Any newly identified critical defect must be addressed before claiming the implementation is ready.

Distinguish a genuine independent review from a self-review if no separate agent is available.

---

## 18. DEFINITION OF DONE

### Corrections

- [ ] R01 is corrected and covered by regression tests.
- [ ] R02 is corrected without removing visible focus.
- [ ] Available verification evidence is recorded.

### Exponentiation

- [ ] `^` implemented.
- [ ] Correct precedence.
- [ ] Right associativity.
- [ ] Negative exponent support.
- [ ] Mathematical domain errors controlled.
- [ ] Frontend button and keyboard support.
- [ ] Backend and frontend tests.

### Square root

- [ ] `sqrt(expression)` implemented.
- [ ] Valid results.
- [ ] Negative argument errors.
- [ ] Integration with other operations.
- [ ] Frontend button.
- [ ] Relevant tests.

### Contextual percentage

- [ ] `%` implemented.
- [ ] `200+10%=220`.
- [ ] `200-10%=180`.
- [ ] `200*10%=20`.
- [ ] `200/10%=2000`.
- [ ] Chained percentage behavior.
- [ ] Frontend controls.
- [ ] Semantics documented.
- [ ] Relevant tests.

### Docker

- [ ] Two-service Docker Compose configuration created.
- [ ] Frontend production container configuration created.
- [ ] Backend container configuration created.
- [ ] Reverse proxy configured.
- [ ] Build contexts and networking reviewed.
- [ ] `.dockerignore` included.
- [ ] Execution and test commands documented.
- [ ] Runtime verification accurately marked as pending if Docker remains unavailable.

### Quality and delivery

- [ ] Original mathematical operations still work.
- [ ] Parser handles advanced expressions.
- [ ] Backend tests pass.
- [ ] Frontend tests pass.
- [ ] TypeScript validation passes.
- [ ] Lint and static analysis pass.
- [ ] Build succeeds.
- [ ] Local integration is verified.
- [ ] Coverage reports updated.
- [ ] README updated.
- [ ] Actual prompts recorded.
- [ ] Git repository preserved.
- [ ] No unauthorized publication performed.

Do not mark unchecked or unverified items as passed.

---

## 19. FINAL REPORT

Provide a concise evidence-based summary covering:

1. The six objectives and their completion status.
2. Important files changed.
3. Parser and percentage-semantics decisions.
4. Tests executed and their real results.
5. Updated backend and frontend coverage.
6. Local full-stack integration results.
7. Docker architecture and verification status.
8. Documentation updates.
9. Independent review findings.
10. Remaining limitations and required user actions.

Use these status distinctions when relevant:

- **Implemented and verified**
- **Implemented but not runtime-verified**
- **Incomplete**
- **Blocked**

Do not invent results.

Do not declare the full application verified merely because individual unit tests pass.

### FINAL DIRECTIVE

Implement all six requested improvements without rebuilding the existing calculator.

Preserve the original architecture wherever it remains sound.

Use Go as the sole mathematical evaluation engine.

Implement the advanced operations with correct precedence, explicit domain handling, and conventional contextual percentages.

Preserve the application's usability, accessibility, existing tests, and API compatibility.

Deliver the two-service Docker Compose configuration without claiming container execution when Docker is unavailable.

**The result must be a more complete technical assessment project, not a more complicated one.**

Begin by inspecting the current repository and establishing the existing test baseline. Then execute the six objectives, validate every available requirement, and report the actual results.