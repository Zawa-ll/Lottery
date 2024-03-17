package main

import (
	"fmt"

	"github.com/Zawa-ll/raffle/bootstrap"
	"github.com/Zawa-ll/raffle/web/middleware/identity"
	"github.com/Zawa-ll/raffle/web/routes"
)

var port = 8080

// Initialize Application
func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("go_lottery_system", "Haydn") // appName, appOwner
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
