# SocialNet Backend

This is the backend server for the SocialNet social media application, built with Go, Gin, and PostgreSQL.

## Project Structure

```
backend/
├── config/        # Application configuration
├── controller/    # HTTP request handlers
├── database/      # Database connection and migrations
├── middleware/    # HTTP middleware
├── model/         # Data models and DTOs
├── util/          # Utility functions
├── main.go        # Application entry point
├── go.mod         # Go module definition
└── .env.example   # Example environment variables
```

## Getting Started

### Prerequisites

- Go (1.18+)
- PostgreSQL
- Docker and Docker Compose (optional)

### Setup

#### Option 1: Local Development

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Create a PostgreSQL database

```bash
cd backend
go mod download
go run main.go
```

#### Option 2: Docker Development

1. Clone the repository
2. Set up environment variables (or use defaults in docker-compose.yml)
3. Run with Docker Compose:

```bash
docker-compose up -d
```

### Environment Variables

| Variable              | Description           | Default            |
|-----------------------|-----------------------|--------------------|
| HOST                  | Server host           | 0.0.0.0            |
| PORT                  | Server port           | 8080               |
| DB_HOST               | Database host         | localhost          |
| DB_PORT               | Database port         | 5432               |
| DB_USER               | Database user         | postgres           |
| DB_PASSWORD           | Database password     | postgres           |
| DB_NAME               | Database name         | socialnet          |
| JWT_SECRET            | JWT secret key        | default_jwt_secret |
| AWS_REGION            | AWS S3 region         | us-east-1          |
| AWS_BUCKET            | AWS S3 bucket name    | socialnet-uploads  |
| AWS_ACCESS_KEY_ID     | AWS access key ID     |                    |
| AWS_SECRET_ACCESS_KEY | AWS secret access key |                    |

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Log in a user

### Users

- `GET /api/v1/users` - Get users list
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/username/:username` - Get user by username
- `PUT /api/v1/users/:id` - Update user (authenticated)
- `GET /api/v1/users/me` - Get current user (authenticated)
- `POST /api/v1/users/follow/:id` - Follow a user (authenticated)
- `DELETE /api/v1/users/follow/:id` - Unfollow a user (authenticated)
- `GET /api/v1/users/followers` - Get followers (authenticated)
- `GET /api/v1/users/following` - Get following (authenticated)

### Posts

- `GET /api/v1/posts` - Get posts
- `GET /api/v1/posts/:id` - Get post by ID
- `POST /api/v1/posts` - Create post (authenticated)
- `PUT /api/v1/posts/:id` - Update post (authenticated)
- `DELETE /api/v1/posts/:id` - Delete post (authenticated)
- `POST /api/v1/posts/:id/like` - Like post (authenticated)
- `DELETE /api/v1/posts/:id/like` - Unlike post (authenticated)
- `GET /api/v1/posts/feed` - Get feed (authenticated)

### Comments

- `GET /api/v1/posts/:id/comments` - Get post comments
- `POST /api/v1/posts/:id/comments` - Create comment (authenticated)
- `PUT /api/v1/posts/comments/:commentId` - Update comment (authenticated)
- `DELETE /api/v1/posts/comments/:commentId` - Delete comment (authenticated)

### Files

- `POST /api/v1/files/upload` - Upload file (authenticated)
