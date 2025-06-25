# Invoice B2B Financing Platform

A comprehensive B2B invoice financing platform built with Nuxt 3, featuring secure authentication, business management, invoice processing, and admin controls.

## 🚀 Features

### ✅ Implemented Features

#### Authentication & Security
- **JWT-based authentication** with refresh tokens
- **Two-factor authentication (2FA)** support
- **Role-based access control** (Business, Admin)
- **Password reset** functionality
- **Security middleware** with rate limiting and CORS
- **Input validation** using Joi schemas
- **Error handling** with custom error classes

#### Business Management
- **Business profile management** (CRUD operations)
- **Business status management** (pending, approved, rejected, suspended)
- **Business analytics** with detailed insights
- **Dashboard with key metrics**

#### Invoice Management
- **Invoice creation** with validation
- **Invoice listing** with filtering and pagination
- **Invoice details** and status tracking
- **Invoice updates** and deletion
- **Invoice status tracking** (pending, funded, overdue)

#### Admin Features
- **Admin dashboard** with platform overview
- **Business approval workflow**
- **Financial institution management**
- **User management** and monitoring
- **Audit logging** for all operations

#### API & Documentation
- **RESTful API** with OpenAPI 3.0 specification
- **Scalar API documentation** integration
- **Comprehensive API endpoints** for all features
- **Health check endpoints** (basic, detailed, readiness, liveness)

#### Database & Infrastructure
- **PostgreSQL database** with comprehensive schema
- **Database migrations** and initialization
- **Connection pooling** and optimization
- **Audit logging** for compliance

#### Development & Testing
- **TypeScript** support throughout
- **ESLint** and **Prettier** configuration
- **Jest** testing framework setup
- **Environment configuration** management

### 🔄 In Progress / Planned Features

#### Advanced Features
- **Real-time notifications** (WebSocket integration)
- **File upload** for invoice documents
- **Email notifications** and alerts
- **Payment processing** integration
- **Advanced reporting** and analytics

#### Security Enhancements
- **API key management** for integrations
- **Advanced rate limiting** strategies
- **Security monitoring** and alerting
- **Compliance reporting** (GDPR, SOX)

#### Business Features
- **Invoice templates** and customization
- **Bulk operations** for invoices
- **Customer management** system
- **Payment scheduling** and automation

## 🛠️ Technology Stack

- **Frontend**: Nuxt 3, Vue 3, TypeScript
- **Backend**: Nuxt 3 Server API, Node.js
- **Database**: PostgreSQL
- **Authentication**: JWT, bcryptjs
- **Validation**: Joi
- **Documentation**: Scalar API Reference
- **Styling**: Tailwind CSS, shadcn/ui
- **Testing**: Jest
- **Linting**: ESLint, Prettier

## 📋 Prerequisites

- Node.js >= 18.0.0
- npm >= 8.0.0
- PostgreSQL >= 12.0
- Docker (optional)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd invoiceB2B/frontend
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp env.example .env
```

Update the `.env` file with your configuration:

```env
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/invoice_b2b

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key
JWT_EXPIRES_IN=15m
JWT_REFRESH_EXPIRES_IN=7d

# Application Configuration
NODE_ENV=development
APP_VERSION=1.0.0
PORT=3000

# Security Configuration
RATE_LIMIT_WINDOW_MS=900000
RATE_LIMIT_MAX=100
AUTH_RATE_LIMIT_WINDOW_MS=900000
AUTH_RATE_LIMIT_MAX=5
SLOW_DOWN_DELAY_MS=1000

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
CORS_CREDENTIALS=true

# Logging Configuration
LOG_LEVEL=info
LOG_SERVICE_URL=

# 2FA Configuration
TOTP_ISSUER=InvoiceB2B
```

### 4. Database Setup

#### Option A: Using Docker (Recommended)

```bash
# Start PostgreSQL container
docker run --name invoice-b2b-db \
  -e POSTGRES_DB=invoice_b2b \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15

# Initialize database
npm run db:init
```

#### Option B: Local PostgreSQL

1. Install PostgreSQL on your system
2. Create a database: `createdb invoice_b2b`
3. Run the initialization script:

```bash
npm run db:init
```

### 5. Start Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

## 📚 API Documentation

### Accessing API Documentation

1. **Scalar Documentation**: Visit `/documentation` for interactive API docs
2. **OpenAPI Spec**: Available at `/api/openapi.json`
3. **Health Checks**: 
   - Basic: `/api/health`
   - Detailed: `/api/health/detailed`
   - Readiness: `/api/health/ready`
   - Liveness: `/api/health/live`

### Key API Endpoints

#### Authentication
- `POST /api/auth/register` - Business registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh JWT token
- `GET /api/auth/profile` - Get user profile
- `POST /api/auth/2fa/verify` - Verify 2FA token

#### Business Management
- `GET /api/business/profile` - Get business profile
- `PUT /api/business/profile` - Update business profile
- `GET /api/business/invoices` - List invoices
- `POST /api/business/invoices` - Create invoice
- `GET /api/business/invoices/{id}` - Get invoice details
- `PUT /api/business/invoices/{id}` - Update invoice
- `DELETE /api/business/invoices/{id}` - Delete invoice
- `GET /api/business/analytics` - Get business analytics

#### Admin Management
- `GET /api/admin/dashboard-summary` - Admin dashboard
- `GET /api/admin/businesses` - List all businesses
- `PUT /api/admin/businesses/{id}/status` - Update business status
- `GET /api/admin/financial-institutions` - List financial institutions

## 🧪 Testing

### Run Tests

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

### Test Structure

- **Unit tests**: Individual function testing
- **Integration tests**: API endpoint testing
- **E2E tests**: Full user workflow testing

## 🏗️ Project Structure

```
frontend/
├── components/          # Vue components
│   ├── ui/             # Reusable UI components
│   ├── User/           # User-specific components
│   └── Nav/            # Navigation components
├── pages/              # Nuxt pages
├── server/             # Server-side code
│   ├── api/            # API routes
│   │   ├── auth/       # Authentication endpoints
│   │   ├── business/   # Business endpoints
│   │   └── admin/      # Admin endpoints
│   ├── middleware/     # Custom middleware
│   └── db.ts           # Database configuration
├── database/           # Database scripts
│   └── init.sql        # Database initialization
├── public/             # Static assets
└── assets/             # Application assets
```

## 🔧 Development

### Code Quality

```bash
# Lint code
npm run lint

# Fix linting issues
npm run lint:fix

# Format code
npm run format
```

### Database Management

```bash
# Initialize database
npm run db:init

# Reset database (development only)
npm run db:reset
```

### Environment Management

- **Development**: Uses `.env` file
- **Production**: Uses environment variables
- **Testing**: Uses `.env.test` file

## 🚀 Deployment

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm run start
```

### Docker Deployment

```bash
# Build Docker image
docker build -t invoice-b2b .

# Run container
docker run -p 3000:3000 invoice-b2b
```

### Environment Variables for Production

Ensure all required environment variables are set in your production environment:

- `DATABASE_URL`
- `JWT_SECRET`
- `JWT_REFRESH_SECRET`
- `NODE_ENV=production`
- `ALLOWED_ORIGINS`

## 📊 Monitoring & Logging

### Health Checks

The application provides multiple health check endpoints:

- `/api/health` - Basic health status
- `/api/health/detailed` - Detailed system information
- `/api/health/ready` - Kubernetes readiness probe
- `/api/health/live` - Kubernetes liveness probe

### Logging

- **Log Levels**: ERROR, WARN, INFO, DEBUG
- **Structured Logging**: JSON format for production
- **Request Logging**: All API requests are logged
- **Error Logging**: Comprehensive error tracking

## 🔒 Security Features

### Authentication & Authorization
- JWT-based authentication
- Role-based access control
- Two-factor authentication
- Password strength requirements
- Session management

### API Security
- Rate limiting
- CORS protection
- Input validation
- SQL injection prevention
- XSS protection

### Data Protection
- Password hashing with bcrypt
- Secure token storage
- Audit logging
- Data encryption at rest

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new features
5. Ensure all tests pass
6. Submit a pull request

### Development Guidelines

- Follow TypeScript best practices
- Write comprehensive tests
- Use conventional commit messages
- Update documentation for new features
- Follow the existing code style

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:

1. Check the documentation
2. Review existing issues
3. Create a new issue with detailed information
4. Contact the development team

## 🔄 Changelog

### Version 1.0.0
- Initial release
- Complete authentication system
- Business and invoice management
- Admin dashboard and controls
- API documentation
- Comprehensive testing suite
