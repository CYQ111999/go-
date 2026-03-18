package process

import (
	"LearnGo/chatroom/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) { //这里来的mes一定是smsMes
	//显示即可
	//第一步先反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
