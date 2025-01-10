package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DBUSER     string
	DBPASSWORD string
	DBPORT     string
	DBHOST     string
	DBNAME     string
	SSLMODE string
}

func ConfigEnv() *Env {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("/home/athul/Documents/Clean-Arch-User-api/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error loading .env:", err)
	}
	fmt.Println("Loaded Config:", viper.AllSettings())

	var env Env
	env.DBUSER = viper.GetString("user")
	env.DBPASSWORD = viper.GetString("password")
	env.DBPORT = viper.GetString("port")
	env.DBHOST = viper.GetString("host")
	env.DBNAME = viper.GetString("dbname")
	env.SSLMODE = viper.GetString("sslmode")
	return &env
}
