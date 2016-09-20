# metrics-fetcher [![Build Status](https://travis-ci.com/Wikia/metrics-fetcher.svg?token=8Hc4nTuxXPoC7GveyjkW&branch=master)](https://travis-ci.com/Wikia/metrics-fetcher)
Tool which pulls metrics from services registered in Consul and send them aggregated to InfluxDB/telegraf

## Sample config
```yaml
filters:
    - path: "io\\.dropwizard\\.db\\.ManagedPooledDataSource\\..*-master\\.idle"
      group: "gauges"
      type: "int64"
      measurement: "http_server"
    - path: "org\\.eclipse\\.jetty\\.util\\.thread\\.QueuedThreadPool\\.dw\\.jobs"
      group: "gauges"
      type: "int64"
      measurement: "http_server"
    - path: "jvm\\.memory\\.pools\\..*\\.usage"
      group: "gauges"
      type: "float64"
      measurement: "jvm_memory"
```