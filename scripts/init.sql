-- MySQL initialization script for go-api-setup
-- This script runs when the MySQL container starts for the first time

-- Use the created database
USE go_api_setup;

-- Create users table (GORM will handle this with AutoMigrate, but this is a fallback)
-- This table structure should match your domain.User struct
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_email (email)
);

-- Insert some test data (optional - remove in production)
INSERT IGNORE INTO users (name, email, password, created_at, updated_at) VALUES
('Test User', 'test@example.com', '$2a$14$XYZ...', NOW(), NOW()),
('Admin User', 'admin@example.com', '$2a$14$ABC...', NOW(), NOW());

-- Print success message
SELECT 'Database initialized successfully!' as message; 