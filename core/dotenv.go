package core

import (
	"os"
	"path"

	"github.com/joho/godotenv"
)

func LoadDotenv() {
	basePath, _ := os.Getwd()
	envFile := path.Join(basePath, ".env")
	_, err := os.Stat(envFile)

	if err == nil {
		AppLog().Infof(".env file exists, loading..")
		err := godotenv.Load()
		if err != nil {
			AppLog().Errorf(".env file exists, loading..")
		}
	}

}
