package routers

import (
	"chat/app/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//restful 路由
    beego.Router("/", &controllers.MainController{})
    beego.Router("/chat", &controllers.ChatController{})
	beego.Router("/chat/ws", &controllers.ServerController{}, "get:WS")
}
