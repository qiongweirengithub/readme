package gconnection

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/nats-io/nats.go"
	"gserver.com/gbus"
)

var (
	Msgbus gbus.GBus
)

func Helloworld() {
	fmt.Println("gconn helloworld")
}

type ClientConn struct {
	// 用户与网关的tcp链接
	TcpCon net.Conn
	// 监听服务器帧消息的nat订阅关系，用来用户关闭连接时，解除订阅关系
	FramServer *nats.Subscription
	// 监听nats的hello消息，
	FramHello *nats.Subscription
	// 用户id
	Userid string
	// 指向gconnmanager
	gm *GconManagers
}

func (gc *ClientConn) SendTcp(msg string) error {
	_, err := gc.TcpCon.Write([]byte(msg))
	return err
}

func (gc *ClientConn) Close() {
	userid := gc.Userid
	tcpCon := gc.TcpCon
	tcpCon.Close()
	fmt.Println("close tcp con")
	framServer := gc.FramServer
	framServer.Unsubscribe()
	fmt.Println("close framServer sub")
	framHello := gc.FramHello
	framHello.Unsubscribe()
	fmt.Println("close framHello sub")

	delete(gc.gm.conmap, userid)
	fmt.Println("delte tcp con")
}

type GconManagers struct {
	conmap map[string]*ClientConn
}

func (gm *GconManagers) Init() {
	gm.conmap = make(map[string]*ClientConn)
	Msgbus = gbus.Init("nats://nats01:4333,nats://nats02:4333")
	go gm.tick()
}

func (gm *GconManagers) NewCon(userid string, con net.Conn) (*ClientConn, error) {
	if clentConn, ok := gm.conmap[userid]; ok {
		clentConn.TcpCon = con
		clentConn.Userid = userid
		clentConn.gm = gm
		return clentConn, nil
	} else {
		clentConn := &ClientConn{}
		gm.conmap[userid] = clentConn
	}
	clentConn, ok := gm.conmap[userid]
	if !ok {
		return nil, errors.New("type")
	}
	clentConn.gm = gm
	clentConn.TcpCon = con
	clentConn.Userid = userid
	clentConn.SendTcp("establish succ")

	// 监听游戏服务器消息
	framS, err := gm.frameListener(userid)
	if err != nil {
		fmt.Println("frame listen fail"+userid, err)
		return nil, err
	}
	// 监听其他系统消息
	framH, err := gm.helloListener(userid)
	clentConn.FramServer = framS
	clentConn.FramHello = framH

	return clentConn, nil
}

// 监听此用户关联的游戏服务器消息， 通过conmanager发送到对应用户
func (gm *GconManagers) frameListener(userid string) (*nats.Subscription, error) {
	natsSub, err := Msgbus.Receive("frame_"+userid, userid, func(msg *nats.Msg) {
		frame := msg.Data
		fmt.Println(frame)
		gm.Send2User(userid, string(frame))
	})
	fmt.Println("franme listener succ")
	return natsSub, err
}

func (gm *GconManagers) helloListener(userid string) (*nats.Subscription, error) {

	natsSub, err := Msgbus.Receive("hello_event_"+userid, userid, func(msg *nats.Msg) {
		frame := msg.Data
		fmt.Println("xxxx+" + string(frame))
		gm.Send2User(userid, string(frame))
	})
	fmt.Println("hello listener succ")
	return natsSub, err
}

func (gm *GconManagers) Send2MatchServer(msg string) {
	Msgbus.Send("game_match_event", msg)
}

func (gm *GconManagers) Send2GameServer(tableinfo string, msg string) {
	Msgbus.Send("game_event_"+tableinfo, msg)
}

func (gm *GconManagers) Send2HelloServer(userid string) {
	Msgbus.Send("hello_event_"+userid, userid+"res from connect,time:"+TimeNow())
}

func (gm *GconManagers) Send2User(userid string, msg string) error {
	con := gm.conmap[userid]
	err := con.SendTcp(msg)
	if err != nil {
		fmt.Println("send msg error, userid:"+userid, err)
		return err
	}
	return nil
}

func (gm *GconManagers) tick() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("con tick error:", err)
			go gm.tick()
		}
	}()

	for range time.Tick(5000 * time.Millisecond) {
		// fmt.Println("server ticking")
		for valu := range gm.conmap {
			clientCon := gm.conmap[valu]
			if clientCon != nil {
				clientCon.SendTcp("server tick")
			}
		}
	}
}

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
