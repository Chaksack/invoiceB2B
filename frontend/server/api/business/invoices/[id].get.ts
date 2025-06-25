import { defineEventHandler } from 'h3'
import { query } from '../../../db'
import { authorize } from '../../auth'

/**
 * @openapi
 * /api/business/invoices/{id}:
 *   get:
 *     summary: Get Invoice Details
 *     tags: [Business]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *           format: uuid
 *         description: Invoice ID
 *     responses:
 *       200:
 *         description: Invoice details retrieved successfully
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
 *                     invoice:
 *                       type: object
 *                       properties:
 *                         id:
 *                           type: string
 *                           format: uuid
 *                         invoice_number:
 *                           type: string
 *                         amount:
 *                           type: number
 *                         due_date:
 *                           type: string
 *                           format: date
 *                         customer_name:
 *                           type: string
 *                         customer_email:
 *                           type: string
 *                         description:
 *                           type: string
 *                         terms:
 *                           type: string
 *                         status:
 *                           type: string
 *                         created_at:
 *                           type: string
 *                           format: date-time
 *                         updated_at:
 *                           type: string
 *                           format: date-time
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Business access required
 *       404:
 *         description: Invoice not found
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Get invoice ID from URL
    const invoiceId = getRouterParam(event, 'id')
    if (!invoiceId) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Invoice ID is required'
      })
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

    // Get invoice details
    const invoiceResult = await query(
      `SELECT 
        i.*,
        CASE 
          WHEN i.status = 'funded' THEN fi.name
          ELSE NULL 
        END as funding_institution
       FROM invoices i
       LEFT JOIN financial_institutions fi ON i.funding_institution_id = fi.id
       WHERE i.id = $1 AND i.business_id = $2`,
      [invoiceId, businessId]
    )

    if (invoiceResult.rows.length === 0) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Invoice not found'
      })
    }

    const invoice = invoiceResult.rows[0]

    return {
      success: true,
      message: 'Invoice details retrieved successfully',
      data: { invoice },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Invoice details error:', error)
    throw error
  }
}) 