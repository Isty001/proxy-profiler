package internal

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gookit/slog"
)

type ProxyHandler struct {
	collector MetricsCollector
	config    Config
}

var proxies map[string]*httputil.ReverseProxy = map[string]*httputil.ReverseProxy{}

func NewProxyHandler(storage MetricsCollector, config Config) ProxyHandler {
	return ProxyHandler{collector: storage, config: config}
}

func (h ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	host := r.Host
	proxy, err := h.getProxy(host)

	if err != nil {
		slog.Errorf("Unable to create proxy for host %s - %w", host, err)

		return
	}

	proxy.ServeHTTP(w, r)
}

func (h ProxyHandler) getProxy(sourceHost string) (*httputil.ReverseProxy, error) {
	if p, ok := proxies[sourceHost]; ok {
		return p, nil
	}

	destination, ok := h.config.Proxy.Destination.SourceToDestinationHostMap[sourceHost]

	if !ok {
		slog.Infof("Source " + sourceHost + " not found in sourceToDestinationHostMap, using the defaultDestination")

		destination = h.config.Proxy.Destination.DefaultDestination
	}

	remoteUrl, err := url.Parse(destination.Scheme + "://" + destination.Host)

	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
	proxy.Transport = newProxyTransport(h.collector, h.config)
	proxies[sourceHost] = proxy

	return proxy, nil
}

type proxyTransport struct {
	transport *http.Transport
	collector MetricsCollector
}

func newProxyTransport(collector MetricsCollector, config Config) proxyTransport {
	transport := &http.Transport{}

	if config.Proxy.Destination.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return proxyTransport{
		transport: transport,
		collector: collector,
	}
}

func (t proxyTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	start := time.Now()
	response, err = t.transport.RoundTrip(request)

	t.collector.SaveRequestDuration(request, response, time.Since(start))

	return response, err
}
