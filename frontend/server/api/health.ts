import { defineEventHandler } from 'h3';
import { query } from '../db';

// Simple success response helper
const successResponse = (data: any, message = 'Success', statusCode = 200) => {
    return {
        success: true,
        message,
        data,
        timestamp: new Date().toISOString()
    };
};

// Async handler wrapper
const asyncHandler = (fn: Function) => {
    return async (event: any) => {
        try {
            return await fn(event);
        } catch (error) {
            console.error('Health check error:', error);
            throw createError({
                statusCode: 500,
                statusMessage: 'Internal server error'
            });
        }
    };
};

export default defineEventHandler(asyncHandler(async (event: any) => {
    const startTime = Date.now();
    
    // Basic health check with proper typing
    const health: any = {
        status: 'healthy',
        timestamp: new Date().toISOString(),
        uptime: process.uptime(),
        environment: process.env.NODE_ENV || 'development',
        version: process.env.APP_VERSION || '1.0.0'
    };

    // Database health check
    try {
        const dbResult = await query('SELECT NOW() as current_time');
        health.database = {
            status: 'connected',
            currentTime: dbResult.rows[0]?.current_time
        };
    } catch (error) {
        health.database = {
            status: 'disconnected',
            error: error instanceof Error ? error.message : 'Unknown error'
        };
        health.status = 'degraded';
    }

    // Memory usage
    const memUsage = process.memoryUsage();
    health.memory = {
        rss: Math.round(memUsage.rss / 1024 / 1024) + ' MB',
        heapTotal: Math.round(memUsage.heapTotal / 1024 / 1024) + ' MB',
        heapUsed: Math.round(memUsage.heapUsed / 1024 / 1024) + ' MB',
        external: Math.round(memUsage.external / 1024 / 1024) + ' MB'
    };

    // Response time
    const responseTime = Date.now() - startTime;
    health.responseTime = `${responseTime}ms`;

    // Set appropriate status code
    const statusCode = health.status === 'healthy' ? 200 : 503;
    
    return successResponse(health, 'Health check completed', statusCode);
})); 