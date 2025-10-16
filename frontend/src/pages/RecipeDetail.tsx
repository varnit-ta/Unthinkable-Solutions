import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../auth'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Skeleton } from '../components/ui/skeleton'
import { ArrowLeft, Clock, Users, ChefHat, Heart, Star } from 'lucide-react'
import { toast } from 'sonner'

export default function RecipeDetail() {
  const { id } = useParams()
  const { token } = useAuth()
  const navigate = useNavigate()
  const [recipe, setRecipe] = useState<any | null>(null)
  const [loading, setLoading] = useState(true)
  const [isFavorite, setIsFavorite] = useState(false)
  const [rating, setRating] = useState<number>(0)
  const [submittingRating, setSubmittingRating] = useState(false)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    api.getRecipe(Number(id))
      .then(setRecipe)
      .catch(() => setRecipe(null))
      .finally(() => setLoading(false))
  }, [id])

  const handleFavorite = async () => {
    if (!token) {
      toast.error('Please login to add favorites')
      navigate('/login')
      return
    }

    try {
      if (isFavorite) {
        await api.removeFavorite(token, Number(id))
        setIsFavorite(false)
        toast.success('Removed from favorites')
      } else {
        await api.addFavorite(token, Number(id))
        setIsFavorite(true)
        toast.success('Added to favorites')
      }
    } catch (error) {
      toast.error('Failed to update favorite')
    }
  }

  const handleRating = async (value: number) => {
    if (!token) {
      toast.error('Please login to rate recipes')
      navigate('/login')
      return
    }

    setSubmittingRating(true)
    try {
      await api.rate(token, Number(id), value)
      setRating(value)
      toast.success('Rating submitted!')
      // Refresh recipe to get updated average rating
      const updated = await api.getRecipe(Number(id))
      setRecipe(updated)
    } catch (error) {
      toast.error('Failed to submit rating')
    } finally {
      setSubmittingRating(false)
    }
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Button variant="ghost" asChild className="mb-6">
          <Link to="/recipes">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Recipes
          </Link>
        </Button>
        <Card>
          <CardHeader>
            <Skeleton className="h-8 w-3/4 mb-2" />
            <Skeleton className="h-4 w-full" />
          </CardHeader>
          <CardContent className="space-y-4">
            <Skeleton className="h-32 w-full" />
            <Skeleton className="h-32 w-full" />
          </CardContent>
        </Card>
      </div>
    )
  }

  if (!recipe) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <ChefHat className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">Recipe not found</h3>
            <Button asChild className="mt-4">
              <Link to="/recipes">Browse Recipes</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Button variant="ghost" asChild className="mb-6">
        <Link to="/recipes">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Recipes
        </Link>
      </Button>

      <div className="grid lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <CardTitle className="text-3xl mb-2">{recipe.title}</CardTitle>
                  <CardDescription className="text-base">{recipe.description}</CardDescription>
                </div>
                <Button
                  variant={isFavorite ? 'default' : 'outline'}
                  size="icon"
                  onClick={handleFavorite}
                  className="ml-4"
                >
                  <Heart className={`h-5 w-5 ${isFavorite ? 'fill-current' : ''}`} />
                </Button>
              </div>

              <div className="flex flex-wrap gap-2 mt-4">
                {recipe.cuisine && <Badge>{recipe.cuisine}</Badge>}
                {recipe.difficulty && <Badge variant="outline">{recipe.difficulty}</Badge>}
                {recipe.dietType && <Badge variant="secondary">{recipe.dietType}</Badge>}
              </div>

              <div className="flex items-center gap-6 mt-4 text-muted-foreground">
                <div className="flex items-center gap-2">
                  <Clock className="h-5 w-5" />
                  <span className="font-medium">
                    {recipe.totalTime || recipe.cookTime || 30} min
                  </span>
                </div>
                {recipe.servings && (
                  <div className="flex items-center gap-2">
                    <Users className="h-5 w-5" />
                    <span className="font-medium">{recipe.servings} servings</span>
                  </div>
                )}
                {recipe.averageRating && (
                  <div className="flex items-center gap-2">
                    <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
                    <span className="font-medium">{recipe.averageRating.toFixed(1)}</span>
                  </div>
                )}
              </div>
            </CardHeader>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Ingredients</CardTitle>
            </CardHeader>
            <CardContent>
              {recipe.ingredients && recipe.ingredients.length > 0 ? (
                <ul className="space-y-2">
                  {recipe.ingredients.map((ingredient: string, index: number) => (
                    <li key={index} className="flex items-start">
                      <span className="mr-2">•</span>
                      <span>{ingredient}</span>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-muted-foreground">No ingredients listed</p>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Instructions</CardTitle>
            </CardHeader>
            <CardContent>
              {recipe.instructions && recipe.instructions.length > 0 ? (
                <ol className="space-y-4">
                  {recipe.instructions.map((instruction: string, index: number) => (
                    <li key={index} className="flex gap-4">
                      <span className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center font-semibold">
                        {index + 1}
                      </span>
                      <span className="pt-1">{instruction}</span>
                    </li>
                  ))}
                </ol>
              ) : (
                <p className="text-muted-foreground">No instructions available</p>
              )}
            </CardContent>
          </Card>
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Rate this Recipe</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-2 justify-center">
                {[1, 2, 3, 4, 5].map((value) => (
                  <button
                    key={value}
                    onClick={() => handleRating(value)}
                    disabled={submittingRating}
                    className="transition-transform hover:scale-110 disabled:opacity-50"
                  >
                    <Star
                      className={`h-8 w-8 ${
                        value <= rating
                          ? 'fill-yellow-400 text-yellow-400'
                          : 'text-gray-300'
                      }`}
                    />
                  </button>
                ))}
              </div>
              {rating > 0 && (
                <p className="text-center mt-4 text-sm text-muted-foreground">
                  You rated this {rating} star{rating !== 1 ? 's' : ''}
                </p>
              )}
            </CardContent>
          </Card>

          {recipe.tags && recipe.tags.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Tags</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {recipe.tags.map((tag: string, index: number) => (
                    <Badge key={index} variant="secondary">
                      {tag}
                    </Badge>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
