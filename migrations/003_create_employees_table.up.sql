-- Create employees table
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_code VARCHAR(10) UNIQUE NOT NULL,
    department VARCHAR(50) NOT NULL,
    position VARCHAR(50) NOT NULL,
    hire_date DATE NOT NULL,
    salary DECIMAL(10,2) NOT NULL CHECK (salary > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('ACTIVE', 'INACTIVE', 'TERMINATED', 'ON_LEAVE')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_employees_user_id ON employees(user_id);
CREATE INDEX idx_employees_employee_code ON employees(employee_code);
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_employees_status ON employees(status);
CREATE INDEX idx_employees_hire_date ON employees(hire_date);
CREATE INDEX idx_employees_created_at ON employees(created_at);
CREATE INDEX idx_employees_updated_at ON employees(updated_at);

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_employees_updated_at 
    BEFORE UPDATE ON employees 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();