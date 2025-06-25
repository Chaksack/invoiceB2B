import { defineEventHandler, createError } from 'h3';
import { protect, authorize } from '../auth';

export default defineEventHandler(async (event) => {
    // Protect the route
    const user = protect(event);
    
    // Authorize admin users only
    authorize('admin')(event);
    
    return {
        success: true,
        message: 'Admin API endpoint',
        user: user,
        timestamp: new Date().toISOString()
    };
}); 