# metrics-fetcher [![Build Status](https://travis-ci.org/Wikia/metrics-fetcher.svg?branch=master)](https://travis-ci.org/Wikia/metrics-fetcher) [![Coverage Status](https://coveralls.io/repos/github/Wikia/metrics-fetcher/badge.svg?branch=master)](https://coveralls.io/github/Wikia/metrics-fetcher?branch=master)
Tool which pulls metrics from services registered in Consul and send them aggregated to InfluxDB/telegraf

## Sample config
```yaml
filters:
    - path: "io\\.dropwizard\\.db\\.ManagedPooledDataSource\\..*-master\\.idle"
      group: "gauges"
      measurement: "http_server"
    - path: "org\\.eclipse\\.jetty\\.util\\.thread\\.QueuedThreadPool\\.dw\\.jobs"
      group: "gauges"
      measurement: "http_server"
    - path: "jvm\\.memory\\.pools\\..*\\.usage"
      group: "gauges"
      measurement: "jvm_memory"
```

## Running
`metrics-fetcher fetch --label metrics --marathon http://marathon.service.consul:8080 --influx http://influx.service.consul:8086 --database test`

## Releasing
Do it only on **master** branch!

* install [bumpversion](https://github.com/peritus/bumpversion)
* install [github-changelog-generator](https://github.com/skywinder/github-changelog-generator)
* run `github-changelog-generator -u Wikia -p metrics-fetcher`
* `git add CHANGELOG.md`
* commit changes
* run `bumpversion patch` (or replace `patch` with either `minor` or `major`)
* `git push --tags`
* `git push`

