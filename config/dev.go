package config

import "time"

var (
    DEBUG                = true
    ENGINE_LOOP_PAUSE    = 10 * time.Millisecond
    PAUSE_SHORT          = 1 * time.Millisecond
    PAUSE_MEDIUM         = 2 * time.Millisecond
    PAUSE_LONG           = 5000 * time.Millisecond
    MAX_PLAYERS          = 9
)
