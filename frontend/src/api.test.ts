import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './api'

afterEach(() => vi.unstubAllGlobals())
function respond(body: unknown, status = 200) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(body), { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

describe('HTTP boundary', () => {
  it('posts JSON and returns the server result', async () => {
    const fetchMock = respond({ result: 11 })
    const signal = new AbortController().signal
    expect(await calculate('5+3*2', signal)).toBe(11)
    expect(fetchMock).toHaveBeenCalledWith('/api/calculate', {
      method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{"expression":"5+3*2"}', signal,
    })
  })
  it('preserves structured server errors', async () => {
    respond({ error: { code: 'division_by_zero', message: 'Cannot divide by zero' } }, 400)
    await expect(calculate('5/0', new AbortController().signal)).rejects.toThrow('Cannot divide by zero')
  })
  it.each([{}, null, { error: null }, { error: { message: 5 } }])('handles unstructured HTTP failure %j', async body => {
    respond(body, 500)
    await expect(calculate('2', new AbortController().signal)).rejects.toThrow('The calculation failed')
  })
  it.each([{}, null, { result: '2' }, { result: null }])('rejects malformed success %j', async body => {
    respond(body)
    await expect(calculate('2', new AbortController().signal)).rejects.toThrow('invalid result')
  })
  it('rejects non-finite parsed values', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('{"result":1e999}')))
    await expect(calculate('2', new AbortController().signal)).rejects.toThrow('invalid result')
  })
  it('handles a non-JSON response', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('unavailable', { status: 502 })))
    await expect(calculate('2', new AbortController().signal)).rejects.toThrow('invalid response')
  })
  it('explains network failures', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))
    await expect(calculate('2', new AbortController().signal)).rejects.toThrow('Cannot reach the calculator')
  })
  it('preserves cancellation', async () => {
    const controller = new AbortController()
    controller.abort()
    const failure = new DOMException('Aborted', 'AbortError')
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(failure))
    await expect(calculate('2', controller.signal)).rejects.toBe(failure)
  })
})
