/**
 * API Module
 * 
 * This module provides a centralized API client for interacting with the backend services.
 * It handles HTTP requests, authentication, error handling, and type-safe responses.
 * 
 * @module api
 */

/**
 * JSON type alias for flexible data structures
 */
export type Json = Record<string, unknown> | unknown[] | string | number | boolean | null

/**
 * Standard API error structure
 */
export type ApiError = {
  status: number
  message: string
}

/**
 * Base URL for API requests, configurable via environment variable
 */
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081'

/**
 * Generic HTTP request handler with authentication and error handling
 * 
 * @template T - Expected response type
 * @param {string} path - API endpoint path
 * @param {RequestInit} opts - Fetch options
 * @param {string} token - Optional authentication token
 * @returns {Promise<T>} Parsed response data
 * @throws {ApiError} API error with status and message
 */
async function request<T>(
  path: string,
  opts: RequestInit = {},
  token?: string,
): Promise<T> {
  const headers: Record<string, string> = {
    ...(opts.headers as Record<string, string> | undefined),
  }

  if (!(opts.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...opts, headers })
  const text = await res.text()
  const data = text ? (JSON.parse(text) as T) : (undefined as unknown as T)

  if (!res.ok) {
    const message = (data as any)?.message || res.statusText
    throw { status: res.status, message } as ApiError
  }

  return data
}

/**
 * API client with typed methods for all backend endpoints
 */
export const api = {
  /**
   * Authenticate user with email and password
   * @param {string} email - User email
   * @param {string} password - User password
   * @returns {Promise<{token: string}>} Authentication token
   */
  login: (email: string, password: string) =>
    request<{ token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  /**
   * Register a new user account
   * @param {string} username - Desired username
   * @param {string} email - User email
   * @param {string} password - User password
   * @returns {Promise<{token: string}>} Authentication token
   */
  register: (username: string, email: string, password: string) =>
    request<{ token: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    }),

  /**
   * Fetch list of recipes with optional filters
   * @param {URLSearchParams} params - Query parameters for filtering
   * @returns {Promise<any[]>} Array of recipe objects
   */
  listRecipes: (params: URLSearchParams) =>
    request<any[]>(`/recipes?${params.toString()}`),

  /**
   * Get detailed information for a specific recipe
   * @param {number} id - Recipe ID
   * @returns {Promise<any>} Recipe details
   */
  getRecipe: (id: number) =>
    request<any>(`/recipes/${id}`),

  /**
   * Find recipes matching given ingredients
   * @param {string[]} ingredients - List of ingredients to match
   * @param {URLSearchParams} params - Optional filters
   * @returns {Promise<any>} Matching recipes
   */
  match: (ingredients: string[], params?: URLSearchParams) =>
    request(`/match${params && params.toString() ? `?${params.toString()}` : ''}`, {
      method: 'POST',
      body: JSON.stringify({ detectedIngredients: ingredients }),
    }),

  /**
   * Detect ingredients from an uploaded image using AI
   * @param {File} file - Image file to analyze
   * @returns {Promise<object>} Detection results with ingredients and metadata
   */
  detectIngredients: (file: File) => {
    const form = new FormData()
    form.append('image', file)

    return request<{
      detectedIngredients: string[]
      confidence?: number
      provider?: string
      caption?: string
      message?: string
    }>('/detect-ingredients', { method: 'POST', body: form })
  },

  /**
   * Submit a rating for a recipe
   * @param {string} token - Authentication token
   * @param {number} recipeId - Recipe ID to rate
   * @param {number} rating - Rating value (1-5)
   * @returns {Promise<any>} Rating confirmation
   */
  rate: (token: string, recipeId: number, rating: number) =>
    request('/ratings', {
      method: 'POST',
      body: JSON.stringify({ recipeId, rating }),
    }, token),

  /**
   * Add a recipe to user's favorites
   * @param {string} token - Authentication token
   * @param {number} recipeId - Recipe ID to favorite
   * @returns {Promise<any>} Success confirmation
   */
  addFavorite: (token: string, recipeId: number) =>
    request(`/favorites/${recipeId}`, { method: 'POST' }, token),

  /**
   * Remove a recipe from user's favorites
   * @param {string} token - Authentication token
   * @param {number} recipeId - Recipe ID to unfavorite
   * @returns {Promise<any>} Success confirmation
   */
  removeFavorite: (token: string, recipeId: number) =>
    request(`/favorites/${recipeId}`, { method: 'DELETE' }, token),

  /**
   * Get list of user's favorite recipes
   * @param {string} token - Authentication token
   * @returns {Promise<any>} Array of favorite recipes
   */
  listFavorites: (token: string) =>
    request('/favorites', {}, token),

  /**
   * Check if a specific recipe is in user's favorites
   * @param {string} token - Authentication token
   * @param {number} recipeId - Recipe ID to check
   * @returns {Promise<{isFavorite: boolean}>} Favorite status
   */
  isFavorite: (token: string, recipeId: number) =>
    request<{ isFavorite: boolean }>(`/favorites/${recipeId}`, {}, token),

  /**
   * Get personalized recipe suggestions based on user preferences
   * @param {string} token - Authentication token
   * @returns {Promise<any>} Array of suggested recipes
   */
  suggestions: (token: string) =>
    request('/suggestions', {}, token),
}


