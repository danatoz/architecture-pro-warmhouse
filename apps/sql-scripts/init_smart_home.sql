-- Create the database if it doesn't exist
CREATE DATABASE smarthome;

-- Connect to the database
\c smarthome;

-- Create the sensors table
CREATE TABLE IF NOT EXISTS sensors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    value FLOAT DEFAULT 0,
    unit VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sensors_type ON sensors(type);
CREATE INDEX IF NOT EXISTS idx_sensors_location ON sensors(location);
CREATE INDEX IF NOT EXISTS idx_sensors_status ON sensors(status);

BEGIN;

INSERT INTO sensors (name, type, location, value, unit, status, last_updated, created_at)
SELECT 'Температурный датчик', 'temperature', 'Комната 1', 22.5, '°C', 'active', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sensors WHERE name = 'Температурный датчик' AND location = 'Комната 1'
);

INSERT INTO sensors (name, type, location, value, unit, status, last_updated, created_at)
SELECT 'Датчик влажности', 'humidity', 'Склад', 60.0, '%', 'active', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM sensors WHERE name = 'Датчик влажности' AND location = 'Склад'
);

COMMIT;


