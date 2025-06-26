import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../../db'
import { authorize } from '../../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const financialInstitutionId = event.context.params.id
    const body = await readBody(event)
    const { product_id, terms } = body
    if (!terms) {
      throw createError({ statusCode: 400, statusMessage: 'Terms are required' })
    }
    const result = await query(
      `INSERT INTO financial_institution_terms (financial_institution_id, product_id, terms)
       VALUES ($1, $2, $3) RETURNING *`,
      [financialInstitutionId, product_id, terms]
    )
    return {
      success: true,
      message: 'Terms created',
      data: result.rows[0],
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 