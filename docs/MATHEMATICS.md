# Mathematical behavior

This reference retains the supported numeric policy and full examples. See the [README](../README.md#mathematical-behavior) for a concise overview and the [API reference](API.md) for request/error formats.

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
