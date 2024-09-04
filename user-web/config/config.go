package config

type UserServerConfig struct {
	Host string `mapstructure:"target-host"`
	Port int    `mapstructure:"target-port"`
}

type ServerConfig struct {
	ServerConfig   string           `mapstructure:"name"`
	ServerPort     int              `mapstructure:"server-port"`
	UserServerInfo UserServerConfig `mapstructure:"user-srv"`
}
