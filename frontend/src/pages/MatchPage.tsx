import { useState } from 'react'
import { api } from '../api'

export default function MatchPage() {
  const [file, setFile] = useState<File | null>(null)
  const [detected, setDetected] = useState<string[] | null>(null)
  const [loading, setLoading] = useState(false)

  async function submit() {
    if (!file) return
    setLoading(true)
    try {
      const res = await api.detectIngredients(file)
      setDetected(res.detectedIngredients)
    } catch (e) {
      setDetected([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ padding: 16 }}>
      <h2>Match Ingredients</h2>
      <input type="file" accept="image/*" onChange={(e) => setFile(e.target.files?.[0] ?? null)} />
      <button onClick={submit} disabled={!file || loading} style={{ marginLeft: 8 }}>
        {loading ? 'Detecting...' : 'Detect'}
      </button>
      {detected && (
        <div>
          <h3>Detected</h3>
          <ul>
            {detected.map((d) => (
              <li key={d}>{d}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
