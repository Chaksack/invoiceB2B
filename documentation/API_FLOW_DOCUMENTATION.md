# Centralized API Flow Documentation

## Overview
This document outlines the centralized and cleaned API architecture that provides smooth user flows and consistent API patterns.

## Centralized API Architecture

### 1. Unified Middleware Stack
All API routes now use a consistent middleware stack applied globally:
- **Request ID Generation**: Every request gets a unique identifier for traceability
- **Rate Limiting**: Prevents API abuse and ensures fair usage
- **Authentication**: JWT-based authentication for protected routes
- **CSRF Protection**: Cross-site request forgery protection for state-changing operations
- **Admin Authorization**: Role-based access control for admin operations

### 2. Standardized Response Format
All API endpoints now return consistent response structures:

```json
{
  "status": "success|error",
  "message": "Human readable message",
  "data": {},
  "error": {},
  "timestamp": "2025-01-09T17:35:00Z",
  "request_id": "uuid-string"
}
```

### 3. Paginated Responses
List endpoints include consistent pagination metadata:

```json
{
  "status": "success",
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

## Centralized User Flows

### 1. Authentication Flow
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/login/2fa/verify
POST /api/v1/auth/refresh-token
POST /api/v1/auth/logout (authenticated)
POST /api/v1/auth/2fa/toggle (authenticated)
```

### 2. User Profile Management Flow
```
GET /api/v1/user/profile (authenticated)
PUT /api/v1/user/profile (authenticated, CSRF protected)

GET /api/v1/user/kyc (authenticated)
POST /api/v1/user/kyc (authenticated, CSRF protected)
```

### 3. Centralized Loan Application Flow
The loan application process is now centralized under `/loan-applications`:

#### Create & Manage Applications
```
POST /api/v1/loan-applications (authenticated, CSRF protected)
POST /api/v1/loan-applications/manual (authenticated, CSRF protected)
GET /api/v1/loan-applications (authenticated) - List user's applications
GET /api/v1/loan-applications/:id (authenticated) - Get specific application
GET /api/v1/loan-applications/stats (authenticated) - User's loan statistics
```

#### KYB (Know Your Business) Submission
```
POST /api/v1/loan-applications/:id/kyb (authenticated, CSRF protected)
```

#### Financial Statement Management
```
POST /api/v1/loan-applications/:id/financial-statements (authenticated, CSRF protected)
POST /api/v1/loan-applications/:id/financial-statements/upload (authenticated, CSRF protected)
```

#### Document Upload
```
POST /api/v1/loan-applications/:id/documents/upload (authenticated, CSRF protected)
```

### 4. Invoice Management Flow
```
GET /api/v1/invoices (authenticated)
GET /api/v1/invoices/:id (authenticated)
POST /api/v1/invoices (authenticated, admin, CSRF protected)
PUT /api/v1/invoices/:id (authenticated, admin, CSRF protected)
DELETE /api/v1/invoices/:id (authenticated, admin, CSRF protected)
POST /api/v1/invoices/:id/upload (authenticated, admin, CSRF protected)
```

### 5. Admin Operations Flow
All admin operations are centralized under `/admin` with proper role-based access:

#### User Management
```
GET /api/v1/admin/users
GET /api/v1/admin/users/:id
GET /api/v1/admin/users/:id/kyc
PUT /api/v1/admin/users/:id/kyc/review
```

#### Loan Application Management (Admin)
```
GET /api/v1/admin/loan-applications
GET /api/v1/admin/loan-applications/stats
POST /api/v1/admin/loan-applications/:id/review
POST /api/v1/admin/loan-applications/:id/send-to-financial-institution
```

## Key Improvements Made

### 1. Centralized User-Facing Loan APIs
- **Before**: Loan functionality was only available through admin routes
- **After**: Complete user-facing loan application flow under `/loan-applications`
- **Benefit**: Users can now manage their loan applications directly

### 2. Consistent Middleware Application
- **Before**: Inconsistent middleware usage across different route groups
- **After**: Unified middleware stack applied globally with proper layering
- **Benefit**: Consistent security, rate limiting, and request tracking

### 3. Standardized Response Format
- **Before**: Different response formats across endpoints
- **After**: Consistent JSON structure with metadata and request tracking
- **Benefit**: Easier API consumption and debugging

### 4. Proper Service Architecture
- **Before**: Missing repository layer and validation services
- **After**: Complete service layer with proper dependency injection
- **Benefit**: Better separation of concerns and testability

### 5. Enhanced Security
- **Before**: Basic security measures
- **After**: Comprehensive security with CSRF, rate limiting, and request tracking
- **Benefit**: Better protection against common attacks

## API Usage Examples

### Creating a Loan Application
```bash
curl -X POST /api/v1/loan-applications \
  -H "Authorization: Bearer <jwt-token>" \
  -H "X-CSRF-Token: <csrf-token>" \
  -d '{
    "source": "manual",
    "requested_amount": 25000.00,
    "currency": "USD",
    "purpose": "Business expansion"
  }'
```

### Getting User's Loan Applications
```bash
curl -X GET /api/v1/loan-applications?page=1&pageSize=10 \
  -H "Authorization: Bearer <jwt-token>"
```

### Submitting KYB Information
```bash
curl -X POST /api/v1/loan-applications/123/kyb \
  -H "Authorization: Bearer <jwt-token>" \
  -H "X-CSRF-Token: <csrf-token>" \
  -d '{
    "business_name": "My Business Ltd",
    "business_registration_no": "REG123456",
    "business_type": "Limited Company",
    "industry_type": "Technology",
    "business_address": "123 Business Street",
    "tax_identification_number": "TAX123456",
    "years_in_operation": 5,
    "business_description": "Technology consulting business"
  }'
```

## Error Handling

All endpoints now return consistent error responses:

### Validation Error (422)
```json
{
  "status": "error",
  "message": "Validation failed",
  "error": {
    "field_name": ["error message"]
  },
  "timestamp": "2025-01-09T17:35:00Z",
  "request_id": "uuid-string"
}
```

### Authentication Error (401)
```json
{
  "status": "error",
  "message": "Authentication required",
  "timestamp": "2025-01-09T17:35:00Z",
  "request_id": "uuid-string"
}
```

### Not Found Error (404)
```json
{
  "status": "error",
  "message": "Resource not found",
  "timestamp": "2025-01-09T17:35:00Z",
  "request_id": "uuid-string"
}
```

## Benefits of Centralized API Architecture

1. **Improved Developer Experience**: Consistent patterns and responses
2. **Better Error Tracking**: Request IDs for debugging and monitoring
3. **Enhanced Security**: Unified security measures across all endpoints
4. **Streamlined User Flows**: Logical grouping of related functionality
5. **Easier Maintenance**: Centralized middleware and response handling
6. **Better Performance**: Proper rate limiting and caching strategies
7. **Scalability**: Clean architecture supports easy feature additions

This centralized approach ensures that the API is maintainable, secure, and provides a smooth experience for both developers and end users.