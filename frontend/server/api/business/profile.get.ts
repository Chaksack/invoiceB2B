import { defineEventHandler } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

/**
 * @openapi
 * /api/business/profile:
 *   get:
 *     summary: Get Business Profile
 *     tags: [Business]
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: Business profile retrieved successfully
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
 *                     business:
 *                       type: object
 *                       properties:
 *                         id:
 *                           type: string
 *                         company_name:
 *                           type: string
 *                         status:
 *                           type: string
 *                         industry:
 *                           type: string
 *                         annual_revenue:
 *                           type: number
 *                         employee_count:
 *                           type: integer
 *       401:
 *         description: Unauthorized
 *       404:
 *         description: Business profile not found
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Get business profile
    const result = await query(
      'SELECT * FROM businesses WHERE user_id = $1',
      [user.id]
    )

    if (result.rows.length === 0) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Business profile not found'
      })
    }

    return {
      success: true,
      message: 'Business profile retrieved successfully',
      data: { business: result.rows[0] },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Business profile error:', error)
    throw error
  }
}) 