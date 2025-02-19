package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
    DEBUG bool
    ENGINE_LOOP_PAUSE time.Duration
    PAUSE_SHORT time.Duration
    PAUSE_MEDIUM time.Duration
    PAUSE_LONG time.Duration
    MAX_PLAYERS int
	BACKEND_URL string
}

var AppConfig Config

func Load(env string) error {
	switch env {
	case "dev":
		err := godotenv.Load()
		if err != nil {
			return fmt.Errorf("error loading .env file")
		}
		AppConfig = Config{
			DEBUG: true,
			ENGINE_LOOP_PAUSE: 10 * time.Millisecond,
			PAUSE_SHORT: 1 * time.Millisecond,
			PAUSE_MEDIUM: 2 * time.Millisecond,
			PAUSE_LONG: 5000 * time.Millisecond,
			MAX_PLAYERS: 9,
			BACKEND_URL: os.Getenv("BACKEND_URL"),
		}
	case "prod":
		// prod env vars will be loaded into docker container at runtime
		AppConfig = Config{
			DEBUG: false,
			ENGINE_LOOP_PAUSE: 10 * time.Millisecond,
			PAUSE_SHORT: 1000 * time.Millisecond,
			PAUSE_MEDIUM: 1500 * time.Millisecond,
			PAUSE_LONG: 2000 * time.Millisecond,
			MAX_PLAYERS: 9,
			BACKEND_URL: os.Getenv("BACKEND_URL"),
		}
	default:
		return fmt.Errorf("unknown environment: %s", env)
	}
	return nil
}
