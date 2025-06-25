import { defineEventHandler, readBody, createError } from 'h3';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { query } from '../../db';

// Validate JWT secret exists
const JWT_SECRET = process.env.JWT_SECRET;
if (!JWT_SECRET || JWT_SECRET === 'your_jwt_secret_key_please_change_this') {
    throw new Error('JWT_SECRET environment variable must be set to a secure value');
}

// Success response helper
const successResponse = (data: any, message = 'Success', statusCode = 200) => {
    return {
        success: true,
        message,
        data,
        timestamp: new Date().toISOString()
    };
};

// Validation schemas
const schemas = {
    register: {
        email: (value: string) => {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Invalid email format'
                });
            }
            if (value.length > 255) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Email too long'
                });
            }
            return value;
        },
        password: (value: string) => {
            if (value.length < 8) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Password must be at least 8 characters long'
                });
            }
            if (value.length > 128) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Password too long'
                });
            }
            const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/;
            if (!passwordRegex.test(value)) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character'
                });
            }
            return value;
        },
        role: (value: string = 'business') => {
            if (!['business', 'admin'].includes(value)) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Invalid role'
                });
            }
            return value;
        }
    },
    login: {
        email: (value: string) => {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(value)) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Invalid email format'
                });
            }
            return value;
        },
        password: (value: string) => {
            if (!value) {
                throw createError({
                    statusCode: 400,
                    statusMessage: 'Password is required'
                });
            }
            return value;
        }
    }
};

// Validation helper
const validate = (schema: any, data: any) => {
    const validated: any = {};
    for (const [key, validator] of Object.entries(schema)) {
        if (data[key] !== undefined) {
            validated[key] = (validator as Function)(data[key]);
        }
    }
    return validated;
};

// Register endpoint
export const register = defineEventHandler(async (event) => {
    try {
        const body = await readBody(event);
        const { email, password, role } = validate(schemas.register, body);
        
        const hashedPassword = await bcrypt.hash(password, 12);
        const isApproved = role === 'admin';

        const result = await query(
            'INSERT INTO users (email, password_hash, role, is_approved) VALUES ($1, $2, $3, $4) RETURNING id, email, role, is_approved, created_at',
            [email, hashedPassword, role, isApproved]
        );
        
        const user = result.rows[0];

        console.log(`User registered: ${user.email}. Approval status: ${user.is_approved}`);
        
        return successResponse(
            { 
                id: user.id, 
                email: user.email, 
                role: user.role, 
                is_approved: user.is_approved,
                created_at: user.created_at
            },
            'User registered successfully. Awaiting admin approval.',
            201
        );
    } catch (error: any) {
        if (error.statusCode) {
            throw error;
        }
        if (error.code === '23505') {
            throw createError({
                statusCode: 409,
                statusMessage: 'Email already registered'
            });
        }
        console.error('Registration error:', error);
        throw createError({
            statusCode: 500,
            statusMessage: 'Internal server error',
            data: {
                details: error.details?.map((detail: any) => ({
                    field: detail.path?.join('.'),
                    message: detail.message,
                    value: detail.context?.value
                })) || []
            }
        });
    }
});

// Login endpoint
export const login = defineEventHandler(async (event) => {
    try {
        const body = await readBody(event);
        const { email, password } = validate(schemas.login, body);
        
        const result = await query(
            'SELECT id, email, password_hash, role, is_approved, created_at FROM users WHERE email = $1',
            [email]
        );
        
        const user = result.rows[0];

        if (!user) {
            throw createError({
                statusCode: 401,
                statusMessage: 'Invalid credentials'
            });
        }

        const isPasswordValid = await bcrypt.compare(password, user.password_hash);
        if (!isPasswordValid) {
            throw createError({
                statusCode: 401,
                statusMessage: 'Invalid credentials'
            });
        }

        if (!user.is_approved && user.role === 'business') {
            throw createError({
                statusCode: 403,
                statusMessage: 'Account awaiting admin approval'
            });
        }

        const token = jwt.sign(
            { 
                id: user.id, 
                email: user.email, 
                role: user.role,
                iat: Math.floor(Date.now() / 1000)
            },
            Buffer.from(JWT_SECRET),
            { 
                expiresIn: process.env.JWT_EXPIRES_IN || '1h',
                issuer: process.env.JWT_ISSUER || 'invoice-b2b-api',
                audience: process.env.JWT_AUDIENCE || 'invoice-b2b-client'
            }
        );

        return successResponse({
            token,
            user: { 
                id: user.id, 
                email: user.email, 
                role: user.role,
                is_approved: user.is_approved,
                created_at: user.created_at
            }
        }, 'Logged in successfully');
    } catch (error: any) {
        if (error.statusCode) {
            throw error;
        }
        console.error('Login error:', error);
        throw createError({
            statusCode: 500,
            statusMessage: 'Internal server error',
            data: {
                details: error.details?.map((detail: any) => ({
                    field: detail.path?.join('.'),
                    message: detail.message,
                    value: detail.context?.value
                })) || []
            }
        });
    }
});

// Profile endpoint (protected)
export const profile = defineEventHandler(async (event) => {
    try {
        const authHeader = getHeader(event, 'authorization');
        if (!authHeader || !authHeader.startsWith('Bearer ')) {
            throw createError({
                statusCode: 401,
                statusMessage: 'Authorization header required'
            });
        }

        const token = authHeader.substring(7);
        const decoded = jwt.verify(token, JWT_SECRET) as any;

        const result = await query(
            'SELECT id, email, role, is_approved, created_at, updated_at FROM users WHERE id = $1',
            [decoded.id]
        );
        
        const user = result.rows[0];
        if (!user) {
            throw createError({
                statusCode: 404,
                statusMessage: 'User not found'
            });
        }

        return successResponse(user, 'User profile retrieved successfully');
    } catch (error: any) {
        if (error.statusCode) {
            throw error;
        }
        console.error('Profile error:', error);
        throw createError({
            statusCode: 500,
            statusMessage: 'Internal server error',
            data: {
                details: error.details?.map((detail: any) => ({
                    field: detail.path?.join('.'),
                    message: detail.message,
                    value: detail.context?.value
                })) || []
            }
        });
    }
});

// Refresh token endpoint
export const refresh = defineEventHandler(async (event) => {
    try {
        const body = await readBody(event);
        const { token } = body;
        
        if (!token) {
            throw createError({
                statusCode: 400,
                statusMessage: 'Token is required'
            });
        }

        const decoded = jwt.verify(token, JWT_SECRET) as any;
        
        const result = await query(
            'SELECT id, email, role, is_approved FROM users WHERE id = $1',
            [decoded.id]
        );
        
        const user = result.rows[0];
        if (!user) {
            throw createError({
                statusCode: 404,
                statusMessage: 'User not found'
            });
        }

        if (!user.is_approved && user.role === 'business') {
            throw createError({
                statusCode: 403,
                statusMessage: 'Account awaiting admin approval'
            });
        }

        const newToken = jwt.sign(
            { 
                id: user.id, 
                email: user.email, 
                role: user.role,
                iat: Math.floor(Date.now() / 1000)
            },
            Buffer.from(JWT_SECRET),
            { 
                expiresIn: process.env.JWT_EXPIRES_IN || '1h',
                issuer: process.env.JWT_ISSUER || 'invoice-b2b-api',
                audience: process.env.JWT_AUDIENCE || 'invoice-b2b-client'
            }
        );

        return successResponse({
            token: newToken,
            user: { 
                id: user.id, 
                email: user.email, 
                role: user.role,
                is_approved: user.is_approved
            }
        }, 'Token refreshed successfully');
    } catch (error: any) {
        if (error.statusCode) {
            throw error;
        }
        console.error('Token refresh error:', error);
        throw createError({
            statusCode: 500,
            statusMessage: 'Internal server error',
            data: {
                details: error.details?.map((detail: any) => ({
                    field: detail.path?.join('.'),
                    message: detail.message,
                    value: detail.context?.value
                })) || []
            }
        });
    }
}); 