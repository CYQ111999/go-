package utils

import (
	"LearnGo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据....")
	// 确保读取完整的4字节头部
	n, err := io.ReadFull(this.Conn, this.Buf[:4])
	if err != nil {
		if err == io.EOF {
			return mes, errors.New("客户端已关闭连接")
		}
		err = errors.New("read pkg header error: " + err.Error())
		return
	}
	if n != 4 {
		return mes, errors.New("读取头部长度不足")
	}
	// 根据buf 转化成一个uint32类型
	pkgLen := binary.BigEndian.Uint32(this.Buf[0:4])
	// 确保读取完整的消息体
	n, err = io.ReadFull(this.Conn, this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		if err == io.EOF {
			return mes, errors.New("客户端在发送数据过程中关闭连接")
		}
		return mes, errors.New("读取消息体错误: " + err.Error())
	}
	// 把pkgLen反序列化成message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给客户端
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	//发送data数据本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	return
}
