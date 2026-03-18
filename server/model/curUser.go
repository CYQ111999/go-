package model

import (
	"LearnGo/chatroom/common/message"
	"net"
)

// CurUser 因为在客户端有很多地方会使用到curUser,我们将其做成全局
type CurUser struct {
	Conn net.Conn
	message.User
}
