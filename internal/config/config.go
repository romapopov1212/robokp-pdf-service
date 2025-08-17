package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Env        string `mapstructure:"env"`
	Database   `mapstructure:"database"`
	HttpServer `mapstructure:"http_server"`
	AWS        AWSConfig `mapstructure:"aws"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type HttpServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type AWSConfig struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	UploadDir       string `mapstructure:"upload_dir"`
	AccessKeyID     string `mapstructure:"access_key_id"`     // опционально
	SecretAccessKey string `mapstructure:"secret_access_key"` // опционально
	EndpointUri     string `mapstructure:"endpoint_uri"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	
	viper.SetConfigFile(path)
	
	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}
	
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("error unmarshal config")
	}
	
	return cfg, nil
}
