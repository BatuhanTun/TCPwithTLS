package lib

import (
	"os"

	"github.com/astaxie/beego"
	"github.com/joho/godotenv"
)

func ConfigGet(config string) string {
	if err := godotenv.Load(); err != nil {
		return beego.AppConfig.String(config)
	}
	return os.Getenv(config)
}
