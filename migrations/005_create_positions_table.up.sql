-- Create positions table
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    department_id UUID,
    requirements TEXT NOT NULL,
    min_salary DECIMAL(10,2) NOT NULL,
    max_salary DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_positions_department_id FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX idx_positions_title ON positions(title);
CREATE INDEX idx_positions_department_id ON positions(department_id);
CREATE INDEX idx_positions_created_at ON positions(created_at);
CREATE INDEX idx_positions_updated_at ON positions(updated_at);

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_positions_updated_at 
    BEFORE UPDATE ON positions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();