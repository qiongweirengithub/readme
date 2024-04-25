package mynatesclient

import (
	"fmt"
	"time"

	nats "github.com/nats-io/nats.go"
)

func Helloworld() {
	fmt.Println("hello nates client")
}

type MyClient struct {
	nates_server string
	cn           *nats.Conn
}

func (c *MyClient) Helloworld() {
	fmt.Println("my client hello")
}

func (c *MyClient) Init(nates_server string) {
	c.nates_server = nates_server
	cn, err := nats.Connect(nates_server)
	if err != nil {
		fmt.Println("nat cient connect err")
		panic(err)
	}
	fmt.Println("nat cient connect succ:" + nates_server)
	c.cn = cn
}

func (c *MyClient) Receive(topic string, cb nats.MsgHandler) {
	c.cn.Subscribe(topic, cb)
}

func (c *MyClient) Send(topic string, data string) {
	// 不需要消费端返回结果
	err := c.cn.Publish(topic, []byte(data))
	if err != nil {
		fmt.Println("publish err:"+data, err)
	}
}

func (c *MyClient) Request(topic string, data string) {
	// 需要消费返回结果，等于是同步请求
	_, err := c.cn.Request(topic, []byte(data), 20*time.Millisecond)
	if err != nil {
		fmt.Println("send err:"+data, err)
	}
}
