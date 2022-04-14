## Troubleshooting

### Elasticsearch is unhealthy

Make sure `sysctl -w vm.max_map_count=262144` got run. Elasticsearch needs it to function properly.