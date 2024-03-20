package main

import (
	"fmt"

	"github.com/Zawa-ll/raffle/bootstrap"
	"github.com/Zawa-ll/raffle/conf"
	"github.com/Zawa-ll/raffle/web/middleware/identity"
	"github.com/Zawa-ll/raffle/web/routes"
)

var port = 8080

// Initialize Application
func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("lottery_system", "Haydn") // appName, appOwner
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)
	return app
}

func main() {
	if port == 8080 {
		conf.RunningCrontabService = true
	}

	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
