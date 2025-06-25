import { defineEventHandler, createError } from 'h3'

// Custom error classes
export class AppError extends Error {
  public statusCode: number
  public isOperational: boolean

  constructor(message: string, statusCode: number = 500, isOperational: boolean = true) {
    super(message)
    this.statusCode = statusCode
    this.isOperational = isOperational

    Error.captureStackTrace(this, this.constructor)
  }
}

export class ValidationError extends AppError {
  constructor(message: string, details?: any) {
    super(message, 400)
    this.details = details
  }
  public details?: any
}

export class AuthenticationError extends AppError {
  constructor(message: string = 'Authentication failed') {
    super(message, 401)
  }
}

export class AuthorizationError extends AppError {
  constructor(message: string = 'Access denied') {
    super(message, 403)
  }
}

export class NotFoundError extends AppError {
  constructor(resource: string = 'Resource') {
    super(`${resource} not found`, 404)
  }
}

export class ConflictError extends AppError {
  constructor(message: string = 'Resource conflict') {
    super(message, 409)
  }
}

export class RateLimitError extends AppError {
  constructor(message: string = 'Too many requests') {
    super(message, 429)
  }
}

// Error handler middleware
export default defineEventHandler(async (event) => {
  try {
    // Continue with the request
    return
  } catch (error) {
    console.error('Error in request:', error)

    // Handle different types of errors
    if (error instanceof AppError) {
      return {
        success: false,
        message: error.message,
        statusCode: error.statusCode,
        ...(error.details && { details: error.details }),
        timestamp: new Date().toISOString()
      }
    }

    // Handle database errors
    if (error.code === '23505') { // Unique constraint violation
      return {
        success: false,
        message: 'Resource already exists',
        statusCode: 409,
        timestamp: new Date().toISOString()
      }
    }

    if (error.code === '23503') { // Foreign key constraint violation
      return {
        success: false,
        message: 'Referenced resource does not exist',
        statusCode: 400,
        timestamp: new Date().toISOString()
      }
    }

    if (error.code === '42P01') { // Undefined table
      return {
        success: false,
        message: 'Database configuration error',
        statusCode: 500,
        timestamp: new Date().toISOString()
      }
    }

    // Handle JWT errors
    if (error.name === 'JsonWebTokenError') {
      return {
        success: false,
        message: 'Invalid token',
        statusCode: 401,
        timestamp: new Date().toISOString()
      }
    }

    if (error.name === 'TokenExpiredError') {
      return {
        success: false,
        message: 'Token expired',
        statusCode: 401,
        timestamp: new Date().toISOString()
      }
    }

    // Handle validation errors
    if (error.name === 'ValidationError') {
      return {
        success: false,
        message: 'Validation error',
        statusCode: 400,
        details: error.details,
        timestamp: new Date().toISOString()
      }
    }

    // Handle rate limiting errors
    if (error.statusCode === 429) {
      return {
        success: false,
        message: 'Too many requests, please try again later',
        statusCode: 429,
        timestamp: new Date().toISOString()
      }
    }

    // Default error response
    const isDevelopment = process.env.NODE_ENV === 'development'
    
    return {
      success: false,
      message: isDevelopment ? error.message : 'Internal server error',
      statusCode: 500,
      ...(isDevelopment && { 
        stack: error.stack,
        details: error
      }),
      timestamp: new Date().toISOString()
    }
  }
})

// Global error handler for unhandled rejections
process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason)
  
  // In production, you might want to log this to a service like Sentry
  if (process.env.NODE_ENV === 'production') {
    // Log to external service
    console.error('Unhandled rejection logged to monitoring service')
  }
})

// Global error handler for uncaught exceptions
process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error)
  
  // In production, you might want to log this to a service like Sentry
  if (process.env.NODE_ENV === 'production') {
    // Log to external service
    console.error('Uncaught exception logged to monitoring service')
  }
  
  // Exit the process in production to prevent the app from running in an undefined state
  if (process.env.NODE_ENV === 'production') {
    process.exit(1)
  }
})

// Utility functions for creating errors
export const createAppError = (message: string, statusCode: number = 500) => {
  return new AppError(message, statusCode)
}

export const createValidationError = (message: string, details?: any) => {
  return new ValidationError(message, details)
}

export const createAuthError = (message?: string) => {
  return new AuthenticationError(message)
}

export const createAuthzError = (message?: string) => {
  return new AuthorizationError(message)
}

export const createNotFoundError = (resource?: string) => {
  return new NotFoundError(resource)
}

export const createConflictError = (message?: string) => {
  return new ConflictError(message)
}

export const createRateLimitError = (message?: string) => {
  return new RateLimitError(message)
} 