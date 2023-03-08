package uniottrans

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	MaxBufferSize = 1024
	MaxThreads    = 8
)

type StartTracerOption func(tracer *Tracer)

func WithBufferSize(size int) StartTracerOption {
	return func(tracer *Tracer) {

	}
}

func WithThreads(num int) StartTracerOption {
	return func(tracer *Tracer) {

	}
}

func WithService(name string) StartTracerOption {
	return func(tracer *Tracer) {

	}
}

func WithMeta(meta map[string]string) StartTracerOption {
	return func(tracer *Tracer) {

	}
}

func WithSampleRatio(ratio int) StartTracerOption {
	return func(tracer *Tracer) {

	}
}

func NewTracer(opts ...StartTracerOption) *Tracer {

}

type Tracer struct {
	buffer, threads int
	service         string
	meta            map[string]string
	sampleRatio     int
	flush           chan struct{}
	close           chan struct{}
}

// Create, start, and return a new Span with the given `operationName` and
// incorporate the given StartSpanOption `opts`. (Note that `opts` borrows
// from the "functional options" pattern, per
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
//
// A Span with no SpanReference options (e.g., opentracing.ChildOf() or
// opentracing.FollowsFrom()) becomes the root of its own trace.
//
// Examples:
//
//     var tracer opentracing.Tracer = ...
//
//     // The root-span case:
//     sp := tracer.StartSpan("GetFeed")
//
//     // The vanilla child span case:
//     sp := tracer.StartSpan(
//         "GetFeed",
//         opentracing.ChildOf(parentSpan.Context()))
//
//     // All the bells and whistles:
//     sp := tracer.StartSpan(
//         "GetFeed",
//         opentracing.ChildOf(parentSpan.Context()),
//         opentracing.Tag{"user_agent", loggedReq.UserAgent},
//         opentracing.StartTime(loggedReq.Timestamp),
//     )
//
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	var (
		now = time.Now().UnixNano()
		id  = strconv.FormatInt(newID(now), 10)
	)
	sp := &Span{
		TraceID:   id,
		ParentID:  "0",
		SpanID:    id,
		StartTime: now,
	}

	return sp
}

// Inject() takes the `sm` SpanContext instance and injects it for
// propagation within `carrier`. The actual type of `carrier` depends on
// the value of `format`.
//
// OpenTracing defines a common set of `format` values (see BuiltinFormat),
// and each has an expected carrier type.
//
// Other packages may declare their own `format` values, much like the keys
// used by `context.Context` (see https://godoc.org/context#WithValue).
//
// Example usage (sans error handling):
//
//     carrier := opentracing.HTTPHeadersCarrier(httpReq.Header)
//     err := tracer.Inject(
//         span.Context(),
//         opentracing.HTTPHeaders,
//         carrier)
//
// NOTE: All opentracing.Tracer implementations MUST support all
// BuiltinFormats.
//
// Implementations may return opentracing.ErrUnsupportedFormat if `format`
// is not supported by (or not known by) the implementation.
//
// Implementations may return opentracing.ErrInvalidCarrier or any other
// implementation-specific error if the format is supported but injection
// fails anyway.
//
// See Tracer.Extract().
func (t *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return nil
}

// Extract() returns a SpanContext instance given `format` and `carrier`.
//
// OpenTracing defines a common set of `format` values (see BuiltinFormat),
// and each has an expected carrier type.
//
// Other packages may declare their own `format` values, much like the keys
// used by `context.Context` (see
// https://godoc.org/golang.org/x/net/context#WithValue).
//
// Example usage (with StartSpan):
//
//
//     carrier := opentracing.HTTPHeadersCarrier(httpReq.Header)
//     clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
//
//     // ... assuming the ultimate goal here is to resume the trace with a
//     // server-side Span:
//     var serverSpan opentracing.Span
//     if err == nil {
//         span = tracer.StartSpan(
//             rpcMethodName, ext.RPCServerOption(clientContext))
//     } else {
//         span = tracer.StartSpan(rpcMethodName)
//     }
//
//
// NOTE: All opentracing.Tracer implementations MUST support all
// BuiltinFormats.
//
// Return values:
//  - A successful Extract returns a SpanContext instance and a nil error
//  - If there was simply no SpanContext to extract in `carrier`, Extract()
//    returns (nil, opentracing.ErrSpanContextNotFound)
//  - If `format` is unsupported or unrecognized, Extract() returns (nil,
//    opentracing.ErrUnsupportedFormat)
//  - If there are more fundamental problems with the `carrier` object,
//    Extract() may return opentracing.ErrInvalidCarrier,
//    opentracing.ErrSpanContextCorrupted, or implementation-specific
//    errors.
//
// See Tracer.Inject().
func (t *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return nil, nil
}

func (t *Tracer) Start() {

}

func (t *Tracer) Flush() {

}

func (t *Tracer) Close() {

}

func newID(start int64) int64 {
	return rand.Int63() ^ start
}

func getEnvPairs() map[string]string {
	l := os.Environ()
	m := make(map[string]string)
	for _, v := range l {
		if strings.HasPrefix(v, Prefix) {
			if i := strings.IndexByte(v, '='); i != -1 {
				m[v[:i]] = v[i+1:]
			}
		}
	}

	return m
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
