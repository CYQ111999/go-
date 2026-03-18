package process

import (
	"LearnGo/chatroom/client/utils"
	"LearnGo/chatroom/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

// SendGroupMes 发送群聊的消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//创建一个smsMes 实例
	var smsMes message.SmsMes
	smsMes.Content = content //内容
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fial =", err.Error())
		return
	}
	mes.Data = string(data)
	//对mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fial =", err.Error())
		return
	}

	//将mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes WritePkg fial =", err.Error())
		return
	}
	return
}
