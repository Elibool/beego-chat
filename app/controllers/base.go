package controllers

import (
	"chat/app/repositories"
	"github.com/astaxie/beego"
	"time"
)

type NestPreparer interface {
	NestPrepare()
}

type baseController struct {
	beego.Controller
	repositories.MemberRepository
}

func (this *baseController) Prepare() {
	this.Data["pageStartTime"] = time.Now()
	this.Layout = "layout/container.html"

	////断言判断,是否实现了  NestPreparer 接口
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

