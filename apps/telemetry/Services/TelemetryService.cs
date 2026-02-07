using System;
using System.Text.Json;
using System.Threading;
using Confluent.Kafka;

public class TelemetryService
{
    private readonly ConsumerConfig _config = default!;
    private readonly IServiceProvider _serviceProvider;
    private readonly ILogger<TelemetryService> _logger;
    public TelemetryService(IConfiguration configuration, IServiceProvider serviceProvider, ILogger<TelemetryService> logger)
    {
        var consumerOptions = configuration.GetRequiredSection("ConsumerConfig").Get<ConsumerOptions>();
        var config = new ConsumerConfig();
        config.BootstrapServers = consumerOptions!.BootstrapServers;
        config.SaslMechanism = SaslMechanism.Plain;
        config.SecurityProtocol = SecurityProtocol.Plaintext;
        config.GroupId = consumerOptions.GroupId;
        _config = config;
        _logger = logger;
        _serviceProvider = serviceProvider;
    }

    public async Task StartConsume(CancellationToken cancellationToken)
    {
        using var consumer = new ConsumerBuilder<Ignore, string>(_config!).Build();
        using var scope = _serviceProvider.CreateScope();
        var repo = scope.ServiceProvider.GetRequiredService<TimescaleRepository>();
        consumer.Subscribe("smart-home.sensors-datas.telemetry");
        try
        {
            while (true)
            {
                var result = consumer.Consume(cancellationToken);
                var str = result.Message.Value;
                await Handle(repo, str);
            }
        }
        catch (OperationCanceledException)
        {
            // Ctrl+C
        }
        finally
        {
            consumer.Close(); // корректный commit offsets
        }
    }

    private async Task Handle(TimescaleRepository repo, string message)
    {
        try
        {
            var data = JsonSerializer.Deserialize<SensorData>(message, new JsonSerializerOptions());
            //{"SensorID":123,"Type":"temperature","Unit":"C","Value":26.1,"Status":"PopPit","Timestamp":"2026-02-03T18:59:14.1115728+00:00"}
            if(data is null){
                _logger.LogError("Data is empty");
                return;
            }
            _logger.LogInformation("New telemetry {@Data}", data);

            await repo.InsertMetricAsync(data!);
        }
        catch(JsonException ex)
        {
            _logger.LogError(ex, "Error parse");
        }
        catch(Exception ex)
        {
            _logger.LogError(ex, "Enexpected error");
        }
    }
}
