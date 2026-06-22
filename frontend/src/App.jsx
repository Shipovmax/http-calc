import { useState } from 'react'
import './App.css'

const OPERATORS = ['+', '-', '*', '/']

export default function App() {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [op, setOp] = useState('+')
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(false)

  async function handleCalculate() {
    const numA = parseFloat(a)
    const numB = parseFloat(b)

    if (isNaN(numA) || isNaN(numB)) {
      setError('Введи два числа')
      setResult(null)
      return
    }

    setLoading(true)
    setError(null)
    setResult(null)

    try {
      const res = await fetch('http://localhost:8080/calculate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ a: numA, op, b: numB }),
      })
      const data = await res.json()
      if (data.error) {
        setError(data.error)
      } else {
        setResult(data.result)
      }
    } catch {
      setError('Сервер недоступен. Запусти: go run .')
    } finally {
      setLoading(false)
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter') handleCalculate()
  }

  return (
    <div className="page">
      <div className="card">
        <div className="card-header">
          <span className="badge">Go + React</span>
          <h1>HTTP Калькулятор</h1>
          <p className="subtitle">REST API на чистом net/http</p>
        </div>

        <div className="inputs">
          <input
            className="input"
            type="number"
            placeholder="Число A"
            value={a}
            onChange={e => setA(e.target.value)}
            onKeyDown={handleKeyDown}
          />

          <div className="ops">
            {OPERATORS.map(o => (
              <button
                key={o}
                className={`op-btn ${op === o ? 'active' : ''}`}
                onClick={() => setOp(o)}
              >
                {o}
              </button>
            ))}
          </div>

          <input
            className="input"
            type="number"
            placeholder="Число B"
            value={b}
            onChange={e => setB(e.target.value)}
            onKeyDown={handleKeyDown}
          />
        </div>

        <button
          className="calc-btn"
          onClick={handleCalculate}
          disabled={loading}
        >
          {loading ? 'Считаю...' : 'Посчитать'}
        </button>

        {result !== null && (
          <div className="result">
            <span className="result-label">Результат</span>
            <span className="result-value">{result}</span>
          </div>
        )}

        {error && (
          <div className="error">
            {error}
          </div>
        )}

        <div className="request-preview">
          <span className="preview-label">POST /calculate</span>
          <code>{`{"a":${a || '?'},"op":"${op}","b":${b || '?'}}`}</code>
        </div>
      </div>
    </div>
  )
}
