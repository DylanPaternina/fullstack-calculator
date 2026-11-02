# API reference

The same contract is available directly at `127.0.0.1:8080` during local development and through Nginx at `127.0.0.1:8088` with Docker. React uses relative `/api/calculate` requests. See the [README](../README.md#api-usage) for setup and representative basic/advanced examples.

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

## Representative cases

| Expression | HTTP | Result or error code |
| --- | ---: | --- |
| `5+3` | 200 | 8 |
| `5+3*2` | 200 | 11 |
| `2^3` | 200 | 8 |
| `sqrt(81)` | 200 | 9 |
| `2^3+sqrt(16)` | 200 | 12 |
| `200+10%` | 200 | 220 |
| `0/5` | 200 | `{"result":0}`, with no error property |
| `5/0` | 400 | division_by_zero |
| `5+*2` | 400 | invalid_expression |

## Through the Docker frontend

After starting Compose, use the published frontend port:

```sh
curl http://127.0.0.1:8088/
curl http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"2^3+sqrt(16)"}'
curl http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"200+10%"}'
curl -i http://127.0.0.1:8088/api/calculate -H 'Content-Type: application/json' -d '{"expression":"sqrt(-1)"}'
```

Expected: HTML 200, result 12, result 220 and HTTP 400/domain_error. Request bodies use the same contract as the direct API. The shell examples use POSIX quoting; use the Invoke-RestMethod example above in PowerShell.
