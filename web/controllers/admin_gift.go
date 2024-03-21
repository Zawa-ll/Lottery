package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Zawa-ll/raffle/comm"
	"github.com/Zawa-ll/raffle/models"
	"github.com/Zawa-ll/raffle/services"
	"github.com/Zawa-ll/raffle/web/utils"
	"github.com/Zawa-ll/raffle/web/viewmodels"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type AdminGiftController struct {
	Ctx             iris.Context
	ServicesUser    services.UserService
	ServicesGift    services.GiftService
	ServicesCode    services.CodeService
	ServicesResult  services.ResultService
	ServicesUserday services.UserdayService
	ServicesBlackip services.BlackipService
}

func (c *AdminGiftController) Get() mvc.Result {
	// TODO: datalist, total
	datalist := c.ServicesGift.GetAll(false) // useCache: false
	total := len(datalist)
	for i, giftInfo := range datalist {
		// Prize Distribution Program Data
		prizedata := make([][2]int, 0)
		err := json.Unmarshal([]byte(giftInfo.PrizeData), &prizedata) // unmarshal each gift into a slice of integer pairs (prizedata)
		if err != nil || len(prizedata) < 1 {
			datalist[i].PrizeData = "[]"
		} else {
			newpd := make([]string, len(prizedata))
			for index, pd := range prizedata { // pd represents prizedata
				ct := comm.FormatFromUnixTime(int64(pd[0])) // ct represents current time
				newpd[index] = fmt.Sprintf("[%s] : %d", ct, pd[1])
			}
			str, err := json.Marshal(newpd)
			if err == nil && len(str) > 0 {
				datalist[i].PrizeData = string(str)
			} else {
				datalist[i].PrizeData = "[]"
			}
		}
	}
	return mvc.View{
		Name: "admin/gift.html", // Default Admin Page
		Data: iris.Map{
			"Title":    "Management Backstage",
			"Channel":  "gift",
			"Datalist": datalist,
			"Total":    total,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) GetEdit() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServicesGift.Get(id, false)
		giftInfo.Id = data.Id
		giftInfo.Title = data.Title
		giftInfo.PrizeNum = data.PrizeNum
		giftInfo.PrizeCode = data.PrizeCode
		giftInfo.PrizeTime = data.PrizeTime
		giftInfo.Img = data.Img
		giftInfo.Displayorder = data.Displayorder
		giftInfo.Gtype = data.Gtype
		giftInfo.Gdata = data.Gdata
		giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
		giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
	}
	// TODO: giftInfo
	return mvc.View{
		Name: "admin/gift.html", // Default Admin Page
		Data: iris.Map{
			"Title":   "Management Backstage",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) PostSave() mvc.Result {
	data := viewmodels.ViewGift{}
	// Read form data into a viewmodels.ViewGift struct (data)
	err := c.Ctx.ReadForm(&data) // Fill in value based on view_gift.go
	if err != nil {
		fmt.Println("admin_gift.PostSave ReadForm error = ", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadForm conversion exception, err= %s", err),
		}
	}
	giftInfo := models.LtGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.Img = data.Img
	giftInfo.Displayorder = data.Displayorder
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	t1, err1 := comm.ParseTime(data.TimeBegin)
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		return mvc.Response{
			Text: fmt.Sprintf("Error in start time or end time format, err1=%s, err2=%s", err1, err2),
		}
	}
	giftInfo.TimeBegin = int(t1.Unix())
	giftInfo.TimeEnd = int(t2.Unix())
	if giftInfo.Id > 0 {
		datainfo := c.ServicesGift.Get(giftInfo.Id, false)
		if datainfo != nil {
			giftInfo.SysUpdated = int(time.Now().Unix())
			giftInfo.SysIp = comm.ClientIP(c.Ctx.Request())
			// Comparison of modified content items
			if datainfo.PrizeNum != giftInfo.PrizeNum {
				// Total number of prizes changed
				giftInfo.LeftNum = datainfo.LeftNum - (datainfo.PrizeNum - giftInfo.PrizeNum)
				if giftInfo.LeftNum < 0 || giftInfo.PrizeNum <= 0 {
					giftInfo.LeftNum = 0
				}
				giftInfo.SysStatus = datainfo.SysStatus
				// Reset the prize cycle information for a prize
				// Potentially Re-update the prize plan according to the prize cycle.
				utils.ResetGiftPrizeData(&giftInfo, c.ServicesGift)
			}
			if datainfo.PrizeTime != giftInfo.PrizeTime {
				// The prize distribution cycle has changed
				utils.ResetGiftPrizeData(&giftInfo, c.ServicesGift)
			}
			c.ServicesGift.Update(&giftInfo, []string{"title", "prize_num", "left_num", "prize_code", "prize_time",
				"img", "displayorder", "gtype", "gdata", "time_begin", "time_end", "sys_updated"})
		} else {
			giftInfo.Id = 0
		}
	}
	if giftInfo.Id <= 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIP(c.Ctx.Request())
		giftInfo.SysCreated = int(time.Now().Unix())
		c.ServicesGift.Create(&giftInfo)
		// Update the prize distribution plan for prizes
		utils.ResetGiftPrizeData(&giftInfo, c.ServicesGift)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetDelete() mvc.Result {

	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		c.ServicesGift.Delete(id) // Delete: Status changed to 1
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetReset() mvc.Result {
	// TODO:
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		c.ServicesGift.Update(&models.LtGift{Id: id, SysStatus: 0}, // Reset: Status=1 represents deleted
			[]string{"sys_status"}) // Update the table to a empty table (with only 'sys_status' as title)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}
