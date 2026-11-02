# FULL-STACK CALCULATOR — CODEX MASTER PROMPT

## ROLE AND OBJECTIVE

Act as a Senior Full-Stack Software Engineer, Software Architect, and Technical Code Reviewer.

Build a complete, professional full-stack calculator for a **Software Engineer Intern technical assessment**.

The employer explicitly permits AI-assisted development and requires the actual prompts used during implementation to be shared.

Your objective is to deliver a small, correct, maintainable, thoroughly tested, and well-documented application that a technical interviewer can understand, run, evaluate, and discuss.

**Core engineering principle:**

> Prefer the simplest correct and maintainable solution that fully satisfies the requirements. Do not confuse complexity with engineering quality.

Engineering priorities, in order:

1. Functional correctness and edge-case handling.
2. Clean, idiomatic, maintainable code.
3. Meaningful automated tests and verifiable behavior.
4. Simple, well-defined architecture.
5. Reliable frontend-backend integration.
6. Clear, reproducible documentation.
7. Responsive, accessible user experience.
8. Performance improvements with demonstrated value.
9. Optional functionality.

The employer suggests approximately 2–4 hours of development effort. Treat this as guidance for **appropriate project scope**, not a strict runtime deadline.

Additional time may be spent validating and improving the mandatory implementation. Do not use it as justification for unnecessary functionality or architectural complexity.

---

## 1. OFFICIAL REQUIREMENTS

Build a full-stack calculator with:

**Frontend**
- React.
- TypeScript, using strict mode.
- Intuitive interface for performing calculations.
- Input validation and error handling.
- Basic responsive support for mobile devices.

**Backend**
- Go.
- REST API.
- JSON requests and responses.
- Mathematical operations performed by the backend.
- Input validation and mathematical error handling.

**Mandatory operations**
- Addition.
- Subtraction.
- Multiplication.
- Division.

**Optional operations**
- Exponentiation.
- Square root.
- Percentage.

**Required deliverables**
- Git repository containing frontend and backend.
- README with setup, execution, API examples, assumptions, and design decisions.
- Unit tests for both layers.
- Real test coverage reports.
- Actual prompts used during AI-assisted development.

**Optional deliverable**
- Docker configuration.

The employer prefers Go, and Go has been selected for this implementation.

Do not replace it without explicit authorization.

---

## 2. EXECUTION AUTONOMY AND SAFETY

You are authorized to:

- Inspect the current working directory and available development tools.
- Create an isolated project.
- Implement and modify the application.
- Install necessary project-local dependencies.
- Run formatting, linting, tests, builds, and integration checks.
- Refactor code when justified by identifiable problems.
- Create and update documentation.
- Initialize Git and prepare logical local commits.

You must not:

- Modify unrelated projects or personal files.
- Install unnecessary global dependencies.
- Delete or overwrite unrelated data.
- Expose credentials or secrets.
- Create, publish, or push to a remote repository without explicit approval.
- Change repository visibility without authorization.
- Claim that tests passed when they were not executed.

**Work autonomously within these boundaries.**

Do not repeatedly request approval for routine implementation decisions.

Ask questions only when an essential requirement is genuinely ambiguous and cannot be resolved through a safe, documented engineering assumption.

If a tool is unavailable, use a reasonable alternative when possible. Report anything that could not be verified.

---

## 3. INITIAL ENVIRONMENT INSPECTION

Before creating files:

1. Inspect the working directory and Git status.
2. Check available Go, Node.js, and npm versions.
3. Confirm the environment supports the selected stack.
4. Identify any existing files that must be preserved.
5. Select compatible dependency versions.
6. Identify available testing and browser-validation tools.

If appropriate, create an isolated directory named `fullstack-calculator`.

Use current compatible, stable dependencies. Do not introduce a dependency solely because it is newer.

Produce a concise implementation plan identifying the major stages, technical decisions, and foreseeable risks.

Then proceed without waiting for approval unless an essential blocker exists.

---

## 4. TECHNOLOGY STACK

### Frontend

Use:

- React.
- TypeScript with strict type checking.
- Vite.
- Native Fetch API.
- CSS.
- Vitest.
- React Testing Library.

Use React's built-in state management.

Avoid Redux, unnecessary state libraries, large UI frameworks, and excessive dependencies.

Do not use `any`, suppress TypeScript errors, or introduce memoization hooks without a concrete reason.

### Backend

Use:

- Go.
- `net/http`.
- `encoding/json`.
- Standard Go testing tools.
- `net/http/httptest` for HTTP handler tests.

Prefer the Go standard library.

Avoid introducing a web framework, ORM, database, or third-party dependency without a clear and measurable benefit.

### Architecture

Keep the application divided into:

- Frontend presentation and interaction.
- Frontend HTTP communication.
- Backend HTTP transport and validation.
- Backend mathematical expression evaluation.

The mathematical engine must be independently testable without HTTP.

The backend must run independently as a single service.

Do not create additional microservices.

Use the simplest sensible project structure. Do not introduce unnecessary interfaces, repositories, factories, service layers, or wrapper files.

---

## 5. CALCULATOR USER EXPERIENCE

Implement a **classic calculator interface**, rather than a form containing two numeric fields.

The interface must include:

- A readable mathematical expression display.
- A result display.
- Numeric buttons.
- Decimal point.
- Addition, subtraction, multiplication, and division buttons.
- Equals button.
- Clear action.
- Delete-last-character action.
- Visible loading and error states.

### Keyboard support

Support:

- `0–9` for numbers.
- `+`, `-`, `*`, `/` for operators.
- `.` for decimal input.
- `Enter` or `=` to calculate.
- `Backspace` to remove the last character.
- `Escape` to clear.

Handle keyboard input without unnecessarily interfering with unrelated interactive elements.

### Expected behavior

- Users can construct mathematical expressions.
- The expression remains visible and editable.
- Pressing equals sends the expression to the backend.
- The backend returns the calculated result.
- The frontend displays the returned result.
- Errors do not permanently block further calculations.
- Subsequent input behaves consistently after successful calculations or errors.
- Repeated submissions are prevented while a request is pending.
- Outdated responses must not overwrite newer calculator state.

Support both mouse/touch and physical keyboard interaction.

### Critical restriction

**The frontend must never independently evaluate mathematical expressions.**

React may construct, display, edit, and perform basic syntactic validation on expressions.

All arithmetic evaluation, including operator precedence, must occur in Go.

Do not duplicate the mathematical engine in TypeScript.

---

## 6. MATHEMATICAL EXPRESSION ENGINE

Implement a small, deterministic mathematical expression evaluator in Go.

### Mandatory capabilities

Support:

- Integer values.
- Decimal values.
- Negative numbers.
- Unary minus.
- Addition.
- Subtraction.
- Multiplication.
- Division.
- Multiple operators in one expression.
- Standard mathematical operator precedence.

Required examples:

| Expression | Expected result |
|---|---|
| `5 + 3 * 2` | `11` |
| `10 - 6 / 2` | `7` |
| `-3 * 2` | `-6` |
| `2 * -3` | `-6` |
| `10 / 4` | `2.5` |
| `5 / 0` | Controlled error |
| `5 + * 2` | Invalid expression error |

Operators with equal precedence should follow standard left-to-right associativity.

Support ordinary mathematical whitespace.

### Parser implementation

Choose the simplest correct approach appropriate for the supported grammar, such as a small recursive-descent parser.

Do not implement an unnecessarily general-purpose expression engine.

Do not use:

- Arbitrary code execution.
- Dynamic evaluation facilities.
- Unrestricted expression libraries.
- Large parsing frameworks without justification.
- Regex-only evaluation shortcuts that cannot reliably enforce precedence.

Separate expression parsing/evaluation from HTTP handling.

### Parentheses

Parentheses are optional.

Implement them only after all mandatory functionality, tests, coverage, documentation, and integration checks are satisfactory.

If implemented, they must support correct grouping and have meaningful tests, including nested expressions.

Never introduce incomplete parentheses support merely to increase functionality.

### Numeric behavior

Use `float64` for calculations.

Implement a consistent policy for:

- Floating-point limitations.
- Decimal result presentation.
- Division by zero.
- Numeric overflow.
- Non-finite results.
- Malformed numeric input.
- Excessively large expressions.

Never return NaN or Infinity in JSON.

Choose a reasonable result-formatting policy and document it.

Presentation rounding must not silently change the underlying value used for later calculations.

Do not introduce arbitrary-precision arithmetic unless a verified requirement justifies it.

---

## 7. REST API CONTRACT

Prefer one simple calculation endpoint:

`POST /api/calculate`

### Request

```json
{
  "expression": "5 + 3 * 2"
}
```

### Successful response

```json
{
  "result": 11
}
```

### Error response

```json
{
  "error": {
    "code": "division_by_zero",
    "message": "Cannot divide by zero"
  }
}
```

Use consistent structured error responses.

### Required validation

Handle:

- Invalid JSON.
- Empty request bodies.
- Missing expression.
- Incorrect JSON field types.
- Unsupported or malformed syntax.
- Division by zero.
- Non-finite results.
- Excessively large request bodies.
- Unsupported HTTP methods.
- Multiple unexpected JSON values.

Use appropriate HTTP status codes:

- `200` for successful calculations.
- `400` for invalid requests or invalid mathematical expressions.
- `405` for unsupported methods.
- Appropriate client-error status for oversized requests.
- `500` only for unexpected internal failures.

Avoid exposing internal implementation details in client-facing errors.

Use reasonable body-size limits, request timeouts, and expression complexity limits.

The default backend port should be `8080`, with straightforward configuration where appropriate.

Use a Vite development proxy for `/api` requests.

Avoid unnecessarily permissive CORS configuration.

The API should be simple to run, test, and understand.

---

## 8. OPTIONAL MATHEMATICAL FEATURES

Do not implement optional features until the mandatory application is complete and verified.

Only consider:

- Exponentiation.
- Square root.
- Percentage.
- Parentheses.

Before implementing any optional feature, confirm:

1. All mandatory operations work.
2. Operator precedence is correct.
3. API validation and error handling work.
4. Backend and frontend tests pass.
5. Coverage reporting works.
6. Frontend and backend integration has been verified.
7. The build succeeds.
8. Required documentation is substantially complete.

Optional features must not compromise these conditions.

### Percentage semantics

If percentage is implemented, use explicitly documented conventional calculator behavior.

At minimum, consider these examples:

- `200 + 10% = 220`
- `200 * 10% = 20`

Because percentage behavior depends on context, define and test the exact supported semantics before implementation.

If the expected behavior cannot be implemented simply and consistently, omit percentage and document the decision.

Do not interpret percentage as modulo unless explicitly required.

### Other optional operations

If exponentiation or square root is implemented:

- Define clear syntax.
- Define operator precedence where applicable.
- Validate mathematical domains.
- Test errors and edge cases.
- Ensure the frontend and backend support the operation consistently.

Do not prioritize optional operations over code quality or test completeness.

---

## 9. FRONTEND QUALITY

Implement a clean, modern, restrained visual design.

Prefer:

- A centered calculator layout.
- Clear numeric display.
- Consistent spacing.
- Readable typography.
- Distinguishable operator buttons.
- Visible hover and focus states.
- Useful error feedback.
- Appropriate disabled states.

The application must work on desktop and mobile screens.

Check representative viewport widths around:

- 1440 px.
- 768 px.
- 375 px.

Avoid horizontal overflow.

Provide basic accessibility:

- Semantic buttons.
- Accessible names and labels.
- Keyboard navigation.
- Visible focus indicators.
- Readable contrast.
- Accessible status and error announcements.

Do not spend disproportionate effort on animations, decorative graphics, or unnecessary visual effects.

Use English for all interface text.

---

## 10. TESTING STRATEGY

Testing is a mandatory deliverable.

Tests must verify real behavior rather than merely increase a coverage percentage.

### Backend unit tests

Use table-driven tests where appropriate.

Cover:

- Every mandatory arithmetic operation.
- Operator precedence.
- Multiple operators.
- Negative numbers.
- Unary operators.
- Decimal values.
- Zero.
- Division by zero.
- Invalid expression syntax.
- Unsupported characters.
- Empty expressions.
- Excessively large or invalid numeric inputs.
- Non-finite results.
- Parentheses, if implemented.

Handle floating-point comparisons appropriately rather than assuming exact decimal representation for every operation.

### HTTP handler tests

Use `httptest`.

Verify:

- Valid calculation requests.
- JSON response structure.
- Correct status codes.
- Missing fields.
- Incorrect field types.
- Invalid JSON.
- Unexpected multiple JSON values.
- Empty bodies.
- Wrong HTTP method.
- Expression errors.
- Mathematical errors.
- Request-size limits.

### Frontend tests

Use Vitest and React Testing Library.

Verify:

- Initial render.
- Number-button input.
- Operator selection.
- Expression display.
- Keyboard input.
- Backspace and clear.
- API request submission.
- Correct display of backend results.
- Input errors.
- HTTP failures.
- Network failures.
- Loading state.
- Prevention of duplicate submissions.
- Recovery after errors.
- Handling of outdated responses where relevant.

Mock the HTTP boundary appropriately in component tests.

Do not couple tests to private implementation details.

### Coverage

Generate real coverage reports for both layers.

For Go, use commands such as:

- `go test ./... -coverprofile=coverage.out`
- `go tool cover -func=coverage.out`

For React, configure and execute a coverage command such as:

- `npm run test:coverage`

Ensure the coverage configuration includes relevant application source files, not only files imported by tests.

Create `docs/COVERAGE.md` containing:

- Commands executed.
- Actual coverage results.
- Relevant tested areas.
- Important gaps and limitations.

No arbitrary 100% coverage requirement.

Never fabricate coverage results or create superficial tests to inflate metrics.

---

## 11. DEVELOPMENT AND VALIDATION WORKFLOW

Work in controlled stages.

### Phase 1 — Project setup

- Inspect environment.
- Create project structure.
- Configure TypeScript, Go, development scripts, and tests.
- Confirm that the initial projects can run.

### Phase 2 — Backend mathematical engine

- Implement the restricted expression evaluator.
- Add mathematical unit tests.
- Validate precedence, error handling, and numerical edge cases.

Do not proceed with unverified mathematical logic.

### Phase 3 — HTTP API

- Implement and validate the calculation endpoint.
- Add handler tests.
- Verify representative requests independently of React.

### Phase 4 — React calculator

- Implement the classic calculator interface.
- Add keyboard support.
- Connect the frontend to the backend through Fetch.
- Handle loading, errors, and changing input.
- Add frontend behavioral tests.

### Phase 5 — Full-stack integration

Run the real frontend and backend together.

Verify calculations across the actual HTTP boundary.

Test at least:

- Basic arithmetic.
- Operator precedence.
- Negative numbers.
- Decimal numbers.
- Division by zero.
- Invalid expressions.
- Backend unavailable.

Use browser testing when available.

If browser testing is unavailable, perform the strongest feasible integration verification and explicitly state the limitation.

Do not present mocked component tests as proof of end-to-end integration.

### Phase 6 — Code review and refactoring

Review the implementation as a senior technical interviewer.

Inspect:

- Correctness.
- Readability.
- Function responsibilities.
- Code duplication.
- Error handling.
- Type safety.
- Parser correctness.
- Test quality.
- Architectural simplicity.
- Unnecessary dependencies.
- Performance concerns supported by evidence.

Correct important issues.

Do not refactor functioning code solely to create more files, abstractions, or design patterns.

### Phase 7 — Final quality gates

Run the applicable checks.

**Backend**

- `gofmt`
- `go vet ./...`
- `go test ./...`
- Coverage generation

**Frontend**

- Lint
- Tests
- Coverage generation
- TypeScript checking
- Production build

Ensure the project scripts support these checks.

If a command fails:

1. Diagnose the cause.
2. Correct the underlying problem.
3. Rerun the relevant validation.
4. Report the actual outcome.

Do not hide errors by weakening validation rules without a justified reason.

---

## 12. DOCUMENTATION

Create a professional, concise README in English.

Include:

- Project overview.
- Technology stack.
- Project structure.
- Requirements and dependency versions.
- Installation instructions.
- Backend startup commands.
- Frontend startup commands.
- API documentation.
- Working request/response examples.
- Curl examples.
- Testing instructions.
- Coverage results or report locations.
- Relevant architectural decisions.
- Assumptions and trade-offs.
- Known limitations.

Explain why the implementation uses:

- React and TypeScript.
- Go.
- Standard-library HTTP handling.
- A restricted backend expression evaluator.
- Separated UI, transport, and mathematical logic.

Document important decisions without turning the README into an extensive engineering report.

Another developer must be able to clone the repository, install dependencies, run the application, execute the tests, and understand the basic architecture using the README alone.

---

## 13. AI USAGE TRANSPARENCY

The employer explicitly requires the prompts used during development.

Create `PROMPTS.md`.

Preserve the **actual instructions used**, including:

- The original master prompt.
- Relevant follow-up prompts.
- Correction requests.
- Important AI-assisted review instructions.
- Tool and model names when known.

Record prompts faithfully, not reconstructed or invented versions.

If exact prompts are unavailable, identify the missing information instead of fabricating it.

Distinguish actual prompts from summaries of actions performed.

Do not claim manual authorship for AI-generated work.

Document how AI assistance was used and how the resulting implementation was validated.

Do not include private internal model reasoning.

Keep the documentation factual and professional.

---

## 14. GIT AND PUBLICATION

Use Git for version control.

Create a suitable `.gitignore`.

Do not commit:

- `node_modules`.
- Secrets or credentials.
- Tokens.
- Unnecessary temporary files.
- Local machine-specific information.
- Unrelated project files.

Review the final diff and repository contents.

Prepare clear local commits where appropriate.

The repository destination has not been selected yet.

**Do not create a remote repository, change repository visibility, or push code without my explicit authorization.**

Once implementation, validation, and documentation are complete, request authorization to publish.

After authorization, use the chosen Git hosting service and verify the resulting repository URL.

Never invent a repository link.

---

## 15. OPTIONAL DOCKER SUPPORT

Docker is not part of the initial implementation scope.

Do not begin development with Docker.

Only reconsider it after all mandatory requirements are complete, tested, and documented.

Add Docker only if it provides a meaningful improvement to reproducibility without introducing disproportionate complexity.

If omitted, a brief explanation is sufficient.

Do not introduce Kubernetes or unrelated deployment infrastructure.

---

## 16. DEFINITION OF DONE

The implementation is complete when all mandatory requirements have been verified.

### Functional

- [ ] Addition works.
- [ ] Subtraction works.
- [ ] Multiplication works.
- [ ] Division works.
- [ ] Operator precedence is correct.
- [ ] Negative and decimal inputs work.
- [ ] Invalid expressions are handled.
- [ ] Division by zero is handled.
- [ ] Arithmetic executes in Go, not React.

### Frontend

- [ ] Classic calculator interface works.
- [ ] Buttons work.
- [ ] Physical keyboard works.
- [ ] Responsive layout works.
- [ ] Loading and error states work.
- [ ] Accessibility basics are implemented.

### Backend

- [ ] REST endpoint works.
- [ ] JSON contract is consistent.
- [ ] Request validation works.
- [ ] Error responses are consistent.
- [ ] Mathematical logic is independent of HTTP.
- [ ] No non-finite values are returned in JSON.

### Quality

- [ ] Backend tests pass.
- [ ] Frontend tests pass.
- [ ] Coverage reports contain actual results.
- [ ] Lint and static checks pass.
- [ ] TypeScript checks pass.
- [ ] Frontend production build succeeds.
- [ ] Full-stack integration has been verified.
- [ ] Important review findings have been addressed.

### Documentation and delivery

- [ ] README is complete.
- [ ] API examples work.
- [ ] Design decisions are documented.
- [ ] Coverage is documented.
- [ ] Actual AI prompts are recorded.
- [ ] Git repository is prepared locally.
- [ ] No secrets or unrelated files are included.
- [ ] Publication awaits explicit authorization.

If a requirement cannot be verified, mark it as incomplete and explain why.

Do not silently remove mandatory requirements.

---

## 17. FINAL REPORT

After implementation, provide a concise factual report containing:

1. Implemented functionality and architecture.
2. Important engineering decisions.
3. Commands executed and their actual results.
4. Backend and frontend coverage figures.
5. Full-stack integration verification.
6. Optional features implemented or deliberately omitted.
7. Known issues or unverified behavior.
8. Documentation and Git status.
9. Whether publication authorization is required.

Do not reproduce the entire README in the final message.

Use precise statements supported by actual execution results.

---

## FINAL ENGINEERING DIRECTIVE

This is an internship technical assessment, not an enterprise software platform.

Demonstrate professional engineering judgment through:

- Correct behavior.
- Clear architecture.
- Readable, idiomatic code.
- Meaningful tests.
- Appropriate scope.
- Reliable integration.
- Honest documentation.
- Explainable decisions.

Avoid overengineering, premature optimization, unnecessary dependencies, and features added solely to appear impressive.

**Every architectural decision must have a concrete justification.**

**Every claim about correctness must be supported by appropriate verification.**

**Every optional feature must earn its complexity.**

A small, reliable, well-tested solution is better than an unnecessarily ambitious one.

Begin by inspecting the local environment, presenting a short implementation plan, and then executing the mandatory work autonomously within the authorization boundaries defined above.