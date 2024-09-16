package config

type UserServerConfig struct {
	Host string `mapstructure:"targetHost"`
	Port int    `mapstructure:"targetPort"`
}

type ServerConfig struct {
	ServerConfig   string           `mapstructure:"name"`
	ServerPort     int              `mapstructure:"serverPort"`
	UserServerInfo UserServerConfig `mapstructure:"userSrv"`
	JWTInfo        JwtConfig        `mapstructure:"jwt"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"SigningKey"`
}
