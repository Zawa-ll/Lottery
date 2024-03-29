package routes

import (
	"github.com/Zawa-ll/raffle/bootstrap"
	"github.com/Zawa-ll/raffle/services"
	"github.com/Zawa-ll/raffle/web/controllers"
	"github.com/Zawa-ll/raffle/web/middleware"
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

	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(
		userService,
		giftService,
		codeService,
		resultService,
		userdayService,
		blackipService)
	admin.Handle(new(controllers.AdminController))

	adminGift := admin.Party("/gift")
	adminGift.Register(giftService)
	adminGift.Handle(new(controllers.AdminGiftController))
}
