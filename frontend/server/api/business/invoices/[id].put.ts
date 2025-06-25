import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../db'
import { authorize } from '../../auth'
import Joi from 'joi'

const updateInvoiceSchema = Joi.object({
  invoiceNumber: Joi.string().min(1).max(50).optional(),
  amount: Joi.number().positive().optional(),
  dueDate: Joi.date().greater('now').optional(),
  customerName: Joi.string().min(2).max(100).optional(),
  customerEmail: Joi.string().email().optional(),
  description: Joi.string().max(500).optional(),
  terms: Joi.string().max(200).optional()
}).min(1) // At least one field must be provided

/**
 * @openapi
 * /api/business/invoices/{id}:
 *   put:
 *     summary: Update Invoice
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
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               invoiceNumber:
 *                 type: string
 *                 minLength: 1
 *                 maxLength: 50
 *               amount:
 *                 type: number
 *                 minimum: 0.01
 *               dueDate:
 *                 type: string
 *                 format: date
 *                 description: Must be in the future
 *               customerName:
 *                 type: string
 *                 minLength: 2
 *                 maxLength: 100
 *               customerEmail:
 *                 type: string
 *                 format: email
 *               description:
 *                 type: string
 *                 maxLength: 500
 *               terms:
 *                 type: string
 *                 maxLength: 200
 *     responses:
 *       200:
 *         description: Invoice updated successfully
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
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Business access required
 *       404:
 *         description: Invoice not found
 *       409:
 *         description: Invoice number already exists
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

    // Validate request body
    const body = await readBody(event)
    const { error, value } = updateInvoiceSchema.validate(body)
    
    if (error) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Validation error',
        data: {
          details: error.details.map(detail => ({
            field: detail.path.join('.'),
            message: detail.message,
            value: detail.context?.value
          }))
        }
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

    const currentInvoice = existingInvoice.rows[0]

    // Check if invoice can be updated (not funded)
    if (currentInvoice.status === 'funded') {
      throw createError({
        statusCode: 400,
        statusMessage: 'Cannot update funded invoice'
      })
    }

    // Check if new invoice number already exists (if being updated)
    if (value.invoiceNumber && value.invoiceNumber !== currentInvoice.invoice_number) {
      const duplicateInvoice = await query(
        'SELECT id FROM invoices WHERE business_id = $1 AND invoice_number = $2 AND id != $3',
        [businessId, value.invoiceNumber, invoiceId]
      )

      if (duplicateInvoice.rows.length > 0) {
        throw createError({
          statusCode: 409,
          statusMessage: 'Invoice number already exists'
        })
      }
    }

    // Build update query
    const updateFields: string[] = []
    const updateValues: any[] = []
    let paramIndex = 1

    Object.keys(value).forEach(key => {
      if (value[key] !== undefined) {
        updateFields.push(`${key.replace(/([A-Z])/g, '_$1').toLowerCase()} = $${paramIndex}`)
        updateValues.push(value[key])
        paramIndex++
      }
    })

    // Add updated_at timestamp
    updateFields.push(`updated_at = $${paramIndex}`)
    updateValues.push(new Date())
    paramIndex++

    // Update invoice
    const updateResult = await query(
      `UPDATE invoices SET ${updateFields.join(', ')} WHERE id = $${paramIndex} RETURNING *`,
      [...updateValues, invoiceId]
    )

    const updatedInvoice = updateResult.rows[0]

    // Log the invoice update
    await query(
      `INSERT INTO audit_logs (user_id, action, table_name, record_id, old_values, new_values)
       VALUES ($1, 'UPDATE', 'invoices', $2, $3, $4)`,
      [
        user.id,
        invoiceId,
        JSON.stringify(currentInvoice),
        JSON.stringify(value)
      ]
    )

    return {
      success: true,
      message: 'Invoice updated successfully',
      data: { invoice: updatedInvoice },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Invoice update error:', error)
    throw error
  }
}) 