FROM debian:latest
COPY ethereum-validator-metrics-exporter* /ethereum-validator-metrics-exporter
ENTRYPOINT ["/ethereum-validator-metrics-exporter"]
