# Flux Server

A simple Go HTTP server with Supabase authentication middleware.

## Features

- ✅ HTTP server with health check endpoint
- ✅ Supabase JWT authentication middleware
- ✅ Protected and public routes
- ✅ Environment configuration
- ✅ User context management

## Setup

### 1. Environment Configuration

Copy the example environment file and configure your Supabase credentials:

```bash
cp env.example .env
```

Edit `.env` with your Supabase project details:

```env
# Supabase Configuration
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_ANON_KEY=your-supabase-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key

# Server Configuration
PORT=8080
ENVIRONMENT=development
```

### 2. Get Supabase Credentials

1. Go to your [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Go to Settings → API
4. Copy the following:
   - **Project URL** → `SUPABASE_URL`
   - **anon public** key → `SUPABASE_ANON_KEY`
   - **service_role** key → `SUPABASE_SERVICE_ROLE_KEY`

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run the Server

```bash
# Build the server
go build -o flux-server ./cmd/main.go

# Run the server
./flux-server
```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

## API Endpoints

### Public Endpoints

- **GET /** - Welcome message
- **GET /health** - Health check

### Protected Endpoints (Require Authentication)

- **GET /profile** - User profile information
- **GET /protected** - Protected data

## Authentication

The server uses JWT token validation to authenticate requests. The middleware validates Supabase JWT tokens and extracts user information from the token claims.

### How to Test Protected Endpoints

1. **Get a JWT token from Supabase** (through your frontend or Supabase Auth)
2. **Include the token in the Authorization header:**

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/profile
```

**Note:** The JWT token must be a valid Supabase JWT token. You can get one by:

- Using Supabase Auth in your frontend application
- Using the Supabase CLI or API to generate test tokens
- Creating a user through Supabase Auth and using their session token

### Example Response

```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "user_email": "user@example.com",
  "message": "Welcome to your profile",
  "timestamp": "2025-10-04T19:13:27.783978+01:00"
}
```

## Project Structure

```
flux-server/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── middleware/
│   │   └── auth.go            # Authentication middleware
│   └── transport/
│       └── rest/
│           └── handlers.go    # HTTP handlers and routes
├── env.example                # Environment variables template
├── go.mod                     # Go module file
└── README.md                  # This file
```

## Middleware

### AuthMiddleware

Validates Supabase JWT tokens and adds user information to the request context.

**Usage:**

```go
mux.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(handler)))
```

### OptionalAuthMiddleware

Similar to AuthMiddleware but doesn't require authentication. Adds user info to context if a valid token is provided.

**Usage:**

```go
mux.Handle("/optional", middleware.OptionalAuthMiddleware(http.HandlerFunc(handler)))
```

## Configuration

The application uses environment variables for configuration. All required variables are validated at startup.

**Required:**

- `SUPABASE_URL`
- `SUPABASE_ANON_KEY`

**Optional:**

- `PORT` (default: 8080)
- `ENVIRONMENT` (default: development)
- `SUPABASE_SERVICE_ROLE_KEY`

## Development

### Adding New Protected Routes

1. Create a handler function
2. Wrap it with `middleware.AuthMiddleware`
3. Add it to the router in `handlers.go`

```go
func MyProtectedHandler(w http.ResponseWriter, r *http.Request) {
    userID, _ := middleware.GetUserID(r)
    // Your handler logic here
}

// In NewRouter()
mux.Handle("/my-protected-route", middleware.AuthMiddleware(http.HandlerFunc(MyProtectedHandler)))
```

### Accessing User Information

```go
userID, ok := middleware.GetUserID(r)
userEmail, ok := middleware.GetUserEmail(r)
```

## Error Handling

- **401 Unauthorized**: Invalid or missing JWT token
- **500 Internal Server Error**: Configuration errors or server issues

## Production Security

### Required Environment Variables

For production deployment, ensure these security configurations:

```env
# REQUIRED: JWT Secret from Supabase
JWT_SECRET=your-actual-jwt-secret-from-supabase

# REQUIRED: Use HTTPS in production
ENVIRONMENT=production

# OPTIONAL: Configure CORS origins
CORS_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

### Security Features

- ✅ **JWT Secret Validation**: Uses proper JWT secret instead of anon key
- ✅ **Security Headers**: X-Content-Type-Options, X-Frame-Options, etc.
- ✅ **CORS Protection**: Configurable allowed origins
- ✅ **Rate Limiting**: 60 requests per minute per IP
- ✅ **Input Validation**: Proper request validation
- ✅ **Error Handling**: Secure error messages

### Getting Your JWT Secret

1. Go to your [Supabase Dashboard](https://supabase.com/dashboard)
2. Select your project
3. Go to Settings → API
4. Copy the **JWT Secret** (not the anon key)

### Production Checklist

- [ ] Set `JWT_SECRET` environment variable
- [ ] Use HTTPS in production
- [ ] Configure proper CORS origins
- [ ] Set up proper rate limiting (Redis-based for multiple instances)
- [ ] Use environment-specific database URLs
- [ ] Enable Supabase RLS (Row Level Security)
- [ ] Set up monitoring and logging
- [ ] Use a reverse proxy (nginx/Cloudflare)

## License

MIT
