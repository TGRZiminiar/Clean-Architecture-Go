package config

import (
	"log"
	"os"
	"strconv"

	godotenv "github.com/joho/godotenv"
)

type (
	Config struct {
		App App
		Db  Db
		Jwt Jwt
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

	Jwt struct {
		AccessSecretKey  string
		RefreshSecretKey string
		ApiSecretKey     string
		PrivateKeyPem    string
		PublicKeyPem     string
		AccessDuration   int64
		RefreshDuration  int64
		ApiDuration      int64
	}

	Db struct {
		Url string
	}
)

func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal("Error loading .env file : %s", err.Error())
	}

	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},

		Db: Db{
			Url: os.Getenv("DB_URL"),
		},
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			PrivateKeyPem:    os.Getenv("PrivateKeyPem"),
			PublicKeyPem:     os.Getenv("PublicKeyPem"),

			AccessDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_ACCESS_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading access duration file")
				}
				return result
			}(),

			RefreshDuration: func() int64 {
				result, err := strconv.ParseInt(os.Getenv("JWT_REFRESH_DURATION"), 10, 64)
				if err != nil {
					log.Fatal("Error loading refresh duration file")
				}
				return result
			}(),
		},
	}

}
