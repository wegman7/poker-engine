package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
}

var (
    DEBUG                = true
    ENGINE_LOOP_PAUSE    = 100 * time.Millisecond
    PAUSE_SHORT          = 1 * time.Millisecond
    PAUSE_MEDIUM         = 2 * time.Millisecond
    PAUSE_LONG           = 3 * time.Millisecond
	MYFAKESECRET 		 = os.Getenv("MYFAKESECRET")
)
