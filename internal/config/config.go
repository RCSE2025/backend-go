package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
}

type EmailConfig struct {
	AuthEmail    string `env:"SMTP_EMAIL" env-required:"true"`
	AuthPassword string `env:"SMTP_PASSWORD" env-required:"true"`
	Host         string `env:"SMTP_HOST" env-required:"true"`
	Port         int    `env:"SMTP_PORT" env-required:"true"`
	From         string `env:"MAIL_FROM" env-required:"true"`
}

type Config struct {
	Port       string `env:"PORT"           env-default:"80"`
	Host       string `env:"HOST"           env-default:"0.0.0.0"`
	Version    string `env:"VERSION"        env-default:"1"`
	Production bool   `env:"PRODUCTION"     env-default:"true"`
	Database   DatabaseConfig
	Email      EmailConfig
}

var (
	config Config    //nolint:gochecknoglobals,lll // Global config is initialized once and accessed throughout the application.
	once   sync.Once //nolint:gochecknoglobals,lll // Ensures the config is initialized only once, which requires a global sync.Once.
)

func Get() *Config {
	once.Do(func() {
		err := godotenv.Load()

		if err != nil {
			log.Println("error loading .env file")
		}
		err = cleanenv.ReadEnv(&config)
		if err != nil {
			panic(fmt.Sprintf("Failed to get config: %s", err))
		}
	})

	return &config
}

func (c *Config) GetDSN() string {
	//dsn := "host=localhost user=postgres password=postgres dbname=hack-2025-backend-go port=5432 sslmode=disable TimeZone=Europe/Moscow"
	//
	//hostPort := net.JoinHostPort(c.MongoHost, c.MongoPort)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow",
		c.Database.Host, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.Port)
}
