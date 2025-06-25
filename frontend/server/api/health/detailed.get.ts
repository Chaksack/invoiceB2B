import { defineEventHandler } from 'h3'
import { query } from '../../db'

/**
 * @openapi
 * /api/health/detailed:
 *   get:
 *     summary: Detailed Health Check
 *     tags: [System]
 *     responses:
 *       200:
 *         description: Detailed health status
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
 *                       enum: [healthy, degraded, unhealthy]
 *                     timestamp:
 *                       type: string
 *                       format: date-time
 *                     uptime:
 *                       type: number
 *                     environment:
 *                       type: string
 *                     version:
 *                       type: string
 *                     database:
 *                       type: object
 *                       properties:
 *                         status:
 *                           type: string
 *                         responseTime:
 *                           type: number
 *                         currentTime:
 *                           type: string
 *                           format: date-time
 *                     memory:
 *                       type: object
 *                       properties:
 *                         rss:
 *                           type: string
 *                         heapTotal:
 *                           type: string
 *                         heapUsed:
 *                           type: string
 *                         external:
 *                           type: string
 *                     system:
 *                       type: object
 *                       properties:
 *                         nodeVersion:
 *                           type: string
 *                         platform:
 *                           type: string
 *                         arch:
 *                           type: string
 *       503:
 *         description: Service unavailable
 */
export default defineEventHandler(async (event) => {
  const startTime = Date.now()
  
  try {
    // Basic health check
    const health: any = {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      environment: process.env.NODE_ENV || 'development',
      version: process.env.APP_VERSION || '1.0.0'
    }

    // Database health check
    try {
      const dbStartTime = Date.now()
      const dbResult = await query('SELECT NOW() as current_time, version() as version')
      const dbResponseTime = Date.now() - dbStartTime
      
      health.database = {
        status: 'connected',
        responseTime: dbResponseTime,
        currentTime: dbResult.rows[0]?.current_time,
        version: dbResult.rows[0]?.version?.split(' ')[1] // Extract PostgreSQL version
      }
    } catch (error) {
      health.database = {
        status: 'disconnected',
        error: error instanceof Error ? error.message : 'Unknown error'
      }
      health.status = 'degraded'
    }

    // Memory usage
    const memUsage = process.memoryUsage()
    health.memory = {
      rss: Math.round(memUsage.rss / 1024 / 1024) + ' MB',
      heapTotal: Math.round(memUsage.heapTotal / 1024 / 1024) + ' MB',
      heapUsed: Math.round(memUsage.heapUsed / 1024 / 1024) + ' MB',
      external: Math.round(memUsage.external / 1024 / 1024) + ' MB'
    }

    // System information
    health.system = {
      nodeVersion: process.version,
      platform: process.platform,
      arch: process.arch,
      pid: process.pid
    }

    // Response time
    const responseTime = Date.now() - startTime
    health.responseTime = `${responseTime}ms`

    // Set appropriate status code
    const statusCode = health.status === 'healthy' ? 200 : 503
    
    return {
      success: true,
      message: 'Detailed health check completed',
      data: health,
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Detailed health check error:', error)
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