# Flux Server API Documentation

## Base URL

```
http://localhost:8080
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## Endpoints

### Health Check

```http
GET /health
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2025-10-04T20:34:11.000Z",
  "service": "flux-server"
}
```

### Root

```http
GET /
```

**Response:**

```json
{
  "message": "Welcome to the Flux Server",
  "version": "1.0.0"
}
```

### Signup

```http
POST /auth/signup
```

**Request Body:**

```json
{
  "username": "farmer123",
  "password": "mypassword123",
  "role": "farmer", // Optional: "farmer" | "extension_officer", defaults to "farmer"
  "phone_number": "+1234567890", // Optional
  "crop_type": "corn", // Optional (farmer only)
  "location_id": 1, // Optional (farmer only)
  "language": "en", // Optional (farmer only), defaults to "en"
  "assigned_location_id": 2 // Optional (extension_officer only)
}
```

**Response:**

```json
{
  "user": {
    "id": "uuid",
    "username": "farmer123",
    "created_at": "2025-10-04T19:34:28.886822Z"
  },
  "access_token": "",
  "token_type": "",
  "expires_in": 0,
  "message": "User created successfully. Profile will be created automatically. Please sign in to get your access token."
}
```

### Signin

```http
POST /auth/signin
```

**Request Body:**

```json
{
  "username": "farmer123",
  "password": "mypassword123"
}
```

**Response:**

```json
{
  "user": {
    "id": "uuid",
    "username": "farmer123",
    "created_at": "2025-10-04T19:34:28.886822Z"
  },
  "access_token": "jwt_token_here",
  "token_type": "bearer",
  "expires_in": 3600,
  "message": "Sign in successful"
}
```

### Profile (Protected)

```http
GET /profile
```

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "user_id": "uuid",
  "user_email": "farmer123",
  "message": "Welcome to your profile",
  "timestamp": "2025-10-04T20:34:11.000Z",
  "profile": {
    "id": "uuid",
    "auth_user_id": "uuid",
    "role_id": "uuid",
    "display_name": "farmer123",
    "phone": null,
    "metadata": {},
    "created_at": "2025-10-04T20:34:11.000Z"
  }
}
```

### Protected Data (Protected)

```http
GET /protected
```

**Headers:** `Authorization: Bearer <token>`

**Response:**

```json
{
  "data": "This is protected data",
  "user_id": "uuid",
  "timestamp": "2025-10-04T20:34:11.000Z"
}
```

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "message": "Error message",
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": "Additional details"
  },
  "timestamp": "2025-10-04T20:34:11.000Z"
}
```

## User Roles

- **farmer**: Can access farmer-specific features
- **extension_officer**: Can access extension officer features

## Database Tables Created

When a user signs up, the system automatically creates:

1. **auth.users** record (via Supabase)
2. **user_profiles** record with role assignment
3. **farmers** or **extension_officers** record (based on role)

## Environment Variables

```env
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_anon_key
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
JWT_SECRET=your_jwt_secret
PORT=8080
ENVIRONMENT=development
```
