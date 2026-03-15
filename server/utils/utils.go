package utils

import (
	"LearnGo/chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据....")

	// 读取4字节头部
	n, err := this.Conn.Read(this.Buf[:4])
	fmt.Printf("读取头部: 实际读取 %d 字节, 错误: %v\n", n, err)

	if err != nil {
		err = errors.New("read pkg header error: " + err.Error())
		return
	}

	if n != 4 {
		fmt.Printf("头部长度不足: 期望4字节, 实际%d字节\n", n)
		return mes, errors.New("read pkg header error: 头部长度不足")
	}

	// 将4字节转换为消息长度
	pkgLen := binary.BigEndian.Uint32(this.Buf[0:4])
	fmt.Printf("消息长度字段: %d 字节\n", pkgLen)

	// 读取消息体
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	fmt.Printf("读取消息体: 实际读取 %d 字节, 错误: %v\n", n, err)

	if n != int(pkgLen) || err != nil {
		return
	}

	// 打印原始数据
	fmt.Printf("原始消息体: %s\n", string(this.Buf[:pkgLen]))

	// 反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Printf("JSON反序列化失败: %v\n", err)
		fmt.Printf("无法解析的数据: %s\n", string(this.Buf[:pkgLen]))
		return
	}

	fmt.Printf("消息解析成功: 类型=%v\n", mes.Type)
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
