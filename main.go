package main

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	// ...we can set default context logger like this
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &log

	// ... and then access it like this
	zerolog.Ctx(ctx).Info().Msg("registering request handlers")

	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8090", nil)
}

func hello(_ http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

	// ...we need some value that can be used as a correlation id across log entries - this might be an SQS message id or any other
	// identifier that we get passed.
	// We also might generate an UUID and use it for correlation of log messages within a request processing process.
	requestID := uuid.NewString()

	// we probably don't want to pass the request id around and log it every time manually
	// instead we can create a sublogger of the context logger and associate the request id with it

	requestLogger := zerolog.Ctx(ctx).With().Str("requestID", requestID).Logger()

	// ... now everytime we use the request logger the request id will be written to the output
	// but we probably don't want to pass the logger around either, instead we can pass it by context

	requestContext := requestLogger.WithContext(ctx)
	doProcess(requestContext)

}

func doProcess(ctx context.Context) {
	// we just take the logger from the context and assume that it has all necessary infos associated with it (request id, other business ids if any)
	zerolog.Ctx(ctx).Info().Msg("we are processing a request")
}
