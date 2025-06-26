import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../../../db'
import { authorize } from '../../../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const financialInstitutionId = event.context.params.id
    const productId = event.context.params.productId
    const body = await readBody(event)
    const fields = []
    const values = []
    let idx = 1
    for (const key of ['name', 'description']) {
      if (body[key] !== undefined) {
        fields.push(`${key} = $${idx}`)
        values.push(body[key])
        idx++
      }
    }
    if (fields.length === 0) {
      throw createError({ statusCode: 400, statusMessage: 'No fields to update' })
    }
    values.push(productId, financialInstitutionId)
    const result = await query(
      `UPDATE financial_institution_products SET ${fields.join(', ')}, updated_at = NOW() WHERE id = $${idx} AND financial_institution_id = $${idx+1} RETURNING *`,
      values
    )
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