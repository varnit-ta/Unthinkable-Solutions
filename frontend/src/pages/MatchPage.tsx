/**
 * Match Page Component
 * 
 * AI-powered ingredient detection and recipe matching page.
 * Features include:
 * - Image upload with preview
 * - AI-based ingredient detection using vision API
 * - Manual ingredient editing (add/remove)
 * - Recipe matching based on detected ingredients
 * - Filter matching results by diet, difficulty, cuisine, and time
 * - Display confidence scores and AI captions
 * 
 * @module MatchPage
 */

import { useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../components/ui/card'
import { Button } from '../components/ui/button'
import { Badge } from '../components/ui/badge'
import { Input } from '../components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select'
import { Upload, X, Search, Clock, Users, Camera } from 'lucide-react'
import { toast } from 'sonner'

/**
 * MatchPage Component
 * 
 * Allows users to upload images of ingredients and find matching recipes.
 * Uses AI vision API to detect ingredients and provides confidence metrics.
 * 
 * @returns {JSX.Element} Match page component
 */
export default function MatchPage() {
  const [file, setFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const [detected, setDetected] = useState<string[]>([])
  const [loading, setLoading] = useState(false)
  const [matchedRecipes, setMatchedRecipes] = useState<any[]>([])
  const [loadingRecipes, setLoadingRecipes] = useState(false)
  const [confidence, setConfidence] = useState<number | null>(null)
  const [caption, setCaption] = useState<string | null>(null)
  
  // Filters
  const [diet, setDiet] = useState('')
  const [difficulty, setDifficulty] = useState('')
  const [cuisine, setCuisine] = useState('')
  const [maxTime, setMaxTime] = useState('')

  /**
   * Handle file input change and generate preview
   * 
   * @param {React.ChangeEvent<HTMLInputElement>} e - File input change event
   */
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] ?? null
    setFile(selectedFile)
    
    if (selectedFile) {
      const reader = new FileReader()
      reader.onloadend = () => {
        setPreview(reader.result as string)
      }
      reader.readAsDataURL(selectedFile)
    } else {
      setPreview(null)
    }
  }

  /**
   * Clear uploaded file and reset all state
   */
  const clearFile = () => {
    setFile(null)
    setPreview(null)
    setDetected([])
    setMatchedRecipes([])
    setConfidence(null)
    setCaption(null)
  }

  /**
   * Use AI to detect ingredients from uploaded image
   */
  const detectIngredients = async () => {
    if (!file) return
    setLoading(true)
    try {
      const res = await api.detectIngredients(file)
      setDetected(res.detectedIngredients || [])
      setConfidence(res.confidence || null)
      setCaption(res.caption || null)
      
      if (res.detectedIngredients?.length > 0) {
        const confidencePercent = res.confidence ? Math.round(res.confidence * 100) : 0
        toast.success(`Detected ${res.detectedIngredients.length} ingredients! (${confidencePercent}% confidence)`)
      } else if (res.message) {
        toast.warning(res.message)
      } else {
        toast.info('No ingredients detected. Try another image or add manually.')
      }
    } catch (e: any) {
      const errorMsg = e?.message || 'Failed to detect ingredients'
      toast.error(errorMsg)
      setDetected([])
      setConfidence(null)
      setCaption(null)
    } finally {
      setLoading(false)
    }
  }

  /**
   * Find recipes that match detected ingredients with optional filters
   */
  const findMatchingRecipes = async () => {
    if (detected.length === 0) {
      toast.error('No ingredients detected')
      return
    }

    setLoadingRecipes(true)
    try {
      const params = new URLSearchParams()
      if (diet) params.set('diet', diet)
      if (difficulty) params.set('difficulty', difficulty)
      if (cuisine) params.set('cuisine', cuisine)
      if (maxTime) params.set('maxTime', maxTime)
      params.set('limit', '12')

      const res = await api.match(detected, params)
      setMatchedRecipes(res as any[])
      toast.success(`Found ${(res as any[]).length} matching recipes!`)
    } catch (e) {
      toast.error('Failed to find matching recipes')
      setMatchedRecipes([])
    } finally {
      setLoadingRecipes(false)
    }
  }

  /**
   * Remove an ingredient from detected list
   * 
   * @param {string} ingredient - Ingredient to remove
   */
  const removeIngredient = (ingredient: string) => {
    setDetected(detected.filter(d => d !== ingredient))
  }

  /**
   * Manually add an ingredient to the list
   * 
   * @param {string} ingredient - Ingredient to add
   */
  const addIngredient = (ingredient: string) => {
    if (ingredient && !detected.includes(ingredient)) {
      setDetected([...detected, ingredient])
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold mb-2">Match Ingredients</h1>
          <p className="text-muted-foreground">Upload an image of your ingredients and find matching recipes</p>
        </div>

        <div className="grid lg:grid-cols-2 gap-6">
          {/* Upload Section */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Camera className="h-5 w-5" />
                Upload Image
              </CardTitle>
              <CardDescription>Upload a photo of your ingredients</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {!preview ? (
                <label className="flex flex-col items-center justify-center border-2 border-dashed rounded-lg p-12 cursor-pointer hover:border-primary transition-colors">
                  <Upload className="h-12 w-12 text-muted-foreground mb-4" />
                  <span className="text-sm text-muted-foreground mb-2">Click to upload image</span>
                  <span className="text-xs text-muted-foreground">PNG, JPG, GIF up to 10MB</span>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleFileChange}
                    className="hidden"
                  />
                </label>
              ) : (
                <div className="space-y-4">
                  <div className="relative">
                    <img
                      src={preview}
                      alt="Preview"
                      className="w-full h-64 object-cover rounded-lg"
                    />
                    <Button
                      variant="destructive"
                      size="icon"
                      className="absolute top-2 right-2"
                      onClick={clearFile}
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                  <Button
                    onClick={detectIngredients}
                    disabled={loading}
                    className="w-full"
                  >
                    {loading ? 'Detecting...' : 'Detect Ingredients'}
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Detected Ingredients */}
          <Card>
            <CardHeader>
              <CardTitle>Detected Ingredients</CardTitle>
              <CardDescription>
                {detected.length > 0
                  ? `${detected.length} ingredient${detected.length !== 1 ? 's' : ''} detected`
                  : 'No ingredients detected yet'}
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {detected.length > 0 ? (
                <>
                  {/* Show AI detection metadata */}
                  {(confidence !== null || caption) && (
                    <div className="mb-4 p-3 bg-muted rounded-lg space-y-2">
                      {confidence !== null && (
                        <div className="flex items-center justify-between text-sm">
                          <span className="text-muted-foreground">AI Confidence:</span>
                          <Badge variant={confidence > 0.8 ? "default" : confidence > 0.6 ? "secondary" : "outline"}>
                            {Math.round(confidence * 100)}%
                          </Badge>
                        </div>
                      )}
                      {caption && (
                        <div className="text-xs text-muted-foreground">
                          <span className="font-medium">AI Caption:</span> "{caption}"
                        </div>
                      )}
                    </div>
                  )}
                  
                  <div className="flex flex-wrap gap-2">
                    {detected.map((ingredient) => (
                      <Badge key={ingredient} variant="secondary" className="text-sm py-1 px-3">
                        {ingredient}
                        <button
                          onClick={() => removeIngredient(ingredient)}
                          className="ml-2 hover:text-destructive"
                        >
                          <X className="h-3 w-3" />
                        </button>
                      </Badge>
                    ))}
                  </div>
                  
                  <div className="pt-4 border-t">
                    <label className="text-sm font-medium mb-2 block">Add ingredient manually</label>
                    <div className="flex gap-2">
                      <Input
                        placeholder="e.g., tomato"
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            addIngredient((e.target as HTMLInputElement).value)
                            ;(e.target as HTMLInputElement).value = ''
                          }
                        }}
                      />
                    </div>
                  </div>
                </>
              ) : (
                <div className="text-center py-12 text-muted-foreground">
                  <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
                  <p>Upload an image to detect ingredients</p>
                </div>
              )}
            </CardContent>
            {detected.length > 0 && (
              <CardFooter>
                <Button onClick={findMatchingRecipes} disabled={loadingRecipes} className="w-full">
                  {loadingRecipes ? 'Finding Recipes...' : 'Find Matching Recipes'}
                </Button>
              </CardFooter>
            )}
          </Card>
        </div>

        {/* Filters */}
        {detected.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Filter Results</CardTitle>
            </CardHeader>
            <CardContent>
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
            </CardContent>
          </Card>
        )}

        {/* Matched Recipes */}
        {matchedRecipes.length > 0 && (
          <div>
            <h2 className="text-2xl font-bold mb-4">
              Matching Recipes ({matchedRecipes.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {matchedRecipes.map((recipe: any) => (
                <Card key={recipe.id} className="hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <div className="flex items-start justify-between mb-2">
                      <CardTitle className="text-lg line-clamp-2">{recipe.title}</CardTitle>
                      {(recipe.score || recipe.matchCount) && (
                        <Badge variant="secondary" className="ml-2">
                          {recipe.score || recipe.matchCount} matches
                        </Badge>
                      )}
                    </div>
                    <CardDescription className="line-clamp-2">
                      {recipe.description || 'Delicious recipe to try'}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-2">
                    <div className="flex flex-wrap gap-2">
                      {recipe.cuisine && <Badge variant="outline">{recipe.cuisine}</Badge>}
                      {recipe.difficulty && <Badge variant="outline">{recipe.difficulty}</Badge>}
                      {recipe.dietType && <Badge variant="outline">{recipe.dietType}</Badge>}
                      {recipe.diet_type && <Badge variant="outline">{recipe.diet_type}</Badge>}
                    </div>
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        <span>{recipe.total_time_minutes || recipe.totalTime || recipe.cook_time_minutes || recipe.cookTime || 30} min</span>
                      </div>
                      {recipe.servings && (
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
          </div>
        )}
      </div>
    </div>
  )
}
