# ilog.go

[![license](https://img.shields.io/github/license/kunitsucom/ilog.go)](LICENSE)
[![pkg](https://pkg.go.dev/badge/github.com/kunitsucom/ilog.go)](https://pkg.go.dev/github.com/kunitsucom/ilog.go)
[![goreportcard](https://goreportcard.com/badge/github.com/kunitsucom/ilog.go)](https://goreportcard.com/report/github.com/kunitsucom/ilog.go)
[![workflow](https://github.com/kunitsucom/ilog.go/workflows/go-lint/badge.svg)](https://github.com/kunitsucom/ilog.go/tree/main)
[![workflow](https://github.com/kunitsucom/ilog.go/workflows/go-test/badge.svg)](https://github.com/kunitsucom/ilog.go/tree/main)
[![codecov](https://codecov.io/gh/kunitsucom/ilog.go/branch/main/graph/badge.svg?token=4UML9FB7BX)](https://codecov.io/gh/kunitsucom/ilog.go)
[![sourcegraph](https://sourcegraph.com/github.com/kunitsucom/ilog.go/-/badge.svg)](https://sourcegraph.com/github.com/kunitsucom/ilog.go)

`ilog.go` is a simple logging interface library for Go. By defining only the logger interface `ilog.Logger`, users can easily swap out the underlying logging implementation without changing their application code.

## Features

- **Flexibility**: Allows users to choose or switch between different logging implementations.
- **Simplicity**: Only defines the logging interface, leaving the choice of logging implementation to the user.

## Interface

`ilog.Logger` interface is defined here:

- [ilog.go](ilog.go)

## Reference Implementations

We provide reference implementations for a custom logger and some of the popular logging libraries:

- [ilog_default_implementation.go](ilog_default_implementation.go): custom logger implementation (Default. You can use this as a blueprint to integrate other loggers)
- [implementations/zap/zap.go](implementations/zap/zap.go): implementation for [go.uber.org/zap](https://github.com/uber-go/zap)
- [implementations/zerolog/zerolog.go](implementations/zerolog/zerolog.go): implementation for [github.com/rs/zerolog](https://github.com/rs/zerolog)

## Usage

First, go get `ilog.go` in your Go application:

```bash
go get -u github.com/kunitsucom/ilog.go
```

Then, define a variable of the `ilog.Logger` interface type and initialize it with your chosen logger implementation.

 For example, if using the default implementation:

```go
import (
    "github.com/kunitsucom/ilog.go"
)

func main() {
    l := ilog.NewBuilder(ilog.DebugLevel, os.Stdout).Build()
}
```

if zap:

```bash
go get -u github.com/kunitsucom/ilog.go/implementations/zap
```

```go
import (
    "github.com/kunitsucom/ilog.go"
    "go.uber.org/zap"
    ilogzap "github.com/kunitsucom/ilog.go/implementations/zap"
)

func main() {
    l := ilogzap.New(ilog.DebugLevel, zap.NewProduction())
}
```

if zerolog:

```bash
go get -u github.com/kunitsucom/ilog.go/implementations/zerolog
```

```go
import (
    "github.com/kunitsucom/ilog.go"
    ilogzerolog "github.com/kunitsucom/ilog.go/implementations/zerolog"
    "github.com/rs/zerolog"
)

func main() {
    l := ilogzerolog.New(ilog.DebugLevel, zerolog.New(os.Stdout))
}
```

Now, you can use the `l` as ilog.Logger for logging in your application:

```go
func main() {
    l := ... // your chosen logger implementation

    l.String("key", "value").Infof("This is an info message")
    l.Err(err).Errorf("This is an error message")
}
```

If you wish to switch to another logger, simply change the initialization of the `l` variable.

## Implementing a Custom Logger

If the provided reference implementations do not meet your requirements, you can easily implement the `Logger` interface with your desired logging package. Ensure that your custom logger adheres to the methods defined in the `ilog.go` interface.

## License

[here.](LICENSE)
