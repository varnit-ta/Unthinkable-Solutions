import { useEffect, useState } from 'react'
import { api } from '../api'
import { useAuth } from '../auth'

export default function FavoritesPage() {
  const { token } = useAuth()
  const [list, setList] = useState<any[]>([])
  useEffect(() => {
    if (!token) return
    api.listFavorites(token).then((l) => setList(l as any[])).catch(() => setList([]))
  }, [token])
  if (!token) return <div style={{ padding: 16 }}>Please login to view favorites</div>
  return (
    <div style={{ padding: 16 }}>
      <h2>Favorites</h2>
      <ul>
        {list.map((f) => (
          <li key={f.recipeId?.toString() ?? f.id}>{f.title ?? f.recipeId}</li>
        ))}
      </ul>
    </div>
  )
}
