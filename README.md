# Chirpy

A Twitter-like social media API built with Go as part of the Boot.dev guided project.

## Overview

Chirpy is a RESTful API server that allows users to create accounts, post short messages called "chirps" (max 140 characters), and manage their content. The project demonstrates building a production-ready web server with authentication, database operations, and webhooks.

## Features

- **User Management**: Create accounts with email and password (hashed with Argon2id)
- **Authentication**: JWT-based authentication with access and refresh tokens
- **Chirps**: Post, retrieve, and delete short messages
- **Filtering & Sorting**: Query chirps by author and sort in ascending/descending order
- **Webhooks**: Integration with external payment service (Polka) for premium user upgrades
- **Metrics**: Built-in admin metrics tracking for monitoring server usage

## Tech Stack

- **Language**: Go 1.25.4
- **Database**: PostgreSQL15 with `sqlc` for type-safe SQL queries
- **Migrations**: Goose for database schema management
- **Authentication**: JWT tokens with `golang-jwt/jwt/v5`
- **Password Hashing**: Argon2id via `alexedwards/argon2id`

## Project Structure

```text
chirpy/
├── handlers/          # HTTP request handlers
├── internal/          # Internal packages (database access)
├── sql/
│   ├── queries/      # SQL query definitions
│   └── schema/       # Database migrations
├── assets/           # Static files
├── main.go           # Application entry point
└── sqlc.yaml         # sqlc configuration
```

## Setup

1. **Prerequisites**
   - Go 1.25+
   - PostgreSQL15 database
   - sqlc (for generating Go code from SQL)
   - Goose (for database migrations)

2. **Environment Variables**

   Create a `.env` file based on `.env.sample`:

   ```
   DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=your-secret-key
   POLKA_KEY=your-polka-api-key
   ```

3. **Database Setup**

   Run migrations:
   ```bash
   goose -dir sql/schema postgres "your-db-url" up
   ```

4. **Run the Server**

   ```bash
   go build -o chirpy && ./chirpy
   ```

   The server will start on `http://localhost:8080`

## API Endpoints

### Health & Metrics
- `GET /api/healthz` - Health check
- `GET /admin/metrics` - View request metrics
- `POST /admin/reset` - Reset metrics (dev only)

### Users
- `POST /api/users` - Create new user
- `POST /api/login` - Login and receive JWT tokens
- `PUT /api/users` - Update user information (requires auth)

### Authentication
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Chirps
- `POST /api/chirps` - Create a new chirp (requires auth)
- `GET /api/chirps` - Get all chirps (optional `?author_id=uuid` and `?sort=asc|desc`)
- `GET /api/chirps/{chirpID}` - Get specific chirp by ID
- `DELETE /api/chirps/{chirpID}` - Delete chirp (requires auth, owner only)

### Webhooks
- `POST /api/polka/webhooks` - Polka payment webhook for premium upgrades

## Database Schema

### Users Table
- `id` (UUID, primary key)
- `email` (unique)
- `hashed_password`
- `is_chirpy_red` (premium status)
- `created_at`, `updated_at`

### Chirps Table
- `id` (UUID, primary key)
- `body` (VARCHAR 140)
- `user_id` (foreign key to users)
- `created_at`, `updated_at`

### Refresh Tokens Table
- `token` (primary key)
- `user_id` (foreign key to users)
- `expires_at`
- `revoked_at`
- `created_at`, `updated_at`

## Development

This project is part of the [Boot.dev](https://boot.dev) Learn Web Servers course, which teaches building production-grade APIs in Go from scratch.

## License

Educational project - part of Boot.dev curriculum
