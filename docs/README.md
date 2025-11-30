# Pretty handler for structured Logging with slog

The handler is ready to use:
```go
package main

import (
	"log/slog"

	logger "gopkg.in/slog-handler.v1"
)

func main() {
	slog.Info("test", "example 1", map[int]string{
		0: "test",
	})
	slog.Info("example of raw output", "raw;example 2", map[int]string{
		0: "test",
	})
	slog.Error("example of raw output", "raw;example 3", map[int]string{
		0: "test",
	})
	slog.Warn("test 2", "example 4", map[int]string{
		0: "test",
	})
}

func init() {
	logger.SetLogger(
		logger.Options{
			AddSource: true,
			Format:    "text",
			Level:     "debug",
			Pretty:    true,
		},
	)
}
```

Output:
![handler output](output.png?raw=true)

## NullHandler

The `NullHandler` is a special handler that discards all log records. It's useful for:
- Testing where you don't want log output
- Disabling logging in production without code changes
- Benchmarking code without logging overhead

### Usage

```go
package main

import (
	"log/slog"

	logger "gopkg.in/slog-handler.v1"
)

func main() {
	// Create logger with NullHandler
	log := logger.NewLogger(logger.Options{
		Null: true, // Enable NullHandler
	})

	// These logs will be discarded
	log.Info("this will not be logged")
	log.Error("this will also be discarded")
	log.With("key", "value").Debug("nothing here")
}
```

You can also set it as the global logger:

```go
func init() {
	logger.SetGlobalLogger(logger.Options{
		Null: true, // All slog calls will be discarded
	})
}
```
