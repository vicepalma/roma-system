package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Port           string `env:"PORT" envDefault:"8080"`
	DSN            string `env:"DATABASE_DSN,required"` // postgres://roma:roma@localhost:5432/roma?sslmode=disable
	DefaultTZ      string `env:"DEFAULT_TZ" envDefault:"America/Santiago"`
	AccessTTLMin   int    `env:"ACCESS_TTL_MIN" envDefault:"15"`
	RefreshTTLDays int    `env:"REFRESH_TTL_DAYS" envDefault:"7"`
	CorsOrigins    string `env:"CORS_ORIGINS" envDefault:"*"`
	JWTSecret      string `env:"JWT_SECRET,required"`
}

func Load() (Config, error) {
	var c Config
	err := env.Parse(&c)
	return c, err
}
