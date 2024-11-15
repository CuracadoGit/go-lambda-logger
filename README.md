# Golang Lambda Logger
This module provides logging functions that are compatible with AWS lambda advanced logging controls.

It supports both `Text` format as well as `JSON` format.

## Features

* supports structured logging with key-value data
* automatically adds Lambda requestId to log entries
* automatically adds XRay traceId to log entries

## Configuration
This module does not provide any configuration or special settings.
It uses Lambda settings to configure itself.

## Installation and Usage

Require the module using `import logger "github.com/CuracadoGit/go-lambda-logger"`, then use the alias `logger.xyz` in your code.

This module provides functions corresponding to these log levels:
* Debug
* Info
* Error
* Warning

Example:
```go
package main

import (
	"context"
	logger "github.com/CuracadoGit/go-lambda-logger"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func main() {
	ctx := context.TODO()

	ctx = lambdacontext.NewContext(ctx, &lambdacontext.LambdaContext{AwsRequestID: "8f507cfc-xmpl-4697-b07a-ac58fc914c95"})
	ctx = context.WithValue(ctx, "x-amzn-trace-id", "Root=1-5759e988-bd862e3fe1be46a994272793;Sampled=1")

	user := struct {
		ID    string
		Name  string
		Email string
	}{
		"user-123-abc", "User Name", "user-email@email.test",
	}

	logger.Info(ctx, "user has been created successfully", "userID", user.ID, "user_data", user)
}

```

Using the default Text format, this will in a result like this:
`2024-11-15T10:50:09.01298877+01:00      8f507cfc-xmpl-4697-b07a-ac58fc914c95    INFO    user has been created successfully {"traceId":"1-5759e988-bd862e3fe1be46a994272793","userID":"user-123-abc","user_data":{"Email":"user-email@email.test","ID":"user-123-abc","Name":"User Name"}}`

When JSON formatting is enabled in your Lambda then the following output should be visible:
`{"timestamp":"2024-11-15T10:52:12.212202617+01:00","level":"INFO","message":"user has been created successfully","userID":"user-123-abc","user_data":{"ID":"user-123-abc","Name":"User Name","Email":"user-email@email.test"},"requestId":"8f507cfc-xmpl-4697-b07a-ac58fc914c95","traceId":"1-5759e988-bd862e3fe1be46a994272793"}`

## Restrictions

Currently, this module only supports the same log level as slog does. This causes a mismatch between this module anc AWS.
Notably, these two levels are not supported:
* TRACE
* FATAL