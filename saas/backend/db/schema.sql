-- ============================================================================
-- TS2GO SaaS Database Schema
-- ============================================================================
-- PostgreSQL 15+ required for GENERATED ALWAYS AS IDENTITY
-- Created: 2025-11-16
-- Description: Core database schema for ts2go SaaS platform
-- ============================================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- USERS & AUTHENTICATION
-- ============================================================================

-- User accounts
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url VARCHAR(500),
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API keys for programmatic access
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    key_prefix VARCHAR(20) NOT NULL, -- First 8 chars for display (e.g., "ts2go_1a2b...")
    scopes TEXT[] DEFAULT '{}', -- Array of permissions: ['read:projects', 'write:transpile']
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- OAuth providers (GitHub, Google, etc.)
CREATE TABLE oauth_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'github', 'google', 'microsoft'
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

-- ============================================================================
-- TEAMS & COLLABORATION
-- ============================================================================

-- Teams/Organizations
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    avatar_url VARCHAR(500),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_type VARCHAR(20) DEFAULT 'free', -- 'free', 'pro', 'team', 'enterprise'
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Team memberships
CREATE TABLE team_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member', 'viewer'
    invited_by UUID REFERENCES users(id),
    invited_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    joined_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(team_id, user_id)
);

-- ============================================================================
-- PROJECTS & TRANSPILATION
-- ============================================================================

-- Projects (collections of TypeScript code)
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    description TEXT,
    visibility VARCHAR(20) DEFAULT 'private', -- 'private', 'team', 'public'
    settings JSONB DEFAULT '{}', -- Project-specific transpilation settings
    total_files INTEGER DEFAULT 0,
    total_size_bytes BIGINT DEFAULT 0,
    last_transpiled_at TIMESTAMP WITH TIME ZONE,
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, slug)
);

-- Transpilation jobs
CREATE TABLE transpilations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    job_id VARCHAR(100) UNIQUE NOT NULL, -- Queue job identifier
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed', 'cancelled'
    input_files JSONB NOT NULL, -- Array of file metadata
    output_files JSONB, -- Array of generated files
    settings JSONB DEFAULT '{}', -- Transpilation options
    error_message TEXT,
    input_size_bytes BIGINT,
    output_size_bytes BIGINT,
    processing_time_ms INTEGER,
    worker_id VARCHAR(100), -- Which worker processed this
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- File storage metadata
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    transpilation_id UUID REFERENCES transpilations(id) ON DELETE SET NULL,
    file_type VARCHAR(20) NOT NULL, -- 'input', 'output'
    original_name VARCHAR(500) NOT NULL,
    storage_path VARCHAR(1000) NOT NULL, -- S3/MinIO path
    mime_type VARCHAR(100),
    size_bytes BIGINT NOT NULL,
    checksum VARCHAR(64), -- SHA256
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- BILLING & USAGE
-- ============================================================================

-- Subscription plans
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    stripe_subscription_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    plan_type VARCHAR(20) NOT NULL, -- 'free', 'pro', 'team', 'enterprise'
    status VARCHAR(20) DEFAULT 'active', -- 'active', 'cancelled', 'past_due', 'trialing'
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (user_id IS NOT NULL OR team_id IS NOT NULL)
);

-- Usage tracking
CREATE TABLE usage_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    resource_type VARCHAR(50) NOT NULL, -- 'transpilation', 'api_call', 'storage'
    quantity INTEGER DEFAULT 1,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Monthly usage aggregates (for quota enforcement)
CREATE TABLE usage_quotas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    transpilations_count INTEGER DEFAULT 0,
    api_calls_count INTEGER DEFAULT 0,
    storage_bytes BIGINT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CHECK (user_id IS NOT NULL OR team_id IS NOT NULL),
    UNIQUE(user_id, period_start),
    UNIQUE(team_id, period_start)
);

-- ============================================================================
-- ANALYTICS & AUDIT
-- ============================================================================

-- Audit log for security/compliance
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- 'user.login', 'project.create', 'transpilation.start'
    resource_type VARCHAR(50),
    resource_id UUID,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- API Keys
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_is_active ON api_keys(is_active) WHERE is_active = TRUE;

-- Teams
CREATE INDEX idx_teams_owner_id ON teams(owner_id);
CREATE INDEX idx_teams_slug ON teams(slug);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);

-- Projects
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_team_id ON projects(team_id);
CREATE INDEX idx_projects_slug ON projects(user_id, slug);
CREATE INDEX idx_projects_visibility ON projects(visibility);

-- Transpilations
CREATE INDEX idx_transpilations_user_id ON transpilations(user_id);
CREATE INDEX idx_transpilations_project_id ON transpilations(project_id);
CREATE INDEX idx_transpilations_status ON transpilations(status);
CREATE INDEX idx_transpilations_created_at ON transpilations(created_at DESC);
CREATE INDEX idx_transpilations_job_id ON transpilations(job_id);

-- Files
CREATE INDEX idx_files_project_id ON files(project_id);
CREATE INDEX idx_files_transpilation_id ON files(transpilation_id);

-- Usage
CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_team_id ON usage_records(team_id);
CREATE INDEX idx_usage_records_created_at ON usage_records(created_at DESC);
CREATE INDEX idx_usage_quotas_user_id_period ON usage_quotas(user_id, period_start);
CREATE INDEX idx_usage_quotas_team_id_period ON usage_quotas(team_id, period_start);

-- Audit
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- ============================================================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_keys_updated_at BEFORE UPDATE ON api_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transpilations_updated_at BEFORE UPDATE ON transpilations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Development Only)
-- ============================================================================

-- Default admin user (password: "admin123" - CHANGE IN PRODUCTION!)
-- Password hash generated with bcrypt cost 12
INSERT INTO users (email, username, password_hash, full_name, email_verified, is_active)
VALUES (
    'admin@ts2go.dev',
    'admin',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIVj3MKmVS', -- "admin123"
    'System Administrator',
    TRUE,
    TRUE
) ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- VIEWS FOR COMMON QUERIES
-- ============================================================================

-- Active user subscriptions with plan details
CREATE VIEW v_user_subscriptions AS
SELECT 
    u.id as user_id,
    u.email,
    u.username,
    COALESCE(s.plan_type, t.plan_type, 'free') as current_plan,
    s.status as subscription_status,
    s.current_period_end as subscription_expires,
    t.id as team_id,
    t.name as team_name
FROM users u
LEFT JOIN subscriptions s ON u.id = s.user_id AND s.status = 'active'
LEFT JOIN team_members tm ON u.id = tm.user_id AND tm.is_active = TRUE
LEFT JOIN teams t ON tm.team_id = t.id AND t.is_active = TRUE;

-- Project statistics
CREATE VIEW v_project_stats AS
SELECT 
    p.id,
    p.name,
    p.user_id,
    p.team_id,
    COUNT(DISTINCT t.id) as total_transpilations,
    COUNT(DISTINCT CASE WHEN t.status = 'completed' THEN t.id END) as successful_transpilations,
    COUNT(DISTINCT CASE WHEN t.status = 'failed' THEN t.id END) as failed_transpilations,
    SUM(t.processing_time_ms) as total_processing_time_ms,
    MAX(t.completed_at) as last_transpiled_at
FROM projects p
LEFT JOIN transpilations t ON p.id = t.project_id
GROUP BY p.id, p.name, p.user_id, p.team_id;

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE users IS 'User accounts with authentication credentials';
COMMENT ON TABLE api_keys IS 'API keys for programmatic access to the platform';
COMMENT ON TABLE teams IS 'Organizations/teams for collaborative work';
COMMENT ON TABLE projects IS 'User projects containing TypeScript code to transpile';
COMMENT ON TABLE transpilations IS 'Transpilation job records and results';
COMMENT ON TABLE usage_records IS 'Individual usage events for billing/analytics';
COMMENT ON TABLE usage_quotas IS 'Monthly aggregated usage for quota enforcement';
COMMENT ON TABLE audit_logs IS 'Security audit trail for compliance';
