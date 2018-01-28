package main

import (
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func main() {
	tracer, closer := tracing.Init("formatter")
	defer closer.Close()

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		// FIXME: baggage item is not working ...
		greeting := span.BaggageItem("greeting")
		if greeting == "" {
			greeting = "Hello"
		}
		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("%s, %s!", greeting, helloTo)
		w.Write([]byte(helloStr))

		span.LogFields(
			otlog.String("event", "string-format"),
			otlog.String("value", helloStr),
		)
	})

	log.Println("listen on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
