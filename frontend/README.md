# 🎨 Smart Recipe Generator - Frontend

A modern, responsive React frontend for the Smart Recipe Generator application.

## Features

- 🤖 **AI-Powered Ingredient Detection** - Upload images to detect ingredients
- 🔍 **Smart Recipe Search** - Find recipes based on available ingredients
- 🎯 **Advanced Filtering** - Filter by diet, cuisine, difficulty, and cooking time
- ⭐ **Favorites Management** - Save and manage favorite recipes
- 💡 **Personalized Suggestions** - Get recipe recommendations
- 🔐 **User Authentication** - Secure login and registration
- 📱 **Fully Responsive** - Works on desktop, tablet, and mobile
- 🎨 **Beautiful UI** - Built with Tailwind CSS and Shadcn/ui components

## Tech Stack

- **React 18** - Modern React with hooks
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **Shadcn/ui** - Beautiful, accessible component library
- **React Router** - Client-side routing
- **Lucide React** - Modern icon library

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
cd frontend
npm install
```

### Development

```bash
# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Configuration

The frontend connects to the backend API. Update `src/api.ts` if your backend URL differs:

```typescript
const API_BASE_URL = 'http://localhost:8081';
```

## Project Structure

```
frontend/
├── src/
│   ├── components/         # Reusable UI components
│   │   └── ui/            # Shadcn/ui components
│   ├── pages/             # Page components
│   │   ├── LoginPage.tsx
│   │   ├── RegisterPage.tsx
│   │   ├── RecipesList.tsx
│   │   ├── RecipeDetail.tsx
│   │   ├── MatchPage.tsx
│   │   ├── FavoritesPage.tsx
│   │   └── SuggestionsPage.tsx
│   ├── api.ts             # API client with typed endpoints
│   ├── auth.tsx           # Authentication context & hooks
│   ├── App.tsx            # Main app component with routing
│   └── main.tsx           # Application entry point
├── public/                # Static assets
├── Dockerfile             # Container definition
└── package.json           # Dependencies and scripts
```

## Key Features

### Authentication
- JWT-based authentication
- Protected routes
- Persistent login (localStorage)
- Automatic token refresh

### Recipe Discovery
- Browse all recipes
- Advanced filtering (diet, cuisine, difficulty, time)
- Search by ingredients
- Recipe details with nutrition info

### AI Image Upload
- Drag & drop or click to upload
- Real-time ingredient detection
- Automatic recipe matching

### Favorites
- Save favorite recipes
- Quick access to saved recipes
- One-click add/remove

### Personalized Suggestions
- Algorithm-based recommendations
- Based on favorites and preferences
- Refreshable suggestions

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Docker

Build and run with Docker:

```bash
docker build -t recipe-frontend .
docker run -p 3000:80 recipe-frontend
```

Or use Docker Compose from project root:

```bash
docker-compose up frontend
```

## Contributing

1. Follow the existing code style
2. Use TypeScript for type safety
3. Keep components small and focused
4. Write meaningful commit messages

## License

Same as parent project - see root LICENSE file.
