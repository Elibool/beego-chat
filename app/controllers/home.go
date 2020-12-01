package controllers

type MainController struct {
	baseController
}


func (view *MainController) Get() {
	view.Data["title"] = "go chat"
	view.TplName = "home/index.tpl"
}
