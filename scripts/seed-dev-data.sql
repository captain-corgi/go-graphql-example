-- Development seed data script
-- Run this script to populate the database with sample data for development
-- Usage: psql -d graphql_service_dev -f scripts/seed-dev-data.sql

-- Ensure the users table exists
\echo 'Seeding development data...'

-- Insert sample users with realistic data
INSERT INTO users (id, email, name, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'john.doe@example.com', 'John Doe', '2024-07-17 10:00:00+00', '2024-08-11 15:30:00+00'),
    ('550e8400-e29b-41d4-a716-446655440002', 'jane.smith@example.com', 'Jane Smith', '2024-07-22 14:20:00+00', '2024-08-13 09:45:00+00'),
    ('550e8400-e29b-41d4-a716-446655440003', 'bob.wilson@example.com', 'Bob Wilson', '2024-07-27 08:15:00+00', '2024-08-15 12:20:00+00'),
    ('550e8400-e29b-41d4-a716-446655440004', 'alice.johnson@example.com', 'Alice Johnson', '2024-08-01 16:45:00+00', '2024-08-14 07:10:00+00'),
    ('550e8400-e29b-41d4-a716-446655440005', 'charlie.brown@example.com', 'Charlie Brown', '2024-08-06 11:30:00+00', '2024-08-15 21:30:00+00'),
    ('550e8400-e29b-41d4-a716-446655440006', 'diana.prince@example.com', 'Diana Prince', '2024-08-11 13:00:00+00', '2024-08-16 07:50:00+00'),
    ('550e8400-e29b-41d4-a716-446655440007', 'edward.norton@example.com', 'Edward Norton', '2024-08-13 09:20:00+00', '2024-08-16 07:55:00+00'),
    ('550e8400-e29b-41d4-a716-446655440008', 'fiona.gallagher@example.com', 'Fiona Gallagher', '2024-08-14 17:40:00+00', '2024-08-16 07:59:00+00'),
    ('550e8400-e29b-41d4-a716-446655440009', 'george.washington@example.com', 'George Washington', '2024-08-15 12:00:00+00', '2024-08-16 08:00:00+00'),
    ('550e8400-e29b-41d4-a716-446655440010', 'helen.keller@example.com', 'Helen Keller', '2024-08-15 20:00:00+00', '2024-08-16 08:00:00+00')
ON CONFLICT (id) DO NOTHING;

-- Display inserted data
\echo 'Development data seeded successfully!'
\echo 'Sample users:'
SELECT id, email, name, created_at FROM users ORDER BY created_at;