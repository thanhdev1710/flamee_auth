package settings

type Config struct {
	Port                      string `mapstructure:"PORT"`
	Env                       string `mapstructure:"ENV"`
	ApiKey                    string `mapstructure:"API_KEY"`
	JwtSecret                 string `mapstructure:"JWT_SECRET"`
	JwtExpirationTimeDefault  string `mapstructure:"JWT_EXPIRATION_TIME_DEFAULT"`
	JwtExpirationTimeRemember string `mapstructure:"JWT_EXPIRATION_TIME_REMEMBER"`
	Postgre                   Postgre
	Email                     Email
	Logger                    Logger
}

type Url struct {
	UrlPostService string `mapstructure:"URL_POST_SERVICE"`
}

type Logger struct {
	Level      string `mapstructure:"LOG_LEVEL"`
	Filename   string `mapstructure:"LOG_FILE"`
	MaxSize    int    `mapstructure:"LOG_MAXSIZE"`
	MaxBackups int    `mapstructure:"LOG_MAXBACKUPS"`
	MaxAge     int    `mapstructure:"LOG_MAXAGE"`
	Compress   bool   `mapstructure:"LOG_COMPRESS"`
}

type Email struct {
	Username string `mapstructure:"EMAIL_FROM"`
	Password string `mapstructure:"EMAIL_PASSWORD"`
	Host     string `mapstructure:"EMAIL_SMTPHOST"`
	Port     string `mapstructure:"EMAIL_SMTPPORT"`
	Secret   string `mapstructure:"EMAIL_SECRET"`
}

type Postgre struct {
	Host            string `mapstructure:"HOST_DB"`
	Port            string `mapstructure:"PORT_DB"`
	Username        string `mapstructure:"USERNAME_DB"`
	Password        string `mapstructure:"PASSWORD_DB"`
	Database        string `mapstructure:"DATABASE_DB"`
	SslMode         string `mapstructure:"SSL_MODE_DB"`
	ConnMaxIdleTime int    `mapstructure:"CONN_MAX_IDLE_TIME_DB"`
	ConnMaxLifeTime int    `mapstructure:"CONN_MAX_LIFE_TIME"`
	ConnMaxOpen     int    `mapstructure:"CONN_MAX_OPEN"`
}
