create database smarthome_telemetry;

CREATE TABLE sensor_telemetry (
    id BIGSERIAL PRIMARY KEY,
    sensor_id INT NOT NULL,
    sensor_type TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    unit TEXT,
    status TEXT,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

