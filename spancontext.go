package uniottrans

import "context"

type spctxkey struct{}

func (spctx *SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	for k, v := range spctx.Meta {
		handler(k, v)
	}
}

func SpanContextFromContext(ctx context.Context) *SpanContext {
	if ctx != nil {
		if v := ctx.Value(spctxkey{}); v != nil {
			switch t := v.(type) {
			case *SpanContext:
				return t
			case SpanContext:
				return &t
			}
		}
	}

	return nil
}
