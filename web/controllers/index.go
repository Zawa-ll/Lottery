package controllers

import (
	"fmt"

	"github.com/Zawa-ll/raffle/comm"
	"github.com/Zawa-ll/raffle/models"
	"github.com/Zawa-ll/raffle/services"
	"github.com/kataras/iris"
)

type IndexController struct {
	Ctx             iris.Context
	ServicesUser    services.UserService
	ServicesGift    services.GiftService
	ServicesCode    services.CodeService
	ServicesResult  services.ResultService
	ServicesUserday services.UserdayService
	ServicesBlackip services.BlackipService
}

// http://localhost:8080/
func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to Go Lottery, <a href='/public.index.html'>Button<a/>"
}

// localhost:8080/gifts
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	datalist := c.ServicesGift.GetAll()
	list := make([]models.LtGift, 0)

	for _, data := range datalist {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}

	rs["gifts"] = list
	return rs
}

// localhost:8080/newprize
func (c *IndexController) GetNewprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""

	//TODO:

	return rs
}

func (c *IndexController) GetLogin() {
	uid := comm.Random(100000)
	loginuser := models.ObjLoginuser{
		Uid:      uid,
		Username: fmt.Sprint("admin-%d", uid),
		Now:      comm.NowUnix(),
		Ip:       comm.ClientIP(c.Ctx.Request()),
	}

	comm.SetLoginuser(c.Ctx.ResponseWriter(), &loginuser)
	comm.Redirect(c.Ctx.ResponseWriter(),
		"/public/index.html?from=login")
}

func (c *IndexController) GetLogout() {
	comm.SetLoginuser(c.Ctx.ResponseWriter(), nil)
	comm.Redirect(c.Ctx.ResponseWriter(),
		"/public/index.html?from=logout")
}
