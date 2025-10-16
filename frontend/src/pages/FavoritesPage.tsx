import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../auth'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Skeleton } from '../components/ui/skeleton'
import { Heart, Clock, Users, LogIn } from 'lucide-react'

export default function FavoritesPage() {
  const { token } = useAuth()
  const [list, setList] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) {
      setLoading(false)
      return
    }
    
    setLoading(true)
    api.listFavorites(token)
      .then((l) => setList(l as any[]))
      .catch(() => setList([]))
      .finally(() => setLoading(false))
  }, [token])

  if (!token) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Card>
          <CardContent className="py-12 text-center">
            <LogIn className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">Login Required</h3>
            <p className="text-muted-foreground mb-4">Please login to view your favorite recipes</p>
            <Button asChild>
              <Link to="/login">Login</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold mb-2 flex items-center gap-2">
            <Heart className="h-8 w-8 fill-current text-red-500" />
            My Favorites
          </h1>
          <p className="text-muted-foreground">Your saved recipes in one place</p>
        </div>

        {loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <Card key={i}>
                <CardHeader>
                  <Skeleton className="h-6 w-3/4 mb-2" />
                  <Skeleton className="h-4 w-full" />
                </CardHeader>
                <CardContent>
                  <Skeleton className="h-4 w-full mb-2" />
                  <Skeleton className="h-4 w-2/3" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : list.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Heart className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
              <h3 className="text-lg font-semibold mb-2">No favorites yet</h3>
              <p className="text-muted-foreground mb-4">
                Start adding recipes to your favorites to see them here
              </p>
              <Button asChild>
                <Link to="/recipes">Browse Recipes</Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {list.map((recipe: any) => (
              <Card key={recipe.recipeId || recipe.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex items-start justify-between mb-2">
                    <CardTitle className="text-lg line-clamp-2">
                      {recipe.title || `Recipe ${recipe.recipeId}`}
                    </CardTitle>
                    {recipe.averageRating && (
                      <Badge variant="secondary" className="ml-2">
                        ⭐ {recipe.averageRating.toFixed(1)}
                      </Badge>
                    )}
                  </div>
                  <CardDescription className="line-clamp-2">
                    {recipe.description || 'Delicious recipe to try'}
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-2">
                  <div className="flex flex-wrap gap-2">
                    {recipe.cuisine && (
                      <Badge variant="outline">{recipe.cuisine}</Badge>
                    )}
                    {recipe.difficulty && (
                      <Badge variant="outline">{recipe.difficulty}</Badge>
                    )}
                    {recipe.dietType && (
                      <Badge variant="outline">{recipe.dietType}</Badge>
                    )}
                  </div>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Clock className="h-4 w-4" />
                      <span>{recipe.totalTime || recipe.cookTime || 30} min</span>
                    </div>
                    {recipe.servings && (
                      <div className="flex items-center gap-1">
                        <Users className="h-4 w-4" />
                        <span>{recipe.servings} servings</span>
                      </div>
                    )}
                  </div>
                </CardContent>
                <CardFooter>
                  <Button asChild className="w-full">
                    <Link to={`/recipes/${recipe.recipeId || recipe.id}`}>View Recipe</Link>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
