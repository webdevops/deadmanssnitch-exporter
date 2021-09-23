# DeadMansSnitch Exporter

[![license](https://img.shields.io/github/license/webdevops/deadmanssnitch-exporter.svg)](https://github.com/webdevops/deadmanssnitch-exporter/blob/master/LICENSE)
[![DockerHub](https://img.shields.io/badge/DockerHub-webdevops%2Fdeadmanssnitch--exporter-blue)](https://hub.docker.com/r/webdevops/deadmanssnitch-exporter/)
[![Quay.io](https://img.shields.io/badge/Quay.io-webdevops%2Fdeadmanssnitch--exporter-blue)](https://quay.io/repository/webdevops/deadmanssnitch-exporter)

Prometheus exporter for [DeadMansSnitch](https://deadmanssnitch.com/) information (snitches)

## Configuration

```
Usage:
  deadmanssnitch-exporter [OPTIONS]

Application Options:
      --debug                 debug mode [$DEBUG]
  -v, --verbose               verbose mode [$VERBOSE]
      --log.json              Switch log output to json format [$LOG_JSON]
      --deadmanssnitch.token= DeadMansSnitch access token [$DEADMANSSNITCH_TOKEN]
      --bind=                 Server address (default: :8080) [$SERVER_BIND]
      --scrape.time=          Scrape time (time.duration) (default: 5m) [$SCRAPE_TIME]

Help Options:
  -h, --help                  Show this help message
```

## Metrics

| Metric                                | Scraper            | Description                                                                           |
|---------------------------------------|--------------------|---------------------------------------------------------------------------------------|
| `deadmanssnitch_stats`                | Collector          | Collector stats                                                                       |
| `deadmanssnitch_api_counter`          | Collector          | Api call counter                                                                      |
| `deadmanssnitch_snitch_info`          | Snitch             | Basic snitch information                                                              |
| `deadmanssnitch_snitch_heartbeat`     | Snitch             | Last heartbeat timestamp                                                              |
| `deadmanssnitch_snitch_status`        | Snitch             | Current health status                                                                 |
