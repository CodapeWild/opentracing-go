package optcgo

// Prefix for universal-opentracing-transformer key in k-v paris
const Prefix = "uni-ot-"

// Service key
const (
	ServiceNameKey    = "uni-ot-service"
	DataSourceKey     = "uni-ot-data-source"
	SamplePriorityKey = "uni-ot-smaple-priority"
	SampleRatioKey    = "uni-ot-sample-ratio"
	ExternalTraceID   = "uni-ot-external-trace-id"
)

// Trace key
const (
	TraceIDKey  = "uni-ot-trace-id"
	ParentIDKey = "uni-ot-parent-id"
)

type FormatExternalTraceID func(tid interface{}) int64
