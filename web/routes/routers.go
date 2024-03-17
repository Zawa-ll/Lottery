package routes

import (
	"github.com/Zawa-ll/raffle/bootstrap"
	"github.com/Zawa-ll/raffle/services"
	"github.com/Zawa-ll/raffle/web/controllers"
	"github.com/kataras/iris/mvc"
)

func Configure(b *bootstrap.Bootstrapper) {
	userService := services.NewUserService()
	giftService := services.NewGiftService()
	codeService := services.NewCodeService()
	resultService := services.NewResultService()
	userdayService := services.NewUserdayService()
	blackipService := services.NewBlackipService()

	index := mvc.New(b.Party("/"))

	// register all services in index controller
	index.Register(
		userService,
		giftService,
		codeService,
		resultService,
		userdayService,
		blackipService)
	index.Handle(new(controllers.IndexController))
}
