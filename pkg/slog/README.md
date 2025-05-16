## 📘 Custom Compatible Version of `slog` for Go <= 1.18

### Background

Go 1.21 introduced a new structured logging package: `log/slog`. It supports modular handlers, leveled logging, and context-aware features. However, in some production environments or legacy systems, you may need similar functionality while using an older version of Go, such as Go 1.18.

### Approach

To make `slog` usable with Go 1.18, we follow these steps:

1. **Copy the source code of `log/slog`** into a custom module directory (e.g., `slog`).
2. **Refactor or remove features that require Go 1.21+**, including:

    * Replacing `any` with `interface{}`.

3. **Customize the `defaultLogger` behavior**:

    * Allow multiple independent logger instances (instead of relying on a global logger).
    * Configurable output format, log levels, and filters per instance.

### Key Refactoring Examples

#### 1. Replace `any` with `interface{}`

```go
// Original slog interface
func (l *Logger) Info(msg string, args ...any) {}

// Refactored for Go 1.18
func (l *Logger) Info(msg string, args ...interface{}) {}
```

#### 2. Replace `init`

```go
// Original slog init
func init() {
    defaultLogger.Store(New(newDefaultHandler(loginternal.DefaultOutput)))
}

// Setup the default logger
func init() {
    if os.Getenv("LOG_ROTATE_ENABLED") == "true" {
        logFileName := "goagent/fi-otel-goagent.log"
        baseLogDirOnPlatform := "/applogs"
        baseLogDirOnLocal := filepath.Join(os.Getenv("HOME"), "applogs")
        baseLogDirOnCurrent := filepath.Join(".", "applogs")
        
        baseDir := baseLogDirOnPlatform
        if _, err := os.Stat(baseDir); os.IsNotExist(err) {
            baseDir = baseLogDirOnLocal
            if _, err := os.Stat(baseDir); os.IsNotExist(err) {
                baseDir = baseLogDirOnCurrent
            }
        }

        appId := os.Getenv("APP_ID")
        var logFilePath string
        if appId == "" {
            logFilePath = filepath.Join(baseDir, logFileName)
        } else {
            logFilePath = filepath.Join(baseDir, appId, logFileName)
        }

        if os.Getenv("LOG_FILE_PATH") != "" {
            logFilePath = os.Getenv("LOG_FILE_PATH")
        }

        maxSize := 10
        if v := os.Getenv("LOG_MAX_SIZE_MB"); v != "" {
            if n, err := strconv.Atoi(v); err == nil {
                maxSize = n
            }
        }

        maxBackups := 5
        if v := os.Getenv("LOG_MAX_BACKUPS"); v != "" {
            if n, err := strconv.Atoi(v); err == nil {
                maxBackups = n
            }
        }

        maxAge := 30
        if v := os.Getenv("LOG_MAX_AGE_DAYS"); v != "" {
            if n, err := strconv.Atoi(v); err == nil {
                maxAge = n
            }
        }

        compress := os.Getenv("LOG_COMPRESS") == "true"

        logFile := &lumberjack.Logger{
            Filename:   logFilePath,
            MaxSize:    maxSize,    // megabytes
            MaxBackups: maxBackups,
            MaxAge:     maxAge,      // days
            Compress:   compress,
        }

        defaultLogger.Store(New(
            NewJSONHandler(logFile, &HandlerOptions{Level: LevelDebug}),
        ))
    }
}
```