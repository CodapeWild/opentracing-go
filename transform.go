package uniottrans

import (
	"errors"
)

type TraceNative interface {
	Foreach(handler func(native SpanNative))
}

type SpanNative interface {
	GetTraceID() string
	GetParentID() string
	GetSpanID() string
	GetService() string
	GetOperation() string
	GetMeta() map[string]string
	GetMetrics() map[string]*Numeric
	GetSpanStatus() SpanStatus
	GetStart() int64
	GetEnd() int64
}

type Decoder func(bts []byte) ([]TraceNative, error)

type UniversalTransformer interface {
	SetDecoder(deocder Decoder)
	BeforeBuildSpan(native SpanNative)
	BuildSpan(native SpanNative) *Span
	AfterBuildSpan(native SpanNative, span *Span)
	Transform(bts []byte) (*Traces, error)
}

type RawTransformer struct {
	Decoder
}

func (def *RawTransformer) SetDecoder(decoder Decoder) {
	def.Decoder = decoder
}

func (def *RawTransformer) BeforeBuildSpan(native SpanNative) {}

func (def *RawTransformer) BuildSpan(native SpanNative) *Span {
	return &Span{
		TraceID:   native.GetTraceID(),
		ParentID:  native.GetParentID(),
		SpanID:    native.GetSpanID(),
		Service:   native.GetService(),
		Operation: native.GetOperation(),
		Meta:      native.GetMeta(),
		Metrics:   native.GetMetrics(),
		Status:    native.GetSpanStatus(),
		Start:     native.GetStart(),
		End:       native.GetEnd(),
	}
}

func (def *RawTransformer) AfterBuildSpan(native SpanNative, span *Span) {}

func (def *RawTransformer) Transform(bts []byte, before, after func()) (*Traces, error) {
	if def.Decoder == nil {
		return nil, errors.New("decoder: nil")
	}

	tracesNative, err := def.Decoder(bts)
	if err != nil {
		return nil, err
	}

	var traces *Traces = &Traces{}
	for _, traceNative := range tracesNative {
		var trace *Trace = &Trace{}
		traceNative.Foreach(func(native SpanNative) {
			def.BeforeBuildSpan(native)
			span := def.BuildSpan(native)
			def.AfterBuildSpan(native, span)
			trace.Trace = append(trace.Trace, span)
		})
		traces.Traces = append(traces.Traces, trace)
	}

	return traces, nil
}
