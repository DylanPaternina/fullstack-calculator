// Independent HTTP integration check. Run with the real Go and Vite servers up.
import assert from 'node:assert/strict'

const cases = [
  ['5+3', 8], ['5-3', 2], ['5*3', 15], ['10/4', 2.5],
  ['5+3*2', 11], ['10-6/2', 7], ['-3*2', -6], ['2*-3', -6], ['.5+1.25', 1.75],
]
for (const base of ['http://127.0.0.1:8080', 'http://127.0.0.1:5173']) {
  for (const [expression, expected] of cases) {
    const response = await fetch(`${base}/api/calculate`, {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression }), signal: AbortSignal.timeout(5000),
    })
    assert.equal(response.status, 200)
    assert.deepEqual(await response.json(), { result: expected })
  }
  for (const [expression, code] of [['5/0', 'division_by_zero'], ['5+*2', 'invalid_expression'], ['1e308*10', 'numeric_overflow']]) {
    const response = await fetch(`${base}/api/calculate`, {
      method: 'POST', body: JSON.stringify({ expression }), signal: AbortSignal.timeout(5000),
    })
    assert.equal(response.status, 400)
    const body = await response.json()
    assert.equal(body.error.code, code)
    assert.equal(typeof body.error.message, 'string')
  }
  console.log(`PASS ${base}: 9 calculations and 3 controlled errors`)
}
