using Npgsql;
using System;
using System.Reflection.Metadata;
using System.Threading.Tasks;

public class TimescaleRepository
{
    private readonly string _connectionString;

    public TimescaleRepository(IConfiguration configuration)
    {
        _connectionString = configuration.GetConnectionString("Default") ??
            throw new ArgumentException("Connection string is empty");
    }

    public async Task InsertMetricAsync(SensorData sensorData)
    {
        const string sql = @"
            INSERT INTO sensor_telemetry (sensor_id, sensor_type, value, unit, status, recorded_at)
            VALUES (@sensor_id, @sensor_type, @value, @unit, @status, @recorded_at);
        ";

        await using var connection = new NpgsqlConnection(_connectionString);
        await connection.OpenAsync();

        await using var command = new NpgsqlCommand(sql, connection);
        command.Parameters.AddWithValue("sensor_id", sensorData.SensorID);
        command.Parameters.AddWithValue("sensor_type", sensorData.Type);
        command.Parameters.AddWithValue("value", sensorData.Value);
        command.Parameters.AddWithValue("unit", sensorData.Unit);
        command.Parameters.AddWithValue("status", sensorData.Status);
        command.Parameters.AddWithValue("recorded_at", sensorData.Timestamp.UtcDateTime);

        await command.ExecuteNonQueryAsync();
    }
}
