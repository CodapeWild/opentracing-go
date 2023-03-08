package uniottrans

func (spctx *SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	for k, v := range spctx.Meta {
		handler(k, v)
	}
}
