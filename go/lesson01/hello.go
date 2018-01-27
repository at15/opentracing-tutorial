package main

import (
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
	return tracer, closer
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need one argument")
	}
	to := os.Args[1]
	// str := os.Args[2]
	// fmt.Println(str)
	// tracer := opentracing.GlobalTracer() // default global tracer is a nop
	tracer, closer := initJager("hello-s")
	span := tracer.StartSpan("say-hello")
	span.SetTag("hello-to", to) // tag is for filter span in the future
	// log is for event inside a span
	span.LogFields(
		tlog.String("event", "string-format"),
	)
	fmt.Printf("Hello, %s!\n", to)
	span.LogKV("event", "println")
	span.Finish()
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}
}
