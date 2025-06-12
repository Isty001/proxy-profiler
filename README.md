This is a simple reverse proxy that captures each requests and measures the the it takes for the destination to server a response, then collect this data into 
Prometheus compatible metrics to be displayed in Grafana.

# Usage

You can provide your own config in `config/proxy/config.yml`.

## Source - Destination

You can map the different source hosts to destinations if that's necessary

```yml
proxy:
  sourceToDestinationHostMap:
    source1-example.com:
      scheme: 'https'
      host: 'destination1-example.com'
    source2-example.com:
      scheme: 'https'
      host: 'destination2-example.com'
```

If this is not set, or the source is not found, then the default will be use

```yml
proxy:
  destination:
    defaultDestination:
      scheme: 'https'
      host: 'destination-example.com'
```

In case the proxy is running on the same host machine as the destination app, then you can just use `host.docker.internal` as the destination host.


## HTTPS

If you want to run the proxy to run on https, then you need to provide a cert and key in `config/proxy/config.yml`:

```yml
proxy:
  tls:
    cert: './path/to/proxy.cert'
    key: './path/to/proxy.key'
```

If you are using a self-signed cert, then you need to add the following to `prometheus.yml` for the `profiler-metrics` job:

```yml
scheme: https
tls_config:
    insecure_skip_verify: true
```

In case the destination server uses self-signed certs, then you can need to add the following to config.yml:

```yml
proxy:
  destination:
    insecureSkipVerify: true
```

## Running

`make up` will start and instance of Grafana and Prometheus additionally, by default you can acces them on the host at:

Prometheus: `http://127.0.0.1:9090`
Grafana: `http://127.0.0.1:3000`

`make up-standalone` will only start the proxy, in case you already have your own Prometheus

## Grafana

By default, you can add the data source at `http://prometheus:9090`
