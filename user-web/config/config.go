package config

type UserServerConfig struct {
	Host        string
	Port        int
	ServiceName string
}

type Cfg struct {
	ServerConfig   string `mapstructure:"serverNme"`
	ServerPort     int    `mapstructure:"serverPort"`
	TimeZone       string `mapstructure:"timeZone"`
	UserServerInfo UserServerConfig
	JWTInfo        JwtConfig     `mapstructure:"jwt"`
	CaptchaInfo    CaptchaConfig `mapstructure:"captcha"`
	RedisInfo      RedisConfig   `mapstructure:"redis"`
	ConsulInfo     ConsulConfig  `mapstructure:"consul"`
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

type ConsulConfig struct {
	Host             string `mapstructure:"host" json:"host"`
	Port             int    `mapstructure:"port" json:"port"`
	Target           string `mapstructure:"target" json:"target"`
	TargetServerName string `mapstructure:"targetServerName" json:"targetServerName"`
}
