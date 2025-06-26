import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../db'
import { authorize } from '../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const id = event.context.params.id
    const body = await readBody(event)
    const fields = []
    const values = []
    let idx = 1
    for (const key of ['name', 'type', 'is_active', 'funding_capacity', 'interest_rate_range', 'contact_email']) {
      if (body[key] !== undefined) {
        fields.push(`${key} = $${idx}`)
        values.push(body[key])
        idx++
      }
    }
    if (fields.length === 0) {
      throw createError({ statusCode: 400, statusMessage: 'No fields to update' })
    }
    values.push(id)
    const result = await query(`UPDATE financial_institutions SET ${fields.join(', ')}, updated_at = NOW() WHERE id = $${idx} RETURNING *`, values)
    if (result.rows.length === 0) {
      throw createError({ statusCode: 404, statusMessage: 'Not found' })
    }
    return {
      success: true,
      data: result.rows[0],
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 