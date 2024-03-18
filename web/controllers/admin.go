package controllers

import (
	"github.com/Zawa-ll/raffle/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type AdminController struct {
	Ctx             iris.Context
	ServicesUser    services.UserService
	ServicesGift    services.GiftService
	ServicesCode    services.CodeService
	ServicesResult  services.ResultService
	ServicesUserday services.UserdayService
	ServicesBlackip services.BlackipService
}

func (c *AdminController) Get() mvc.Result {
	return mvc.View{
		Name: "admin/index.html", // Default Admin Page
		Data: iris.Map{
			"Title":"Management Backstage",
			"Channel": "",
		},
		Layout: "admin/layout.html"
	}
}
