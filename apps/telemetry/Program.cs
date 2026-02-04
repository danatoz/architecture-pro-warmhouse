using System.Text.Json;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

var builder = WebApplication.CreateBuilder(args);
var service = builder.Services;
var configuration = builder.Configuration;


service.AddScoped<TimescaleRepository>();
service.AddSingleton<TelemetryService>();
service.AddHostedService<ConsumerBgService>();
service.AddHealthChecks();

var app = builder.Build();

app.MapHealthChecks("/healtz");

await app.RunAsync();
