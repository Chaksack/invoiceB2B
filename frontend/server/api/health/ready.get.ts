import { defineEventHandler } from 'h3'
import { query } from '../../db'

/**
 * @openapi
 * /api/health/ready:
 *   get:
 *     summary: Kubernetes Readiness Probe
 *     tags: [System]
 *     responses:
 *       200:
 *         description: Service is ready to receive traffic
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 status:
 *                   type: string
 *                   example: "ready"
 *                 timestamp:
 *                   type: string
 *                   format: date-time
 *       503:
 *         description: Service is not ready
 */
export default defineEventHandler(async (event) => {
  try {
    // Check database connectivity
    await query('SELECT 1')
    
    // Check if application is ready (you can add more checks here)
    const isReady = true // Add your readiness logic here
    
    if (!isReady) {
      throw createError({
        statusCode: 503,
        statusMessage: 'Service not ready'
      })
    }

    return {
      status: 'ready',
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Readiness check failed:', error)
    throw createError({
      statusCode: 503,
      statusMessage: 'Service not ready'
    })
  }
}) 