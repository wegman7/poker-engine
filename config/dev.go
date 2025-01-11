package config

import "time"

var (
    DEBUG                = true
    ENGINE_LOOP_PAUSE    = 100 * time.Millisecond
    PAUSE_SHORT          = 1 * time.Millisecond
    PAUSE_MEDIUM         = 2 * time.Millisecond
    PAUSE_LONG           = 3 * time.Millisecond
)
