package uniottrans

import (
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	DefService       = "UniversalOpenTracingTransform"
	DefFlushBuffer   = 1024
	DefFlushInterval = 3 * time.Second
)

var defGlobalTracer *Tracer

type StartTracerOption func(tracer *Tracer)

func WithSampleRatio(ratio float64) StartTracerOption {
	return func(tracer *Tracer) {
		tracer.sampler = CommonSampler(ratio)
	}
}

func WithGlobalTags(tags map[string]interface{}) StartTracerOption {
	return func(tracer *Tracer) {
		if &tracer.tags == &tags {
			return
		}
		if tracer.tags == nil {
			tracer.tags = make(map[string]interface{})
		}
		for k, v := range tags {
			tracer.tags[k] = v
		}
	}
}

func WithFlushBuffer(size int) StartTracerOption {
	return func(tracer *Tracer) {
		if size <= 0 {
			size = DefFlushBuffer
		}
		tracer.finished = make(chan *Span, size)
	}
}

func WithFlushInterval(d time.Duration) StartTracerOption {
	return func(tracer *Tracer) {
		tracer.flushInterval = d
	}
}

func NewTracer(service string, opts ...StartTracerOption) *Tracer {
	envs := getEnvPairs()
	if s, ok := envs[ServiceNameKey]; ok {
		service = s
	}

	if service == "" {
		service = DefService
	}
	tracer := &Tracer{service: service}
	for i := range opts {
		opts[i](tracer)
	}
	if tracer.finished == nil {
		tracer.finished = make(chan *Span, DefFlushBuffer)
	}
	if tracer.flushInterval <= 0 {
		tracer.flushInterval = DefFlushInterval
	}
	tracer.flush = make(chan struct{})
	tracer.close = make(chan struct{})

	if tracer.sampler == nil {
		if p, ok := envs[SampleRatioKey]; ok {
			if ratio, err := strconv.ParseFloat(p, 10); err == nil {
				tracer.sampler = CommonSampler(ratio)
			}
		}
	}

	return tracer
}

type Tracer struct {
	service       string
	sampler       Sampler
	tags          map[string]interface{}
	finished      chan *Span
	flush         chan struct{}
	flushInterval time.Duration
	close         chan struct{}
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
func (tcr *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	ssopts := &opentracing.StartSpanOptions{}
	for i := range opts {
		opts[i].Apply(ssopts)
	}

	var spctx *SpanContext
	if len(ssopts.References) > 0 {
		spctx = ssopts.References[0].ReferencedContext.(*SpanContext)
	}

	var start int64
	if ssopts.StartTime.IsZero() {
		start = time.Now().UnixNano()
	} else {
		start = ssopts.StartTime.UnixNano()
	}

	sp := &Span{
		Service:   tcr.service,
		StartTime: start,
	}
	sp.SetTags(tcr.tags)
	sp.SetTags(ssopts.Tags)

	if spctx != nil {
		sp.TraceID = spctx.TraceID
		sp.ParentID = spctx.ParentID
		for k, v := range spctx.Meta {
			sp.SetTag(k, v)
		}
		if sp.ParentID == 0 {
			sp.SetTag(SamplePriorityKey, &Numeric_Int32Value{Int32Value: int32(spctx.SamplePriority)})
			sp.SetTag(SampleRatioKey, &Numeric_Doublevalue{Doublevalue: spctx.SampleRatio})
		}
	} else {
		sp.TraceID = newID(sp.StartTime)
		sp.ParentID = 0
		if tcr.sampler != nil {
			sp.SetTag(SamplePriorityKey, &Numeric_Int32Value{Int32Value: int32(SamplePriority_AutoKeep)})
			sp.SetTag(SampleRatioKey, &Numeric_Doublevalue{Doublevalue: tcr.sampler.Ratio()})
		}
	}
	sp.SpanID = newID(time.Now().UnixNano())

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
func (tcr *Tracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
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
func (tcr *Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return nil, nil
}

func (tcr *Tracer) Start() {
	go func() {
		ticker := time.NewTicker(tcr.flushInterval)
		for {
			select {
			case <-tcr.close:
				return
			default:
			}

			select {
			case <-tcr.flush:
			case <-ticker.C:
			case <-tcr.close:
				return
			}
		}
	}()
}

func (tcr *Tracer) Flush() {
	tcr.flush <- struct{}{}
}

func (tcr *Tracer) Close() {
	select {
	case <-tcr.close:
	default:
		close(tcr.close)
	}
}

func (tcr *Tracer) finishSpan(span *Span) error {
	timeout := time.NewTimer(time.Second)
	for {
		select {
		case tcr.finished <- span:
			return nil
		case <-timeout.C:
			return errors.New("finish span timeout")
		default:
			tcr.Flush()
		}
	}
}

func (tcr *Tracer) doFlush() {

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

func newID(start int64) int64 {
	return rand.Int63() ^ start
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
