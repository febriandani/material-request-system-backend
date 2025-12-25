package config

import "github.com/spf13/viper"

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
	JWT 	Secret
}

type AppConfig struct {
	Port string
}

type AuthConfig struct {
	Master struct {
		Username string
		Password string
	}
}

type Secret struct {
	Secret string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.SetDefault("app.port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		App: AppConfig{
			Port: viper.GetString("app.port"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		Auth: AuthConfig{
			Master: struct {
				Username string
				Password string
			}{
				Username: viper.GetString("auth.master.username"),
				Password: viper.GetString("auth.master.password"),
			},
		},
		JWT: Secret{
			Secret: viper.GetString("jwt.secret"),
		},
	}
}
