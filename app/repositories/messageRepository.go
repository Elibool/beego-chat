package repositories

type MessageRepository struct {
}

//定义接收结构体, 由于是 html 端传入，所有 value 都转为字符串类型
type AcceptMessage struct {
	AcceptId int64  `json:"accept_id"`
	Type     int64  `json:"type"`
	Msg      string `json:"msg"`
}

