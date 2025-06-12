package main

import (
	"net/http"
	profiler "proxy-profiler/internal"
	"strconv"

	"github.com/gookit/slog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := profiler.ReadConfig("./config/proxy/config.yml", "./config/proxy/config.default.yml")

	if err != nil {
		slog.Fatalf("Unable to load the Config, exiting: %s", err)

		return;
	}

	proxyHandler := profiler.NewProxyHandler(
		profiler.NewMetricsCollector(config),
		config,
	)

	http.Handle("/", proxyHandler)
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: ":" + strconv.Itoa(config.Proxy.Port),
	}

	slog.Infof("Listening on %d", config.Proxy.Port)

	tlsConfig := config.Proxy.Tls

	if "" != tlsConfig.Cert && "" != tlsConfig.Key  {
		slog.Fatal(server.ListenAndServeTLS(tlsConfig.Cert, tlsConfig.Key))
	} else {
		slog.Fatal(server.ListenAndServe())
	}

}
