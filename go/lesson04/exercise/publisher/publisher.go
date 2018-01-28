package main

import (
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func main() {
	tracer, closer := tracing.Init("publisher")
	defer closer.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
		println(helloStr)

		span.LogFields(
			otlog.String("event", "string-publish"),
			otlog.String("value", helloStr),
		)
	})

	log.Println("listen on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
