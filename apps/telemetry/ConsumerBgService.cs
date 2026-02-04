public class ConsumerBgService : BackgroundService
{
    private readonly ILogger<ConsumerBgService> _logger;
    private readonly TelemetryService _consumer;
    public ConsumerBgService(ILogger<ConsumerBgService> logger, TelemetryService consumer)
    {
        _logger = logger;
        _consumer = consumer;
    }

    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        _logger.LogInformation("Bg service started");

        Task.Run(() => RunConsumer(stoppingToken), stoppingToken);

        return Task.CompletedTask;
    }

    private async Task RunConsumer(CancellationToken token)
    {
        try
        {
            await _consumer.StartConsume(token);
        }
        catch (OperationCanceledException)
        {
            _logger.LogInformation("Consumer cancelled");
        }
        catch (Exception ex)
        {
            _logger.LogCritical(ex, "Consumer crashed");
        }
    }
}