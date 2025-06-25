import { defineEventHandler } from 'h3'
import { query } from '../../../db'
import { authorize } from '../../auth'

/**
 * @openapi
 * /api/business/invoices/{id}:
 *   delete:
 *     summary: Delete Invoice
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
 *         description: Invoice deleted successfully
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 success:
 *                   type: boolean
 *                 message:
 *                   type: string
 *       400:
 *         description: Cannot delete funded invoice
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

    // Check if invoice exists and belongs to this business
    const existingInvoice = await query(
      'SELECT id, status, invoice_number FROM invoices WHERE id = $1 AND business_id = $2',
      [invoiceId, businessId]
    )

    if (existingInvoice.rows.length === 0) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Invoice not found'
      })
    }

    const invoice = existingInvoice.rows[0]

    // Check if invoice can be deleted (not funded)
    if (invoice.status === 'funded') {
      throw createError({
        statusCode: 400,
        statusMessage: 'Cannot delete funded invoice'
      })
    }

    // Soft delete the invoice (set deleted_at timestamp)
    await query(
      'UPDATE invoices SET deleted_at = $1 WHERE id = $2',
      [new Date(), invoiceId]
    )

    // Log the invoice deletion
    await query(
      `INSERT INTO audit_logs (user_id, action, table_name, record_id, old_values)
       VALUES ($1, 'DELETE', 'invoices', $2, $3)`,
      [
        user.id,
        invoiceId,
        JSON.stringify({
          invoiceNumber: invoice.invoice_number,
          status: invoice.status
        })
      ]
    )

    return {
      success: true,
      message: 'Invoice deleted successfully',
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Invoice deletion error:', error)
    throw error
  }
}) 