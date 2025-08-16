-- Development seed data for users table
-- This migration should only be run in development environments

INSERT INTO users (id, email, name, created_at, updated_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'john.doe@example.com', 'John Doe', NOW() - INTERVAL '30 days', NOW() - INTERVAL '5 days'),
    ('550e8400-e29b-41d4-a716-446655440002', 'jane.smith@example.com', 'Jane Smith', NOW() - INTERVAL '25 days', NOW() - INTERVAL '3 days'),
    ('550e8400-e29b-41d4-a716-446655440003', 'bob.wilson@example.com', 'Bob Wilson', NOW() - INTERVAL '20 days', NOW() - INTERVAL '1 day'),
    ('550e8400-e29b-41d4-a716-446655440004', 'alice.johnson@example.com', 'Alice Johnson', NOW() - INTERVAL '15 days', NOW() - INTERVAL '2 hours'),
    ('550e8400-e29b-41d4-a716-446655440005', 'charlie.brown@example.com', 'Charlie Brown', NOW() - INTERVAL '10 days', NOW() - INTERVAL '30 minutes'),
    ('550e8400-e29b-41d4-a716-446655440006', 'diana.prince@example.com', 'Diana Prince', NOW() - INTERVAL '5 days', NOW() - INTERVAL '10 minutes'),
    ('550e8400-e29b-41d4-a716-446655440007', 'edward.norton@example.com', 'Edward Norton', NOW() - INTERVAL '3 days', NOW() - INTERVAL '5 minutes'),
    ('550e8400-e29b-41d4-a716-446655440008', 'fiona.gallagher@example.com', 'Fiona Gallagher', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 minute'),
    ('550e8400-e29b-41d4-a716-446655440009', 'george.washington@example.com', 'George Washington', NOW() - INTERVAL '1 day', NOW()),
    ('550e8400-e29b-41d4-a716-446655440010', 'helen.keller@example.com', 'Helen Keller', NOW() - INTERVAL '12 hours', NOW())
ON CONFLICT (id) DO NOTHING;