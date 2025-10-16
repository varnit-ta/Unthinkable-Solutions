import { Routes, Route, Link, Navigate, useNavigate, useLocation } from 'react-router-dom'
import './App.css'
import { ChefHat, Home, Search, Heart, Lightbulb, LogIn, UserPlus, LogOut, User } from 'lucide-react'
import { Button } from './components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from './components/ui/dropdown-menu'
import { useAuth } from './auth'

import RecipesList from './pages/RecipesList'
import RecipeDetail from './pages/RecipeDetail'
import MatchPage from './pages/MatchPage'
import FavoritesPage from './pages/FavoritesPage'
import SuggestionsPage from './pages/SuggestionsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

function Navigation() {
  const { token, setToken } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    setToken(null)
    navigate('/login')
  }

  const isActive = (path: string) => location.pathname === path

  return (
    <nav className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-8">
            <Link to="/" className="flex items-center space-x-2 text-xl font-bold text-primary">
              <ChefHat className="h-6 w-6" />
              <span>Smart Recipe</span>
            </Link>
            
            <div className="hidden md:flex items-center space-x-1">
              <Button
                variant={isActive('/') ? 'secondary' : 'ghost'}
                size="sm"
                asChild
              >
                <Link to="/">
                  <Home className="mr-2 h-4 w-4" />
                  Home
                </Link>
              </Button>
              
              <Button
                variant={isActive('/recipes') ? 'secondary' : 'ghost'}
                size="sm"
                asChild
              >
                <Link to="/recipes">
                  <ChefHat className="mr-2 h-4 w-4" />
                  Recipes
                </Link>
              </Button>
              
              <Button
                variant={isActive('/match') ? 'secondary' : 'ghost'}
                size="sm"
                asChild
              >
                <Link to="/match">
                  <Search className="mr-2 h-4 w-4" />
                  Match
                </Link>
              </Button>
              
              {token && (
                <>
                  <Button
                    variant={isActive('/favorites') ? 'secondary' : 'ghost'}
                    size="sm"
                    asChild
                  >
                    <Link to="/favorites">
                      <Heart className="mr-2 h-4 w-4" />
                      Favorites
                    </Link>
                  </Button>
                  
                  <Button
                    variant={isActive('/suggestions') ? 'secondary' : 'ghost'}
                    size="sm"
                    asChild
                  >
                    <Link to="/suggestions">
                      <Lightbulb className="mr-2 h-4 w-4" />
                      Suggestions
                    </Link>
                  </Button>
                </>
              )}
            </div>
          </div>

          <div className="flex items-center space-x-2">
            {token ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm">
                    <User className="h-4 w-4 mr-2" />
                    Account
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuLabel>My Account</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link to="/favorites">
                      <Heart className="mr-2 h-4 w-4" />
                      Favorites
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link to="/suggestions">
                      <Lightbulb className="mr-2 h-4 w-4" />
                      Suggestions
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={handleLogout}>
                    <LogOut className="mr-2 h-4 w-4" />
                    Logout
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <>
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/login">
                    <LogIn className="mr-2 h-4 w-4" />
                    Login
                  </Link>
                </Button>
                <Button size="sm" asChild>
                  <Link to="/register">
                    <UserPlus className="mr-2 h-4 w-4" />
                    Register
                  </Link>
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}

function HomePage() {
  const { token } = useAuth()
  
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto text-center space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl md:text-6xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            Smart Recipe Generator
          </h1>
          <p className="text-xl text-muted-foreground">
            Discover amazing recipes, match ingredients, and create culinary masterpieces
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-6 mt-12">
          <Link
            to="/recipes"
            className="group p-6 border rounded-lg hover:border-primary transition-all hover:shadow-lg"
          >
            <ChefHat className="h-12 w-12 mx-auto mb-4 text-primary" />
            <h3 className="text-xl font-semibold mb-2">Browse Recipes</h3>
            <p className="text-muted-foreground">
              Explore thousands of recipes from various cuisines
            </p>
          </Link>

          <Link
            to="/match"
            className="group p-6 border rounded-lg hover:border-primary transition-all hover:shadow-lg"
          >
            <Search className="h-12 w-12 mx-auto mb-4 text-primary" />
            <h3 className="text-xl font-semibold mb-2">Match Ingredients</h3>
            <p className="text-muted-foreground">
              Upload an image and find recipes with your ingredients
            </p>
          </Link>

          <Link
            to={token ? "/suggestions" : "/login"}
            className="group p-6 border rounded-lg hover:border-primary transition-all hover:shadow-lg"
          >
            <Lightbulb className="h-12 w-12 mx-auto mb-4 text-primary" />
            <h3 className="text-xl font-semibold mb-2">Get Suggestions</h3>
            <p className="text-muted-foreground">
              Personalized recipe recommendations just for you
            </p>
          </Link>
        </div>
      </div>
    </div>
  )
}

function App() {
  return (
    <div className="min-h-screen bg-background">
      <Navigation />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/recipes" element={<RecipesList />} />
        <Route path="/recipes/:id" element={<RecipeDetail />} />
        <Route path="/match" element={<MatchPage />} />
        <Route path="/favorites" element={<FavoritesPage />} />
        <Route path="/suggestions" element={<SuggestionsPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </div>
  )
}

export default App
