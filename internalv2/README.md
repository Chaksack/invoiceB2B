# Internal API v2

This is the internal API v2 service for the InvoiceB2B application, built with GoFiber.

## Features

- **Authentication**: JWT-based authentication with role-based access control
- **Business Management**: Business profile management and invoice handling
- **Admin Dashboard**: Admin-only endpoints for managing businesses and viewing analytics
- **Database Integration**: PostgreSQL database with connection pooling
- **Error Handling**: Comprehensive error handling and logging
- **CORS Support**: Cross-origin resource sharing enabled

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `GET /api/auth/profile` - Get user profile (requires authentication)

### Business (requires business role)
- `GET /api/business/invoices` - Get business invoices with pagination and filtering
- `GET /api/business/profile` - Get business profile
- `PUT /api/business/profile` - Update business profile

### Admin (requires admin role)
- `GET /api/admin/businesses` - Get all businesses with pagination and filtering
- `GET /api/admin/dashboard-summary` - Get admin dashboard summary
- `PUT /api/admin/businesses/:id/status` - Update business status

### Health Check
- `GET /health` - Health check endpoint
- `GET /` - API information endpoint

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Application
APP_ENV=development
INTERNAL_API_PORT=3001

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=invoice_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret_key_please_change_this
JWT_ACCESS_TOKEN_EXPIRATION_MINUTES=15
JWT_REFRESH_TOKEN_EXPIRATION_DAYS=7

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EVENT_EXCHANGE_NAME=invoice_events_exchange
RABBITMQ_USER_REGISTERED_ROUTING_KEY=user.registered
RABBITMQ_INVOICE_UPLOADED_ROUTING_KEY=invoice.uploaded
RABBITMQ_INVOICE_STATUS_UPDATED_ROUTING_KEY=invoice.status.updated
RABBITMQ_KYC_STATUS_UPDATED_ROUTING_KEY=kyc.status.updated

# SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTPSenderEmail=Invoice App <no-reply@syentia.io>

# Other
OTP_EXPIRATION_MINUTES=5
UPLOADS_DIR=./uploads
MAX_UPLOAD_SIZE_MB=10
INTERNAL_API_KEY=your_internal_api_key_please_change_this
```

## Running the Service

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run the service**:
   ```bash
   go run main.go
   ```

The service will start on port 3001 (or the port specified in `INTERNAL_API_PORT`).

## Database Schema

The service expects the following database tables:

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'business',
    is_approved BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Businesses Table
```sql
CREATE TABLE businesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    industry VARCHAR(255) NOT NULL,
    annual_revenue DECIMAL(15,2) NOT NULL,
    employee_count INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Invoices Table
```sql
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID REFERENCES businesses(id) ON DELETE CASCADE,
    invoice_number VARCHAR(255) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    description TEXT,
    file_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

## API Response Format

All API responses follow this format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

For paginated responses:
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "totalPages": 10,
    "hasNext": true,
    "hasPrev": false
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Error Handling

Errors are returned in the following format:

```json
{
  "success": false,
  "message": "Error description",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Authentication

The service uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Development

### Project Structure
```
internal@v2/
├── main.go              # Application entry point
├── config/              # Configuration management
├── database/            # Database connection and utilities
├── handlers/            # HTTP request handlers
├── middleware/          # HTTP middleware
├── models/              # Data models and DTOs
├── routes/              # Route definitions
└── README.md           # This file
```

### Adding New Endpoints

1. Create a new handler in the `handlers/` directory
2. Add the route in `routes/routes.go`
3. Update the models if needed
4. Test the endpoint

### Testing

To run tests:
```bash
go test ./...
```

## License

This project is part of the InvoiceB2B application. 