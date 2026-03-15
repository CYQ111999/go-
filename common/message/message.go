package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
)

// 这里我们定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"`
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code    int    `json:"code"` //返回状态码
	UsersId []int  `json:"userIds"`
	Error   string `json:"error"` //返回错误信息
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码400表示该用户已经占有200表示注册成功
	Error string `json:"error"` //返回错误信息
}

// NotifyUserStatusMes 为了配合服务器端推送用户状态变化消息通知
type NotifyUserStatusMes struct {
	UserId int `json:"UserId"`
	Status int `json:"Status"` //用户状态
}
