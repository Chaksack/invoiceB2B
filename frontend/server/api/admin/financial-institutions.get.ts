import { defineEventHandler, getQuery } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

/**
 * @openapi
 * /api/admin/financial-institutions:
 *   get:
 *     summary: List Financial Institutions
 *     tags: [Admin]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *           minimum: 1
 *           default: 1
 *         description: Page number
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *           minimum: 1
 *           maximum: 100
 *           default: 10
 *         description: Number of items per page
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *           enum: [active, inactive]
 *         description: Filter by status
 *       - in: query
 *         name: search
 *         schema:
 *           type: string
 *         description: Search by institution name
 *     responses:
 *       200:
 *         description: Financial institutions retrieved successfully
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
 *                   type: array
 *                   items:
 *                     type: object
 *                     properties:
 *                       id:
 *                         type: string
 *                         format: uuid
 *                       name:
 *                         type: string
 *                       type:
 *                         type: string
 *                       is_active:
 *                         type: boolean
 *                       funding_capacity:
 *                         type: number
 *                       interest_rate_range:
 *                         type: string
 *                       created_at:
 *                         type: string
 *                         format: date-time
 *                 pagination:
 *                   type: object
 *                   properties:
 *                     page:
 *                       type: integer
 *                     limit:
 *                       type: integer
 *                     total:
 *                       type: integer
 *                     totalPages:
 *                       type: integer
 *                     hasNext:
 *                       type: boolean
 *                     hasPrev:
 *                       type: boolean
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Admin access required
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify admin authentication
    const user = authorize('admin')(event)

    // Get query parameters
    const queryParams = getQuery(event)
    const page = parseInt(queryParams.page as string) || 1
    const limit = Math.min(parseInt(queryParams.limit as string) || 10, 100)
    const status = queryParams.status as string
    const search = queryParams.search as string
    const offset = (page - 1) * limit

    // Build query with filters
    let whereClause = 'WHERE 1=1'
    let queryParams2: any[] = []
    let paramIndex = 1

    if (status) {
      whereClause += ` AND fi.is_active = $${paramIndex}`
      queryParams2.push(status === 'active')
      paramIndex++
    }

    if (search) {
      whereClause += ` AND fi.name ILIKE $${paramIndex}`
      queryParams2.push(`%${search}%`)
      paramIndex++
    }

    // Get total count
    const countResult = await query(
      `SELECT COUNT(*) as total FROM financial_institutions fi ${whereClause}`,
      queryParams2
    )
    const total = parseInt(countResult.rows[0].total)
    const totalPages = Math.ceil(total / limit)

    // Get financial institutions with statistics
    const institutionsResult = await query(
      `SELECT 
        fi.*,
        COUNT(DISTINCT i.id) as total_invoices_funded,
        COALESCE(SUM(i.amount), 0) as total_amount_funded
       FROM financial_institutions fi
       LEFT JOIN invoices i ON fi.id = i.funding_institution_id AND i.status = 'funded'
       ${whereClause}
       GROUP BY fi.id
       ORDER BY fi.created_at DESC
       LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
      [...queryParams2, limit, offset]
    )

    return {
      success: true,
      message: 'Financial institutions retrieved successfully',
      data: institutionsResult.rows,
      pagination: {
        page,
        limit,
        total,
        totalPages,
        hasNext: page < totalPages,
        hasPrev: page > 1
      },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Admin financial institutions error:', error)
    throw error
  }
}) 