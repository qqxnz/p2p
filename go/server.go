package main

import (
    "fmt"
    "log"
    "net"
    "net/url"
)

var (
    // session存储
    sessionMap map[string]*net.UDPAddr
    listener *net.UDPConn
)

func init() {
    sessionMap = make(map[string]*net.UDPAddr)
}

func main() {
    var (
        err error
    )
    listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9981})
    if err != nil {
        fmt.Println(err)
        return
    }
    log.Printf("本地地址: %s", listener.LocalAddr().String())

    data := make([]byte, 1024)
    for {
        n, remoteAddr, err := listener.ReadFromUDP(data)
        if err != nil {
            log.Printf("error during read: %s", err)
        }
        log.Printf("[%s] %s", remoteAddr.String(), data[:n])

        // 创建连接
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(remoteAddr),
		Port: n,
	})
	if err != nil {
		log.Print("连接失败", err)
		return
	}
	defer conn.Close()

	log.Printf("反向建立连接成功, remote: %s", conn.RemoteAddr().String())

	    // 登录
	    sendMsg("000000", conn)

        // 接收到消息
        request := string(data[:n])

        // 消息处理
        handler(request, remoteAddr)
    }
}

/*
  handler
  @Desc: 消息处理
  @param: request get 请求
  @param: remoteAddr 对端信心
*/
func handler(request string, remoteAddr *net.UDPAddr){
    // 解析消息
    reqUri, err := url.Parse(request)
    if err != nil {
        log.Printf("err-> %v", err)
        listener.WriteToUDP([]byte("url invalid"), remoteAddr)
    }
    log.Printf("uri: %s", reqUri)
    switch reqUri.Path {
    // 登录认证
    case "/login":
        // 存储session信息
        uid := reqUri.Query().Get("uid")
        setSession(uid, remoteAddr)
        // response
        listener.WriteToUDP([]byte("login success"), remoteAddr)
        log.Printf("respnse login success, uid: %s", uid)
        return
    case "/get_user":
        // 获取session信息
        userId := reqUri.Query().Get("user_id")
        session := getSession(userId)
        // response
        listener.WriteToUDP([]byte(session.String()), remoteAddr)
        log.Printf("respnse uid: %s session is %s", userId, session.String())
        return
    case "/heartbeat":
        // 存储session信息
        uid := reqUri.Query().Get("uid")
        setSession(uid, remoteAddr)
        // response
        listener.WriteToUDP([]byte("heartbeat"), remoteAddr)
        log.Printf("respnse heartbeat, uid: %s, session: %s", uid, remoteAddr.String())
        return
    }
    log.Printf("not found services")
}

/*
  setSession
  @Desc: 存储用户session
  @param: uid 用户id
  @param: addr 用户连接地址
*/
func setSession(uid string, addr *net.UDPAddr) {
    if len(uid) <=0 {
        log.Print("ERR: uid invalid")
        return
    }
    sessionMap[uid]= addr
}

/*
  getSession
  @Desc: 获取用户session
  @param: uid
  @return: string
*/
func getSession(uid string) *net.UDPAddr{
    return sessionMap[uid]
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