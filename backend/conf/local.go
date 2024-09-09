// nolint:gochecknoglobals
package conf

import (
	"dish_as_a_service/bot"
	"github.com/Falokut/go-kit/config"
	"time"
)

type LocalConfig struct {
	App struct {
		Debug       bool   `yaml:"debug" env:"APP_DEBUG"`
		AdminSecret string `yaml:"admin_secret" env:"APP_ADMIN_SECRET"`
	} `yaml:"app"`
	HealthcheckPort uint32          `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT" env-default:"8081"`
	Bot             bot.Config      `yaml:"tg_bot"`
	DB              config.Database `yaml:"db"`
	Listen          config.Listen   `yaml:"listen"`
	Images          struct {
		Addr          string `yaml:"addr" env:"IMAGES_SERVICE_ADDR"`
		BaseImagePath string `yaml:"base_image_path" env:"BASE_IMAGE_PATH"`
	} `yaml:"images"`
	Payment struct {
		ExpirationDelay time.Duration `yaml:"expiration_delay" env:"PAYMENT_EXPIRATION_DELAY"`
	} `yaml:"payment"`
}
