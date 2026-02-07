#!/bin/bash

echo "Создаем топики..."

kafka-topics --create --topic smart-home.sensors-datas.telemetry --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
