function record(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

// Validate the wire response: TypeScript types alone cannot validate JSON.
export async function calculate(expression: string, signal: AbortSignal): Promise<number> {
  let response: Response
  try {
    response = await fetch('/api/calculate', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ expression }), signal,
    })
  } catch (error) {
    if (signal.aborted) throw error
    throw new Error('Cannot reach the calculator. Check that the server is running and try again.', { cause: error })
  }
  let data: unknown
  try { data = await response.json() } catch { throw new Error('The server returned an invalid response. Try again.') }
  if (!response.ok) {
    if (record(data) && record(data.error) && typeof data.error.message === 'string') {
      throw new Error(data.error.message)
    }
    throw new Error('The calculation failed. Try again.')
  }
  if (!record(data) || typeof data.result !== 'number' || !Number.isFinite(data.result)) {
    throw new Error('The server returned an invalid result. Try again.')
  }
  return data.result
}
