package uniottrans

type SpanContext struct {
}

func (spctx *SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {

}
