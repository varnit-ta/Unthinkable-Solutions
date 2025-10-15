import { useEffect, useState } from 'react'
import { api } from '../api'

export default function RecipesList() {
  const [recipes, setRecipes] = useState<any[]>([])
  useEffect(() => {
    const params = new URLSearchParams({ limit: '50' })
    api.listRecipes(params).then(setRecipes).catch(() => setRecipes([]))
  }, [])
  return (
    <div style={{ padding: 16 }}>
      <h2>Recipes</h2>
      <ul>
        {recipes.map((r: any) => (
          <li key={r.id}>{r.title}</li>
        ))}
      </ul>
    </div>
  )
}
