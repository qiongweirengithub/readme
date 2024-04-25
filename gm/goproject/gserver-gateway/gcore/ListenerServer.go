package gcore

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	gconnection "gserver-gateway.com/gconn"
)

func ListenerHelloworld() {
	fmt.Println("hello gateway server")
}

type GListenerServer struct {
	gm *gconnection.GconManagers
}

func (gls *GListenerServer) Init(gm *gconnection.GconManagers) {
	gls.gm = gm
}

func (gls *GListenerServer) StartListener(protocel string, addr string) {
	go gls.connListener(protocel, addr)

}

// 监听用户链接事件，并进行establish，进行链接的建立
func (gls *GListenerServer) connListener(protocel string, addr string) {

	listener, err := net.Listen(protocel, addr)
	defer listener.Close()

	if err != nil {
		fmt.Println("listen port 8888 err", err)
		panic(err)
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("con accept fail", err)
			panic(err)
		}
		go gls.establish(con)
	}

}

func (gls *GListenerServer) establish(con net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	for {
		n, err := con.Read(buffer)
		if err != nil {
			fmt.Println("conn read error", err)
			return "", err
		}
		msg := string(buffer[0:n])
		fmt.Println("client msg:" + msg)

		var data map[string]string
		if err := json.Unmarshal(buffer[0:n], &data); err != nil {
			fmt.Println("wrong msg:"+msg, err)
			continue
		}
		msgType := data["type"]
		switch msgType {
		// 链接成功
		case "establish":
			userid := data["userid"]
			// 管理用户链接
			clientConn, err := gls.gm.NewCon(userid, con)
			if err != nil {
				fmt.Println("client:"+userid+" new conn error", err)
				con.Write([]byte("fail"))
				return "fail", err
			}
			// 监听用户输入
			go gls.playEventListener(clientConn)
			return userid, nil
		default:
			fmt.Println("test", con)
		}
	}
}

func (gls *GListenerServer) playEventListener(con *gconnection.ClientConn) error {
	defer con.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := con.TcpCon.Read(buffer)
		if err != nil {
			fmt.Println("conn read error", err)
			return err
		}
		msg := string(buffer[0:n])
		fmt.Println("client msg:" + msg)

		var data map[string]string
		if err := json.Unmarshal(buffer[0:n], &data); err != nil {
			fmt.Println("wrong msg:"+msg, err)
			continue
		}
		msgType := data["type"]
		switch msgType {
		case "match":
			// 发送到匹配服务器
			gls.gm.Send2MatchServer(msg)
		case "event":
			roomid := data["roomid"]
			tableid := data["tableid"]
			// 发送消息到游戏服务器线程
			gls.gm.Send2GameServer(roomid+tableid, msg)
		case "hello":
			userid := data["userid"]
			gls.gm.Send2HelloServer(userid)

		default:
			fmt.Println("test", con)
		}
	}
}

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
