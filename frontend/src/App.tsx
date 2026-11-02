import { useEffect, useRef, useState } from 'react'
import { calculate } from './api'

const MAX_LENGTH = 1024
const keys = [
  { value: 'clear', label: 'AC', name: 'Clear', kind: 'utility' },
  { value: 'delete', label: '⌫', name: 'Delete last character', kind: 'utility' },
  { value: '/', label: '÷', name: 'Divide', kind: 'operator' },
  { value: '*', label: '×', name: 'Multiply', kind: 'operator' },
  ...['7', '8', '9'].map(value => ({ value, label: value, name: value, kind: '' })),
  { value: '-', label: '−', name: 'Subtract', kind: 'operator' },
  ...['4', '5', '6'].map(value => ({ value, label: value, name: value, kind: '' })),
  { value: '+', label: '+', name: 'Add', kind: 'operator' },
  ...['1', '2', '3'].map(value => ({ value, label: value, name: value, kind: '' })),
  { value: '=', label: '=', name: 'Calculate', kind: 'equals' },
  { value: '0', label: '0', name: '0', kind: 'zero' },
  { value: '.', label: '.', name: 'Decimal point', kind: '' },
]

function formatResult(value: number): string {
  return Number(value.toPrecision(12)).toString()
}

export default function App() {
  const [expression, setExpression] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const active = useRef<AbortController | null>(null)
  const completed = useRef(false)
  const keyboard = useRef<(key: string) => void>(() => {})

  useEffect(() => {
    function onKey(event: KeyboardEvent) {
      const target = event.target
      if (event.ctrlKey || event.metaKey || event.altKey) return
      if (target instanceof HTMLInputElement && target.id === 'expression') {
        if (['Enter', '=', 'Escape'].includes(event.key)) {
          event.preventDefault()
          keyboard.current(event.key)
        }
        return
      }
      if (target instanceof HTMLElement && (target.isContentEditable || target.closest('input, textarea, select, a, [contenteditable]'))) return
      // Enter calculates within the keypad; unrelated buttons retain activation.
      if (event.key === 'Enter' && target instanceof HTMLButtonElement && !target.closest('.calculator')) return
      if (/^[0-9+*/.=-]$/.test(event.key) || ['Enter', 'Escape', 'Backspace'].includes(event.key)) {
        event.preventDefault()
        keyboard.current(event.key)
      }
    }
    window.addEventListener('keydown', onKey)
    return () => {
      window.removeEventListener('keydown', onKey)
      active.current?.abort()
      active.current = null
    }
  }, [])

  function cancel() {
    active.current?.abort()
    active.current = null
    setLoading(false)
  }

  function edit(next: string) {
    cancel()
    completed.current = false
    setResult(null)
    if (next.length > MAX_LENGTH) { setError('Expression must be at most 1024 characters.'); return }
    setError('')
    setExpression(next)
  }

  async function submit() {
    if (active.current) return
    if (!expression.trim()) { setError('Enter an expression first.'); return }
    const controller = new AbortController()
    active.current = controller
    completed.current = false
    setLoading(true)
    setError('')
    setResult(null)
    const timeout = window.setTimeout(() => {
      if (active.current === controller) {
        cancel()
        setError('The calculation timed out. Try again.')
      }
    }, 10_000)
    try {
      const value = await calculate(expression, controller.signal)
      if (active.current !== controller) return
      setResult(value)
      completed.current = true
    } catch (failure) {
      if (active.current !== controller) return
      setError(failure instanceof Error ? failure.message : 'The calculation failed. Try again.')
    } finally {
      window.clearTimeout(timeout)
      if (active.current === controller) { active.current = null; setLoading(false) }
    }
  }

  function press(key: string) {
    if (key === '=' || key === 'Enter') { void submit(); return }
    if (key === 'clear' || key === 'Escape') { edit(''); return }
    if (key === 'delete' || key === 'Backspace') { edit(expression.slice(0, -1)); return }
    // Chaining uses the full backend value, never the rounded presentation.
    const base = completed.current ? (/^[+*/-]$/.test(key) && result !== null ? String(result) : '') : expression
    edit(base + key)
  }

  useEffect(() => { keyboard.current = press })

  return (
    <main>
      <header className="page-heading">
        <span className="brand-mark" aria-hidden="true">=</span>
        <span>EVERYDAY TOOLS</span>
      </header>
      <section className="workspace" aria-labelledby="title">
        <div className="intro">
          <p className="eyebrow">A LITTLE CLARITY</p>
          <h1 id="title">Calculator<span>.</span></h1>
          <p>The numbers, made simple.</p>
        </div>
        <section className="calculator" aria-label="Calculator">
          <div className="display">
            <label htmlFor="expression">EXPRESSION</label>
            <input
              id="expression" value={expression} onChange={event => edit(event.target.value)}
              placeholder="0" autoComplete="off" spellCheck={false} inputMode="decimal"
              aria-describedby="input-help calculation-error" aria-invalid={Boolean(error)}
            />
            <div className="result-line">
              <span aria-hidden="true">=</span>
              <output aria-label="Result" title={result === null ? undefined : String(result)}>
                {result === null ? '—' : formatResult(result)}
              </output>
            </div>
            <p className="status" role="status" aria-label="Calculation status" aria-live="polite">
              {loading ? 'Calculating…' : result === null ? 'Ready when you are' : 'Calculated'}
            </p>
          </div>
          <div className="keypad">
            {keys.map(key => (
              <button
                key={key.value} type="button" className={key.kind} aria-label={key.name}
                disabled={key.value === '=' && loading} onClick={() => press(key.value)}
              >
                {key.label}
              </button>
            ))}
          </div>
          <div className="feedback"><p id="calculation-error" role="alert">{error}</p></div>
        </section>
        <p id="input-help" className="help">Type an expression or use the keys above.<br /><kbd>Enter</kbd> calculate <span>·</span> <kbd>Esc</kbd> clear <span>·</span> <kbd>⌫</kbd> delete</p>
      </section>
      <footer>Four operations. Endless possibilities.</footer>
    </main>
  )
}
