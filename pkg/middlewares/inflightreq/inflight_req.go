package inflightreq

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cgws/traefik/v2/pkg/config/dynamic"
	"github.com/cgws/traefik/v2/pkg/log"
	"github.com/cgws/traefik/v2/pkg/middlewares"
	"github.com/cgws/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/vulcand/oxy/connlimit"
)

const (
	typeName = "InFlightReq"
)

type inFlightReq struct {
	handler http.Handler
	name    string
}

// New creates a max request middleware.
func New(ctx context.Context, next http.Handler, config dynamic.InFlightReq, name string) (http.Handler, error) {
	ctxLog := log.With(ctx, log.Str(log.MiddlewareName, name), log.Str(log.MiddlewareType, typeName))
	log.FromContext(ctxLog).Debug("Creating middleware")

	if config.SourceCriterion == nil ||
		config.SourceCriterion.IPStrategy == nil &&
			config.SourceCriterion.RequestHeaderName == "" && !config.SourceCriterion.RequestHost {
		config.SourceCriterion = &dynamic.SourceCriterion{
			RequestHost: true,
		}
	}

	sourceMatcher, err := middlewares.GetSourceExtractor(ctxLog, config.SourceCriterion)
	if err != nil {
		return nil, fmt.Errorf("error creating requests limiter: %v", err)
	}

	handler, err := connlimit.New(next, sourceMatcher, config.Amount)
	if err != nil {
		return nil, fmt.Errorf("error creating connection limit: %v", err)
	}

	return &inFlightReq{handler: handler, name: name}, nil
}

func (i *inFlightReq) GetTracingInformation() (string, ext.SpanKindEnum) {
	return i.name, tracing.SpanKindNoneEnum
}

func (i *inFlightReq) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	i.handler.ServeHTTP(rw, req)
}
