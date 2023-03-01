package uniottrans

import "github.com/opentracing/opentracing-go"

type TraceNative interface {
	Foreach(handler func(spn SpanNative))
}

type SpanNative interface {
	GetTraceID() string
	GetParentID() string
	GetSpanID() string
	GetService() string
	GetOperation() string
	GetMeta() map[string]string
	GetMetrics() map[string]Numeric
	GetSpanStatus() SpanStatus
	GetStart() int64
	GetEnd() int64
}

type UniversalTransformer interface {
	Parse(bts []byte) ([]TraceNative, error)
	BuildSpan(native SpanNative) *Span
	opentracing.Tracer
}
