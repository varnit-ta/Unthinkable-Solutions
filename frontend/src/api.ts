export type Json = Record<string, unknown> | unknown[] | string | number | boolean | null

export type ApiError = { status: number; message: string }

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081'

async function request<T>(
  path: string,
  opts: RequestInit = {},
  token?: string,
): Promise<T> {
  const headers: Record<string, string> = {
    ...(opts.headers as Record<string, string> | undefined),
  }
  // If body is not FormData, default to JSON content-type
  if (!(opts.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }
  if (token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(`${BASE_URL}${path}`, { ...opts, headers })
  const text = await res.text()
  const data = text ? (JSON.parse(text) as T) : (undefined as unknown as T)
  if (!res.ok) {
    const message = (data as any)?.message || res.statusText
    throw { status: res.status, message } as ApiError
  }
  return data
}

export const api = {
  login: (email: string, password: string) =>
    request<{ token: string }>(`/auth/login`, { method: 'POST', body: JSON.stringify({ email, password }) }),
  register: (username: string, email: string, password: string) =>
    request<{ token: string }>(`/auth/register`, { method: 'POST', body: JSON.stringify({ username, email, password }) }),
  listRecipes: (params: URLSearchParams) => request<any[]>(`/recipes?${params.toString()}`),
  getRecipe: (id: number) => request<any>(`/recipes/${id}`),
  match: (ingredients: string[], params?: URLSearchParams) =>
    request(`/match${params && params.toString() ? `?${params.toString()}` : ''}`, {
      method: 'POST',
      body: JSON.stringify({ detectedIngredients: ingredients }),
    }),
  detectIngredients: (file: File) => {
    const form = new FormData()
    form.append('image', file)
    // uses request which will not set JSON content-type for FormData
    return request<{ 
      detectedIngredients: string[]
      confidence?: number
      provider?: string
      caption?: string
      message?: string
    }>(`/detect-ingredients`, { method: 'POST', body: form })
  },
  rate: (token: string, recipeId: number, rating: number) =>
    request(`/ratings`, { method: 'POST', body: JSON.stringify({ recipeId, rating }) }, token),
  addFavorite: (token: string, recipeId: number) => request(`/favorites/${recipeId}`, { method: 'POST' }, token),
  removeFavorite: (token: string, recipeId: number) => request(`/favorites/${recipeId}`, { method: 'DELETE' }, token),
  listFavorites: (token: string) => request(`/favorites`, {}, token),
  isFavorite: (token: string, recipeId: number) => request<{ isFavorite: boolean }>(`/favorites/${recipeId}`, {}, token),
  suggestions: (token: string) => request(`/suggestions`, {}, token),
}


