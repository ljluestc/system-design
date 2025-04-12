package monitoring

import (
    "log"
    "os"
)

// Setup initializes logging
func Setup() {
    log.SetPrefix("CDN: ")
    log.SetOutput(os.Stdout)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("Logging initialized")
}

// Expand with log rotation
// ...