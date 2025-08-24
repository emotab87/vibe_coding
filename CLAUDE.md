# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RealWorld Conduit app implementing **Vibe Coding** techniques and **Armin Ronacher's development philosophy**.

### Core Development Philosophy
- **Minimize dependencies**: Direct implementation over new libraries
- **Explicit over implicit**: Direct SQL over ORM, clear code over complex abstractions  
- **AI-assisted development**: Leverage generative AI for code generation and problem solving

## Architecture

### Backend (Go)
```
backend/
├── cmd/              # Application entry point
├── internal/
│   ├── handlers/     # HTTP handlers (routing logic)
│   ├── models/       # Data model structs
│   ├── middleware/   # Auth, logging middleware
│   └── database/     # Database logic (direct SQL)
└── migrations/       # SQLite schema migrations
```

### Frontend (React + TypeScript)
```
frontend/src/
├── components/       # Reusable UI components
├── pages/           # Route-specific page components
├── hooks/           # Custom React hooks
├── services/        # API communication (axios-based)
└── types/           # TypeScript type definitions
```

## Development Commands

### Docker (Recommended)
```bash
docker-compose up          # Run full environment
docker-compose up -d       # Run in background
docker-compose up backend  # Run specific service
```

### Local Development
```bash
# Backend
cd backend && go mod tidy && go run cmd/main.go

# Frontend  
cd frontend && npm install && npm run dev
npm run build  # Production build
```

## Tech Stack

### Backend
- **Go 1.21+**: Web server
- **Gorilla Mux**: HTTP router (minimal dependencies)
- **SQLite**: Database (direct SQL, no ORM)
- **JWT**: Authentication (direct implementation)
- **bcrypt**: Password hashing

### Frontend
- **React 18+ + TypeScript**: UI framework
- **Vite**: Build tool
- **React Router**: Client-side routing
- **Tailwind CSS**: Utility-first styling
- **Axios**: HTTP client
- **Context API**: Global state (auth only)

## API Design

### Authentication
- `POST /api/users` - User registration
- `POST /api/users/login` - Login
- `GET /api/user` - Current user info
- `PUT /api/user` - Update user info

### Articles
- `GET /api/articles` - List articles (paginated)
- `GET /api/articles/:slug` - Article details
- `POST /api/articles` - Create article (auth required)
- `PUT /api/articles/:slug` - Update article (author only)
- `DELETE /api/articles/:slug` - Delete article (author only)

### Comments
- `GET /api/articles/:slug/comments` - List comments
- `POST /api/articles/:slug/comments` - Create comment (auth required)
- `DELETE /api/articles/:slug/comments/:id` - Delete comment (author only)

## Database Schema

### Core Tables
- **users**: id, username, email, password_hash, bio, image_url
- **articles**: id, slug, title, description, body, author_id, favorites_count
- **comments**: id, body, author_id, article_id
- **tags**: id, name, usage_count (future)
- **favorites**: user_id, article_id (future)

### Indexing Strategy
- articles: author_id, created_at DESC
- comments: article_id

## Authentication & Security

### JWT Token Auth
- Store tokens in `localStorage` (frontend)
- API header: `Authorization: Token jwt.token.here`
- Implement token expiry/refresh logic

### Authorization
- Article/comment edit/delete: author only
- Auth middleware on protected endpoints
- Validate permissions client and server-side

## Development Guidelines

### Code Principles
- **Simplicity first**: Clear, simple code over complex abstractions
- **Explicit naming**: Long, descriptive function names
- **Generate over import**: Prefer code generation to new dependencies

### Error Handling
- Go: Explicit error returns
- React: Error boundaries for component errors
- API: Standard HTTP status codes (400, 401, 403, 404, 500)

### State Management
- **Global state**: Auth state via Context API only
- **Local state**: useState/useReducer for component state
- **Server state**: Consider React Query for future

## Documentation

### Key Docs
- `docs/design.md`: System architecture with Mermaid diagrams
- `docs/tasks.md`: MVP implementation tasks and order
- `docs/PRD.md`: Product requirements

### External References
- [RealWorld Spec](https://realworld-docs.netlify.app/): API and frontend requirements
- [Demo API](https://api.realworld.build/api): Test API

## Communication Rules

- Use Korean for all communication and documentation