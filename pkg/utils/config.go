package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Cloudinary          string        `mapstructure:"CLOUDINARY_API_ENV"`
	TokenKey            string        `mapstructure:"TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MongoUserName       string        `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoPassword       string        `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	MongoUri            string        `mapstructure:"MONGODB_LOCAL_URI"`
	Port                string        `mapstructure:"PORT"`
	ProductCol          string        `mapstructure:"PRODUCT_COL"`
	UserCol             string        `mapstructure:"USER_COl"`
	DbName              string        `mapstructure:"DB_NAME"`
	TokenCol            string        `mapstructure:"TOKEN_COL"`
	RedisUri            string        `mapstructure:"REDIS_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
