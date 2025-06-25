import { defineEventHandler, createError } from 'h3';
import { protect, authorize } from '../auth';

export default defineEventHandler(async (event) => {
    // Protect the route
    const user = protect(event);
    
    // Authorize business users only
    authorize('business')(event);
    
    return {
        success: true,
        message: 'Business API endpoint',
        user: user,
        timestamp: new Date().toISOString()
    };
}); 