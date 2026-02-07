FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

COPY ["telemetry.csproj", "./"]

COPY . .
RUN dotnet build "telemetry.csproj" -c Release -o /app/build && \
    dotnet publish "telemetry.csproj" --no-restore -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
WORKDIR /app
COPY --from=build /app/publish .

RUN apt-get update && apt-get install -y libgssapi-krb5-2

ENV ASPNETCORE_URLS=http://+:8080

EXPOSE 8080

ENTRYPOINT ["dotnet", "telemetry.dll"]
