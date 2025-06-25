import { defineEventHandler, readBody } from 'h3'
import Joi from 'joi'

// Common validation schemas
export const commonSchemas = {
  uuid: Joi.string().uuid().required(),
  email: Joi.string().email().required(),
  password: Joi.string().min(8).pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/).required(),
  phone: Joi.string().pattern(/^\+?[1-9]\d{1,14}$/).optional(),
  pagination: {
    page: Joi.number().integer().min(1).default(1),
    limit: Joi.number().integer().min(1).max(100).default(10)
  }
}

// Authentication schemas
export const authSchemas = {
  register: Joi.object({
    email: commonSchemas.email,
    password: commonSchemas.password,
    firstName: Joi.string().min(2).max(50).required(),
    lastName: Joi.string().min(2).max(50).required(),
    phone: commonSchemas.phone,
    companyName: Joi.string().min(2).max(100).required(),
    industry: Joi.string().min(2).max(50).required(),
    annualRevenue: Joi.number().positive().required(),
    employeeCount: Joi.number().integer().positive().required(),
    address: Joi.object({
      street: Joi.string().min(5).max(200).required(),
      city: Joi.string().min(2).max(50).required(),
      state: Joi.string().min(2).max(50).required(),
      zipCode: Joi.string().pattern(/^\d{5}(-\d{4})?$/).required(),
      country: Joi.string().min(2).max(50).required()
    }).required()
  }),

  login: Joi.object({
    email: commonSchemas.email,
    password: Joi.string().required()
  }),

  refreshToken: Joi.object({
    refreshToken: Joi.string().required()
  }),

  verify2FA: Joi.object({
    token: Joi.string().length(6).pattern(/^\d+$/).required()
  })
}

// Business schemas
export const businessSchemas = {
  updateProfile: Joi.object({
    companyName: Joi.string().min(2).max(100).optional(),
    industry: Joi.string().min(2).max(50).optional(),
    annualRevenue: Joi.number().positive().optional(),
    employeeCount: Joi.number().integer().positive().optional(),
    phone: commonSchemas.phone,
    address: Joi.object({
      street: Joi.string().min(5).max(200).optional(),
      city: Joi.string().min(2).max(50).optional(),
      state: Joi.string().min(2).max(50).optional(),
      zipCode: Joi.string().pattern(/^\d{5}(-\d{4})?$/).optional(),
      country: Joi.string().min(2).max(50).optional()
    }).optional()
  }).min(1), // At least one field must be provided

  createInvoice: Joi.object({
    invoiceNumber: Joi.string().min(1).max(50).required(),
    amount: Joi.number().positive().required(),
    dueDate: Joi.date().greater('now').required(),
    customerName: Joi.string().min(2).max(100).required(),
    customerEmail: commonSchemas.email,
    description: Joi.string().max(500).optional(),
    terms: Joi.string().max(200).optional()
  }),

  updateInvoice: Joi.object({
    invoiceNumber: Joi.string().min(1).max(50).optional(),
    amount: Joi.number().positive().optional(),
    dueDate: Joi.date().greater('now').optional(),
    customerName: Joi.string().min(2).max(100).optional(),
    customerEmail: commonSchemas.email.optional(),
    description: Joi.string().max(500).optional(),
    terms: Joi.string().max(200).optional()
  }).min(1)
}

// Admin schemas
export const adminSchemas = {
  updateBusinessStatus: Joi.object({
    status: Joi.string().valid('approved', 'rejected', 'suspended').required(),
    reason: Joi.string().max(500).optional()
  }),

  listBusinesses: Joi.object({
    page: commonSchemas.pagination.page,
    limit: commonSchemas.pagination.limit,
    status: Joi.string().valid('pending', 'approved', 'rejected', 'suspended').optional(),
    search: Joi.string().min(1).max(100).optional()
  })
}

// Validation middleware factory
export function validateBody(schema: Joi.ObjectSchema) {
  return defineEventHandler(async (event) => {
    try {
      const body = await readBody(event)
      const { error, value } = schema.validate(body, {
        abortEarly: false,
        stripUnknown: true
      })

      if (error) {
        throw createError({
          statusCode: 400,
          statusMessage: 'Validation error',
          data: {
            details: error.details.map(detail => ({
              field: detail.path.join('.'),
              message: detail.message,
              value: detail.context?.value
            }))
          }
        })
      }

      // Replace the body with validated data
      event.context.validatedBody = value
    } catch (error) {
      if (error.statusCode === 400) {
        throw error
      }
      throw createError({
        statusCode: 400,
        statusMessage: 'Invalid request body'
      })
    }
  })
}

// Validation middleware factory for query parameters
export function validateQuery(schema: Joi.ObjectSchema) {
  return defineEventHandler(async (event) => {
    try {
      const query = getQuery(event)
      const { error, value } = schema.validate(query, {
        abortEarly: false,
        stripUnknown: true
      })

      if (error) {
        throw createError({
          statusCode: 400,
          statusMessage: 'Validation error',
          data: {
            details: error.details.map(detail => ({
              field: detail.path.join('.'),
              message: detail.message,
              value: detail.context?.value
            }))
          }
        })
      }

      // Replace the query with validated data
      event.context.validatedQuery = value
    } catch (error) {
      if (error.statusCode === 400) {
        throw error
      }
      throw createError({
        statusCode: 400,
        statusMessage: 'Invalid query parameters'
      })
    }
  })
}

// Validation middleware factory for path parameters
export function validateParams(schema: Joi.ObjectSchema) {
  return defineEventHandler(async (event) => {
    try {
      const params = getRouterParams(event)
      const { error, value } = schema.validate(params, {
        abortEarly: false,
        stripUnknown: true
      })

      if (error) {
        throw createError({
          statusCode: 400,
          statusMessage: 'Validation error',
          data: {
            details: error.details.map(detail => ({
              field: detail.path.join('.'),
              message: detail.message,
              value: detail.context?.value
            }))
          }
        })
      }

      // Replace the params with validated data
      event.context.validatedParams = value
    } catch (error) {
      if (error.statusCode === 400) {
        throw error
      }
      throw createError({
        statusCode: 400,
        statusMessage: 'Invalid path parameters'
      })
    }
  })
} 