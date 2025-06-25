import { defineEventHandler, getHeader, setHeader } from 'h3'
import rateLimit from 'express-rate-limit'
import slowDown from 'express-slow-down'
import helmet from 'helmet'
import cors from 'cors'

// Rate limiting configuration
const rateLimitConfig = {
  windowMs: parseInt(process.env.RATE_LIMIT_WINDOW_MS || '900000'), // 15 minutes
  max: parseInt(process.env.RATE_LIMIT_MAX || '100'), // limit each IP to 100 requests per windowMs
  message: {
    success: false,
    message: 'Too many requests from this IP, please try again later.',
    timestamp: new Date().toISOString()
  },
  standardHeaders: true,
  legacyHeaders: false,
}

// Authentication rate limiting
const authRateLimitConfig = {
  windowMs: parseInt(process.env.AUTH_RATE_LIMIT_WINDOW_MS || '900000'), // 15 minutes
  max: parseInt(process.env.AUTH_RATE_LIMIT_MAX || '5'), // limit each IP to 5 failed auth attempts per windowMs
  message: {
    success: false,
    message: 'Too many failed authentication attempts, please try again later.',
    timestamp: new Date().toISOString()
  },
  standardHeaders: true,
  legacyHeaders: false,
  skipSuccessfulRequests: true, // Only count failed requests
}

// Speed limiting configuration
const speedLimitConfig = {
  windowMs: parseInt(process.env.AUTH_RATE_LIMIT_WINDOW_MS || '900000'), // 15 minutes
  delayAfter: 50, // Allow 50 requests per windowMs without delay
  delayMs: parseInt(process.env.SLOW_DOWN_DELAY_MS || '1000'), // Add 1 second delay per request after delayAfter
  maxDelayMs: 20000, // Maximum delay of 20 seconds
}

// CORS configuration
const corsOptions = {
  origin: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:3000'],
  credentials: process.env.CORS_CREDENTIALS === 'true',
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
  maxAge: 86400 // 24 hours
}

// Create rate limiters
const generalLimiter = rateLimit(rateLimitConfig)
const authLimiter = rateLimit(authRateLimitConfig)
const speedLimiter = slowDown(speedLimitConfig)

// Security headers configuration
const helmetConfig = {
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"],
      scriptSrc: ["'self'"],
      imgSrc: ["'self'", "data:", "https:"],
      connectSrc: ["'self'"],
      fontSrc: ["'self'"],
      objectSrc: ["'none'"],
      mediaSrc: ["'self'"],
      frameSrc: ["'none'"],
    },
  },
  hsts: {
    maxAge: 31536000, // 1 year
    includeSubDomains: true,
    preload: true
  },
  noSniff: true,
  xssFilter: true,
  frameguard: {
    action: 'deny'
  }
}

export default defineEventHandler(async (event) => {
  // Set security headers
  setHeader(event, 'X-Content-Type-Options', 'nosniff')
  setHeader(event, 'X-Frame-Options', 'DENY')
  setHeader(event, 'X-XSS-Protection', '1; mode=block')
  setHeader(event, 'Referrer-Policy', 'strict-origin-when-cross-origin')
  setHeader(event, 'Permissions-Policy', 'geolocation=(), microphone=(), camera=()')

  // CORS headers
  const origin = getHeader(event, 'origin')
  if (origin && corsOptions.origin.includes(origin)) {
    setHeader(event, 'Access-Control-Allow-Origin', origin)
  }
  setHeader(event, 'Access-Control-Allow-Credentials', corsOptions.credentials.toString())
  setHeader(event, 'Access-Control-Allow-Methods', corsOptions.methods.join(', '))
  setHeader(event, 'Access-Control-Allow-Headers', corsOptions.allowedHeaders.join(', '))
  setHeader(event, 'Access-Control-Max-Age', corsOptions.maxAge.toString())

  // Handle preflight requests
  if (event.method === 'OPTIONS') {
    return { status: 'ok' }
  }

  // Apply rate limiting based on endpoint
  const path = event.path || ''
  
  // Apply stricter rate limiting to authentication endpoints
  if (path.startsWith('/api/auth/')) {
    // This would be handled by the auth-specific middleware
    // For now, we'll just log the request
    console.log(`Auth request to: ${path}`)
  }

  // Apply general rate limiting to all API endpoints
  if (path.startsWith('/api/')) {
    // This would be handled by the general rate limiting middleware
    // For now, we'll just log the request
    console.log(`API request to: ${path}`)
  }

  // Continue with the request
  return
}) 