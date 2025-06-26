import { defineEventHandler } from 'h3'
import { query } from '../../../../db'
import { authorize } from '../../../auth'

export default defineEventHandler(async (event) => {
  try {
    const user = authorize('admin')(event)
    const financialInstitutionId = event.context.params.id
    const result = await query(
      'SELECT * FROM financial_institution_products WHERE financial_institution_id = $1',
      [financialInstitutionId]
    )
    return {
      success: true,
      data: result.rows,
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 