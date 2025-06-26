-- Financial Institutions Table
CREATE TABLE IF NOT EXISTS financial_institutions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    funding_capacity NUMERIC,
    interest_rate_range VARCHAR(100),
    contact_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Products Table
CREATE TABLE IF NOT EXISTS financial_institution_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    financial_institution_id UUID REFERENCES financial_institutions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Terms and Conditions Table
CREATE TABLE IF NOT EXISTS financial_institution_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    financial_institution_id UUID REFERENCES financial_institutions(id) ON DELETE CASCADE,
    product_id UUID REFERENCES financial_institution_products(id) ON DELETE CASCADE,
    terms TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
); 