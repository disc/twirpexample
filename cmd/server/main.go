package main

import (
	"github.com/disc/twirpexample/hooks/metrics"
	"github.com/disc/twirpexample/hooks/sentry"
	"github.com/disc/twirpexample/internal/haberdasherserver"
	"github.com/disc/twirpexample/rpc/haberdasher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twitchtv/twirp"
	"net/http"
)

func main() {
	srv := &haberdasherserver.Server{} // implements Haberdasher interface

	//FIXME: Move to .env file
	sentryDsn := "https://xxx@sentry.io/yyy"
	sentryHook := sentry.NewSentryServerHooks(sentryDsn)
	metricsHook := metrics.NewMetricsServerHooks(nil)

	hooks := twirp.ChainHooks(sentryHook, metricsHook)

	twirpHandler := haberdasher.NewHaberdasherServer(srv, hooks)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	mux.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
