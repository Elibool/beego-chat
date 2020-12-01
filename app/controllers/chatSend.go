package controllers

import (
	"chat/app/models"
	"chat/app/repositories"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"log"
)

//消息消费模块， 无关 socket 请求

type Client struct {
	conn *websocket.Conn //用户 websocket 链接
	id   int64           //用户id
}

//普通系统单条消息
type Message struct {
	EventType byte           `json:"type"`    // 0 发布消息， 1 用户进入， 2 用户退出， 3有新消息
	Id        int64          `json:"id"`      //用户名称
	Message   string         `json:"message"` //消息内容
	Member    models.Members `json:"member"`
	AcceptId  int64			`json:"accept_id"`
}

//聊天消息，多条
type ChatMessage struct {
	Status   int    `json:"status"`
	AcceptId int64  `json:"accept_id"`
	MemberId int64  `json:"member_id"`
	Messages string `json:"messages"`
	CreateAt string `json:"create_at"`
}

var (
	// 此处要设置有缓冲的通道。因为这是goroutine自己从通道中发送并接受数据。
	// 若是无缓冲的通道， 该goroutine发送数据到通道后就被锁定，
	// 需要数据被接受后才能解锁， 而恰恰接受数据的又只能是它自己
	join             = make(chan Client, 10)      // 用户加入通道
	leave            = make(chan Client, 10)      // 用户退出通道
	message          = make(chan Message, 10)     // 系统通知消息通道
	chatMessagesChan = make(chan ChatMessage, 20) // 聊天消息通道

	//用户映射,用户 conn 列表
	clientList = make(map[int64]Client)

	//对应聊天人映射表
	chatServices = make(map[int64]int64)
)

func init() {
	//必须启动单独 go 程消费该通道
	go broadcaster()
	log.Println(" broadcaster init success ")
}

//后端广播功能 ， 发消息， 用户加入， 用户退出三种情况广播给所有用户 ; 后两种转换后，以第一种方式发送
func broadcaster() {
	for {
		//
		select {
		case msg := <-chatMessagesChan:
			//有新消息
			str := fmt.Sprintf("聊天类型码: %d , 信息 message: %s \n", msg.Status, msg.Messages)
			logs.Info(str)

			//将 struct 个数据转为 json 数据
			data, err := json.Marshal(msg)
			if err != nil {
				logs.Error("fail to json message: ", err)
				return
			}

			if 200 == msg.Status {
				//返回历史消息,仅发送给请求人
				client, _ := clientList[msg.AcceptId]
				if client.conn.WriteMessage(websocket.TextMessage, data) != nil {
					logs.Error("fail to write chat messages to socket")
				}
			} else if 201 == msg.Status {
				//单对单聊天
				memberClient, _ := clientList[msg.MemberId]
				if memberClient.conn.WriteMessage(websocket.TextMessage, data) != nil {
					//发送给自己
					logs.Error("member -  fail to write chat messages to socket")
				}

				//接收人是否在线
				if acceptClient, ok := clientList[msg.AcceptId]; ok {
					//当前聊天人是否为发送人
					if chatActiveId := chatServices[msg.AcceptId]; chatActiveId == msg.MemberId {
						if acceptClient.conn.WriteMessage(websocket.TextMessage, data) != nil {
							logs.Error("accept - fail to write chat messages to socket")
						}
					} else {
						//未在同个聊天频道， 发送系统消息
						var msg Message
						msg.Id = memberClient.id
						msg.EventType = 3
						msg.Message = "未在同聊天频道消息通知"
						msg.AcceptId = acceptClient.id
						msg.Member = models.Members{}

						message <- msg
					}
				}

			} else if 202 == msg.Status {
				//群聊， 所有在线的都发送
				for _, client := range clientList {
					//发送数据
					if client.conn.WriteMessage(websocket.TextMessage, data) != nil {
						logs.Error("public - fail to write chat messages to socket")
					}
				}
			}
		case msg := <-message:
			//有新系统消息
			str := fmt.Sprintf("用户: %d send message: %s \n", msg.Id, msg.Message)
			logs.Info(str)
			//将该系统消息发送给所有用户
			for _, client := range clientList {
				//将 struct 个数据转为 json 数据
				data, err := json.Marshal(msg)
				if err != nil {
					logs.Error("fail to json message: ", err)
					return
				}

				//发送数据
				if client.conn.WriteMessage(websocket.TextMessage, data) != nil {
					logs.Error("fail to write message to socket")
				}
			}
		case client := <-join:
			//有新的聊天人加入，上线通知
			str := fmt.Sprintf("用户: %d  上线, \n", client.id)
			logs.Info(str)

			//用户加入在线 ws 组
			clientList[client.id] = client

			//将用户消息放入消息通道
			var msg Message
			msg.Id = client.id
			msg.EventType = 1
			msg.Message = fmt.Sprintf("用户: %d 上线 , 有 %d 在线", client.id, len(clientList))

			var memberRepository repositories.MemberRepository
			msg.Member = memberRepository.GetMemberInfo(client.id)

			//写入通道， 待下次 for 处理（所以通道必须是缓冲通道， 否则写入通道后就锁定该 go 程，待读取后才继续下次 for 循环）
			message <- msg

		case client := <-leave:
			//用户退出
			str := fmt.Sprintf("用户下线: %d   \n", client.id)
			logs.Info(str)

			var msg Message
			msg.Id = client.id
			msg.EventType = 2
			msg.Message = fmt.Sprintf("用户下线: %d , 还有 %d 人在线", client.id, len(clientList))
			var memberRepository repositories.MemberRepository
			msg.Member = memberRepository.GetMemberInfo(client.id)

			//待下个 for 循环消费
			message <- msg
		}
	}
}
