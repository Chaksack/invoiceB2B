import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../../db'
import { authorize } from '../../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const financialInstitutionId = event.context.params.id
    const body = await readBody(event)
    const { name, description } = body
    if (!name) {
      throw createError({ statusCode: 400, statusMessage: 'Product name is required' })
    }
    const result = await query(
      `INSERT INTO financial_institution_products (financial_institution_id, name, description)
       VALUES ($1, $2, $3) RETURNING *`,
      [financialInstitutionId, name, description]
    )
    return {
      success: true,
      message: 'Product created',
      data: result.rows[0],
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 