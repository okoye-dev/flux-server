-- Migration: Add farmer_crops junction table for multiple crops per farmer
-- This allows farmers to register multiple crops instead of just one

-- Create farmer_crops junction table
CREATE TABLE IF NOT EXISTS farmer_crops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farmer_id BIGINT NOT NULL,
    crop_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Foreign key constraints
    CONSTRAINT fk_farmer_crops_farmer FOREIGN KEY (farmer_id) REFERENCES farmers(id) ON DELETE CASCADE,
    CONSTRAINT fk_farmer_crops_crop FOREIGN KEY (crop_id) REFERENCES crops(id) ON DELETE CASCADE,
    
    -- Ensure unique farmer-crop combinations
    CONSTRAINT unique_farmer_crop UNIQUE (farmer_id, crop_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_farmer_crops_farmer_id ON farmer_crops(farmer_id);
CREATE INDEX IF NOT EXISTS idx_farmer_crops_crop_id ON farmer_crops(crop_id);
CREATE INDEX IF NOT EXISTS idx_farmer_crops_created_at ON farmer_crops(created_at);

-- Enable Row Level Security
ALTER TABLE farmer_crops ENABLE ROW LEVEL SECURITY;

-- Create policy for service role access
CREATE POLICY "Service role can access all farmer_crops" ON farmer_crops
    FOR ALL USING (auth.role() = 'service_role');

-- Create policy for farmers to access their own crop data
CREATE POLICY "Farmers can access own crop data" ON farmer_crops
    FOR ALL USING (
        farmer_id IN (
            SELECT f.id FROM farmers f 
            WHERE f.auth_user_id = auth.uid()
        )
    );

-- Insert some common crops if they don't exist
-- First, add a unique constraint on the name column (ignore error if it already exists)
DO $$ 
BEGIN
    BEGIN
        ALTER TABLE crops ADD CONSTRAINT unique_crop_name UNIQUE (name);
    EXCEPTION
        WHEN duplicate_object THEN null;
    END;
END $$;

-- Now insert the crops with conflict handling
INSERT INTO crops (id, name, scientific_name, created_at) VALUES
    (gen_random_uuid(), 'Maize', 'Zea mays', NOW()),
    (gen_random_uuid(), 'Rice', 'Oryza sativa', NOW()),
    (gen_random_uuid(), 'Wheat', 'Triticum aestivum', NOW()),
    (gen_random_uuid(), 'Sorghum', 'Sorghum bicolor', NOW()),
    (gen_random_uuid(), 'Millet', 'Pennisetum glaucum', NOW()),
    (gen_random_uuid(), 'Beans', 'Phaseolus vulgaris', NOW()),
    (gen_random_uuid(), 'Groundnuts', 'Arachis hypogaea', NOW()),
    (gen_random_uuid(), 'Cassava', 'Manihot esculenta', NOW()),
    (gen_random_uuid(), 'Sweet Potato', 'Ipomoea batatas', NOW()),
    (gen_random_uuid(), 'Tomatoes', 'Solanum lycopersicum', NOW()),
    (gen_random_uuid(), 'Onions', 'Allium cepa', NOW()),
    (gen_random_uuid(), 'Cabbage', 'Brassica oleracea', NOW()),
    (gen_random_uuid(), 'Carrots', 'Daucus carota', NOW()),
    (gen_random_uuid(), 'Peppers', 'Capsicum annuum', NOW()),
    (gen_random_uuid(), 'Okra', 'Abelmoschus esculentus', NOW())
ON CONFLICT (name) DO NOTHING;
