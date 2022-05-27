package main

import (
	"github.com/toel-app/template-server/src/app"
)

func main() {
	app.WrapWithGracefulShutdown(app.StartApplication)
}
