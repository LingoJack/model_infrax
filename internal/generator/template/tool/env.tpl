package tool

import (
	"github.com/joho/godotenv"
)

func LoadEnv(dotEnvFilePath string) (err error) {
	err = godotenv.Load(dotEnvFilePath)
	return
}
