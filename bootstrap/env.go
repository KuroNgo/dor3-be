package bootstrap

import (
	"github.com/spf13/viper"
	"log"
)

type Database struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBName     string `mapstructure:"DB_NAME"`

	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Prefix string `mapstructure:"prefix"`

	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`

	Secret                string `mapstructure:"secret"`
	AccessTokenExpiresIn  int64  `mapstructure:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int64  `mapstructure:"refreshTokenExpiresIn"`
}

func NewEnv() *Database {
	env := Database{}
	viper.SetConfigFile("app.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file app.env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	} else {
		log.Println("The App is running in deployment env")
	}
	return &env
}
