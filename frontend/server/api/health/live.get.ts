import { defineEventHandler } from 'h3'

/**
 * @openapi
 * /api/health/live:
 *   get:
 *     summary: Kubernetes Liveness Probe
 *     tags: [System]
 *     responses:
 *       200:
 *         description: Service is alive
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 status:
 *                   type: string
 *                   example: "alive"
 *                 timestamp:
 *                   type: string
 *                   format: date-time
 *       503:
 *         description: Service is not alive
 */
export default defineEventHandler(async (event) => {
  try {
    // Simple liveness check - just verify the process is running
    // You can add more sophisticated checks here if needed
    
    return {
      status: 'alive',
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Liveness check failed:', error)
    throw createError({
      statusCode: 503,
      statusMessage: 'Service not alive'
    })
  }
}) 