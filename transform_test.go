package uniottrans

import (
	"log"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestOpenTracing(t *testing.T) {
	opentracing.SetGlobalTracer(mocktracer.New())

	span := opentracing.StartSpan("test_start_span")
	time.Sleep(time.Second)
	span.Finish()

	log.Printf("%#v", span)
}
