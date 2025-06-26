import { defineEventHandler } from 'h3'
import { query } from '../../../../../db'
import { authorize } from '../../../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const financialInstitutionId = event.context.params.id
    const productId = event.context.params.productId
    const result = await query(
      'DELETE FROM financial_institution_products WHERE id = $1 AND financial_institution_id = $2 RETURNING *',
      [productId, financialInstitutionId]
    )
    if (result.rows.length === 0) {
      throw createError({ statusCode: 404, statusMessage: 'Not found' })
    }
    return {
      success: true,
      message: 'Deleted',
      data: result.rows[0],
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 