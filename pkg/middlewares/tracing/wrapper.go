package tracing

import (
	"context"
	"net/http"

	"github.com/containous/alice"
	"github.com/cgws/traefik/v2/pkg/log"
	"github.com/cgws/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

// Tracable embeds tracing information.
type Tracable interface {
	GetTracingInformation() (name string, spanKind ext.SpanKindEnum)
}

// Wrap adds tracability to an alice.Constructor.
func Wrap(ctx context.Context, constructor alice.Constructor) alice.Constructor {
	return func(next http.Handler) (http.Handler, error) {
		if constructor == nil {
			return nil, nil
		}
		handler, err := constructor(next)
		if err != nil {
			return nil, err
		}

		if tracableHandler, ok := handler.(Tracable); ok {
			name, spanKind := tracableHandler.GetTracingInformation()
			log.FromContext(ctx).WithField(log.MiddlewareName, name).Debug("Adding tracing to middleware")
			return NewWrapper(handler, name, spanKind), nil
		}
		return handler, nil
	}
}

// NewWrapper returns a http.Handler struct
func NewWrapper(next http.Handler, name string, spanKind ext.SpanKindEnum) http.Handler {
	return &Wrapper{
		next:     next,
		name:     name,
		spanKind: spanKind,
	}
}

// Wrapper is used to wrap http handler middleware.
type Wrapper struct {
	next     http.Handler
	name     string
	spanKind ext.SpanKindEnum
}

func (w *Wrapper) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, err := tracing.FromContext(req.Context())
	if err != nil {
		w.next.ServeHTTP(rw, req)
		return
	}

	var finish func()
	_, req, finish = tracing.StartSpan(req, w.name, w.spanKind)
	defer finish()

	if w.next != nil {
		w.next.ServeHTTP(rw, req)
	}
}
