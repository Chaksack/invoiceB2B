import { defineEventHandler, readBody } from 'h3'
import { query } from '../../../../db'
import nodemailer from 'nodemailer'
import formidable from 'formidable'
import fs from 'fs'
import { authorize } from '~/server/api/auth'

export default defineEventHandler(async (event) => {
  try {
    // Authorize business user
    const user = authorize('business')(event)
    const financialInstitutionId = event.context.params.id

    // Parse form data (fields + files)
    const form = formidable({ multiples: true })
    const { fields, files } = await new Promise((resolve, reject) => {
      form.parse(event.node.req, (err, fields, files) => {
        if (err) reject(err)
        else resolve({ fields, files })
      })
    })

    // Fetch financial institution contact email
    const result = await query('SELECT contact_email, name FROM financial_institutions WHERE id = $1', [financialInstitutionId])
    if (result.rows.length === 0) {
      throw createError({ statusCode: 404, statusMessage: 'Financial institution not found' })
    }
    const { contact_email, name: institutionName } = result.rows[0]

    // Prepare email
    const transporter = nodemailer.createTransport({
      host: process.env.SMTP_HOST,
      port: parseInt(process.env.SMTP_PORT || '587'),
      secure: false,
      auth: {
        user: process.env.SMTP_USER,
        pass: process.env.SMTP_PASS
      }
    })

    const attachments = []
    if (files && Object.keys(files).length > 0) {
      for (const key in files) {
        const file = Array.isArray(files[key]) ? files[key][0] : files[key]
        attachments.push({
          filename: file.originalFilename,
          content: fs.createReadStream(file.filepath)
        })
      }
    }

    const mailOptions = {
      from: process.env.SMTP_FROM || 'noreply@example.com',
      to: contact_email,
      subject: `New request from business user for ${institutionName}`,
      text: `A business user has submitted a request.\n\nDetails: ${JSON.stringify(fields, null, 2)}`,
      attachments
    }

    await transporter.sendMail(mailOptions)

    return {
      success: true,
      message: 'Request sent to financial institution',
      timestamp: new Date().toISOString()
    }
  } catch (error: any) {
    throw createError({ statusCode: 500, statusMessage: error.message })
  }
}) 