// nolint:gochecknoglobals
package conf

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Falokut/go-kit/log"

	"dish_as_a_service/bot"

	"github.com/Falokut/go-kit/client/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type DbConfig struct {
	Host        string `yaml:"host" env:"DB_HOST"`
	Port        int    `yaml:"port" env:"DB_PORT"`
	Database    string `yaml:"database" env:"DB_NAME"`
	Username    string `yaml:"username" env:"DB_USERNAME"`
	Password    string `yaml:"password" env:"DB_PASSWORD"`
	Schema      string `yaml:"schema" env:"DB_SCHEMA_NAME"`
	MaxOpenConn int    `yaml:"max_open_conn" env:"DB_MAX_OPEN_CONN"`
}

func (c *DbConfig) Convert() db.Config {
	return db.Config{
		Host:        c.Host,
		Port:        c.Port,
		Database:    c.Database,
		Username:    c.Username,
		Password:    c.Password,
		Schema:      c.Schema,
		MaxOpenConn: c.MaxOpenConn,
	}
}

type LocalConfig struct {
	Listen struct {
		Addr string `yaml:"addr" env:"LISTEN_ADDR"`
	}
	App struct {
		Id          string `yaml:"id" env:"APP_ID"`
		Version     string `yaml:"version" env:"APP_VERSION"`
		Debug       bool   `yaml:"debug" env:"APP_DEBUG"`
		AdminSecret string `yaml:"admin_secret" env:"APP_ADMIN_SECRET"`
	} `yaml:"app"`
	Bot    bot.Config `yaml:"tg_bot"`
	DB     DbConfig   `yaml:"db"`
	Images struct {
		Addr          string `yaml:"addr" env:"IMAGES_SERVICE_ADDR"`
		BaseImagePath string `yaml:"base_image_path" env:"BASE_IMAGE_PATH"`
	} `yaml:"images"`
	Payment struct {
		ExpirationDelay time.Duration `yaml:"expiration_delay" env:"PAYMENT_EXPIRATION_DELAY"`
	} `yaml:"payment"`
	Log struct {
		LogLevel      string `yaml:"level" env:"LOG_LEVEL"`
		ConsoleOutput bool   `yaml:"console_output" env:"LOG_CONSOLE_OUTPUT"`
		Filepath      string `yaml:"filepath" env:"LOG_FILEPATH"`
	} `yaml:"log"`
}

const configsPath = "conf/"

var once sync.Once
var instance *LocalConfig

func GetLocalConfig() *LocalConfig {
	once.Do(func() {
		instance = &LocalConfig{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger, _ := log.NewFromConfig(log.Config{
				Loglevel: log.FatalLevel,
				Output: log.OutputConfig{
					Console: true,
				},
			})
			logger.Fatal(context.Background(), fmt.Sprint(help, err))
		}
	})
	return instance
}
