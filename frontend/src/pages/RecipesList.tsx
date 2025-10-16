/**
 * Recipes List Page Component
 * 
 * Displays a searchable and filterable list of recipes.
 * Features include:
 * - Search by recipe name
 * - Filter by diet, difficulty, cuisine, and cooking time
 * - Toggle favorites (for authenticated users)
 * - Responsive grid layout
 * - Loading states with skeletons
 * 
 * @module RecipesList
 */

import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../auth'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card'
import { Input } from '../components/ui/input'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select'
import { Skeleton } from '../components/ui/skeleton'
import { Search, Clock, Users, ChefHat, Heart } from 'lucide-react'
import { toast } from 'sonner'

/**
 * RecipesList Component
 * 
 * Main component for browsing and searching recipes with advanced filtering options.
 * Supports favorite management for authenticated users.
 * 
 * @returns {JSX.Element} Recipes list page component
 */
export default function RecipesList() {
  const { token } = useAuth()
  const navigate = useNavigate()
  const [recipes, setRecipes] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [diet, setDiet] = useState('')
  const [difficulty, setDifficulty] = useState('')
  const [cuisine, setCuisine] = useState('')
  const [maxTime, setMaxTime] = useState('')
  const [favorites, setFavorites] = useState<Set<number>>(new Set())

  /**
   * Fetch recipes from API with current filter settings
   * Also fetches favorite status for authenticated users
   */
  const fetchRecipes = async () => {
    setLoading(true)
    const params = new URLSearchParams({ limit: '50' })
    if (search) params.set('q', search)
    if (diet) params.set('diet', diet)
    if (difficulty) params.set('difficulty', difficulty)
    if (cuisine) params.set('cuisine', cuisine)
    if (maxTime) params.set('maxTime', maxTime)
    
    try {
      const data = await api.listRecipes(params)
      setRecipes(data)
      
      if (token && data.length > 0) {
        const favoriteChecks = await Promise.allSettled(
          data.map((recipe: any) => 
            api.isFavorite(token, recipe.id).then(res => ({ id: recipe.id, isFavorite: res.isFavorite }))
          )
        )
        
        const newFavorites = new Set<number>()
        favoriteChecks.forEach((result) => {
          if (result.status === 'fulfilled' && result.value.isFavorite) {
            newFavorites.add(result.value.id)
          }
        })
        setFavorites(newFavorites)
      }
    } catch (error) {
      setRecipes([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRecipes()
  }, [token])

  /**
   * Handle search button click
   */
  const handleSearch = () => {
    fetchRecipes()
  }

  /**
   * Clear all filters and reset to default recipe list
   */
  const clearFilters = () => {
    setSearch('')
    setDiet('')
    setDifficulty('')
    setCuisine('')
    setMaxTime('')
    const params = new URLSearchParams({ limit: '50' })
    api.listRecipes(params).then(setRecipes).catch(() => setRecipes([]))
  }

  /**
   * Toggle favorite status for a recipe
   * 
   * @param {number} recipeId - ID of the recipe to toggle
   * @param {React.MouseEvent} e - Click event
   */
  const handleFavoriteToggle = async (recipeId: number, e: React.MouseEvent) => {
    e.preventDefault()
    
    if (!token) {
      toast.error('Please login to add favorites')
      navigate('/login')
      return
    }

    const isFavorite = favorites.has(recipeId)
    
    try {
      if (isFavorite) {
        await api.removeFavorite(token, recipeId)
        setFavorites(prev => {
          const newSet = new Set(prev)
          newSet.delete(recipeId)
          return newSet
        })
        toast.success('Removed from favorites')
      } else {
        await api.addFavorite(token, recipeId)
        setFavorites(prev => new Set(prev).add(recipeId))
        toast.success('Added to favorites')
      }
    } catch (error) {
      toast.error('Failed to update favorite')
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">Discover Recipes</h1>
          <p className="text-muted-foreground">Browse and search through our collection of delicious recipes</p>
        </div>

        {/* Search and Filters */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Search & Filter</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search recipes..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                  className="pl-9"
                />
              </div>
              <Button onClick={handleSearch}>Search</Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label className="text-sm font-medium mb-2 block">Diet</label>
                <Select value={diet} onValueChange={setDiet}>
                  <SelectTrigger>
                    <SelectValue placeholder="Any" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value=" ">Any</SelectItem>
                    <SelectItem value="vegetarian">Vegetarian</SelectItem>
                    <SelectItem value="vegan">Vegan</SelectItem>
                    <SelectItem value="gluten-free">Gluten-Free</SelectItem>
                    <SelectItem value="keto">Keto</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label className="text-sm font-medium mb-2 block">Difficulty</label>
                <Select value={difficulty} onValueChange={setDifficulty}>
                  <SelectTrigger>
                    <SelectValue placeholder="Any" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value=" ">Any</SelectItem>
                    <SelectItem value="easy">Easy</SelectItem>
                    <SelectItem value="medium">Medium</SelectItem>
                    <SelectItem value="hard">Hard</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label className="text-sm font-medium mb-2 block">Cuisine</label>
                <Select value={cuisine} onValueChange={setCuisine}>
                  <SelectTrigger>
                    <SelectValue placeholder="Any" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value=" ">Any</SelectItem>
                    <SelectItem value="italian">Italian</SelectItem>
                    <SelectItem value="mexican">Mexican</SelectItem>
                    <SelectItem value="indian">Indian</SelectItem>
                    <SelectItem value="chinese">Chinese</SelectItem>
                    <SelectItem value="japanese">Japanese</SelectItem>
                    <SelectItem value="thai">Thai</SelectItem>
                    <SelectItem value="american">American</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <label className="text-sm font-medium mb-2 block">Max Time (min)</label>
                <Input
                  type="number"
                  placeholder="Any"
                  value={maxTime}
                  onChange={(e) => setMaxTime(e.target.value)}
                />
              </div>
            </div>

            <div className="flex gap-2">
              <Button onClick={handleSearch} className="flex-1">Apply Filters</Button>
              <Button onClick={clearFilters} variant="outline">Clear</Button>
            </div>
          </CardContent>
        </Card>

        {/* Recipe Grid */}
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
        ) : recipes.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <ChefHat className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
              <h3 className="text-lg font-semibold mb-2">No recipes found</h3>
              <p className="text-muted-foreground">Try adjusting your filters or search terms</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {recipes.map((recipe: any) => (
              <Card key={recipe.id} className="hover:shadow-lg transition-all hover:-translate-y-1 duration-200 relative">
                <CardHeader>
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <CardTitle className="text-lg line-clamp-2 flex-1">{recipe.title}</CardTitle>
                    <div className="flex items-center gap-2 shrink-0">
                      {recipe.average_rating && parseFloat(recipe.average_rating) > 0 && (
                        <Badge variant="secondary">
                          ⭐ {parseFloat(recipe.average_rating).toFixed(1)}
                        </Badge>
                      )}
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={(e) => handleFavoriteToggle(recipe.id, e)}
                      >
                        <Heart
                          className={`h-5 w-5 ${
                            favorites.has(recipe.id)
                              ? 'fill-red-500 text-red-500'
                              : 'text-muted-foreground'
                          }`}
                        />
                      </Button>
                    </div>
                  </div>
                  {recipe.description && (
                    <CardDescription className="line-clamp-2">
                      {recipe.description}
                    </CardDescription>
                  )}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex flex-wrap gap-2">
                    {recipe.cuisine && (
                      <Badge variant="outline">
                        🍽️ {recipe.cuisine}
                      </Badge>
                    )}
                    {recipe.difficulty && (
                      <Badge variant="secondary">
                        {recipe.difficulty}
                      </Badge>
                    )}
                    {recipe.diet_type && (
                      <Badge variant="outline">
                        {recipe.diet_type}
                      </Badge>
                    )}
                  </div>
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    {(recipe.total_time_minutes || recipe.cook_time_minutes) && (
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        <span>{recipe.total_time_minutes || recipe.cook_time_minutes} min</span>
                      </div>
                    )}
                    {recipe.servings && recipe.servings > 0 && (
                      <div className="flex items-center gap-1">
                        <Users className="h-4 w-4" />
                        <span>{recipe.servings}</span>
                      </div>
                    )}
                  </div>
                </CardContent>
                <CardFooter>
                  <Button asChild className="w-full">
                    <Link to={`/recipes/${recipe.id}`}>View Recipe</Link>
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
