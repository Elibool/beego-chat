package controllers

import (
	"log"
	"time"
)

type ChatController struct {
	baseController
}

func (this ChatController) NestPrepare() {

}


//前端时间戳转换， 数据库时区错误
func (this *ChatController) Get() {
	//获取随机聊天名
	this.Data["chatInfo"] = this.MemberRepository.GetName(&this.Controller)

	list := this.MemberRepository.GetOnlineList(&this.Controller)

	log.Println(time.Now())
	this.Data["title"] = "chat"
	this.Data["lists"] = list
	this.TplName = "chat/index.tpl"
}
