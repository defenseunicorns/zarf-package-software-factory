## Troubleshooting

### Elasticsearch/Kibana/Jaeger is unhealthy

Make sure `sysctl -w vm.max_map_count=262144` got run on the server. Elasticsearch needs it to function properly.
