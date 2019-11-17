package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/twitchtv/twirp"
)

func NewSentryServerHooks(dsn string) *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}

	sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})

	hooks.Error = func(ctx context.Context, err twirp.Error) context.Context {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			for k, v := range err.MetaMap() {
				scope.SetExtra(k, v)
			}
		})
		sentry.CaptureException(err)

		return ctx
	}

	return hooks
}
