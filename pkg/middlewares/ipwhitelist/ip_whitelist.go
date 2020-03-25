package ipwhitelist

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cgws/traefik/v2/pkg/config/dynamic"
	"github.com/cgws/traefik/v2/pkg/ip"
	"github.com/cgws/traefik/v2/pkg/log"
	"github.com/cgws/traefik/v2/pkg/middlewares"
	"github.com/cgws/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	typeName = "IPWhiteLister"
)

// ipWhiteLister is a middleware that provides Checks of the Requesting IP against a set of Whitelists
type ipWhiteLister struct {
	next        http.Handler
	whiteLister *ip.Checker
	strategy    ip.Strategy
	name        string
}

// New builds a new IPWhiteLister given a list of CIDR-Strings to whitelist
func New(ctx context.Context, next http.Handler, config dynamic.IPWhiteList, name string) (http.Handler, error) {
	logger := log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName))
	logger.Debug("Creating middleware")

	if len(config.SourceRange) == 0 {
		return nil, errors.New("sourceRange is empty, IPWhiteLister not created")
	}

	checker, err := ip.NewChecker(config.SourceRange)
	if err != nil {
		return nil, fmt.Errorf("cannot parse CIDR whitelist %s: %v", config.SourceRange, err)
	}

	strategy, err := config.IPStrategy.Get()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Setting up IPWhiteLister with sourceRange: %s", config.SourceRange)

	return &ipWhiteLister{
		strategy:    strategy,
		whiteLister: checker,
		next:        next,
		name:        name,
	}, nil
}

func (wl *ipWhiteLister) GetTracingInformation() (string, ext.SpanKindEnum) {
	return wl.name, tracing.SpanKindNoneEnum
}

func (wl *ipWhiteLister) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := middlewares.GetLoggerCtx(req.Context(), wl.name, typeName)
	logger := log.FromContext(ctx)

	err := wl.whiteLister.IsAuthorized(wl.strategy.GetIP(req))
	if err != nil {
		logMessage := fmt.Sprintf("rejecting request %+v: %v", req, err)
		logger.Debug(logMessage)
		tracing.SetErrorWithEvent(req, logMessage)
		reject(ctx, rw)
		return
	}
	logger.Debugf("Accept %s: %+v", wl.strategy.GetIP(req), req)

	wl.next.ServeHTTP(rw, req)
}

func reject(ctx context.Context, rw http.ResponseWriter) {
	statusCode := http.StatusForbidden

	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(http.StatusText(statusCode)))
	if err != nil {
		log.FromContext(ctx).Error(err)
	}
}
