package process2

import (
	"LearnGo/chatroom/common/message"
	"LearnGo/chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

// SendGroupMes 写方法转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 取出mes的内容SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	// 重新构建完整的消息
	var forwardMes message.Message
	forwardMes.Type = message.SmsMesType // 设置消息类型

	// 序列化 smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	forwardMes.Data = string(data)

	// 序列化完整消息
	finalData, err := json.Marshal(forwardMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	// 转发给其他用户
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(finalData, up.Conn) // 发送完整消息
	}
}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
