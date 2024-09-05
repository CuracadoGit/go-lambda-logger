// Package logger provides basis logging functions for Lambdas written in Golang.
//
// The output is compatible with AWS structured JSON logging.
// All log records will include Lambda requestIDs and X-Ray traceIDs if present in the context.
//
// Use `go get` to include the module in your code.
//
//	go get github.com/CuracadoGit/go-lambda-logger
package logger // import "github.com/CuracadoGit/go-lambda-logger"
