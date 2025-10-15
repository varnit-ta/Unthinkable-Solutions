import { Routes, Route, Link, Navigate } from 'react-router-dom'
import './App.css'

import RecipesList from './pages/RecipesList'
import RecipeDetail from './pages/RecipeDetail'
import MatchPage from './pages/MatchPage'
import FavoritesPage from './pages/FavoritesPage'
import SuggestionsPage from './pages/SuggestionsPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'

function Home() {
  return (
    <div style={{ padding: 16 }}>
      <h1>Smart Recipe Generator</h1>
      <nav style={{ display: 'flex', gap: 12 }}>
        <Link to="/recipes">Recipes</Link>
        <Link to="/match">Match</Link>
        <Link to="/favorites">Favorites</Link>
        <Link to="/suggestions">Suggestions</Link>
        <Link to="/login">Login</Link>
        <Link to="/register">Register</Link>
      </nav>
    </div>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/recipes" element={<RecipesList />} />
      <Route path="/recipes/:id" element={<RecipeDetail />} />
      <Route path="/match" element={<MatchPage />} />
      <Route path="/favorites" element={<FavoritesPage />} />
      <Route path="/suggestions" element={<SuggestionsPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
