import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api } from '../api'

export default function RecipeDetail() {
  const { id } = useParams()
  const [recipe, setRecipe] = useState<any | null>(null)
  useEffect(() => {
    if (!id) return
    api.getRecipe(Number(id)).then(setRecipe).catch(() => setRecipe(null))
  }, [id])
  if (!recipe) return <div style={{ padding: 16 }}>Loading...</div>
  return (
    <div style={{ padding: 16 }}>
      <h2>{recipe.title}</h2>
      <div>{recipe.description}</div>
    </div>
  )
}
