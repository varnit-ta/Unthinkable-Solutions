# 🍳 Smart Recipe Generator

An intelligent recipe matching application with AI-powered ingredient detection from images!

## ✨ Features

- 🤖 **AI Image Recognition** - Upload photos of ingredients and let AI detect them automatically
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
- Hugging Face API Key (FREE - [Get one here](https://huggingface.co/settings/tokens))

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/varnit-ta/Unthinkable-Solutions.git
cd Unthinkable-Solutions
```

2. **Setup environment variables**

Create `backend/.env`:
```env
DATABASE_URL=postgres://unthinkable:unthinkable@localhost:5432/unthinkable_recipes?sslmode=disable
PORT=8081
JWT_SECRET=your-secret-key
HUGGINGFACE_API_KEY=hf_YOUR_TOKEN_HERE
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

See detailed setup guide: **[VISION_SETUP.md](./VISION_SETUP.md)**

Quick version:
1. Get free Hugging Face API key
2. Add to `backend/.env`: `HUGGINGFACE_API_KEY=hf_xxx`
3. Upload ingredient images and watch the magic! ✨

## 📚 Tech Stack

### Backend
- **Go** - Fast, efficient server
- **Chi Router** - Lightweight HTTP routing
- **PostgreSQL** - Reliable data storage
- **SQLC** - Type-safe SQL queries
- **Hugging Face API** - AI image recognition

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

## 📦 Project Structure

```
.
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── auth/           # Authentication logic
│   │   ├── config/         # Configuration
│   │   ├── db/             # Database queries (SQLC)
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Auth & logging middleware
│   │   ├── service/        # Business logic
│   │   └── vision/         # AI image detection 🆕
│   └── migrations/         # Database migrations
├── frontend/
│   └── src/
│       ├── components/     # Reusable UI components
│       ├── pages/          # Page components
│       └── api.ts          # API client
└── docker-compose.yml      # Docker orchestration
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

## 🤝 Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

MIT License

## 👨‍💻 Author

Created by [Varnit TA](https://github.com/varnit-ta)

---

**Need help?** Check out [VISION_SETUP.md](./VISION_SETUP.md) for AI setup or open an issue!
