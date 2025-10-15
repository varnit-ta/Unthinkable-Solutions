import { useEffect, useState } from 'react'
import { api } from '../api'
import { useAuth } from '../auth'

export default function SuggestionsPage() {
  const { token } = useAuth()
  const [list, setList] = useState<any[]>([])
  useEffect(() => {
    if (!token) return
    api.suggestions(token).then((l) => setList(l as any[])).catch(() => setList([]))
  }, [token])
  if (!token) return <div style={{ padding: 16 }}>Please login to view suggestions</div>
  return (
    <div style={{ padding: 16 }}>
      <h2>Suggestions</h2>
      <ul>
        {list.map((r) => (
          <li key={r.id}>{r.title}</li>
        ))}
      </ul>
    </div>
  )
}
