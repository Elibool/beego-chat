package repositories

import (
	"chat/app/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"log"
	"math/rand"
	"time"
)

//随机命名处理件
type MemberRepository struct {
	ChatName string
	models.Members
}

const letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const NameKey = "name"

func init() {
	log.Println("memberRepository init success")
}

//获取所有在线用户名
func (this *MemberRepository) GetOnlineList(app *beego.Controller) []orm.Params {
	o := orm.NewOrm()
	qs := o.QueryTable(this.Members)

	var lists []orm.Params
	qs.Exclude("name", this.ChatName).OrderBy("-last_login_at").Values(&lists)
	for  key , item := range lists {
		color := "green"
		if 1 == key % 3 {
			color = "blue"
		}

		var isOnlien, _ = item["IsOnline"].(int64)
		if isOnlien < 1 {
			color = "aero"
		}

		item["color"] = color
		lists[key] = item
	}

	return lists
}

//获取用户详细信息
func (this *MemberRepository) GetMemberInfo(memberId int64) models.Members  {
	o := orm.NewOrm()
	member := models.Members{Id: memberId}
	o.Read(&member)

	return member
}

//下线处理
func (this *MemberRepository) DownOnline(memberId int64) {
	o := orm.NewOrm()
	member := models.Members{Id: memberId}
	o.Read(&member)
	member.IsOnline = 0
	o.Update(&member)
}

//上线处理
func (this *MemberRepository) UpOnline(memberId int64) {
	o := orm.NewOrm()
	member := models.Members{Id: memberId}
	o.Read(&member)
	member.IsOnline = 1
	o.Update(&member)
}


//获取当前用户名
func (this *MemberRepository) GetName(app *beego.Controller) map[string]interface{} {
	chatName := app.GetSession("chatName")
	if nil == chatName {
		chatName = this.getRandName()
		app.SetSession("chatName", chatName)
	}

	this.ChatName = chatName.(string)


	//数据库查看是否已存在
	o := orm.NewOrm()
	member := models.Members{Name: this.ChatName}

	err := o.Read(&member, "name")
	member.LoginIp = app.Ctx.Input.IP()
	member.IsOnline = 1
	if err == orm.ErrNoRows {
		//插入数据
		id, _ := o.Insert(&member)
		member.Id = id
	} else {
		//更新最后在线时间和 ip 状态
		o.Update(&member)
	}


	//struct 转 map
	var chatInfo = make(map[string]interface{})
	j, _ := json.Marshal(&member)
	json.Unmarshal(j, &chatInfo)

	logs.Info(chatInfo)

	return chatInfo
}

func (member *MemberRepository) getRandName() interface{} {
	b := make([]byte, 5)
	rand.Seed(time.Now().Unix())

	for i := 0; i <= 4; i++ {
		num := rand.Intn(52)
		b[i] = letterBytes[num]
	}
	return string(b)
}
