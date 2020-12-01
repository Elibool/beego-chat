package controllers

import (
	"chat/app/repositories"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"net/http"
)

//后端 websocket 服务端
type ServerController struct {
	beego.Controller
	repositories.ChatRepository
	repositories.MemberRepository
}




func (this *ServerController) WS() {
	id, _ := this.GetInt64("id")
	if id == 0 {
		logs.Error("id is NULL")
		this.Redirect("/?error=id null", 302)
		return
	}

	//检验 http 是否 websocket 请求，且建立 websocket 服务
	conn, err := (&websocket.Upgrader{}).Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if _,ok := err.(websocket.HandshakeError); ok {
		logs.Error("not a websocket connection")
		http.Error(this.Ctx.ResponseWriter, "not a websocket handshake", 400)
		return
	} else if err != nil {
		logs.Error("connection error: ", err)
		return
	}

	//请求成功建立
	var client Client
	client.id = id
	client.conn = conn


	if _, ok := clientList[id]; !ok {
		//如果该id 和链接资源 未在用户链接列表内， 否则该用户加入 wb 组内
		var memberRepository repositories.MemberRepository
		memberRepository.UpOnline(client.id)

		//默认聊天关系，公共
		chatServices[client.id] = 0
		//写入上线系统通知通道
		join <- client
	}

	defer func() {
		logs.Info(client.id, " ---- 下线了")
		//函数返回，该用户退出通道，并断开链接
		//将用户从映射删除
		delete(clientList, client.id)
		//用户状态更新为下线
		var memberRepository repositories.MemberRepository
		memberRepository.DownOnline(client.id)
		client.conn.Close()

		//写入下线消息通知通道
		leave <- client

		//阻止返回
		this.ServeJSON()
	}()

	// 一直读取客户端发送来的消息，直到链接断开
	for {
		//链接建立后，一直 for 读取消息
		_, msgStr , err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		//将字节切片数组通过json 库转结构体
		var acceptMsg = repositories.AcceptMessage{}
		json.Unmarshal(msgStr, &acceptMsg)

		logs.Info("ws 接收消息 :", string(msgStr))

		//获取聊天关系id
		chatServerId := this.GetChatServerId(acceptMsg.AcceptId, client.id)

		//聊天关系映射
		chatServices[client.id] = acceptMsg.AcceptId

		logs.Info("聊天关系建立 :", chatServices)

		//处理消息
		if 1 == acceptMsg.Type {
			//获取历史聊天记录
			historyMessages, num := this.GetHistoryMessage(chatServerId, acceptMsg.AcceptId)
			if num > 0 {
				var chatMessages ChatMessage
				for _, messageItem := range historyMessages {
					chatMessages.Status = 200
					chatMessages.Messages = messageItem.Msg
					//注意接收方
					chatMessages.AcceptId = client.id
					chatMessages.MemberId = messageItem.MemberId
					chatMessages.CreateAt = messageItem.CreatedAt.Format("2006-01-02 15:04:05")

					//logs.Info("历史聊天信息 :", chatMessages)
					chatMessagesChan <- chatMessages
				}
			}
		} else if 2 == acceptMsg.Type || 3 == acceptMsg.Type {
			//在线聊天
			var chatMessages ChatMessage
			message := this.InsertMessageLog(chatServerId, client.id, acceptMsg.Msg)
			chatMessages.Status = 201
			chatMessages.Messages = message.Msg
			chatMessages.AcceptId = acceptMsg.AcceptId
			chatMessages.MemberId = client.id
			chatMessages.CreateAt = message.CreatedAt.Format("2006-01-02 15:04:05")

			if 3 == acceptMsg.Type {
				//群聊天室
				chatMessages.Status = 202
			}

			chatMessagesChan <- chatMessages
		}


	}
}
