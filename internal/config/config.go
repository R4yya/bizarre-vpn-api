package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env              string     `yaml:"env" env-required:"true"`
	StoragePath      string     `yaml:"storage_path" env-required:"true"`
	TelegramBotToken string     `yaml:"telegram_bot_token" env-required:"true"`
	WebAppUrl        string     `yaml:"web_app_url" env-required:"true"`
	JWT              JWT        `yaml:"jwt" env-required:"true"`
	HttpServer       HttpServer `yaml:"http_server" env-required:"true"`
}

type HttpServer struct {
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type JWT struct {
	AccessSecretKey  string `yaml:"accessSecretKey" env-required:"true"`
	RefreshSecretKey string `yaml:"refreshSecretKey" env-required:"true"`
}

func MustLoadConfig() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
