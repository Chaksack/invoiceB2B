import { Pool, QueryResult, PoolConfig } from 'pg';

// Production-ready database configuration
const dbConfig: PoolConfig = {
    connectionString: process.env.DATABASE_URL,
    max: parseInt(process.env.DB_POOL_MAX || '20'), // Maximum number of clients in the pool
    min: parseInt(process.env.DB_POOL_MIN || '4'),  // Minimum number of clients in the pool
    connectionTimeoutMillis: parseInt(process.env.DB_CONNECTION_TIMEOUT || '10000'), // How long to wait for a connection
    idleTimeoutMillis: parseInt(process.env.DB_IDLE_TIMEOUT || '30000'), // How long a connection can be idle
    query_timeout: parseInt(process.env.DB_QUERY_TIMEOUT || '30000'), // Query timeout in milliseconds
    statement_timeout: parseInt(process.env.DB_STATEMENT_TIMEOUT || '30000'), // Statement timeout in milliseconds
};

const pool = new Pool(dbConfig);

// Enhanced error handling
pool.on('error', (err: Error, client: any) => {
    console.error('Unexpected error on idle client', err);
    // Don't exit process in production, just log the error
    if (process.env.NODE_ENV === 'development') {
        process.exit(-1);
    }
});

pool.on('connect', (client: any) => {
    console.log('New client connected to database');
});

pool.on('acquire', (client: any) => {
    console.log('Client acquired from pool');
});

pool.on('release', (client: any) => {
    console.log('Client released back to pool');
});

// Graceful shutdown
process.on('SIGINT', async () => {
    console.log('Received SIGINT, closing database pool...');
    await pool.end();
    process.exit(0);
});

process.on('SIGTERM', async () => {
    console.log('Received SIGTERM, closing database pool...');
    await pool.end();
    process.exit(0);
});

// Health check function
export const checkDatabaseHealth = async (): Promise<boolean> => {
    try {
        const result = await pool.query('SELECT 1 as health_check');
        return result.rows[0]?.health_check === 1;
    } catch (error) {
        console.error('Database health check failed:', error);
        return false;
    }
};

// Enhanced query function with timeout and better error handling
export const query = async (text: string, params?: any[]): Promise<QueryResult> => {
    const start = Date.now();
    try {
        const result = await pool.query(text, params);
        const duration = Date.now() - start;
        console.log('Executed query', { text, duration, rows: result.rowCount });
        return result;
    } catch (error: any) {
        const duration = Date.now() - start;
        console.error('Query error', { text, duration, error: error.message });
        throw error;
    }
};

// Transaction helper
export const transaction = async (callback: (client: any) => Promise<any>) => {
    const client = await pool.connect();
    try {
        await client.query('BEGIN');
        const result = await callback(client);
        await client.query('COMMIT');
        return result;
    } catch (error) {
        await client.query('ROLLBACK');
        throw error;
    } finally {
        client.release();
    }
};

export default {
    query,
    getClient: () => pool.connect(),
    checkHealth: checkDatabaseHealth,
    transaction,
    pool
};