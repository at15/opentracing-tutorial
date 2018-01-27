package main

// copied from my lesson 1 code, not official solution

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	tlog "github.com/opentracing/opentracing-go/log"
	jaconfig "github.com/uber/jaeger-client-go/config"
)

type logger struct {
}

func (l *logger) Error(msg string) {
	log.Printf("[mylog] Error: %s\n", msg)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	log.Print("[mylog] Info:" + fmt.Sprintf(msg, args...))
}

func initJager(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaconfig.Configuration{
		Sampler: &jaconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaconfig.ReporterConfig{
			LogSpans: true,
		},
	}
	lggr := &logger{}
	tracer, closer, err := cfg.New(service, jaconfig.Logger(lggr))
	if err != nil {
		log.Fatal(err)
	}
	// FIXME: we have to set the global tracer to make `StartSpanFromContext()` work, is there better way?
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func formatString(rootSpan opentracing.Span, to string) string {
	// NOTE: this will start a new trace instead of reusing existing one ....
	// span := rootSpan.Tracer().StartSpan("format")
	span := rootSpan.Tracer().StartSpan("format", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()

	s := fmt.Sprintf("Hello, %s!", to)
	// log is for event inside a span
	span.LogFields(
		tlog.String("event", "string-format"),
		tlog.String("value", to),
	)
	return s
}

func formatStringContext(ctx context.Context, to string) string {
	// FIXME: it seems StartSpanFromContext is using GlobalTracer?
	// TODO: what's the use of second context?
	span, _ := opentracing.StartSpanFromContext(ctx, "format")
	defer span.Finish()

	s := fmt.Sprintf("Hello, %s!", to)
	// log is for event inside a span
	span.LogFields(
		tlog.String("event", "string-format"),
		tlog.String("value", to),
	)
	return s
}

func println(rootSpan opentracing.Span, str string) {
	span := rootSpan.Tracer().StartSpan("print", opentracing.ChildOf(rootSpan.Context()))
	defer span.Finish()

	fmt.Println(str)
	span.LogKV("event", "println")
}

func printlnContext(ctx context.Context, str string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "print")
	defer span.Finish()

	fmt.Println(str)
	span.LogKV("event", "println")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need one argument")
	}
	to := os.Args[1]
	// str := os.Args[2]
	// fmt.Println(str)
	// NOTE: default global tracer is a nop
	// tracer := opentracing.GlobalTracer()
	tracer, closer := initJager("hello-s")
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", to) // tag is for filter span in the future
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	// str := formatString(span, to)
	// println(span, str)
	str := formatStringContext(ctx, to)
	printlnContext(ctx, str)
	span.Finish()
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}
}
