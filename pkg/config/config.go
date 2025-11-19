package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database       DatabaseConfig       `mapstructure:"database"`
	Server         ServerConfig         `mapstructure:"server"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	App            AppConfig            `mapstructure:"app"`
	Authentication AuthenticationConfig `mapstructure:"authentication"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expiry time.Duration `mapstructure:"expiry"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type AuthenticationConfig struct {
	EncryptKey           string        `mapstructure:"encrypt_key"`
	AccessSecretKey      string        `mapstructure:"access_secret_key"`
	RefreshSecretKey     string        `mapstructure:"refresh_secret_key"`
	Issuer               string        `mapstructure:"issuer"`
	AccessTokenExpiry    time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry   time.Duration `mapstructure:"refresh_token_expiry"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./pkg/config/files")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.Timezone)
}