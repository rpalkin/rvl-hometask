package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type (
	Config struct {
		App App `yaml:"app" env-prefix:"RVL_APP_"`
		Db  Db  `yaml:"db" env-prefix:"RVL_DB_"`
	}

	App struct {
		HttpPort    string `yaml:"port" env:"PORT" env-default:"8080"`
		MetricsPort string `yaml:"metrics_port" env:"METRICS_PORT" env-default:"8081"`
		LogLevel    string `yaml:"level" env:"LEVEL" env-default:"info"`
	}

	Db struct {
		Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"PORT" env-default:"5432"`
		Username string `yaml:"username" env:"USERNAME" env-default:"rvl"`
		Password string `yaml:"password" env:"PASSWORD" env-default:"rvl"`
		Dbname   string `yaml:"dbname" env:"DBNAME" env-default:"rvl"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		// ok
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, errors.Wrap(err, "error read env")
	}

	return cfg, nil
}

func (c *Config) GetDbConnString() string {
	return "host=" + c.Db.Host + " port=" + c.Db.Port + " user=" + c.Db.Username + " password=" + c.Db.Password + " dbname=" + c.Db.Dbname + " sslmode=disable"
}
