package repositories

import (
	"chat/app/models"
	"github.com/astaxie/beego/orm"
)

type ChatRepository struct {
}

type ChatInfo struct {


}

//录入聊天数据
func (this *ChatRepository) InsertMessageLog(chatId int64,memberId int64, msg string) models.Messages  {
	o := orm.NewOrm()
	var message models.Messages
	message.ChatId = chatId
	message.MemberId = memberId
	message.Msg = msg
	message.IfRead = 0

	o.Insert(&message)

	return message
}


//获取历史聊天记录 5 条, 公共聊天室接收 100 条历史消息
func (this *ChatRepository) GetHistoryMessage(chatId int64, acceptId int64) ([]models.Messages, int64) {
	var message []models.Messages
	o := orm.NewOrm()

	limit := "100"
	if acceptId > 0 {
		limit = "5"
	}
	num, _ := o.Raw("select * from messages where chat_id = ? order by id desc limit "+limit, chatId).QueryRows(&message)

	return message, num
}

//返回聊天关系 id
func (this *ChatRepository) GetChatServerId(acceptId int64, clientId int64) int64 {
	if 0 == acceptId {
		//聊天室群聊
		return 0
	}

	//建立 1 对 1 聊天关系
	o := orm.NewOrm()
	var chat models.Chats
	o.Raw("select * from chats where (member_a = ?  and  member_b = ?) or (member_a = ?  and  member_b = ?)", acceptId, clientId, clientId, acceptId).QueryRow(&chat)
	if 0 == chat.Id {
		//录入数据
		chat.MemberA = acceptId
		chat.MemberB = clientId

		id, _ := o.Insert(&chat)
		return id
	} else {
		o.Update(&chat)
		return chat.Id
	}
}
