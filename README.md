This is a simple golang application that pushes metrics to Prometheus which are availbale in Grafana.

#### Steps to run
1. `docker-compose up -d`
2. Prometheus is available at 9090
3. Application rusn at 2112
4. Open grafana on port 3000
5. Set data source in grafana to Prometheus with strategy of pulling from browser
6. Add dashboards and all requests to application at 2112 are now visible in Grafana dashboard
