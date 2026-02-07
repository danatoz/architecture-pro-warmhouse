
using System.Text.Json.Serialization;

public class SensorData
{
	[JsonPropertyName("sensor_id")]
    public int SensorID { get; set; }

	[JsonPropertyName("type")]
	public string Type { get; set; } = "";

	[JsonPropertyName("unit")]
	public string Unit { get; set; } = "";

	[JsonPropertyName("value")]
	public double Value { get; set; } = 0.0;

	[JsonPropertyName("status")]
	public string Status { get; set; } = "";
	
	[JsonPropertyName("timestamp")]
	public DateTimeOffset Timestamp { get; set; } = DateTimeOffset.UtcNow;
}
