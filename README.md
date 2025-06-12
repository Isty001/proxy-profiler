# HTTPS

If you want to run the proxy to run on https, then you need to provide a cert and key in `config/proxy/config.yml`:

```yml
proxy:
  tls:
    cert: './path/to/proxy.cert'
    key: './path/to/proxy.key'
```

If you are using a self-signed cert, then you need to add the following to `prometheus.yml` for the `profiler-metrics` job:

```
scheme: https
tls_config:
    insecure_skip_verify: true
```

In case the destination server uses self-signed certs, then you can need to add the following to config.yml:

```
proxy:
  destination:
    insecureSkipVerify: true
```

# Grafana

Add data source at `http://prometheus:9090`
