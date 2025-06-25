import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../../db'
import { authorize } from '../../../auth'
import Joi from 'joi'

const businessStatusSchema = Joi.object({
  status: Joi.string().valid('approved', 'rejected', 'suspended').required(),
  reason: Joi.string().max(500).optional()
})

export default defineEventHandler(async (event) => {
  try {
    // Verify admin authentication
    const user = authorize('admin')(event)

    // Get business ID from URL
    const businessId = getRouterParam(event, 'id')
    if (!businessId) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Business ID is required'
      })
    }

    // Validate request body
    const body = await readBody(event)
    const { error, value } = businessStatusSchema.validate(body)
    
    if (error) {
      throw createError({
        statusCode: 400,
        statusMessage: 'Validation error',
        data: {
          details: error.details.map((detail: any) => ({
            field: detail.path.join('.'),
            message: detail.message,
            value: detail.context?.value
          }))
        }
      })
    }

    // Check if business exists
    const businessResult = await query(
      'SELECT id, status FROM businesses WHERE id = $1',
      [businessId]
    )

    if (businessResult.rows.length === 0) {
      throw createError({
        statusCode: 404,
        statusMessage: 'Business not found'
      })
    }

    const currentStatus = businessResult.rows[0].status

    // Update business status
    const updateData: any = {
      status: value.status,
      updated_at: new Date()
    }

    if (value.status === 'approved') {
      updateData.approved_at = new Date()
      updateData.approved_by = user.id
    }

    const updateFields = Object.keys(updateData).map((key, index) => `${key} = $${index + 2}`).join(', ')
    const updateValues = Object.values(updateData)

    await query(
      `UPDATE businesses SET ${updateFields} WHERE id = $1`,
      [businessId, ...updateValues]
    )

    // Log the status change
    await query(
      `INSERT INTO audit_logs (user_id, action, table_name, record_id, old_values, new_values)
       VALUES ($1, 'STATUS_UPDATE', 'businesses', $2, $3, $4)`,
      [
        user.id,
        businessId,
        JSON.stringify({ status: currentStatus }),
        JSON.stringify({ status: value.status, reason: value.reason })
      ]
    )

    // Get updated business
    const updatedResult = await query(
      'SELECT * FROM businesses WHERE id = $1',
      [businessId]
    )

    return {
      success: true,
      message: `Business status updated to ${value.status}`,
      data: { business: updatedResult.rows[0] },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Business status update error:', error)
    throw error
  }
}) 