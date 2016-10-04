# metrics-fetcher [![Build Status](https://travis-ci.org/Wikia/metrics-fetcher.svg?branch=master)](https://travis-ci.org/Wikia/metrics-fetcher) [![Coverage Status](https://coveralls.io/repos/github/Wikia/metrics-fetcher/badge.svg?branch=master)](https://coveralls.io/github/Wikia/metrics-fetcher?branch=master)
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