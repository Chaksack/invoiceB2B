import { defineEventHandler } from 'h3'
import { query } from '../db'

/**
 * @openapi
 * /api/health:
 *   get:
 *     summary: Health Check
 *     tags: [System]
 *     responses:
 *       200:
 *         description: Health status
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *                 data:
 *                   type: object
 *                   properties:
 *                     status:
 *                       type: string
 *                       enum: [healthy, unhealthy]
 *                     timestamp:
 *                       type: string
 *                       format: date-time
 *                     database:
 *                       type: object
 *                       properties:
 *                         status:
 *                           type: string
 *                         responseTime:
 *                           type: number
 */
export default defineEventHandler(async (event) => {
  try {
    // Check database connectivity
    await query('SELECT NOW() as current_time')
    
    return {
      success: true,
      message: 'Health check completed successfully',
      data: {
        status: 'healthy',
        timestamp: new Date().toISOString(),
        uptime: process.uptime(),
        environment: process.env.NODE_ENV || 'development',
        version: process.env.APP_VERSION || '1.0.0'
      },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Health check error:', error)
    return {
      success: false,
      message: 'Health check failed',
      data: {
        status: 'unhealthy',
        timestamp: new Date().toISOString(),
        error: error instanceof Error ? error.message : 'Unknown error'
      },
      timestamp: new Date().toISOString()
    }
  }
}) 