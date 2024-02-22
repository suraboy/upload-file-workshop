package config

type AppConfig struct {
	Server  ServerConfig  `mapstructure:"server"`
	Storage StorageConfig `mapstructure:"storage"`
	Email   EmailConfig   `mapstructure:"mail"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type StorageConfig struct {
	Services []ServiceConfig `mapstructure:"services"`
}

type ServiceConfig struct {
	Name   string     `mapstructure:"name"`
	Type   string     `mapstructure:"type"`
	Bucket string     `mapstructure:"bucket"`
	Config ConfigData `mapstructure:"config"`
}

type ConfigData struct {
	JSONCredential string `mapstructure:"json_credential"`
}

type EmailConfig struct {
	ApiKey string `mapstructure:"api-key"`
}
