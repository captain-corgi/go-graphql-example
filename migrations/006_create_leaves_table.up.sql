-- Create leaves table
CREATE TABLE leaves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    leave_type VARCHAR(20) NOT NULL CHECK (leave_type IN ('VACATION', 'SICK', 'PERSONAL', 'MATERNITY', 'PATERNITY', 'BEREAVEMENT', 'UNPAID', 'OTHER')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'CANCELLED')),
    approved_by UUID,
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_leaves_employee_id FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_leaves_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
    CONSTRAINT check_leaves_date_range CHECK (start_date <= end_date)
);

-- Create indexes for better query performance
CREATE INDEX idx_leaves_employee_id ON leaves(employee_id);
CREATE INDEX idx_leaves_status ON leaves(status);
CREATE INDEX idx_leaves_start_date ON leaves(start_date);
CREATE INDEX idx_leaves_end_date ON leaves(end_date);
CREATE INDEX idx_leaves_approved_by ON leaves(approved_by);
CREATE INDEX idx_leaves_created_at ON leaves(created_at);
CREATE INDEX idx_leaves_updated_at ON leaves(updated_at);

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_leaves_updated_at 
    BEFORE UPDATE ON leaves 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();