import { defineEventHandler, getQuery } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

/**
 * @openapi
 * /api/admin/businesses:
 *   get:
 *     summary: List All Businesses
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
 *           enum: [pending, approved, rejected, suspended]
 *         description: Filter by business status
 *       - in: query
 *         name: search
 *         schema:
 *           type: string
 *         description: Search by company name
 *     responses:
 *       200:
 *         description: Businesses retrieved successfully
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
 *                       company_name:
 *                         type: string
 *                       status:
 *                         type: string
 *                       industry:
 *                         type: string
 *                       annual_revenue:
 *                         type: number
 *                       employee_count:
 *                         type: integer
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
      whereClause += ` AND b.status = $${paramIndex}`
      queryParams2.push(status)
      paramIndex++
    }

    if (search) {
      whereClause += ` AND b.company_name ILIKE $${paramIndex}`
      queryParams2.push(`%${search}%`)
      paramIndex++
    }

    // Get total count
    const countResult = await query(
      `SELECT COUNT(*) as total FROM businesses b ${whereClause}`,
      queryParams2
    )
    const total = parseInt(countResult.rows[0].total)
    const totalPages = Math.ceil(total / limit)

    // Get businesses with user info
    const businessesResult = await query(
      `SELECT 
        b.*,
        u.email,
        COUNT(i.id) as total_invoices,
        COALESCE(SUM(CASE WHEN i.status = 'funded' THEN i.amount ELSE 0 END), 0) as total_funded
       FROM businesses b
       LEFT JOIN users u ON b.user_id = u.id
       LEFT JOIN invoices i ON b.id = i.business_id
       ${whereClause}
       GROUP BY b.id, u.email
       ORDER BY b.created_at DESC
       LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
      [...queryParams2, limit, offset]
    )

    return {
      success: true,
      message: 'Businesses retrieved successfully',
      data: businessesResult.rows,
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
    console.error('Admin businesses error:', error)
    throw error
  }
}) 