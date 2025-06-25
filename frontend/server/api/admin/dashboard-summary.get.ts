import { defineEventHandler } from 'h3'
import { query } from '../../db'
import { authorize } from '../auth'

export default defineEventHandler(async (event) => {
  try {
    // Verify admin authentication
    const user = authorize('admin')(event)

    // Get dashboard statistics
    const [
      businessesResult,
      pendingBusinessesResult,
      invoicesResult,
      fundedAmountResult,
      pendingInvoicesResult,
      financialInstitutionsResult,
      recentActivitiesResult
    ] = await Promise.all([
      query('SELECT COUNT(*) as count FROM businesses'),
      query("SELECT COUNT(*) as count FROM businesses WHERE status = 'pending'"),
      query('SELECT COUNT(*) as count FROM invoices'),
      query("SELECT COALESCE(SUM(amount), 0) as total FROM invoices WHERE status = 'funded'"),
      query("SELECT COUNT(*) as count FROM invoices WHERE status = 'pending'"),
      query('SELECT COUNT(*) as count FROM financial_institutions WHERE is_active = true'),
      query(`
        SELECT 
          'business' as type,
          b.company_name as name,
          b.status,
          b.created_at
        FROM businesses b
        WHERE b.created_at >= NOW() - INTERVAL '7 days'
        UNION ALL
        SELECT 
          'invoice' as type,
          i.invoice_number as name,
          i.status,
          i.created_at
        FROM invoices i
        WHERE i.created_at >= NOW() - INTERVAL '7 days'
        ORDER BY created_at DESC
        LIMIT 10
      `)
    ])

    const summary = {
      total_businesses: parseInt(businessesResult.rows[0].count),
      pending_businesses: parseInt(pendingBusinessesResult.rows[0].count),
      total_invoices: parseInt(invoicesResult.rows[0].count),
      total_funded_amount: parseFloat(fundedAmountResult.rows[0].total),
      pending_invoices: parseInt(pendingInvoicesResult.rows[0].count),
      total_financial_institutions: parseInt(financialInstitutionsResult.rows[0].count),
      recent_activities: recentActivitiesResult.rows
    }

    return {
      success: true,
      message: 'Dashboard summary retrieved successfully',
      data: summary,
      timestamp: new Date().toISOString()
    }
  } catch (error) {
    console.error('Admin dashboard error:', error)
    throw error
  }
}) 