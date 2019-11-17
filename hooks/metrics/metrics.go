package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/twitchtv/twirp"
)

var (
	requestsReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_requests_total",
			Help: "Number of RPC requests received.",
		},
		[]string{"method"},
	)

	responsesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_responses_total",
			Help: "Number of RPC responses sent.",
		},
		[]string{"method", "status"},
	)

	rpcDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "rpc_durations_seconds",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "status"},
	)
)

// MustRegister registers the prometheus with the registerer.
func MustRegister(registerer prometheus.Registerer) {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}
	registerer.MustRegister(requestsReceived)
	registerer.MustRegister(responsesSent)
	registerer.MustRegister(rpcDurations)
}

func NewMetricsServerHooks(registerer prometheus.Registerer) *twirp.ServerHooks {
	MustRegister(registerer)

	hooks := &twirp.ServerHooks{}

	hooks.RequestReceived = func(ctx context.Context) (context.Context, error) {
		ctx = markReqStart(ctx)
		return ctx, nil
	}

	hooks.RequestRouted = func(ctx context.Context) (context.Context, error) {
		method, ok := twirp.MethodName(ctx)
		if !ok {
			return ctx, nil
		}
		requestsReceived.WithLabelValues(method).Inc()
		return ctx, nil
	}

	hooks.ResponseSent = func(ctx context.Context) {
		method, _ := twirp.MethodName(ctx)
		status, _ := twirp.StatusCode(ctx)

		responsesSent.WithLabelValues(method, status).Inc()

		if start, ok := getReqStart(ctx); ok {
			dur := time.Now().Sub(start).Seconds()
			rpcDurations.WithLabelValues(method, status).Observe(dur)
		}
	}
	return hooks
}

var reqStartTimestampKey = new(int)

func markReqStart(ctx context.Context) context.Context {
	return context.WithValue(ctx, reqStartTimestampKey, time.Now())
}

func getReqStart(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(reqStartTimestampKey).(time.Time)
	return t, ok
}
