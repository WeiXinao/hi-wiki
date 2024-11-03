package config

type AppConfig struct {
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	Minio MinioConfig `mapstructure:"minio"`
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecureAccessKey string
	UseSSL          bool
}
