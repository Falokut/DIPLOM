// nolint:gochecknoglobals
package conf

import (
	"time"

	"github.com/Falokut/go-kit/config"
)

type LocalConfig struct {
	App             App                      `yaml:"app"`
	HealthcheckPort uint32                   `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT" env-default:"8081"`
	Bot             config.TelegramBotConfig `yaml:"tg_bot"`
	DB              config.Database          `yaml:"db"`
	Listen          config.Listen            `yaml:"listen"`
	Images          Images                   `yaml:"images"`
	Payment         Payment                  `yaml:"payment"`
	Auth            Auth                     `yaml:"auth"`
}

type App struct {
	AdminSecret string `yaml:"admin_secret" env:"APP_ADMIN_SECRET"`
}

type Images struct {
	BaseServiceUrl string `yaml:"base_service_url" env:"IMAGES_BASE_SERVICE_URL"`
	BaseImagePath  string `yaml:"base_image_path" env:"BASE_IMAGE_PATH"`
}

type Payment struct {
	ExpirationDelay time.Duration `yaml:"expiration_delay" env:"PAYMENT_EXPIRATION_DELAY"`
}

type Auth struct {
	Access                 JwtToken      `yaml:"access"`
	Refresh                JwtToken      `yaml:"refresh"`
	TelegramExpireDuration time.Duration `yaml:"telegram_expire_duration" validate:"required,min=24h"`
}

type JwtToken struct {
	TTL    time.Duration `yaml:"ttl" validate:"required,min=24h"`
	Secret string        `yaml:"secret" validate:"required,min=10"`
}
