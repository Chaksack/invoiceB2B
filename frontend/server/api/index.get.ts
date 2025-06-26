import { defineEventHandler } from 'h3';

export default defineEventHandler(async (event) => {
  return {
    success: true,
    message: 'Profundr API',
    version: process.env.APP_VERSION || '1.0.0',
    environment: process.env.NODE_ENV || 'development',
    timestamp: new Date().toISOString(),
    endpoints: {
      auth: '/api/auth',
      business: '/api/business',
      admin: '/api/admin',
      health: '/api/health',
      docs: '/api/docs'
    }
  }
}) 