package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
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

// ✅ 只改了这里：json:"userId"
type LoginMes struct {
	UserId   int    `json:"userId"` // 改为小写
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

// ✅ 只改了这里：json:"userId" 和 json:"status"
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` // 改为小写
	Status int `json:"status"` // 改为小写
}

// SmsMes 增加一个SmsMes
type SmsMes struct {
	Content string `json:"content"` //内容
	User           //匿名结构体，继承
}
