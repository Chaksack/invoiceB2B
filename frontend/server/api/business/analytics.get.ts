import { defineEventHandler, getQuery } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

/**
 * @openapi
 * /api/business/analytics:
 *   get:
 *     summary: Get Business Analytics
 *     tags: [Business]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: period
 *         schema:
 *           type: string
 *           enum: [7d, 30d, 90d, 1y]
 *           default: 30d
 *         description: Analytics period
 *     responses:
 *       200:
 *         description: Analytics retrieved successfully
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
 *                     overview:
 *                       type: object
 *                       properties:
 *                         total_invoices:
 *                           type: integer
 *                         total_amount:
 *                           type: number
 *                         funded_invoices:
 *                           type: integer
 *                         funded_amount:
 *                           type: number
 *                         pending_invoices:
 *                           type: integer
 *                         pending_amount:
 *                           type: number
 *                         overdue_invoices:
 *                           type: integer
 *                         overdue_amount:
 *                           type: number
 *                     trends:
 *                       type: object
 *                       properties:
 *                         monthly_invoices:
 *                           type: array
 *                           items:
 *                             type: object
 *                         monthly_amounts:
 *                           type: array
 *                           items:
 *                             type: object
 *                     top_customers:
 *                       type: array
 *                       items:
 *                         type: object
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Business access required
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Get query parameters
    const queryParams = getQuery(event)
    const period = queryParams.period as string || '30d'

    // Calculate date range based on period
    let dateRange: string
    switch (period) {
      case '7d':
        dateRange = "NOW() - INTERVAL '7 days'"
        break
      case '90d':
        dateRange = "NOW() - INTERVAL '90 days'"
        break
      case '1y':
        dateRange = "NOW() - INTERVAL '1 year'"
        break
      default:
        dateRange = "NOW() - INTERVAL '30 days'"
    }

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

    // Get overview statistics
    const overviewResult = await query(
      `SELECT 
        COUNT(*) as total_invoices,
        COALESCE(SUM(amount), 0) as total_amount,
        COUNT(CASE WHEN status = 'funded' THEN 1 END) as funded_invoices,
        COALESCE(SUM(CASE WHEN status = 'funded' THEN amount ELSE 0 END), 0) as funded_amount,
        COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_invoices,
        COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as pending_amount,
        COUNT(CASE WHEN status = 'pending' AND due_date < NOW() THEN 1 END) as overdue_invoices,
        COALESCE(SUM(CASE WHEN status = 'pending' AND due_date < NOW() THEN amount ELSE 0 END), 0) as overdue_amount
       FROM invoices 
       WHERE business_id = $1 AND deleted_at IS NULL`,
      [businessId]
    )

    // Get monthly trends
    const monthlyTrendsResult = await query(
      `SELECT 
        DATE_TRUNC('month', created_at) as month,
        COUNT(*) as invoice_count,
        COALESCE(SUM(amount), 0) as total_amount,
        COUNT(CASE WHEN status = 'funded' THEN 1 END) as funded_count,
        COALESCE(SUM(CASE WHEN status = 'funded' THEN amount ELSE 0 END), 0) as funded_amount
       FROM invoices 
       WHERE business_id = $1 AND deleted_at IS NULL AND created_at >= ${dateRange}
       GROUP BY DATE_TRUNC('month', created_at)
       ORDER BY month DESC
       LIMIT 12`,
      [businessId]
    )

    // Get top customers
    const topCustomersResult = await query(
      `SELECT 
        customer_name,
        customer_email,
        COUNT(*) as invoice_count,
        COALESCE(SUM(amount), 0) as total_amount,
        COUNT(CASE WHEN status = 'funded' THEN 1 END) as funded_count,
        COALESCE(SUM(CASE WHEN status = 'funded' THEN amount ELSE 0 END), 0) as funded_amount
       FROM invoices 
       WHERE business_id = $1 AND deleted_at IS NULL
       GROUP BY customer_name, customer_email
       ORDER BY total_amount DESC
       LIMIT 10`,
      [businessId]
    )

    // Get recent activity
    const recentActivityResult = await query(
      `SELECT 
        'invoice' as type,
        invoice_number as identifier,
        status,
        amount,
        created_at,
        updated_at
       FROM invoices 
       WHERE business_id = $1 AND deleted_at IS NULL
       ORDER BY updated_at DESC
       LIMIT 20`,
      [businessId]
    )

    const analytics = {
      overview: overviewResult.rows[0],
      trends: {
        monthly: monthlyTrendsResult.rows
      },
      top_customers: topCustomersResult.rows,
      recent_activity: recentActivityResult.rows
    }

    return {
      success: true,
      message: 'Analytics retrieved successfully',
      data: analytics,
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Business analytics error:', error)
    throw error
  }
}) 