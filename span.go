package uniottrans

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

var _ opentracing.Span = &Span{}

// Sets the end timestamp and finalizes Span state.
//
// With the exception of calls to Context() (which are always allowed),
// Finish() must be the last call made to any span instance, and to do
// otherwise leads to undefined behavior.
func (sp *Span) Finish() {}

// FinishWithOptions is like Finish() but with explicit control over
// timestamps and log data.
func (sp *Span) FinishWithOptions(opts opentracing.FinishOptions) {}

// Context() yields the SpanContext for this Span. Note that the return
// value of Context() is still valid after a call to Span.Finish(), as is
// a call to Span.Context() after a call to Span.Finish().
func (sp *Span) Context() opentracing.SpanContext {
	return nil
}

// Sets or changes the operation name.
//
// Returns a reference to this Span for chaining.
func (sp *Span) SetOperationName(operationName string) opentracing.Span {
	return sp
}

// Adds a tag to the span.
//
// If there is a pre-existing tag set for `key`, it is overwritten.
//
// Tag values can be numeric types, strings, or bools. The behavior of
// other tag value types is undefined at the OpenTracing level. If a
// tracing system does not know how to handle a particular value type, it
// may ignore the tag, but shall not panic.
//
// Returns a reference to this Span for chaining.
func (sp *Span) SetTag(key string, value interface{}) opentracing.Span {
	return sp
}

// LogFields is an efficient and type-checked way to record key:value
// logging data about a Span, though the programming interface is a little
// more verbose than LogKV(). Here's an example:
//
//    span.LogFields(
//        log.String("event", "soft error"),
//        log.String("type", "cache timeout"),
//        log.Int("waited.millis", 1500))
//
// Also see Span.FinishWithOptions() and FinishOptions.BulkLogData.
func (sp *Span) LogFields(fields ...log.Field) {}

// LogKV is a concise, readable way to record key:value logging data about
// a Span, though unfortunately this also makes it less efficient and less
// type-safe than LogFields(). Here's an example:
//
//    span.LogKV(
//        "event", "soft error",
//        "type", "cache timeout",
//        "waited.millis", 1500)
//
// For LogKV (as opposed to LogFields()), the parameters must appear as
// key-value pairs, like
//
//    span.LogKV(key1, val1, key2, val2, key3, val3, ...)
//
// The keys must all be strings. The values may be strings, numeric types,
// bools, Go error instances, or arbitrary structs.
//
// (Note to implementors: consider the log.InterleavedKVToFields() helper)
func (sp *Span) LogKV(alternatingKeyValues ...interface{}) {}

// SetBaggageItem sets a key:value pair on this Span and its SpanContext
// that also propagates to descendants of this Span.
//
// SetBaggageItem() enables powerful functionality given a full-stack
// opentracing integration (e.g., arbitrary application data from a mobile
// app can make it, transparently, all the way into the depths of a storage
// system), and with it some powerful costs: use this feature with care.
//
// IMPORTANT NOTE #1: SetBaggageItem() will only propagate baggage items to
// *future* causal descendants of the associated Span.
//
// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and
// value is copied into every local *and remote* child of the associated
// Span, and that can add up to a lot of network and cpu overhead.
//
// Returns a reference to this Span for chaining.
func (sp *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	return sp
}

// Gets the value for a baggage item given its key. Returns the empty string
// if the value isn't found in this Span.
func (sp *Span) BaggageItem(restrictedKey string) string {
	return ""
}

// Provides access to the Tracer that created this Span.
func (sp *Span) Tracer() opentracing.Tracer {
	return nil
}

// Deprecated: use LogFields or LogKV
func (sp *Span) LogEvent(event string) {}

// Deprecated: use LogFields or LogKV
func (sp *Span) LogEventWithPayload(event string, payload interface{}) {}

// Deprecated: use LogFields or LogKV
func (sp *Span) Log(data opentracing.LogData) {}
