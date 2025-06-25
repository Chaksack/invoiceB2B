import { defineEventHandler, getQuery } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

/**
 * @openapi
 * /api/business/invoices:
 *   get:
 *     summary: Get Business Invoices
 *     tags: [Business]
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
 *           enum: [pending, submitted, approved, funded, rejected, paid]
 *         description: Filter by invoice status
 *     responses:
 *       200:
 *         description: Invoices retrieved successfully
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
 *                       invoice_number:
 *                         type: string
 *                       customer_name:
 *                         type: string
 *                       amount:
 *                         type: number
 *                       status:
 *                         type: string
 *                       issue_date:
 *                         type: string
 *                         format: date
 *                       due_date:
 *                         type: string
 *                         format: date
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
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Get query parameters
    const queryParams = getQuery(event)
    const page = parseInt(queryParams.page as string) || 1
    const limit = Math.min(parseInt(queryParams.limit as string) || 10, 100)
    const status = queryParams.status as string
    const search = queryParams.search as string
    const offset = (page - 1) * limit

    // Get business ID
    const businessResult = await query(
      'SELECT id FROM businesses WHERE user_id = $1',
      [user.id]
    )

    if (businessResult.rows.length === 0) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Business profile not found'
      })
    }

    const businessId = businessResult.rows[0].id

    // Build query with filters
    let whereClause = 'WHERE business_id = $1 AND deleted_at IS NULL'
    let queryParams2: any[] = [businessId]
    let paramIndex = 2

    if (status) {
      whereClause += ` AND status = $${paramIndex}`
      queryParams2.push(status)
      paramIndex++
    }

    if (search) {
      whereClause += ` AND (invoice_number ILIKE $${paramIndex} OR customer_name ILIKE $${paramIndex})`
      queryParams2.push(`%${search}%`)
      paramIndex++
    }

    // Get total count
    const countResult = await query(
      `SELECT COUNT(*) as total FROM invoices ${whereClause}`,
      queryParams2
    )
    const total = parseInt(countResult.rows[0].total)
    const totalPages = Math.ceil(total / limit)

    // Get invoices
    const invoicesResult = await query(
      `SELECT * FROM invoices ${whereClause} ORDER BY created_at DESC LIMIT $${paramIndex} OFFSET $${paramIndex + 1}`,
      [...queryParams2, limit, offset]
    )

    return {
      success: true,
      message: 'Invoices retrieved successfully',
      data: invoicesResult.rows,
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
    console.error('Invoices error:', error)
    throw error
  }
}) 