import { defineEventHandler, getHeader, getQuery, readBody } from 'h3'

// Log levels
export enum LogLevel {
  ERROR = 'error',
  WARN = 'warn',
  INFO = 'info',
  DEBUG = 'debug'
}

// Log interface
export interface LogEntry {
  timestamp: string
  level: LogLevel
  message: string
  method?: string
  url?: string
  statusCode?: number
  responseTime?: number
  userAgent?: string
  ip?: string
  userId?: string
  requestId?: string
  details?: any
}

// Logger class
export class Logger {
  private static instance: Logger
  private logLevel: LogLevel

  private constructor() {
    this.logLevel = (process.env.LOG_LEVEL as LogLevel) || LogLevel.INFO
  }

  public static getInstance(): Logger {
    if (!Logger.instance) {
      Logger.instance = new Logger()
    }
    return Logger.instance
  }

  private shouldLog(level: LogLevel): boolean {
    const levels = Object.values(LogLevel)
    const currentLevelIndex = levels.indexOf(this.logLevel)
    const messageLevelIndex = levels.indexOf(level)
    return messageLevelIndex <= currentLevelIndex
  }

  private formatLog(entry: LogEntry): string {
    const parts = [
      `[${entry.timestamp}]`,
      entry.level.toUpperCase(),
      entry.message
    ]

    if (entry.method && entry.url) {
      parts.push(`${entry.method} ${entry.url}`)
    }

    if (entry.statusCode) {
      parts.push(`Status: ${entry.statusCode}`)
    }

    if (entry.responseTime) {
      parts.push(`Response Time: ${entry.responseTime}ms`)
    }

    if (entry.ip) {
      parts.push(`IP: ${entry.ip}`)
    }

    if (entry.userId) {
      parts.push(`User: ${entry.userId}`)
    }

    if (entry.requestId) {
      parts.push(`Request ID: ${entry.requestId}`)
    }

    return parts.join(' | ')
  }

  public log(entry: LogEntry): void {
    if (!this.shouldLog(entry.level)) {
      return
    }

    const formattedLog = this.formatLog(entry)
    
    switch (entry.level) {
      case LogLevel.ERROR:
        console.error(formattedLog, entry.details || '')
        break
      case LogLevel.WARN:
        console.warn(formattedLog, entry.details || '')
        break
      case LogLevel.INFO:
        console.info(formattedLog, entry.details || '')
        break
      case LogLevel.DEBUG:
        console.debug(formattedLog, entry.details || '')
        break
    }

    // In production, you might want to send logs to a service like CloudWatch, ELK, etc.
    if (process.env.NODE_ENV === 'production') {
      this.sendToLogService(entry)
    }
  }

  private sendToLogService(entry: LogEntry): void {
    // Implementation for sending logs to external service
    // This is a placeholder - implement based on your logging service
    if (process.env.LOG_SERVICE_URL) {
      // Example: Send to external logging service
      // fetch(process.env.LOG_SERVICE_URL, {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(entry)
      // }).catch(err => console.error('Failed to send log to service:', err))
    }
  }

  // Convenience methods
  public error(message: string, details?: any): void {
    this.log({
      timestamp: new Date().toISOString(),
      level: LogLevel.ERROR,
      message,
      details
    })
  }

  public warn(message: string, details?: any): void {
    this.log({
      timestamp: new Date().toISOString(),
      level: LogLevel.WARN,
      message,
      details
    })
  }

  public info(message: string, details?: any): void {
    this.log({
      timestamp: new Date().toISOString(),
      level: LogLevel.INFO,
      message,
      details
    })
  }

  public debug(message: string, details?: any): void {
    this.log({
      timestamp: new Date().toISOString(),
      level: LogLevel.DEBUG,
      message,
      details
    })
  }
}

// Request logging middleware
export default defineEventHandler(async (event) => {
  const logger = Logger.getInstance()
  const startTime = Date.now()
  const requestId = generateRequestId()

  // Add request ID to context
  event.context.requestId = requestId

  // Get request details
  const method = event.method
  const url = event.path || event.node.req.url || ''
  const userAgent = getHeader(event, 'user-agent') || 'Unknown'
  const ip = getClientIP(event) || 'Unknown'
  const query = getQuery(event)
  const body = method !== 'GET' ? await readBody(event).catch(() => ({})) : {}

  // Log request
  logger.info('Incoming request', {
    requestId,
    method,
    url,
    userAgent,
    ip,
    query: Object.keys(query).length > 0 ? query : undefined,
    body: Object.keys(body).length > 0 ? body : undefined
  })

  // Store original send method
  const originalSend = event.node.res.send

  // Override send method to log response
  event.node.res.send = function(data: any) {
    const responseTime = Date.now() - startTime
    const statusCode = event.node.res.statusCode

    // Log response
    logger.info('Outgoing response', {
      requestId,
      method,
      url,
      statusCode,
      responseTime,
      ip,
      userAgent
    })

    // Call original send method
    return originalSend.call(this, data)
  }

  // Continue with the request
  return
})

// Utility functions
function generateRequestId(): string {
  return `req_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

function getClientIP(event: any): string | undefined {
  // Check various headers for client IP
  const headers = [
    'x-forwarded-for',
    'x-real-ip',
    'x-client-ip',
    'cf-connecting-ip', // Cloudflare
    'x-forwarded',
    'forwarded-for',
    'forwarded'
  ]

  for (const header of headers) {
    const value = getHeader(event, header)
    if (value) {
      // Handle comma-separated IPs (take the first one)
      const ip = value.split(',')[0].trim()
      if (ip && ip !== 'unknown') {
        return ip
      }
    }
  }

  // Fallback to connection remote address
  return event.node.req.connection?.remoteAddress || 
         event.node.req.socket?.remoteAddress ||
         undefined
}

// Export logger instance for use in other parts of the application
export const logger = Logger.getInstance() 