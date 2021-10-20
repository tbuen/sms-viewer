package main

import (
	"os"

	"github.com/tbuen/sms-viewer/internal/app"
)

func main() {
	res := app.Run()
	os.Exit(res)
}
