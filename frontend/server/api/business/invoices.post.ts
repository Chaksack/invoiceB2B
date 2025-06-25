import { defineEventHandler, readBody } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'
import Joi from 'joi'

const createInvoiceSchema = Joi.object({
  invoiceNumber: Joi.string().min(1).max(50).required(),
  amount: Joi.number().positive().required(),
  dueDate: Joi.date().greater('now').required(),
  customerName: Joi.string().min(2).max(100).required(),
  customerEmail: Joi.string().email().required(),
  description: Joi.string().max(500).optional(),
  terms: Joi.string().max(200).optional()
})

/**
 * @openapi
 * /api/business/invoices:
 *   post:
 *     summary: Create New Invoice
 *     tags: [Business]
 *     security:
 *       - bearerAuth: []
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
 *             required:
 *               - invoiceNumber
 *               - amount
 *               - dueDate
 *               - customerName
 *               - customerEmail
 *     responses:
 *       201:
 *         description: Invoice created successfully
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
 *                         invoiceNumber:
 *                           type: string
 *                         amount:
 *                           type: number
 *                         status:
 *                           type: string
 *                         created_at:
 *                           type: string
 *                           format: date-time
 *       400:
 *         description: Validation error
 *       401:
 *         description: Unauthorized
 *       403:
 *         description: Forbidden - Business access required
 */
export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Validate request body
    const body = await readBody(event)
    const { error, value } = createInvoiceSchema.validate(body)
    
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

    // Check if invoice number already exists for this business
    const existingInvoice = await query(
      'SELECT id FROM invoices WHERE business_id = $1 AND invoice_number = $2',
      [businessId, value.invoiceNumber]
    )

    if (existingInvoice.rows.length > 0) {
      throw createError({
        statusCode: 409,
        statusMessage: 'Invoice number already exists'
      })
    }

    // Create invoice
    const invoiceResult = await query(
      `INSERT INTO invoices (
        business_id, invoice_number, amount, due_date, customer_name, 
        customer_email, description, terms, status
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending')
      RETURNING *`,
      [
        businessId,
        value.invoiceNumber,
        value.amount,
        value.dueDate,
        value.customerName,
        value.customerEmail,
        value.description || null,
        value.terms || null
      ]
    )

    const invoice = invoiceResult.rows[0]

    // Log the invoice creation
    await query(
      `INSERT INTO audit_logs (user_id, action, table_name, record_id, new_values)
       VALUES ($1, 'CREATE', 'invoices', $2, $3)`,
      [
        user.id,
        invoice.id,
        JSON.stringify({
          invoiceNumber: invoice.invoice_number,
          amount: invoice.amount,
          customerName: invoice.customer_name
        })
      ]
    )

    return {
      success: true,
      message: 'Invoice created successfully',
      data: { invoice },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Invoice creation error:', error)
    throw error
  }
}) 