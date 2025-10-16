# Frontend Documentation

## Overview

The Smart Recipe Generator frontend is a modern React application built with TypeScript, Vite, and Tailwind CSS. It provides an intuitive interface for discovering recipes, matching ingredients using AI vision, and managing personalized recipe collections.

## Technology Stack

- **React 18** - UI framework
- **TypeScript** - Type-safe development
- **Vite** - Fast build tool and dev server
- **React Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **shadcn/ui** - Reusable component library
- **Lucide React** - Icon library
- **Sonner** - Toast notifications

## Core Modules

### 1. API Client (`src/api.ts`)

Centralized API communication layer with type-safe methods.

**Key Features:**
- Generic request handler with authentication support
- Automatic error handling and parsing
- Support for JSON and FormData requests
- Environment-based API URL configuration

**Main Functions:**
- `login(email, password)` - User authentication
- `register(username, email, password)` - Account creation
- `listRecipes(params)` - Fetch recipes with filters
- `getRecipe(id)` - Get recipe details
- `match(ingredients, params)` - Find matching recipes
- `detectIngredients(file)` - AI image analysis
- `rate(token, recipeId, rating)` - Submit recipe rating
- `addFavorite(token, recipeId)` - Add to favorites
- `removeFavorite(token, recipeId)` - Remove from favorites
- `listFavorites(token)` - Get user's favorites
- `isFavorite(token, recipeId)` - Check favorite status
- `suggestions(token)` - Get personalized suggestions

### 2. Authentication (`src/auth.tsx`)

React Context-based authentication management.

**Components:**
- `AuthProvider` - Context provider for authentication state
- `useAuth()` - Hook to access authentication context

**Features:**
- Token persistence in localStorage
- Automatic token sync across tabs
- Protected route support
- Login/logout state management

### 3. Application Entry (`src/main.tsx`)

Application initialization with provider wrappers.

**Provider Stack:**
1. React.StrictMode - Development mode checks
2. BrowserRouter - Client-side routing
3. AuthProvider - Authentication state
4. Toaster - Notification system

### 4. Main Application (`src/App.tsx`)

Root component with routing and navigation.

**Components:**
- `Navigation` - Top navigation bar with authentication-aware menu
- `HomePage` - Landing page with feature overview
- `App` - Main routing component

**Routes:**
- `/` - Home page
- `/recipes` - Recipe listing
- `/recipes/:id` - Recipe details
- `/match` - Ingredient matching
- `/favorites` - User favorites (protected)
- `/suggestions` - Personalized suggestions (protected)
- `/login` - Login page
- `/register` - Registration page

## Pages

### LoginPage (`src/pages/LoginPage.tsx`)

User authentication interface.

**Features:**
- Email and password validation
- Error handling with toast notifications
- Link to registration page
- Redirect to home on success
- Loading state during authentication

### RegisterPage (`src/pages/RegisterPage.tsx`)

New user account creation.

**Features:**
- Username, email, and password fields
- Password minimum length validation (6 characters)
- Automatic login after registration
- Error handling with toast notifications
- Link to login page

### RecipesList (`src/pages/RecipesList.tsx`)

Main recipe browsing and search interface.

**Features:**
- Text search across recipe titles and descriptions
- Multi-filter support:
  - Diet type (vegetarian, vegan, gluten-free, keto)
  - Difficulty level (easy, medium, hard)
  - Cuisine type (Italian, Mexican, Indian, Chinese, Japanese, Thai, American)
  - Maximum cooking time
- Favorite toggle for authenticated users
- Responsive grid layout (1-3 columns)
- Loading skeletons
- Empty state handling
- Recipe cards with metadata (rating, time, servings)

### RecipeDetail (`src/pages/RecipeDetail.tsx`)

Comprehensive recipe information display.

**Features:**
- Full recipe details:
  - Title and description
  - Cuisine, difficulty, and diet badges
  - Cooking time and servings
  - Average rating display
- Complete ingredients list with quantities
- Step-by-step instructions
- Interactive rating system (1-5 stars)
- Favorite toggle
- Tag display
- Back navigation
- Loading states
- 404 handling

### MatchPage (`src/pages/MatchPage.tsx`)

AI-powered ingredient detection and recipe matching.

**Features:**
- Image upload with drag-and-drop
- Image preview
- AI ingredient detection using vision API
- Confidence score display
- AI-generated caption
- Manual ingredient management (add/remove)
- Recipe matching based on ingredients
- Filter matching results by diet, difficulty, cuisine, time
- Match score display
- Responsive grid layout

**Workflow:**
1. Upload image of ingredients
2. Detect ingredients using AI
3. Review and edit detected ingredients
4. Apply optional filters
5. Find matching recipes
6. View results with match scores

### FavoritesPage (`src/pages/FavoritesPage.tsx`)

User's saved favorite recipes.

**Features:**
- Authentication requirement check
- Grid display of favorite recipes
- Quick remove functionality
- Recipe metadata display
- Empty state with call-to-action
- Loading skeletons
- Direct navigation to recipe details

### SuggestionsPage (`src/pages/SuggestionsPage.tsx`)

Personalized recipe recommendations.

**Features:**
- Authentication requirement check
- AI-powered suggestions based on user behavior
- Visual "Suggested" badges
- Recipe metadata display
- Empty state with call-to-action
- Loading skeletons
- Grid layout

## Component Library

The application uses **shadcn/ui** components, a collection of reusable components built with Radix UI and Tailwind CSS.

### Key Components

- **Button** - Variants: default, secondary, outline, ghost, destructive
- **Card** - Container with header, content, and footer sections
- **Input** - Styled text input fields
- **Select** - Dropdown selection component
- **Badge** - Labels and tags with variants
- **Skeleton** - Loading state placeholders
- **Dialog** - Modal dialogs
- **Dropdown Menu** - Context menus and dropdowns
- **Avatar** - User profile images
- **Separator** - Visual dividers
- **Tabs** - Tabbed interface component
- **Label** - Form field labels
- **Toaster** - Toast notifications via Sonner

## State Management

### Local Component State

Each page manages its own state using React hooks:

- `useState` - Local state management
- `useEffect` - Side effects and data fetching
- `useNavigate` - Programmatic navigation
- `useParams` - URL parameter extraction
- `useLocation` - Current route information

### Global State

- **Authentication** - Managed by `AuthContext` via `useAuth()` hook
- **Toast Notifications** - Provided by Sonner's `toast` API

## Styling

### Tailwind CSS

Utility-first CSS framework for rapid UI development.

**Configuration:**
- Custom color scheme
- Responsive breakpoints
- Dark mode support (via CSS variables)
- Custom animations

**Common Patterns:**
- Responsive grid: `grid-cols-1 md:grid-cols-2 lg:grid-cols-3`
- Spacing: `space-y-{n}`, `gap-{n}`, `p-{n}`, `m-{n}`
- Typography: `text-{size}`, `font-{weight}`
- Colors: Theme-aware via CSS variables

### CSS Variables

Located in `src/index.css`:

```css
:root {
  --background
  --foreground
  --primary
  --primary-foreground
  --secondary
  --secondary-foreground
  --muted
  --muted-foreground
  --accent
  --border
  --input
  --ring
  --destructive
  --radius
}
```

## API Integration

### Environment Variables

Configure API endpoint via `.env`:

```env
VITE_API_URL=http://localhost:8081
```

### Error Handling

All API calls include error handling with user-friendly toast notifications:

```typescript
try {
  const data = await api.someMethod()
  toast.success('Operation successful!')
} catch (error) {
  toast.error(error.message || 'Operation failed')
}
```

### Authentication Flow

1. User logs in via `LoginPage`
2. API returns JWT token
3. Token stored in localStorage via `AuthProvider`
4. Token included in subsequent API requests
5. Protected routes check for token presence

## Development

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
cd frontend
npm install
```

### Development Server

```bash
npm run dev
```

Runs on `http://localhost:5173`

### Build

```bash
npm run build
```

Outputs to `dist/` directory

### Type Checking

```bash
npm run type-check
```

### Linting

```bash
npm run lint
```

## Best Practices

### Code Organization

1. **Single Responsibility** - Each component has one clear purpose
2. **Type Safety** - All components use TypeScript for type checking
3. **Documentation** - JSDoc comments for all major functions
4. **Error Handling** - Consistent error handling patterns
5. **Loading States** - Skeleton loaders for better UX

### Component Structure

```typescript
/**
 * Component documentation
 */
export default function MyComponent() {
  // State declarations
  const [state, setState] = useState()

  // Effects
  useEffect(() => {
    // Side effects
  }, [dependencies])

  // Event handlers
  const handleEvent = () => {
    // Handler logic
  }

  // Render
  return (
    // JSX
  )
}
```

### Naming Conventions

- **Components** - PascalCase (`LoginPage`, `RecipeCard`)
- **Functions** - camelCase (`fetchRecipes`, `handleSubmit`)
- **Constants** - UPPER_SNAKE_CASE (`BASE_URL`, `MAX_RESULTS`)
- **Types/Interfaces** - PascalCase (`Recipe`, `ApiError`)

## Performance Optimizations

1. **Code Splitting** - Route-based lazy loading via React Router
2. **Memoization** - `useMemo` for expensive computations
3. **Debouncing** - Search input debouncing (where applicable)
4. **Image Optimization** - Proper image sizing and formats
5. **Build Optimization** - Vite's optimized production builds

## Accessibility

- **Semantic HTML** - Proper use of HTML5 elements
- **ARIA Labels** - Screen reader support
- **Keyboard Navigation** - Full keyboard accessibility
- **Focus Management** - Visible focus indicators
- **Color Contrast** - WCAG compliant color schemes

## Testing Recommendations

### Unit Tests
- Test utility functions in isolation
- Test custom hooks behavior
- Test form validation logic

### Integration Tests
- Test API integration
- Test authentication flows
- Test routing behavior

### E2E Tests
- Test complete user workflows
- Test responsive design
- Test error scenarios

## Troubleshooting

### Common Issues

**Build Errors:**
- Clear `node_modules` and reinstall: `rm -rf node_modules && npm install`
- Clear Vite cache: `rm -rf node_modules/.vite`

**API Connection:**
- Verify `VITE_API_URL` in `.env`
- Check backend server is running
- Check CORS configuration

**Authentication Issues:**
- Clear localStorage: `localStorage.clear()`
- Verify token format
- Check token expiration

## Future Enhancements

1. **Progressive Web App (PWA)** - Offline support
2. **Internationalization (i18n)** - Multi-language support
3. **Advanced Search** - Full-text search with filters
4. **Recipe Collections** - User-created recipe lists
5. **Social Features** - Share recipes, follow users
6. **Recipe Creation** - User-submitted recipes
7. **Meal Planning** - Weekly meal planning feature
8. **Shopping Lists** - Generate shopping lists from recipes

## Resources

- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [shadcn/ui](https://ui.shadcn.com/)
- [React Router](https://reactrouter.com/)

## License

Part of the Smart Recipe Generator project.
