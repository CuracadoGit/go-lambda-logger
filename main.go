// Package logger
package logger

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"log"
	"log/slog"
	"os"
	"strings"
)

var lg *slog.Logger

type loggerContextKey string
var addedContextDataKey loggerContextKey = "addedContextData"

func init() {
	// log prefixes
	log.SetFlags(0)

	// rename some of the attributes to conform to AWS advanced logging
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{Key: "timestamp", Value: a.Value}
		}
		if a.Key == slog.MessageKey {
			return slog.Attr{Key: "message", Value: a.Value}
		}

		return a
	}

	// this is the default in slog
	level := slog.LevelInfo

	switch os.Getenv("AWS_LAMBDA_LOG_LEVEL") {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	switch os.Getenv("AWS_LAMBDA_LOG_FORMAT") {
	case "JSON":
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replace, Level: level}))
		break
	default:
		lg = slog.New(slog.NewJSONHandler(&jsonToTextWriter{w: os.Stdout}, &slog.HandlerOptions{ReplaceAttr: replace, Level: level}))
		break
	}

}

// WithContextValues will add additional values to the context. These values will be added to all log entries using the context.
func WithContextValues(ctx context.Context, data map[string]interface{}) context.Context {

	return context.WithValue(ctx, addedContextDataKey, data)
}

func handlerForContext(ctx context.Context, args ...any) *slog.Logger {
	h := lg

	if lCtx, has := lambdacontext.FromContext(ctx); has == true {
		args = append(args, "requestId", lCtx.AwsRequestID)
	}

	if traceID, has := ctx.Value("x-amzn-trace-id").(string); has {

		idParts := strings.Split(traceID, ";")
		for _, idPart := range idParts {
			if strings.HasPrefix(idPart, "Root=") {
				args = append(args, "traceId", strings.TrimPrefix(idPart, "Root="))
				break
			}
		}
	}

	if values, has := ctx.Value(addedContextDataKey).(map[string]interface{}); has == true {
		for k, v := range values {
			args = append(args, k)
			args = append(args, v)
		}

	}

	return h.With(args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	handlerForContext(ctx, args...).DebugContext(ctx, msg)
}

func Info(ctx context.Context, msg string, args ...any) {
	handlerForContext(ctx, args...).InfoContext(ctx, msg)
}

func Error(ctx context.Context, msg string, args ...any) {
	handlerForContext(ctx, args...).ErrorContext(ctx, msg)
}

func Warning(ctx context.Context, msg string, args ...any) {
	handlerForContext(ctx, args...).WarnContext(ctx, msg)
}
