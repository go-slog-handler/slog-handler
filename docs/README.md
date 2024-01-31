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
