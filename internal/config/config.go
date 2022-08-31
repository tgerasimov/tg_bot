package config

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

const (
	Telegram_token = "TELEGRAM_TOKEN"
	Database_dsn   = "DATABASE_DSN"
	Cfg            = iota >> 1
)

var (
	once sync.Once
	cfg  Config
)

type Config struct {
	sync.Mutex
	vault map[string]string
}

//ParseConfig parse config values
func ParseConfig(ctx context.Context) context.Context {
	var ctxlocal context.Context
	once.Do(func() {
		err := godotenv.Load("internal/config/.config.env")
		if err != nil {
			log.Fatal(err)
		}
		cfg.vault = make(map[string]string)

		cfg.vault[Telegram_token] = os.Getenv(Telegram_token)
		cfg.vault[Database_dsn] = os.Getenv(Database_dsn)

		ctxlocal = context.WithValue(ctx, Cfg, cfg)
	})
	return ctxlocal
}

//GetConfigFromCtx returns config from context
func GetConfigFromCtx(ctx context.Context) (Config, bool) {
	value := ctx.Value(Cfg)
	if j, ok := value.(Config); ok {
		return j, true
	}
	return Config{}, false
}

func (c *Config) GetValue(key string) string {
	c.Lock()
	if value, ok := c.vault[key]; ok {
		return value
	}
	c.Unlock()
	return ""
}
