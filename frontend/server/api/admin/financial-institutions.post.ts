import { defineEventHandler, readBody } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const body = await readBody(event)
    const { name, type, is_active, funding_capacity, interest_rate_range, contact_email } = body
    if (!name || !contact_email) {
      throw createError({ statusCode: 400, statusMessage: 'Name and contact_email are required' })
    }
    const result = await query(
      `INSERT INTO financial_institutions (name, type, is_active, funding_capacity, interest_rate_range, contact_email)
       VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
      [name, type, is_active ?? true, funding_capacity, interest_rate_range, contact_email]
    )
    return {
      success: true,
      message: 'Financial institution created',
      data: result.rows[0],
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 