-- Create the database if it doesn't exist
CREATE DATABASE smarthome_telemetry;

-- Connect to the database
\c smarthome_telemetry;

-- Create the sensor_telemetry table
CREATE TABLE IF NOT EXISTS sensor_telemetry (
    id BIGSERIAL PRIMARY KEY,
    sensor_id INT NOT NULL,
    sensor_type TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit TEXT,
    status TEXT,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

