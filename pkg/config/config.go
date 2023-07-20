package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	DBUrl              string `mapstructure:"DB_URL"`
	JWTSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	TWILIO_ACCOUNT_SID string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TWILIO_AUTH_TOCKEN string `mapstructure:"TWILIO_AUTH_TOCKEN"`
	TWILIO_SERVICE_ID  string `mapstructure:"TWILIO_SERVICE_ID"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
