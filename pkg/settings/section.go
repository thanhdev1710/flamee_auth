package settings

type Config struct {
	Port                      string `mapstructure:"PORT"`
	JwtSecret                 string `mapstructure:"JWT_SECRET"`
	JwtExpirationTimeDefault  string `mapstructure:"JWT_EXPIRATION_TIME_DEFAULT"`
	JwtExpirationTimeRemember string `mapstructure:"JWT_EXPIRATION_TIME_REMEMBER"`
	Postgre                   Postgre
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
