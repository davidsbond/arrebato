# Debugging

This document aims to outline how to run a development server locally and debug it.

To build the project, just run `make build`. For tests, there are 2 commands.
To run the unit test suite, use `make test`, for the end-to-end testing suite
use `make test-e2e`.

## Docker Compose

At the root of this repository is a docker-compose file. You can use this to
spin up a development cluster. Usually, a single call to `docker-compose up`
would be sufficient. However, this can cause some problems with the consensus
of the nodes within the cluster. It is recommended that individual services
within docker compose be spun up individually:

```shell
docker-compose up arrebato-0
# Wait for arrebato-0 to enter leadership state
docker-compose up arrebato-1
docker-compose up arrebato-2
```

## OpenTelemetry

The arrebato project provides the ability to enable distributed tracing via
[OpenTelemetry](https://opentelemetry.io/). The [docker-compose file](../../docker-compose.yaml)
at the root of this repository includes the ability to configure the
OpenTelemetry collector and [Jaeger](https://www.jaegertracing.io/) for viewing
the traces.

If you aren't using the docker-compose file, you just need to provide the
`--tracing-enabled` and `--tracing-endpoint` flags to enable tracing and route
spans via HTTP to an OpenTelemetry collector.

Once Jaeger and the collector are running, you should be able to start a
cluster as described above and view the Jaeger UI at `http://localhost:16686`.
