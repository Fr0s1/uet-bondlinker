
# SocialNet

A full-stack social networking application with React frontend and Go backend.

## Project Structure

```
socialnet/
├── backend/       # Go backend API
│   ├── config/    # Configuration
│   ├── controller/ # Request handlers
│   ├── database/  # Database connection
│   ├── middleware/ # HTTP middleware
│   ├── model/     # Data models
│   ├── util/      # Utility functions
│   └── main.go    # Entry point
├── src/           # React frontend
│   ├── components/ # UI components
│   ├── contexts/   # React context providers
│   ├── hooks/      # Custom React hooks
│   ├── lib/        # Utility functions
│   ├── pages/      # Application pages
│   └── App.tsx     # Main component
└── docker-compose.yml # Docker configuration
```

## Getting Started

### Prerequisites

- Node.js & npm or yarn (for frontend)
- Go 1.18+ (for backend)
- PostgreSQL (database)
- Docker & Docker Compose (optional, for containerized setup)

### Development Setup

#### Option 1: Local Development

**Backend Setup:**

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy the environment example file and configure it:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Run the backend server:
   ```bash
   go run main.go
   ```

**Frontend Setup:**

1. Install npm dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Access the application at http://localhost:5173

#### Option 2: Docker Setup

1. Make sure Docker and Docker Compose are installed
2. Set necessary environment variables or use defaults in docker-compose.yml
3. Build and start the containers:
   ```bash
   docker-compose up -d
   ```
4. Access the frontend at http://localhost:5173
5. The API will be available at http://localhost:8080

### Environment Variables

**Backend:**

| Variable | Description | Default |
|----------|-------------|---------|
| HOST | Server host | 0.0.0.0 |
| PORT | Server port | 8080 |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | socialnet |
| JWT_SECRET | JWT secret key | default_jwt_secret |
| AWS_REGION | AWS S3 region | us-east-1 |
| AWS_BUCKET | AWS S3 bucket name | socialnet-uploads |

**Frontend:**

| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_URL | Backend API URL | http://localhost:8080/api/v1 |

## Features

- User authentication (register, login)
- Profile creation and management
- Social connections (follow/unfollow)
- Content creation and sharing
- Interactions (likes, comments)
- Real-time messaging
- Notifications
- File uploads

## Technologies

### Backend
- Go
- Gin (web framework)
- GORM (ORM)
- PostgreSQL (database)
- JWT (authentication)
- AWS S3 (file storage)

### Frontend
- React
- TypeScript
- Vite (build tool)
- React Router
- React Query
- Tailwind CSS
- Shadcn UI

## Deployment

This project can be deployed using Docker containers. The docker-compose.yml file is configured for production use.

For production deployment:

```bash
docker-compose -f docker-compose.yml up -d
```

## License

[MIT License](LICENSE)
