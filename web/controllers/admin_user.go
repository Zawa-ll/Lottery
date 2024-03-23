package controllers

import (
	"fmt"

	"github.com/Zawa-ll/raffle/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc.com/lottery/comm"
	"imooc.com/lottery/models"
)

type AdminUserController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

// GET: /admin/user/
func (c *AdminUserController) Get() mvc.Result {
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""

	datalist := c.ServiceUser.GetAll(page, size)
	total := (page-1)*size + len(datalist) // total number of data

	if len(datalist) >= size {
		total = c.ServiceUser.CountAll()
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/user.html",
		Data: iris.Map{
			"Title":    "Management Back Office",
			"Channel":  "user",
			"Datalist": datalist,
			"Total":    total,
			"Now":      comm.NowUnix(),
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

// GET: /admin/user/black?id=1&time=0
// Updates a user's blacklist duration based on URL parameters `id` and `time`
// then redirects to `/admin/user`
func (c *AdminUserController) GetBlack() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	t := c.Ctx.URLParamIntDefault("time", 0) // time in param
	if err == nil {
		if t > 0 {
			t = t*86400 + comm.NowUnix()
		}
		c.ServiceUser.Update(&models.LtUser{Id: id, Blacktime: t, SysUpdated: comm.NowUnix()},
			[]string{"blacktime"})
	}
	return mvc.Response{
		Path: "/admin/user",
	}
}
