import { defineEventHandler, readBody } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'
import Joi from 'joi'

const updateProfileSchema = Joi.object({
  companyName: Joi.string().min(2).max(100).optional(),
  industry: Joi.string().min(2).max(50).optional(),
  annualRevenue: Joi.number().positive().optional(),
  employeeCount: Joi.number().integer().positive().optional(),
  phone: Joi.string().pattern(/^\+?[1-9]\d{1,14}$/).optional(),
  address: Joi.object({
    street: Joi.string().min(5).max(200).optional(),
    city: Joi.string().min(2).max(50).optional(),
    state: Joi.string().min(2).max(50).optional(),
    zipCode: Joi.string().pattern(/^\d{5}(-\d{4})?$/).optional(),
    country: Joi.string().min(2).max(50).optional()
  }).optional()
}).min(1) // At least one field must be provided

export default defineEventHandler(async (event) => {
  try {
    // Verify business authentication
    const user = authorize('business')(event)

    // Validate request body
    const body = await readBody(event)
    const { error, value } = updateProfileSchema.validate(body)
    
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

    // Build update query
    const updateFields: string[] = []
    const updateValues: any[] = []
    let paramIndex = 1

    Object.keys(value).forEach(key => {
      if (value[key] !== undefined) {
        const dbField = key.replace(/([A-Z])/g, '_$1').toLowerCase()
        updateFields.push(`${dbField} = $${paramIndex}`)
        updateValues.push(value[key])
        paramIndex++
      }
    })

    // Add updated_at timestamp
    updateFields.push(`updated_at = $${paramIndex}`)
    updateValues.push(new Date())

    // Update business profile
    const updateResult = await query(
      `UPDATE businesses SET ${updateFields.join(', ')} WHERE id = $${paramIndex + 1} RETURNING *`,
      [...updateValues, businessId]
    )

    return {
      success: true,
      message: 'Business profile updated successfully',
      data: { business: updateResult.rows[0] },
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Business profile update error:', error)
    throw error
  }
}) 