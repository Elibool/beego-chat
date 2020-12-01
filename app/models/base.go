package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init()  {
	var mysqlUser  = beego.AppConfig.String("mysqlUser")
	var mysqlPass  = beego.AppConfig.String("mysqlPass")
	var mysqlUrls  = beego.AppConfig.String("mysqlUrls")
	var mysqlDb  = beego.AppConfig.String("mysqlDb")

	var mysqlInfo = mysqlUser + ":" + mysqlPass + "@tcp(" + mysqlUrls + ":3306)/" + mysqlDb + "?charset=utf8&loc=Asia%2FShanghai"

	//启用 debug
	orm.Debug, _ = beego.AppConfig.Bool("mysqlDebug")

	//初始链接
	orm.RegisterDataBase("default", "mysql", mysqlInfo, 30)

	logs.Info("mysql connect init success !")

	orm.RegisterModel(new(Members), new(Chats), new (Messages))
	logs.Info("orm new success!")
}


