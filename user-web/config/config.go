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
	CaptchaInfo    CaptchaConfig    `mapstructure:"captcha"`
	RedisInfo      RedisConfig      `mapstructure:"redis"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"SigningKey"`
}

type CaptchaConfig struct {
	Type          string `mapstructure:"type"`
	SourceChinese string `mapstructure:"sourceChinese"`
	EnableCaptcha bool   `mapstructure:"enableCaptcha"`
}

type RedisConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	DB             int    `mapstructure:"db"`
	ExpirationTime string `mapstructure:"expirationTime"`
}
