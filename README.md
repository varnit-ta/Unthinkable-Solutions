# 🍳 Smart Recipe Generator

An intelligent recipe matching application with AI-powered ingredient detection from images!

## ✨ Features

- 🤖 **AI Image Recognition** - Upload photos and detect ingredients with local AI (no API key!)
- 🔍 **Smart Recipe Matching** - Find recipes based on available ingredients
- 🎯 **Advanced Filtering** - Filter by diet, difficulty, cuisine, and cooking time
- ⭐ **Favorites System** - Save your favorite recipes
- 💡 **Personalized Suggestions** - Get recipe recommendations based on your preferences
- 🔐 **User Authentication** - Secure login and registration
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- PostgreSQL 15+
- Python 3.11+ (for AI service)
- Docker & Docker Compose (recommended)

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/varnit-ta/Unthinkable-Solutions.git
cd Unthinkable-Solutions
```

2. **Setup environment variables**

Create `.env` in the root directory:
```env
# Database Configuration
DATABASE_URL=postgres://unthinkable:unthinkable@db:5432/unthinkable_recipes?sslmode=disable
POSTGRES_USER=unthinkable
POSTGRES_PASSWORD=unthinkable
POSTGRES_DB=unthinkable_recipes

# Backend Configuration
PORT=8081
JWT_SECRET=change-me-to-a-secure-secret

# AI Service Configuration (runs locally, no API key needed!)
AI_SERVICE_URL=http://ai-service:8000

# Image Processing
MAX_IMAGE_SIZE_MB=10

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000,http://localhost:4173
```

3. **Start with Docker Compose** (Recommended)
```bash
docker-compose up --build
```

Or **run locally**:

```bash
# Start PostgreSQL (if not using Docker)
# Then run migrations...

# Backend
cd backend
go run cmd/server/main.go

# Frontend (in another terminal)
cd frontend
npm install
npm run dev
```

4. **Access the application**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8081

## 🎯 AI Image Detection Setup

**Local Python AI Service (No API Key Required!)** ✨
1. The AI service runs automatically with Docker Compose
2. Uses Salesforce BLIP model for ingredient detection
3. Model downloads automatically on first run (~990MB)
4. Fast local inference - no external API calls!
5. Configured via `AI_SERVICE_URL` environment variable

**Architecture:**
```
Frontend → Go Backend → Python AI Service (FastAPI + BLIP)
```

**Configuration:**
- **Docker Compose**: `AI_SERVICE_URL=http://ai-service:8000` (default)
- **Local Development**: `AI_SERVICE_URL=http://localhost:8000`

**Manual Setup** (if not using Docker):
```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

See `ai-service/README.md` for detailed documentation.

Upload ingredient images and watch the magic! ✨

## 📚 Tech Stack

### Backend
- **Go** - Fast, efficient server
- **Chi Router** - Lightweight HTTP routing
- **PostgreSQL** - Reliable data storage
- **SQLC** - Type-safe SQL queries
- **Python + FastAPI** - Local AI microservice
- **Salesforce BLIP** - Image captioning model

### Frontend
- **React + TypeScript** - Modern UI development
- **Vite** - Fast build tool
- **Tailwind CSS** - Utility-first styling
- **Shadcn/ui** - Beautiful components
- **React Router** - Client-side routing

## 📖 API Documentation

### Authentication
- `POST /auth/register` - Create new account
- `POST /auth/login` - Login

### Recipes
- `GET /recipes` - List all recipes (with filters)
- `GET /recipes/:id` - Get recipe details
- `POST /match` - Find recipes matching ingredients

### AI Detection
- `POST /detect-ingredients` - Upload image to detect ingredients

### Favorites (Protected)
- `GET /favorites` - List user favorites
- `POST /favorites/:id` - Add to favorites
- `DELETE /favorites/:id` - Remove from favorites

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## 🌟 Key Features Explained

### AI Ingredient Detection
Upload an image of your ingredients or fridge contents, and our AI will automatically detect what's in the image. The system uses Hugging Face's BLIP-2 model for accurate food recognition.

**Supported:**
- Single ingredients
- Multiple ingredients
- Prepared dishes (extracts ingredients)
- Fridge/pantry photos

### Smart Recipe Matching
Once ingredients are detected (or manually entered), the app finds recipes that:
- Use the most of your available ingredients
- Match your dietary preferences
- Fit your skill level
- Meet your time constraints

### Personalized Suggestions
The app learns from your favorites and ratings to suggest recipes you'll love.