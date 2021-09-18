package main

import (
	"log"
	"net"
	"time"
)

/**
    @date: 2021/9/10
**/

func main() {
	// 创建连接
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(113, 97, 31, 163),
		Port: 64386,
	})
	if err != nil {
		log.Print("连接失败", err)
		return
	}
	defer conn.Close()

	log.Printf("建立连接成功, remote: %s", conn.RemoteAddr().String())

	// 登录
	sendMsg("/login?uid=1", conn)
	time.Sleep(time.Duration(1) * time.Second)
	// 获取指定用户session信息
	sendMsg("/get_user?user_id=1", conn)
	for {
		// heartbeat
		sendMsg("/heartbeat?uid=1", conn)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

/*
  sendMsg
  @Desc: 发送消息
  @param: socket 远端连接
*/
func sendMsg(msg string, conn *net.UDPConn){
	// 发送数据
	senddata := []byte(msg)
	_, err := conn.Write(senddata)
	if err != nil {
		log.Printf("发送数据失败: err-> %v", err)
		return
	}

	// 接收数据
	data := make([]byte, 1024)
	_, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Printf("读取数据失败: err-> %v", err)
		return
	}
	log.Printf("addr: %s, msg: %s",remoteAddr, string(data))
}