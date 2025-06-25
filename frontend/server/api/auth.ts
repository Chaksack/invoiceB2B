import { createError, getHeader } from 'h3';
import jwt from 'jsonwebtoken';

// Validate JWT secret exists
const JWT_SECRET = process.env.JWT_SECRET;
if (!JWT_SECRET || JWT_SECRET === 'your_jwt_secret_key_please_change_this') {
    throw new Error('JWT_SECRET environment variable must be set to a secure value');
}

// JWT verification helper
const verifyToken = (token: string) => {
    try {
        return jwt.verify(token, JWT_SECRET) as any;
    } catch (error) {
        throw createError({
            statusCode: 401,
            statusMessage: 'Invalid or expired token'
        });
    }
};

// Middleware to protect routes
export const protect = (event: any) => {
    const authHeader = getHeader(event, 'authorization');
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        throw createError({
            statusCode: 401,
            statusMessage: 'Authorization header required'
        });
    }

    const token = authHeader.substring(7);
    const decoded = verifyToken(token);
    
    // Attach user to event context
    event.context.user = decoded;
    return decoded;
};

// Middleware to restrict access based on role
export const authorize = (roles: string[] | string = []) => {
    return (event: any) => {
        const user = protect(event);
        const allowedRoles = Array.isArray(roles) ? roles : [roles];
        
        if (allowedRoles.length > 0 && !allowedRoles.includes(user.role)) {
            throw createError({
                statusCode: 403,
                statusMessage: 'Access denied'
            });
        }
        
        return user;
    };
}; 