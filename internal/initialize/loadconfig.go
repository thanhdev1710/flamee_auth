package initialize

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/thanhdev1710/flamee_auth/global"
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, loading from OS ENV")
	}

	viper.BindEnv("PORT")
	viper.BindEnv("API_KEY")
	viper.BindEnv("ENV")
	viper.BindEnv("DOMAIN")

	viper.BindEnv("URL_FRONT_END")
	viper.BindEnv("URL_POST_SERVICE")
	viper.BindEnv("URL_USER_SERVICE")

	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRATION_TIME_DEFAULT")
	viper.BindEnv("JWT_EXPIRATION_TIME_REMEMBER")

	viper.BindEnv("EMAIL_FROM")
	viper.BindEnv("EMAIL_PASSWORD")
	viper.BindEnv("EMAIL_SMTPHOST")

	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("LOG_FILE")
	viper.BindEnv("LOG_MAXSIZE")
	viper.BindEnv("LOG_MAXBACKUPS")
	viper.BindEnv("LOG_MAXAGE")
	viper.BindEnv("LOG_COMPRESS")

	viper.BindEnv("EMAIL_SMTPPORT")
	viper.BindEnv("EMAIL_SECRET")

	viper.BindEnv("HOST_DB")
	viper.BindEnv("PORT_DB")
	viper.BindEnv("USERNAME_DB")
	viper.BindEnv("PASSWORD_DB")
	viper.BindEnv("DATABASE_DB")
	viper.BindEnv("SSL_MODE_DB")
	viper.BindEnv("CONN_MAX_IDLE_TIME_DB")
	viper.BindEnv("CONN_MAX_OPEN")
	viper.BindEnv("CONN_MAX_LIFE_TIME")

	viper.BindEnv("NATS_HOST")
	viper.BindEnv("NATS_PORT")
	viper.BindEnv("NATS_MAX_RECONNECTS")
	viper.BindEnv("NATS_RECONNECT_WAIT")

	err = viper.Unmarshal(&global.Config)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	err = viper.Unmarshal(&global.Config.Postgre)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	err = viper.Unmarshal(&global.Config.Email)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	err = viper.Unmarshal(&global.Config.Logger)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	err = viper.Unmarshal(&global.Url)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %v", err)
	}
	err = viper.Unmarshal(&global.Config.Nats)
	if err != nil {
		log.Fatalf("Cannot unmarshal NATS config: %v", err)
	}
}
